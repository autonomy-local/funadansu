# 鍵の管理（sops と age）

フェーズ0の 5.2（[Issue #18](https://github.com/autonomy-local/funadansu/issues/18)）の手引きです。

> **方針の改訂（[ADR-0003](../decisions/0003-リポジトリと秘密情報とデプロイ.md)、採用）**：秘密情報の正本は、非公開のデプロイ用リポジトリの `secrets/` の 1 か所です。本体（公開）には、暗号化した秘密情報も `.sops.yaml` の実在の宛先も置きません。下の手順は、移行の作業（dev のダミーを確かめる手順）として残しています。公開の前に、`secrets/` と `.sops.yaml` は非公開のデプロイ用リポジトリへ移します（[公開の前の確認](public-repo-checklist.md)）。

## 仕組み

| 場所 | 中身 |
| --- | --- |
| `.sops.yaml` | 暗号化の規則。`secrets/` の下の `*.yaml` を、age の公開鍵で暗号化する。公開鍵なので、リポジトリに入れてよい |
| `secrets/dev.yaml` | dev の秘密情報（sops で暗号化したもの）。**正本** |
| 鍵ファイル（age の秘密鍵） | リポジトリに入れない。既定の場所は `~/.config/sops/age/keys.txt`（`SOPS_AGE_KEY_FILE` で変えられる） |

暗号化した中身は、鍵を持つ人が `sops` で復号します。復号した値は、ファイルに書き出さず、必要な実行環境に直接渡します。

## 復号する（鍵を持つ人）

`nix develop` の中で実行します（`sops` と `age` は開発シェルに入っています）。

```sh
# 鍵ファイルを置く（初回だけ）。秘密鍵は、この場所から出さない
mkdir -p ~/.config/sops/age && cp /安全な場所/keys.txt ~/.config/sops/age/keys.txt && chmod 600 ~/.config/sops/age/keys.txt

sops -d secrets/dev.yaml            # 復号して標準出力に出す（画面に秘密が出るので、人のいない場所では使わない）
sops -d --extract '["jwt_key"]' secrets/dev.yaml   # 1 つの値だけ取り出す
```

## 値を変える（鍵を持つ人）

```sh
sops secrets/dev.yaml               # エディタで開き、保存すると暗号化し直す
```

## 鍵の持ち主を増やす・入れ替える

1. 新しい人の公開鍵（`age1…`）を `.sops.yaml` の `age:` に、カンマ区切りで足す
2. `sops updatekeys secrets/dev.yaml` で、暗号化の宛先を更新する（既存の鍵を持つ人が実行する）
3. 外れた人の鍵で、秘密情報そのものを変えたいときは、値を入れ替えてから `updatekeys` を実行する

## 含むもの（dev のダミー）

| キー | 中身 |
| --- | --- |
| `jwt_kid` | JWT の鍵の識別子（dev のダミー）。proxy の `FUNADANSU_JWT_KID` に渡す |
| `jwt_key` | JWT の検証に使う鍵（dev のダミー）。proxy の `FUNADANSU_JWT_KEY` に渡す |

**これはダミーです。本番の鍵は、この仕組みには入れません**（本番の鍵の扱いは、本番の環境を作るときに決めます）。

## 注意

- 平文の鍵・パスワード・トークンを、コミット、Issue、PR、コミットメッセージに書かない（[AGENTS.md](../../AGENTS.md)）
- 復号した値を、ログや CI の出力に出さない
- 鍵ファイルを失くしたら、`secrets/dev.yaml` は復号できなくなる。鍵の持ち主を 2 人以上にしておく
- 漏えいが疑われたら、[運用の手引き](README.md) の鍵の入れ替えの手順に従う（dev のダミーでも、同じ手順で値を変える）
