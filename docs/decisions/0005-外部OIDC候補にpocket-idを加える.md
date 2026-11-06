# ADR-0005：外部 OIDC の候補として pocket-id を加える（組み込まず、別に運用する）

| 項目 | 内容 |
| --- | --- |
| 状態 | accepted（採用） |
| 日付 | 2026-11-05（検証結果を 2026-11-06 に追記） |
| 関連 Issue | なし |
| 関連 ADR | ADR-0004（決定 2 の候補に加える。0004 の決定は変えない）、ADR-0001（依存の最小化） |

## 背景（Context）

ADR-0004 は、ユーザー認証を外部の OIDC プロバイダーに任せ、基盤は検証だけを行うと決めました。自前で持ちたい地域のために、OSS のプロバイダーを候補に挙げています（最初は Keycloak）。

OSS の候補として pocket-id（https://github.com/pocket-id/pocket-id、BSD 2-Clause）を検討しました。pocket-id はパスキーだけでログインする、軽量な OIDC プロバイダーです。調査の結果は [docs/decisions/0005 の付録](#付録調査の要点) のとおりです。

## 決定（Decision）

1. **pocket-id は Funadansu に組み込まない**。コードを取り込まず、外部の OIDC プロバイダーとして別に運用する
2. ADR-0004 の決定 2 の候補に pocket-id を加える。Firebase Authentication、Keycloak と並べて、地域ごとの設定で選べるようにする
3. 検証は ADR-0004 の決定 3・4 のとおり、Go の標準ライブラリで行う。pocket-id 専用の検証の実装を、プロバイダーの実装として1つ足す（既存の実装は変えない）
4. パスキーは pocket-id 側で扱う。Funadansu は WebAuthn を持たない

## 結果（Consequences）

- Funadansu の依存は増えない。pocket-id のコードも NOTICE の表示も、この決定の時点では不要
- pocket-id を運用する地域は、パスキーの復旧手順と、ドメインの設計（RP ID の範囲）を運用の手引きに書く必要がある
- pocket-id は既定で `email` や `name` をクレームに含む。ADR-0004 の決定 6（個人情報を預けない）に合わせ、スコープとクレームを絞る設定が必要
- 組み込みを選ぶ場合は、この ADR を覆す新しい ADR が必要（依存の増加が ADR-0001 に反する）

## 検証結果（2026-11-06、pocket-id 2.19.2 を手元で起動して確かめた）

ADR-0004 の確認事項に対する結果です。確認した範囲は、ローカルの pocket-id（SQLite、HTTP）と、テスト用のクライアント（PKCE 有効）です。

1. **検証の共通部分**：discovery と JWKS を取得できる。署名方式は RS256（`id_token_signing_alg_values_supported` は RS256 のみ）。RS256 の署名、`kid` の照合、`iss` の完全一致、`exp` を確かめるだけで、検証できた。`alg` は許可リスト（RS256・ES256）で照合できる
2. **旧 API の窓口の `sid`（ADR-0001）**：ID トークンにもアクセストークンにも `sid` がない。旧 API の窓口の例外を採るなら、`sid` は pocket-id から得られない。例外を採るかどうかは、この結果を前提に決める
3. **個人情報のクレーム**：`scope=openid profile email` では、ID トークンに `email`、`name`、`given_name`、`family_name`、`display_name`、`picture`、`preferred_username` が入る。`scope=openid` だけでは、これらは入らず、`sub`・`amr`・`auth_time`・`nonce` などに絞れる。Funadansu は `openid` だけを要求すれば、ADR-0004 の決定 6 に合う
4. **トークンの種類**：ID トークンは `typ: JWT`、アクセストークンは `typ: at+jwt` で、ヘッダーで見分けられる。ID トークンの `aud` は client_id、`nonce` あり。アクセストークンの `aud` は client_id と issuer の両方を含み、`client_id`・`scp`・`scope` を持つ。Funadansu は、トークンを受け取る経路（ID トークンかアクセストークンか）ごとに、`aud` と `typ` の期待値を決める
5. **`sub`**：利用者の UUID（ログインのたびに同じ値）。`(issuer, sub)` から内部 ID を引く方式（ADR-0004 の決定 5）と合う

まだ確かめていないこと：CockroachDB での動作、本番に近い HTTPS の設定（今回は HTTP）、refresh token の扱い、パスキーの復旧手順。

## 付録：調査の要点

- ライセンス：pocket-id は BSD 2-Clause。主な依存は MIT、BSD 系、Apache-2.0。SBOM での確認は未実施（フロントエンドと推移的依存）
- 技術：gin、GORM、fosite のフォーク、埋め込みの SvelteKit。Funadansu の方針（標準ライブラリ、sqlc、DI なし）と違う。組み込むと依存が大きく増える
- 署名：RS256 が既定。ES256 や EdDSA にも対応。Discovery あり。`sub` は利用者の UUID
- 詳しい調査の記録は、プロジェクトの共有フォルダの調査メモ（`pocket-id-review.md`）にある
