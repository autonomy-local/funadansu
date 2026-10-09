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
      in
      {
        packages = { inherit db-start db-stop db-init; };

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
