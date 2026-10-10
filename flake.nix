{
  description = "Funadansu development shell";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};

        # ローカルの PostgreSQL。データは flake.nix のあるディレクトリの .data/postgres に置く。
        # 接続は localhost の TCP だけ（Unix ソケットは使わない）。認証は trust。
        dbEnv = ''
          root="$PWD"
          while [ "$root" != "/" ] && [ ! -f "$root/flake.nix" ]; do
            root="$(dirname "$root")"
          done
          if [ ! -f "$root/flake.nix" ]; then
            echo "flake.nix が見つかりません。リポジトリの中で実行してください。" >&2
            exit 1
          fi
          export PGDATA="''${PGDATA:-$root/.data/postgres}"
          export PGHOST="''${PGHOST:-localhost}"
          export PGPORT="''${PGPORT:-5432}"
          export PGUSER="''${PGUSER:-postgres}"
          export PGDATABASE="''${PGDATABASE:-funadansu}"
        '';

        db-start = pkgs.writeShellApplication {
          name = "db-start";
          runtimeInputs = with pkgs; [ postgresql coreutils ];
          text = dbEnv + ''
            if [ ! -f "$PGDATA/PG_VERSION" ]; then
              mkdir -p "$PGDATA"
              initdb --pgdata="$PGDATA" --username="$PGUSER" --auth=trust \
                --encoding=UTF8 --locale=C >/dev/null
            fi
            if pg_ctl status --pgdata="$PGDATA" >/dev/null 2>&1; then
              echo "PostgreSQL はすでに起動しています（$PGHOST:$PGPORT）。"
            else
              pg_ctl start --pgdata="$PGDATA" --log="$PGDATA/postgres.log" --wait \
                -o "-c listen_addresses=localhost -c port=$PGPORT -c unix_socket_directories=" >/dev/null
              echo "PostgreSQL を起動しました（$PGHOST:$PGPORT）。"
            fi
            if [ -z "$(psql --dbname=postgres --tuples-only --no-align \
              --command="SELECT 1 FROM pg_database WHERE datname = '$PGDATABASE'")" ]; then
              createdb "$PGDATABASE"
            fi
            echo "データベース: $PGDATABASE　ユーザー: $PGUSER　データ: $PGDATA"
          '';
        };

        db-stop = pkgs.writeShellApplication {
          name = "db-stop";
          runtimeInputs = with pkgs; [ postgresql coreutils ];
          text = dbEnv + ''
            if pg_ctl status --pgdata="$PGDATA" >/dev/null 2>&1; then
              pg_ctl stop --pgdata="$PGDATA" --wait >/dev/null
              echo "PostgreSQL を停止しました。"
            else
              echo "PostgreSQL は起動していません。"
            fi
          '';
        };

        # db/<スキーマ>/*.sql を順に流し、旧スキーマ（テーブルは空）を作る。
        # 各 DDL は IF NOT EXISTS なので、何度流しても安全。
        db-init = pkgs.writeShellApplication {
          name = "db-init";
          runtimeInputs = with pkgs; [ postgresql coreutils findutils ];
          text = dbEnv + ''
            if ! pg_isready --quiet; then
              echo "PostgreSQL が起動していません。先に db-start を実行してください。" >&2
              exit 1
            fi
            # DDL の IF NOT EXISTS が出す NOTICE（既存の表の表示）は抑える。
            export PGOPTIONS="-c client_min_messages=warning"
            find "$root/db" -name '*.sql' | sort | while read -r file; do
              psql --set=ON_ERROR_STOP=1 --quiet --file="$file" >/dev/null
            done
            count="$(psql --tuples-only --no-align --command="SELECT count(*) FROM information_schema.schemata WHERE schema_name LIKE 'pxr\_%'")"
            echo "旧スキーマ：$count 個（データベース $PGDATABASE）"
          '';
        };

        # Go のサービス（bootstrap/）。標準ライブラリと pgx だけで、静的なバイナリを作る。
        bootstrap = pkgs.buildGoModule {
          pname = "funadansu-bootstrap";
          version = "0.1.0";
          src = ./bootstrap;
          vendorHash = "sha256-8H2nyRtOHNVngFc77SGccsEWI29adbKWTzL9XfPTX5w=";
          subPackages = [ "cmd/funadansu" ];
          env.CGO_ENABLED = "0";
          ldflags = [ "-s" "-w" ];
        };

        # コンテナ（Go）。root で動かさず、待ち受けは 0.0.0.0:8080。
        # 証明書は cacert（PostgreSQL へ TLS で接続するとき、証明書を検証するため）。
        bootstrap-image = pkgs.dockerTools.buildLayeredImage {
          name = "funadansu-bootstrap";
          contents = [ pkgs.cacert ];
          config = {
            Entrypoint = [ "${bootstrap}/bin/funadansu" ];
            Env = [
              "FUNADANSU_ADDR=0.0.0.0:8080"
              "SSL_CERT_FILE=${pkgs.cacert}/etc/ssl/certs/ca-bundle.crt"
            ];
            ExposedPorts = { "8080/tcp" = { }; };
            User = "65532:65532";
          };
        };

        # proxy（Hono）の依存を入れた作業ディレクトリ。実行時の依存（hono）だけを入れる。
        proxy-app = pkgs.buildNpmPackage {
          pname = "funadansu-proxy";
          version = "0.1.0";
          src = ./bootstrap/proxy;
          npmDepsHash = "sha256-Qj1TSyy3U2ViDbwgkYupQclg8MaBeX9ykGCyrGVUPSY=";
          npmInstallFlags = [ "--omit=dev" ];
          dontNpmBuild = true;
          installPhase = ''
            runHook preInstall
            mkdir -p $out/app
            cp -r src package.json node_modules $out/app/
            runHook postInstall
          '';
        };

        # コンテナ（Bun）。Bun が TypeScript をそのまま動かす（src/entry/bun.ts）。
        proxy-image = pkgs.dockerTools.buildLayeredImage {
          name = "funadansu-proxy";
          contents = [ proxy-app ];
          config = {
            Entrypoint = [ "${pkgs.bun}/bin/bun" "src/entry/bun.ts" ];
            WorkingDir = "/app";
            Env = [ "FUNADANSU_PROXY_ADDR=0.0.0.0:8787" ];
            ExposedPorts = { "8787/tcp" = { }; };
            User = "65532:65532";
          };
        };
      in
      {
        packages = {
          inherit db-start db-stop db-init;
          inherit bootstrap bootstrap-image proxy-app proxy-image;
        };

        # db-start で起動し、接続でき、db-stop で止まることを確かめる。
        checks.db = pkgs.runCommand "funadansu-db-check"
          { nativeBuildInputs = [ db-start db-stop pkgs.postgresql ]; }
          ''
            touch flake.nix
            export PGPORT=54329
            db-start
            db-start
            test "$(psql --host=localhost --username=postgres --dbname=funadansu \
              --tuples-only --no-align --command='SELECT current_database()')" = funadansu
            db-stop
            ! pg_isready --host=localhost --port=$PGPORT
            db-stop
            touch $out
          '';

        # db-init で旧スキーマ（15スキーマ、108 表）が立ち、表が空であることを確かめる。
        checks.db-init = pkgs.runCommand "funadansu-db-init-check"
          { nativeBuildInputs = [ db-start db-init db-stop pkgs.postgresql ]; }
          ''
            cp -r ${self} work && chmod -R u+w work && cd work
            export PGPORT=54330
            db-start
            db-init
            test "$(psql --tuples-only --no-align --command="SELECT count(*) FROM information_schema.schemata WHERE schema_name LIKE 'pxr\_%'")" = 15
            test "$(psql --tuples-only --no-align --command="SELECT count(*) FROM information_schema.tables WHERE table_schema LIKE 'pxr\_%'")" = 108
            test "$(psql --tuples-only --no-align --command="SELECT COALESCE(SUM(n_live_tup), 0) FROM pg_stat_user_tables WHERE schemaname LIKE 'pxr\_%'")" = 0
            db-stop
            touch $out
          '';

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            bun
            nodejs
            wrangler
            pulumi
            sqlc
            postgresql
            sops
            age
            # API のシナリオ（test/api/）。版は flake.lock の nixpkgs で固定する。
            runn
            curl
            db-start
            db-stop
            db-init
          ];

          # psql などが、引数なしでローカルの PostgreSQL につながるようにする。
          PGHOST = "localhost";
          PGUSER = "postgres";
          PGDATABASE = "funadansu";

          shellHook = ''
            echo "Funadansu dev shell activated."
          '';
        };
      }
    );
}
