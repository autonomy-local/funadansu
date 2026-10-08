# 単位（ドメイン）の一覧

移行作業は、1つの作業で1つの単位を扱います。単位ごとの引き継ぎファイルは、[_template.md](_template.md) をもとに `docs/units/<単位の名前>.md` として作ります。

上から順に移します。依存の少ない単位から先に移し、移した単位は、それ以降の単位から直接呼べるようにします。

| 順 | 単位 | 旧実装（サービス） | 引き継ぎファイル |
| --- | --- | --- | --- |
| 1 | proxy（入口、認証） | block-proxy | [proxy.md](proxy.md) |
| 2 | operator（認証・権限） | operator | [operator.md](operator.md) |
| 3 | access-control（アクセス認可） | access-control-manage、access-control | [access-control.md](access-control.md) |
| 4 | catalog（カタログ） | catalog、catalog-update | [catalog.md](catalog.md) |
| 5 | notification（通知） | notification | [notification.md](notification.md) |
| 6 | identity-verify（本人確認） | identity-verificate | [identity-verify.md](identity-verify.md) |
| 7 | certificate（証明書） | certification-authority、certificate-manage | [certificate.md](certificate.md) |
| 8 | binary（バイナリ） | binary-manage | [binary.md](binary.md) |
| 9 | ctoken（CToken 台帳） | ctoken-ledger、local-ctoken | [ctoken.md](ctoken.md) |
| 10 | book-manage（データ通帳の管理） | book-manage | [book-manage.md](book-manage.md) |
| 11 | book-operate（データ通帳の操作） | book-operate | [book-operate.md](book-operate.md) |

- 2つの旧サービスが同じテーブルを持っている単位（証明書、アクセス認可など）は、1つの単位にまとめて移します
- 本人確認・認証・認可は、別の単位のまま保ちます。統合のときにも1つの単位に混ぜません
- 全体の状況は [MIGRATION.md](../../MIGRATION.md) を見てください
