# 記録する ID（見本）

このファイルは、非公開のデプロイ用リポジトリの `ids.md` にコピーして使います。本体には置きません。

| 項目 | 値 | 備考 |
| --- | --- | --- |
| GCP のプロジェクト ID | （記入） | Cloud Run、Artifact Registry、Secret Manager |
| Cloud KMS のキーリング、鍵 | （記入） | sops の本番の宛先、Pulumi の秘密 |
| GCS の状態バケット | （記入） | Pulumi の状態（バージョン管理を有効に） |
| Cloudflare のアカウント ID | （記入） | Workers |
| Cloudflare の API トークンの名前 | （記入。値は書かない） | 最小の権限 |
| 鍵の持ち主（dev、本番） | （記入） | 2 人以上（dev）。本番はオフラインの鍵の保管場所 |
