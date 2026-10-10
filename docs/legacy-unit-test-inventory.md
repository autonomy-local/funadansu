# 旧 E2E 以外の試験の一覧（単体テストへ回すもの）

旧基盤の各サービス（`pxr-*-service`）の試験のうち、supertest を使うものを、E2E（runn）に入れるかどうかで分けた一覧です。作ったのは、[差分台帳](spec-deviations.md) の D-009〜D-017 の対応表でもあります。

## 分類の決め方

- **単体テストへ（D-0xx）**：スタブ（`Stub…`）、`jest.mock`、`jest.spyOn` のどれかを使う試験。下流を差し替えた結果を見るため、E2E では意味が変わります（AGENTS.md の互換性）。Go 側の単体テストに写します。
- **候補（未確認）**：supertest を使い、スタブもモックも使わない試験。下流の実物に依る可能性があるため、旧実装を起動して GREEN になったものだけを E2E に入れる、と決めます。まだ確かめていません。
- **対象外（supertest なし）**：supertest を使わない試験。今回の一覧の対象外です。

件数は、`test(` と `it(` を単語として数えた概数です。ファイルの中身は読んでいません。依存の判定は、名前とインポートの文字列による機械的なものです。

## 単位ごとの集計

| 単位 | 差分台帳 | 試験（件） | 単体テストへ（件） | 候補（未確認、件） |
| --- | --- | --- | --- | --- |
| access-control | D-009 | 321 | 184 | 137 |
| catalog | D-010 | 2,084 | 1,214 | 870 |
| notification | D-011 | 86 | 55 | 31（確認済み。E2E に入れた） |
| identity-verify | D-012 | 177 | 177 | 0 |
| certificate | D-013 | 256 | 256 | 0 |
| binary | D-014 | 127 | 127 | 0 |
| ctoken | D-015 | 442 | 442 | 0 |
| book-manage | D-016 | 2,785 | 1,920 | 4（確認済み。E2E に入れた） |
| book-operate | D-017 | 1,677 | 1,677 | 0 |

単位 proxy は [proxy.md](units/proxy.md)（D-008、350 件すべて単体テストへ）、operator は [operator.md](units/operator.md)（D-006、D-007）で、この一覧の対象外です。

## サービスごとの試験

### pxr-access-control-manage-service（単位：access-control）

試験ファイル 12、試験 200 件。単体テストへ 125 件、E2E の候補（未確認）75 件

| 試験ファイル | 件数 | 分類 | 依存 |
| --- | --- | --- | --- |
| `src/tests/CreateAPIKey.validator.spec.ts` | 75 | E2E（確認済み、GREEN） | なし |
| `src/tests/CreateActorAPIKey.spec.ts` | 6 | 単体テストへ（D-009） | スタブ: StubAccessControlServer, StubCatalogServer, StubOperatorServer 他 |
| `src/tests/CreateAppWfUserAPIKey.spec.ts` | 6 | 単体テストへ（D-009） | スタブ: StubAccessControlServer, StubCatalogServer, StubOperatorServer 他 |
| `src/tests/CreateBlockAPIKey.spec.ts` | 6 | 単体テストへ（D-009） | スタブ: StubAccessControlServer, StubBookManageServer, StubCatalogServer 他 |
| `src/tests/CreateBookAPIKey.spec.ts` | 9 | 単体テストへ（D-009） | スタブ: StubAccessControlServer, StubBookManageServer1, StubBookManageServer2 他 |
| `src/tests/CreateCatalogAPIKey.spec.ts` | 10 | 単体テストへ（D-009） | スタブ: StubAccessControlServer, StubCatalogServer, StubOperatorServer 他 |
| `src/tests/CreateJoinAPIKey.spec.ts` | 6 | 単体テストへ（D-009） | スタブ: StubAccessControlServer, StubCatalogServer, StubOperatorServer 他 |
| `src/tests/CreateOperatorAPIKey.spec.ts` | 7 | 単体テストへ（D-009） | スタブ: StubAccessControlServer, StubCatalogServer, StubOperatorServer 他 |
| `src/tests/CreateSettingAPIKey.spec.ts` | 6 | 単体テストへ（D-009） | スタブ: StubAccessControlServer, StubCatalogServer, StubOperatorServer 他 |
| `src/tests/CreateShareContinuousAPIKey.spec.ts` | 28 | 単体テストへ（D-009） | スタブ: StubAccessControlServer, StubBookManageServer, StubBookManageServer1 他; jest.spyOn 2 件 |
| `src/tests/CreateShareTempAPIKey.spec.ts` | 20 | 単体テストへ（D-009） | スタブ: StubAccessControlServer, StubBookManageServer, StubCTokenLedgerServer 他 |
| `src/tests/CreateStoreAPIKey.spec.ts` | 21 | 単体テストへ（D-009） | スタブ: StubAccessControlServer, StubBookManageServer, StubCatalogServer 他; jest.spyOn 2 件 |

### pxr-access-control-service（単位：access-control）

試験ファイル 17、試験 121 件。単体テストへ 59 件（D-009 の一覧に次の行）、E2E に入れた 62 件（GREEN。`test/e2e/access-control/`）

| 試験ファイル | 件数 | 分類 | 依存 |
| --- | --- | --- | --- |
| `src/tests/00-00.Validation.spec.ts` | 8 | E2E（GREEN。`00-00.Validation.yml`） | なし |
| `src/tests/01-01.Token.spec.ts` | 40 | 単体テストへ（D-009） | スタブ: StubServer |
| `src/tests/01-10.Token.Abnormal3.spec.ts` | 1 | 単体テストへ（D-009） | スタブ: StubServer |
| `src/tests/01-12.Token.AccessNomal1.spec.ts` | 1 | 単体テストへ（D-009） | スタブ: StubServer |
| `src/tests/01-13.Token.AccessNomal2.spec.ts` | 2 | 単体テストへ（D-009） | スタブ: StubServer |
| `src/tests/01-14.Token.Normal2.spec.ts` | 1 | 単体テストへ（D-009） | スタブ: StubServer |
| `src/tests/01-15.Token.Normal3.spec.ts` | 1 | 単体テストへ（D-009） | スタブ: StubServer |
| `src/tests/01-16.Token.dbSelectDateOk.spec.ts` | 1 | 単体テストへ（D-009） | スタブ: StubServer |
| `src/tests/01-17.Token.dbSelectDateError.spec.ts` | 1 | 単体テストへ（D-009） | スタブ: StubServer |
| `src/tests/01-19.Token.ApiTokenReceiveCheckOK.spec.ts` | 1 | 単体テストへ（D-009） | スタブ: StubServer |
| `src/tests/01-28.Token.BinaryUpload.spec.ts` | 3 | 単体テストへ（D-009） | スタブ: StubServer |
| `src/tests/01-29.Token.OperatorError1.spec.ts` | 7 | 単体テストへ（D-009） | スタブ: StubServer |
| `src/tests/02-01.AccessControl.spec.ts` | 30 | E2E（GREEN。`02-01.AccessControl.yml`。応答の token は実行ごとに変わるため、後のステップは `steps[N]` で受け渡す） | なし |
| `src/tests/03-01.Collate.spec.ts` | 12 | E2E（GREEN。`03-01.Collate.yml`） | なし |
| `src/tests/03-02.CollateOk.spec.ts` | 10 | E2E（GREEN。`03-02.CollateOk.yml`） | なし |
| `src/tests/03-08.Collate.NoMacthError.spec.ts` | 1 | E2E（GREEN。`03-08.Collate.NoMacthError.yml`） | なし |
| `src/tests/03-09.Collate.Macth.spec.ts` | 1 | E2E（GREEN。`03-09.Collate.Macth.yml`） | なし |

### pxr-catalog-service（単位：catalog）

試験ファイル 27、試験 1444 件。単体テストへ 636 件、E2E（確認済み、GREEN、`test/e2e/catalog/`）7 ファイル 1144 件

| 試験ファイル | 件数 | 分類 | 依存 |
| --- | --- | --- | --- |
| `src/tests/01-01.NameSpace.spec.ts` | 124 | 単体テストへ（D-010） | スタブ: StubOperatorServer |
| `src/tests/01-02.NameSpace.dbError.spec.ts` | 15 | 単体テストへ（D-010） | jest.mock 1 件 |
| `src/tests/02-01.CatalogCodeScope.spec.ts` | 92 | 単体テストへ（D-010） | スタブ: StubOperatorServer |
| `src/tests/02-02.CatalogCodeScope.dbError.spec.ts` | 13 | 単体テストへ（D-010） | jest.mock 1 件 |
| `src/tests/02-03.CatalogCodeScope.dbError.spec.ts` | 6 | 単体テストへ（D-010） | jest.mock 1 件 |
| `src/tests/03-01.CatalogName.spec.ts` | 46 | 単体テストへ（D-010） | スタブ: StubOperatorServer |
| `src/tests/03-02.CatalogName.dbError.spec.ts` | 3 | 単体テストへ（D-010） | jest.mock 1 件 |
| `src/tests/03-03.CatalogName.dbError.spec.ts` | 1 | 単体テストへ（D-010） | jest.mock 1 件 |
| `src/tests/03-04.CatalogName.dbError.spec.ts` | 1 | 単体テストへ（D-010） | jest.mock 1 件 |
| `src/tests/04-01.Catalog.model.spec.ts` | 303 | E2E（確認済み、GREEN） | なし |
| `src/tests/04-02.Catalog.built_in.spec.ts` | 229 | E2E（確認済み、GREEN） | なし |
| `src/tests/04-03.Catalog.ext.spec.ts` | 229 | E2E（確認済み、GREEN） | なし |
| `src/tests/04-04.Catalog.model.error.spec.ts` | 88 | 単体テストへ（D-010） | スタブ: StubOperatorServer |
| `src/tests/04-05.Catalog.dbError.spec.ts` | 15 | 単体テストへ（D-010） | jest.mock 1 件 |
| `src/tests/04-06.Catalog.dbError.spec.ts` | 3 | 単体テストへ（D-010） | jest.mock 1 件 |
| `src/tests/04-07.Catalog.dbError.spec.ts` | 15 | 単体テストへ（D-010） | jest.mock 1 件 |
| `src/tests/04-08.Catalog.dbError.spec.ts` | 12 | 単体テストへ（D-010） | jest.mock 1 件 |
| `src/tests/04-09.Catalog.bulk.spec.ts` | 1 | E2E（確認済み、GREEN。繰り返しで 337 件） | なし |
| `src/tests/05-01.CatalogInner.spec.ts` | 20 | E2E（確認済み、GREEN） | なし |
| `src/tests/06-01.CatalogFullText.spec.ts` | 21 | 単体テストへ（D-010） | スタブ: StubCloudSearchServer, StubOperatorServer; jest.mock 1 件 |
| `src/tests/07-01.UpdateSet.spec.ts` | 65 | 単体テストへ（D-010） | スタブ: StubOperatorServer |
| `src/tests/07-02.UpdateSet.spec.ts` | 72 | 単体テストへ（D-010） | スタブ: StubOperatorServer |
| `src/tests/07-03.UpdateSet.dbError.spec.ts` | 8 | 単体テストへ（D-010） | jest.mock 3 件 |
| `src/tests/07-04.UpdateSet.dbError.spec.ts` | 5 | 単体テストへ（D-010） | jest.mock 1 件 |
| `src/tests/08-01.CatalogPublic.spec.ts` | 4 | E2E（確認済み、GREEN） | なし |
| `src/tests/09-01.Attribute.spec.ts` | 31 | 単体テストへ（D-010） | スタブ: StubOperatorServer |
| `src/tests/10-01.CatalogHistoryCode.spec.ts` | 22 | E2E（確認済み、GREEN） | なし |

### pxr-catalog-update-service（単位：catalog）

試験ファイル 47、試験 640 件。単体テストへ 597 件（D-010）、E2E に入れた 43 件（GREEN。`test/e2e/catalog-update/`）

| 試験ファイル | 件数 | 分類 | 依存 |
| --- | --- | --- | --- |
| `src/tests/01-01.GetActorAccreditor.spec.ts` | 12 | 単体テストへ（D-010） | スタブ: _StubCatalogServer, _StubCatalogServer2, _StubCatalogServer3 他 |
| `src/tests/02-01.GetActor.spec.ts` | 18 | 単体テストへ（D-010） | スタブ: _StubCatalogServer, _StubOperatorServer, _StubOperatorServer2 他 |
| `src/tests/03-01.PostActor.spec.ts` | 41 | 単体テストへ（D-010） | スタブ: _StubCatalogServer, _StubCatalogServer2, _StubCertificationAuthServer 他 |
| `src/tests/04-01.PostActorApproval.spec.ts` | 17 | 単体テストへ（D-010） | スタブ: _StubCatalogServer, _StubCatalogServer2, _StubCatalogServer3 他 |
| `src/tests/04-02.PostActorApproval.spec.ts` | 1 | 単体テストへ（D-010） | スタブ: _StubCatalogServer2, _StubCertificationAuthServer, _StubNotificationServer 他 |
| `src/tests/04-03.PostActorApproval.spec.ts` | 1 | 単体テストへ（D-010） | スタブ: _StubCatalogServer2, _StubCertificationAuthServer, _StubNotificationServer 他 |
| `src/tests/05-01.Join.DraftOK.spec.ts` | 2 | 単体テストへ（D-010） | スタブ: _StubCatalogServer, _StubNotificationServer, _StubOperatorServer |
| `src/tests/05-02.Join.UpdateOK.spec.ts` | 9 | 単体テストへ（D-010） | スタブ: _StubCatalogServer, _StubNotificationServer, _StubOperatorServer |
| `src/tests/05-03.Join.DraftNG.spec.ts` | 13 | 単体テストへ（D-010） | スタブ: _StubCatalogServerEr |
| `src/tests/05-04.Join.AppOK.spec.ts` | 2 | 単体テストへ（D-010） | スタブ: _StubNotificationServer |
| `src/tests/05-05.Join.ParamNg.spec.ts` | 18 | E2E（GREEN。`05-05.Join.ParamNg.yml`） | なし |
| `src/tests/06-01.JoinRemove.DraftNG.spec.ts` | 13 | 単体テストへ（D-010） | スタブ: _StubCatalogServerOk |
| `src/tests/06-02.JoinRemove.CatalogActorNG.spec.ts` | 1 | 単体テストへ（D-010） | スタブ: _StubCatalogServerOk |
| `src/tests/06-03.JoinRemove.DraftOK.spec.ts` | 1 | 単体テストへ（D-010） | スタブ: _StubCatalogServerOk |
| `src/tests/06-04.JoinRemove.AppOK.spec.ts` | 2 | 単体テストへ（D-010） | スタブ: StubServer, _StubCatalogServerOk, _StubNotificationServer |
| `src/tests/06-05.JoinRemove.WfOK.spec.ts` | 3 | 単体テストへ（D-010） | スタブ: _StubCatalogServerOk, _StubNotificationServer |
| `src/tests/06-06.JoinRemove.OK2.spec.ts` | 1 | 単体テストへ（D-010） | スタブ: StubServer, _StubCatalogServerOk, _StubNotificationServer |
| `src/tests/06-07.JoinRemove.AppOK.spec.ts` | 5 | 単体テストへ（D-010） | スタブ: _StubCatalogServerOk, _StubNotificationServer |
| `src/tests/06-08.JoinRemove.ParamNg.spec.ts` | 18 | E2E（GREEN。`06-08.JoinRemove.ParamNg.yml`） | なし |
| `src/tests/07-01.Join.ApprovalApOK.spec.ts` | 2 | 単体テストへ（D-010） | スタブ: _StubCatalogServerOk, _StubNoticeServer |
| `src/tests/07-02.Join.ApprovalApCatalogError.spec.ts` | 1 | 単体テストへ（D-010） | スタブ: _StubCatalogServerOk, _StubNoticeServer |
| `src/tests/07-03.Join.ApprovalWfOK.spec.ts` | 1 | 単体テストへ（D-010） | スタブ: _StubCatalogServerOk, _StubNoticeServer |
| `src/tests/07-04.Join.ApprovalApNo.spec.ts` | 1 | 単体テストへ（D-010） | スタブ: _StubCatalogServerOk, _StubNoticeServer |
| `src/tests/07-05.Join.ApprovalApPutNG.spec.ts` | 1 | 単体テストへ（D-010） | スタブ: _StubCatalogServerOk, _StubNoticeServer |
| `src/tests/07-06.Join.ApprovalApNoticeNG.spec.ts` | 1 | 単体テストへ（D-010） | スタブ: _StubCatalogServerOk, _StubNoticeServer |
| `src/tests/07-07.JoinApproval.ParamNg.spec.ts` | 5 | E2E（GREEN。`07-07.JoinApproval.ParamNg.yml`。承認コードがない試験は、仕様どおり 404 だけを見る） | なし |
| `src/tests/07-08.Join.ApprovalApOK.spec.ts` | 1 | 単体テストへ（D-010） | スタブ: _StubCatalogServerOk, _StubNoticeServer |
| `src/tests/07-09.Join.ApprovalWfOK.spec.ts` | 1 | 単体テストへ（D-010） | スタブ: _StubCatalogServerOk, _StubNoticeServer |
| `src/tests/07-10.Join.ApprovalWfOK.spec.ts` | 1 | 単体テストへ（D-010） | スタブ: _StubCatalogServerOk, _StubNoticeServer |
| `src/tests/09-01.GetJoin.spec.ts` | 21 | 一部 E2E（GREEN。試験 19・20 の 2 件。`09-01.GetJoin.yml`）。残り 19 件は単体テストへ（D-010） | 18 件はスタブ（`CatalogServer4Get` 他）。1 件は catalog の停止に依る（「カタログサービスへの接続に失敗」） |
| `src/tests/11-01.PostActorRemove.spec.ts` | 32 | 単体テストへ（D-010） | スタブ: StubServer, _StubCatalogServer, _StubNotificationServer 他 |
| `src/tests/12-01.PostActorRemoveApproval.spec.ts` | 23 | 単体テストへ（D-010） | スタブ: _StubCatalogServer, _StubNotificationServer, _StubOperatorServer |
| `src/tests/13-01.PostTermsOfUse.spec.ts` | 35 | 単体テストへ（D-010） | スタブ: _StubCatalogServer, _StubOperatorServer |
| `src/tests/13-02.PutTermsOfUse.spec.ts` | 53 | 単体テストへ（D-010） | スタブ: _StubBookManageServer, _StubCatalogServer, _StubCatalogServer2 他 |
| `src/tests/14-01.PostStoreEvent.spec.ts` | 57 | 単体テストへ（D-010） | スタブ: _StubBookManageServer, _StubCatalogServer, _StubOperatorServer |
| `src/tests/14-02.PutStoreEvent.spec.ts` | 41 | 単体テストへ（D-010） | スタブ: _StubBookManageServer, _StubCatalogServer, _StubOperatorServer |
| `src/tests/16-01.GetRegionStatus.spec.ts` | 17 | 単体テストへ（D-010） | スタブ: _StubCatalogServer, _StubOperatorServer |
| `src/tests/16-02.PostRegionStatusStart.spec.ts` | 22 | 単体テストへ（D-010） | スタブ: _StubCatalogServer, _StubNotificationServer, _StubOperatorServer |
| `src/tests/16-03.PostRegionStatusEnd.spec.ts` | 28 | 単体テストへ（D-010） | スタブ: _StubCatalogServer, _StubNotificationServer, _StubOperatorServer |
| `src/tests/16-04.PostRegionStatusApproval.spec.ts` | 27 | 単体テストへ（D-010） | スタブ: _StubBookManageServer, _StubCatalogServer, _StubNotificationServer 他 |
| `src/tests/17-01.GetRegion.spec.ts` | 16 | 単体テストへ（D-010） | スタブ: _StubCatalogServer, _StubOperatorServer |
| `src/tests/17-02.PostRegion.spec.ts` | 6 | 単体テストへ（D-010） | スタブ: _StubCatalogServer, _StubOperatorServer |
| `src/tests/17-03.PostRegionDelete.spec.ts` | 12 | 単体テストへ（D-010） | スタブ: _StubCatalogServer, _StubOperatorServer |
| `src/tests/18-01.ActorAcquire.spec.ts` | 9 | 単体テストへ（D-010） | スタブ: _StubOperatorServer |
| `src/tests/18-02.PostDataOperation.spec.ts` | 24 | 単体テストへ（D-010） | スタブ: _StubBookManageServer, _StubCatalogServer, _StubOperatorServer |
| `src/tests/19-01.RequestUpdateSet.spec.ts` | 12 | 単体テストへ（D-010） | スタブ: _StubCatalogServer, _StubNotificationServer, _StubOperatorServer |
| `src/tests/19-02.ApprovalUpdateSet.spec.ts` | 12 | 単体テストへ（D-010） | スタブ: _StubCatalogServer, _StubNotificationServer, _StubOperatorServer |

### pxr-notification-service（単位：notification）

試験ファイル 8、試験 86 件。単体テストへ 55 件、E2E（確認済み、GREEN、`test/e2e/notification/`）31 件

| 試験ファイル | 件数 | 分類 | 依存 |
| --- | --- | --- | --- |
| `src/tests/AbnormalNotification.service.spec.ts` | 8 | 単体テストへ（D-011） | スタブ: StubServer |
| `src/tests/AbnormalServices.spec.ts` | 8 | 単体テストへ（D-011） | スタブ: StubServer |
| `src/tests/Notificatioin.approval.spec.ts` | 5 | 単体テストへ（D-011） | スタブ: StubServer |
| `src/tests/Notification.add.spec.ts` | 12 | 単体テストへ（D-011） | スタブ: StubServer |
| `src/tests/Notification.list.spec.ts` | 13 | 単体テストへ（D-011） | スタブ: StubServer |
| `src/tests/Notification.read.spec.ts` | 4 | 単体テストへ（D-011） | スタブ: StubServer |
| `src/tests/Notification.transfer.spec.ts` | 5 | 単体テストへ（D-011） | スタブ: StubBookManageServer, StubServer |
| `src/tests/Notification.validator.spec.ts` | 31 | E2E（確認済み、GREEN） | なし |

### pxr-identity-verificate-service（単位：identity-verify）

試験ファイル 20、試験 177 件。単体テストへ 177 件、E2E の候補（未確認）0 件

| 試験ファイル | 件数 | 分類 | 依存 |
| --- | --- | --- | --- |
| `src/tests/01-01.IssuanceCode.spec.ts` | 36 | 単体テストへ（D-012） | スタブ: StubServer; jest.mock 1 件 |
| `src/tests/01-02.IssuanceCode.abnormal.service.spec.ts` | 1 | 単体テストへ（D-012） | スタブ: StubServer; jest.mock 1 件 |
| `src/tests/02-01.PostUrl.spec.ts` | 14 | 単体テストへ（D-012） | スタブ: StubServer |
| `src/tests/02-02.PostUrl.abnormal.service.spec.ts` | 1 | 単体テストへ（D-012） | スタブ: StubServer; jest.mock 1 件 |
| `src/tests/03-01.Collation.spec.ts` | 17 | 単体テストへ（D-012） | スタブ: StubServer |
| `src/tests/03-02.Collation.abnormal.spec.ts` | 1 | 単体テストへ（D-012） | スタブ: StubServer; jest.mock 1 件 |
| `src/tests/04-01.AcquisitionPersonalIdentificationItems.spec.ts` | 7 | 単体テストへ（D-012） | スタブ: StubServer |
| `src/tests/04-02.AcquisitionPersonalIdentificationItems.abnormal.spec.ts` | 1 | 単体テストへ（D-012） | スタブ: StubServer; jest.mock 1 件 |
| `src/tests/05-01.VerifyPersonalIdentifyByOthers.spec.ts` | 9 | 単体テストへ（D-012） | スタブ: StubServer |
| `src/tests/05-02.VerifyPersonalIdentifyByOthers.abnormal.spec.ts` | 1 | 単体テストへ（D-012） | スタブ: StubServer; jest.mock 1 件 |
| `src/tests/06-01.UserInformationSetting.spec.ts` | 24 | 単体テストへ（D-012） | スタブ: StubServer |
| `src/tests/06-02.UserInformationSetting.abnormal.spec.ts` | 1 | 単体テストへ（D-012） | スタブ: StubServer; jest.mock 1 件 |
| `src/tests/07-01.CodeVerified.spec.ts` | 34 | 単体テストへ（D-012） | スタブ: StubServer; jest.mock 1 件 |
| `src/tests/07-02.CodeVerified.abnormal.service.spec.ts` | 1 | 単体テストへ（D-012） | スタブ: StubServer; jest.mock 1 件 |
| `src/tests/08-01.IndCollation.spec.ts` | 13 | 単体テストへ（D-012） | スタブ: StubServer |
| `src/tests/08-02.IndCollation.abnormal.spec.ts` | 1 | 単体テストへ（D-012） | スタブ: StubServer; jest.mock 1 件 |
| `src/tests/09-01.IssuanceUrl.spec.ts` | 7 | 単体テストへ（D-012） | スタブ: StubServer |
| `src/tests/09-02.IssuanceUrl.abnormal.service.spec.ts` | 1 | 単体テストへ（D-012） | スタブ: StubServer; jest.mock 1 件 |
| `src/tests/10-01.ConfirmUrl.spec.ts` | 6 | 単体テストへ（D-012） | スタブ: StubServer |
| `src/tests/10-02.ConfirmUrl.abnormal.spec.ts` | 1 | 単体テストへ（D-012） | スタブ: StubServer; jest.mock 1 件 |

### pxr-certification-authority-service（単位：certificate）

試験ファイル 36、試験 209 件。単体テストへ 209 件、E2E の候補（未確認）0 件

| 試験ファイル | 件数 | 分類 | 依存 |
| --- | --- | --- | --- |
| `src/tests/01-01.Root.spec.ts` | 18 | 単体テストへ（D-013） | スタブ: StubOperatorServer; jest.mock 1 件 |
| `src/tests/01-02.Root.dbError.spec.ts` | 2 | 単体テストへ（D-013） | jest.mock 1 件 |
| `src/tests/01-03.Root.insertError.spec.ts` | 1 | 単体テストへ（D-013） | jest.mock 2 件 |
| `src/tests/01-04.Root.openSslError.spec.ts` | 1 | 単体テストへ（D-013） | jest.mock 1 件 |
| `src/tests/01-05.Root.openSslError2.spec.ts` | 1 | 単体テストへ（D-013） | jest.mock 1 件 |
| `src/tests/01-06.Root.openSslError3.spec.ts` | 1 | 単体テストへ（D-013） | jest.mock 1 件 |
| `src/tests/01-07.Root.openSslError4.spec.ts` | 1 | 単体テストへ（D-013） | jest.mock 1 件 |
| `src/tests/01-08.Root.openSslError5.spec.ts` | 1 | 単体テストへ（D-013） | jest.mock 1 件 |
| `src/tests/02-01.Server.spec.ts` | 46 | 単体テストへ（D-013） | スタブ: StubOperatorServer; jest.mock 1 件 |
| `src/tests/02-02.Server.dbError.spec.ts` | 3 | 単体テストへ（D-013） | jest.mock 1 件 |
| `src/tests/02-03.Server.insertError.spec.ts` | 1 | 単体テストへ（D-013） | jest.mock 2 件 |
| `src/tests/02-04.Server.openSslError.spec.ts` | 2 | 単体テストへ（D-013） | jest.mock 1 件 |
| `src/tests/02-05.Server.openSslError2.spec.ts` | 1 | 単体テストへ（D-013） | jest.mock 1 件 |
| `src/tests/02-06.Server.openSslError3.spec.ts` | 1 | 単体テストへ（D-013） | jest.mock 1 件 |
| `src/tests/02-07.Server.openSslError4.spec.ts` | 1 | 単体テストへ（D-013） | jest.mock 1 件 |
| `src/tests/02-08.Server.openSslError5.spec.ts` | 1 | 単体テストへ（D-013） | jest.mock 1 件 |
| `src/tests/02-09.Server.deleteError.spec.ts` | 1 | 単体テストへ（D-013） | jest.mock 2 件 |
| `src/tests/03-01.Client.spec.ts` | 49 | 単体テストへ（D-013） | スタブ: StubOperatorServer; jest.mock 1 件 |
| `src/tests/03-02.Client.dbError.spec.ts` | 3 | 単体テストへ（D-013） | jest.mock 1 件 |
| `src/tests/03-03.Client.insertError.spec.ts` | 1 | 単体テストへ（D-013） | jest.mock 2 件 |
| `src/tests/03-04.Client.openSslError.spec.ts` | 2 | 単体テストへ（D-013） | jest.mock 1 件 |
| `src/tests/03-05.Client.openSslError2.spec.ts` | 1 | 単体テストへ（D-013） | jest.mock 1 件 |
| `src/tests/03-06.Client.openSslError3.spec.ts` | 1 | 単体テストへ（D-013） | jest.mock 1 件 |
| `src/tests/03-07.Client.openSslError4.spec.ts` | 1 | 単体テストへ（D-013） | jest.mock 1 件 |
| `src/tests/03-08.Client.openSslError5.spec.ts` | 1 | 単体テストへ（D-013） | jest.mock 1 件 |
| `src/tests/03-09.Client.deleteError.spec.ts` | 1 | 単体テストへ（D-013） | jest.mock 2 件 |
| `src/tests/04-01.Actor.spec.ts` | 11 | 単体テストへ（D-013） | スタブ: StubOperatorServer |
| `src/tests/04-02.Actor.dbError.spec.ts` | 1 | 単体テストへ（D-013） | jest.mock 1 件 |
| `src/tests/05-01.List.spec.ts` | 9 | 単体テストへ（D-013） | スタブ: StubOperatorServer |
| `src/tests/05-02.List.dbError.spec.ts` | 1 | 単体テストへ（D-013） | jest.mock 1 件 |
| `src/tests/05-03.List.ActorCode.spec.ts` | 10 | 単体テストへ（D-013） | スタブ: StubOperatorServer |
| `src/tests/06-01.Distributed.spec.ts` | 16 | 単体テストへ（D-013） | スタブ: StubOperatorServer |
| `src/tests/06-02.Distributed.dbError.spec.ts` | 1 | 単体テストへ（D-013） | jest.mock 1 件 |
| `src/tests/06-03.Distributed.insertError.spec.ts` | 1 | 単体テストへ（D-013） | jest.mock 1 件 |
| `src/tests/07-01.Valid.spec.ts` | 15 | 単体テストへ（D-013） | スタブ: StubOperatorServer |
| `src/tests/07-02.Valid.dbError.spec.ts` | 1 | 単体テストへ（D-013） | jest.mock 1 件 |

### pxr-certificate-manage-service（単位：certificate）

試験ファイル 6、試験 47 件。単体テストへ 47 件、E2E の候補（未確認）0 件

| 試験ファイル | 件数 | 分類 | 依存 |
| --- | --- | --- | --- |
| `src/tests/01-01.CertificateManage.spec.ts` | 34 | 単体テストへ（D-013） | スタブ: StubCertificationAuthorityServer, StubOperatorServer, StubOperatorServer1 他 |
| `src/tests/01-02.CertificateRevokeList.spec.ts` | 5 | 単体テストへ（D-013） | スタブ: StubCertificationAuthorityServer, StubCertificationAuthorityServer1, StubOperatorServer |
| `src/tests/01-03.CertificateManage.insertError.spec.ts` | 1 | 単体テストへ（D-013） | スタブ: StubCertificationAuthorityServer; jest.mock 1 件 |
| `src/tests/01-04.CertificateManage.dbRollback.spec.ts` | 1 | 単体テストへ（D-013） | スタブ: StubCertificationAuthorityServer, StubOperatorServer, StubOperatorServer0 |
| `src/tests/01-05.CertificateCheck.spec.ts` | 3 | 単体テストへ（D-013） | スタブ: StubCertificationAuthorityServer, StubOperatorServer |
| `src/tests/01-06.CertificateManageNotPxrRootBlock.spec.ts` | 3 | 単体テストへ（D-013） | スタブ: StubCertificationAuthorityServer, StubOperatorServer; jest.mock 1 件 |

### pxr-binary-manage-service（単位：binary）

試験ファイル 16、試験 127 件。単体テストへ 127 件、E2E の候補（未確認）0 件

| 試験ファイル | 件数 | 分類 | 依存 |
| --- | --- | --- | --- |
| `src/tests/01-01.UploadStart.spec.ts` | 14 | 単体テストへ（D-014） | スタブ: StubOperatorServer |
| `src/tests/01-02.UploadStart.dbInsertError.spec.ts` | 1 | 単体テストへ（D-014） | jest.spyOn 1 件 |
| `src/tests/02-01.UploadEnd.spec.ts` | 13 | 単体テストへ（D-014） | スタブ: StubOperatorServer |
| `src/tests/03-01.UploadCancel.spec.ts` | 10 | 単体テストへ（D-014） | スタブ: StubOperatorServer |
| `src/tests/03-02.UploadCancel.dbUpdateError.spec.ts` | 1 | 単体テストへ（D-014） | jest.spyOn 1 件 |
| `src/tests/04-01.Upload.spec.ts` | 16 | 単体テストへ（D-014） | スタブ: StubOperatorServer |
| `src/tests/04-02.Upload.dbInsertError.spec.ts` | 1 | 単体テストへ（D-014） | jest.spyOn 3 件 |
| `src/tests/05-01.DownloadStart.spec.ts` | 11 | 単体テストへ（D-014） | スタブ: StubOperatorServer |
| `src/tests/05-02.DownloadStart.dbUpdateError.spec.ts` | 1 | 単体テストへ（D-014） | jest.spyOn 3 件 |
| `src/tests/06-01.DownloadEnd.spec.ts` | 12 | 単体テストへ（D-014） | スタブ: StubOperatorServer |
| `src/tests/07-01.DownloadCancel.spec.ts` | 10 | 単体テストへ（D-014） | スタブ: StubOperatorServer |
| `src/tests/07-02.DownloadCancel.dbUpdateError.spec.ts` | 1 | 単体テストへ（D-014） | jest.spyOn 1 件 |
| `src/tests/08-01.Download.spec.ts` | 13 | 単体テストへ（D-014） | スタブ: StubOperatorServer |
| `src/tests/09-01.GetBinaryManage.spec.ts` | 13 | 単体テストへ（D-014） | スタブ: StubOperatorServer |
| `src/tests/10-01.DeleteBinaryManage.spec.ts` | 9 | 単体テストへ（D-014） | スタブ: StubOperatorServer |
| `src/tests/10-02.DeleteBinaryManage.dbUpdateError.spec.ts` | 1 | 単体テストへ（D-014） | jest.spyOn 1 件 |

### pxr-ctoken-ledger-service（単位：ctoken）

試験ファイル 5、試験 240 件。単体テストへ 240 件、E2E の候補（未確認）0 件

| 試験ファイル | 件数 | 分類 | 依存 |
| --- | --- | --- | --- |
| `src/tests/01-01.LocalCToken.spec.ts` | 167 | 単体テストへ（D-015） | スタブ: StubServer; jest.spyOn 1 件 |
| `src/tests/01-02.LocalCToken.insertError.spec.ts` | 3 | 単体テストへ（D-015） | スタブ: StubServer; jest.mock 1 件 |
| `src/tests/02-01.CTokenLedgerCount.spec.ts` | 23 | 単体テストへ（D-015） | スタブ: StubServer |
| `src/tests/03-01.PxrIdSearch.spec.ts` | 14 | 単体テストへ（D-015） | スタブ: StubServer |
| `src/tests/04-01.CTokenLedger.spec.ts` | 33 | 単体テストへ（D-015） | スタブ: StubServer |

### pxr-local-ctoken-service（単位：ctoken）

試験ファイル 3、試験 202 件。単体テストへ 202 件、E2E の候補（未確認）0 件

| 試験ファイル | 件数 | 分類 | 依存 |
| --- | --- | --- | --- |
| `src/tests/01-01.LocalCToken.spec.ts` | 160 | 単体テストへ（D-015） | スタブ: StubServer |
| `src/tests/01-02.LocalCToken.insertError.spec.ts` | 3 | 単体テストへ（D-015） | スタブ: StubServer; jest.mock 1 件 |
| `src/tests/02-01.CTokenLedger.spec.ts` | 39 | 単体テストへ（D-015） | スタブ: StubServer |

### pxr-book-manage-service（単位：book-manage）

試験ファイル 69、試験 2785 件。単体テストへ 1920 件、E2E（確認済み、GREEN、`test/e2e/book-manage/`）4 件

| 試験ファイル | 件数 | 分類 | 依存 |
| --- | --- | --- | --- |
| `src/tests/01-01.BookOpen.spec.ts` | 136 | 単体テストへ（D-016） | スタブ: StubServer |
| `src/tests/02-01.BookSearch.spec.ts` | 44 | 単体テストへ（D-016） | スタブ: StubCatalogServer, StubOperatorServer, StubOperatorServerBookSearch |
| `src/tests/03-01.IdentitySearch.spec.ts` | 13 | 単体テストへ（D-016） | スタブ: StubCatalogServer, StubOperatorServer |
| `src/tests/04-01.BookCooperate.spec.ts` | 44 | 単体テストへ（D-016） | スタブ: StubCatalogServer, StubCatalogServerNoIdService, StubIdServiceServer 他 |
| `src/tests/05-01.DataStorePost.spec.ts` | 96 | 単体テストへ（D-016） | スタブ: StubCatalogServer, StubCatalogServerDataStore, StubOperatorServer 他 |
| `src/tests/06-01.DataStoreGet.spec.ts` | 39 | 単体テストへ（D-016） | スタブ: StubCatalogServer, StubCatalogServerDataStore, StubOperatorServer 他 |
| `src/tests/07-01.CheckPxrId.spec.ts` | 13 | 単体テストへ（D-016） | スタブ: StubServer |
| `src/tests/08-01.CheckIdentification.spec.ts` | 45 | 単体テストへ（D-016） | スタブ: StubServer |
| `src/tests/09-01.ReleaseCooperate.spec.ts` | 32 | 単体テストへ（D-016） | スタブ: StubCatalogServer, StubCatalogServer05, StubCatalogServerReleaseCooperate 他; jest.spyOn 1 件 |
| `src/tests/10-01.DataStoreDelete.spec.ts` | 10 | 単体テストへ（D-016） | スタブ: StubOperatorServer, StubOperatorServerType0 |
| `src/tests/11-01.BookClose.spec.ts` | 14 | 単体テストへ（D-016） | スタブ: StubCatalogServer, StubCatalogServerBookClose, StubIdServiceServer 他 |
| `src/tests/12-01.GetCooperate.spec.ts` | 20 | 単体テストへ（D-016） | スタブ: StubCatalogServer, StubCatalogServerGetCooperate, StubOperatorServer 他 |
| `src/tests/13-01.LoginCode.spec.ts` | 13 | 単体テストへ（D-016） | スタブ: StubOperatorServer, StubOperatorServerLoginCode |
| `src/tests/14-01.Identification.spec.ts` | 2 | E2E（確認済み、GREEN） | なし |
| `src/tests/14-02.Identification.spec.ts` | 2 | E2E（確認済み、GREEN） | なし |
| `src/tests/15-01.CooperateRequest.spec.ts` | 33 | 単体テストへ（D-016） | スタブ: StubOperatorServer, StubOperatorService |
| `src/tests/16-01.ForceDeletion.spec.ts` | 7 | 単体テストへ（D-016） | スタブ: StubServer; jest.mock 1 件 |
| `src/tests/16-02.ForceDeletion.FailedUpdateOperator.spec.ts` | 2 | 単体テストへ（D-016） | スタブ: StubServer; jest.mock 1 件 |
| `src/tests/16-03.ForceDeletion.NotFoundPxr.spec.ts` | 2 | 単体テストへ（D-016） | スタブ: StubServer; jest.mock 1 件 |
| `src/tests/16-04.ForceDeletion.SearchFailedOperator.spec.ts` | 2 | 単体テストへ（D-016） | スタブ: StubServer; jest.mock 1 件 |
| `src/tests/17-01.TemporarilySharedCode.spec.ts` | 174 | 単体テストへ（D-016） | スタブ: StubCTokenServer, StubCatalogServer, StubCatalogServerShare 他 |
| `src/tests/18-01.indAccessLog.spec.ts` | 74 | 単体テストへ（D-016） | スタブ: StubCatalogForAccessLog, StubCatalogForAccessLogGetCatalogInfosError, StubCatalogServer 他 |
| `src/tests/18-02.AccessLog.spec.ts` | 74 | 単体テストへ（D-016） | スタブ: StubCatalogForAccessLog, StubCatalogForAccessLogGetCatalogInfosError, StubCatalogServer 他 |
| `src/tests/19-01.SendSms.spec.ts` | 18 | 単体テストへ（D-016） | スタブ: StubCatalogServer, StubCatalogServerNotify, StubOperatorServer 他 |
| `src/tests/21-01.DataShare.spec.ts` | 30 | 単体テストへ（D-016） | スタブ: StubProxyServer, StubServer |
| `src/tests/21-02.DataShare.spec.ts` | 27 | 単体テストへ（D-016） | スタブ: StubProxyServer, StubServer |
| `src/tests/22-01.UserInfo.spec.ts` | 18 | 単体テストへ（D-016） | スタブ: StubOperatorServer, StubOperatorServiceUserInfoService, StubServer |
| `src/tests/23-01.Book.spec.ts` | 41 | 単体テストへ（D-016） | スタブ: StubCatalogServer, StubCatalogServerGetBook, StubOperatorServer 他; jest.spyOn 1 件 |
| `src/tests/24-01.CTokenSearch.spec.ts` | 9 | 単体テストへ（D-016） | スタブ: StubCTokenLedgerServer, StubServer |
| `src/tests/26-01.GetIndBook.spec.ts` | 21 | 単体テストへ（D-016） | スタブ: StubCatalogServer, StubCatalogServerGetIndBook, StubOperatorServer 他 |
| `src/tests/26-02.GetBook.spec.ts` | 21 | 単体テストへ（D-016） | スタブ: StubCatalogServer, StubCatalogServerGetIndBook, StubOperatorServer 他 |
| `src/tests/29-01.IndTermsOfUse.spec.ts` | 19 | 単体テストへ（D-016） | スタブ: StubCatalogServer, StubCatalogServerIndTermsOfUse, StubOperatorServer |
| `src/tests/29-02.TermsOfUse.spec.ts` | 19 | 単体テストへ（D-016） | スタブ: StubCatalogServer, StubCatalogServerIndTermsOfUse, StubOperatorServer |
| `src/tests/31-01.StoreEvent.spec.ts` | 47 | 単体テストへ（D-016） | スタブ: StubCatalogServer, StubCatalogServerStoreEvent, StubCatalogServerStoreEventNotificate 他; jest.mock 1 件; jest.spyOn 6 件 |
| `src/tests/32-01.UpdateStoreEvent.spec.ts` | 36 | 単体テストへ（D-016） | スタブ: StubCatalogServer, StubCatalogServerStoreEvent, StubOperatorServer 他 |
| `src/tests/33-01.PostTermsOfUse.spec.ts` | 11 | 単体テストへ（D-016） | スタブ: StubCatalogServer, StubCatalogServerTermsOfUse, StubOperatorServer 他 |
| `src/tests/33-02.PostTermsOfUse.spec.ts` | 11 | 単体テストへ（D-016） | スタブ: StubCatalogServer, StubCatalogServerTermsOfUse, StubOperatorServer 他 |
| `src/tests/35-01.ReserveDeletetion.spec.ts` | 9 | 単体テストへ（D-016） | スタブ: StubOperatorServer, StubOperatorServerType0 |
| `src/tests/36-01.OutputCondition.spec.ts` | 11 | 単体テストへ（D-016） | スタブ: StubOperatorServer, StubOperatorServerType0 |
| `src/tests/37-01.TermsOfUseTargetFind.spec.ts` | 28 | 単体テストへ（D-016） | スタブ: StubOperatorServer, StubOperatorServerType0 |
| `src/tests/37-02.TermsOfUseTarget.spec.ts` | 28 | 単体テストへ（D-016） | スタブ: StubOperatorServer, StubOperatorServerType0 |
| `src/tests/38-01.TermsOfUseRegionTarget.spec.ts` | 18 | 単体テストへ（D-016） | スタブ: StubCatalogServer, StubOperatorServer, StubOperatorServerType0 他 |
| `src/tests/38-02.TermsOfUsePlatformTarget.spec.ts` | 18 | 単体テストへ（D-016） | スタブ: StubCatalogServer, StubOperatorServer, StubOperatorServerType0 他 |
| `src/tests/39-01.TermsOfUseDeletionTarget.spec.ts` | 39 | 単体テストへ（D-016） | スタブ: StubCatalogServer, StubOperatorServer, StubOperatorServerType0 他 |
| `src/tests/40-01.TermOfUseUpdate.spec.ts` | 28 | 単体テストへ（D-016） | スタブ: StubOperatorServer, StubOperatorServerType0 |
| `src/tests/41-01.OutputConditionDataManage.spec.ts` | 94 | 単体テストへ（D-016） | スタブ: StubOperatorServer, StubOperatorServerType0 |
| `src/tests/42-01.OutputMcdData.spec.ts` | 74 | 単体テストへ（D-016） | スタブ: StubCatalogServer, StubCatalogServerOutputMcdData, StubOperatorServer 他 |
| `src/tests/43-01.OutputPrepare.spec.ts` | 40 | 単体テストへ（D-016） | スタブ: StubCatalogServer, StubCatalogServerPrepare, StubOperatorServer 他 |
| `src/tests/44-01.GetDeleteTargetBook.spec.ts` | 16 | 単体テストへ（D-016） | スタブ: StubCatalogServer, StubCatalogServerDelete, StubOperatorServer 他 |
| `src/tests/45-01.BookForceDelete.spec.ts` | 25 | 単体テストへ（D-016） | スタブ: StubCatalogServer, StubCatalogServerBookForceDelete, StubIdServiceServer 他; jest.mock 1 件 |
| `src/tests/46-01.TermsOfUse.NotificationComplete.spec.ts` | 30 | 単体テストへ（D-016） | スタブ: StubOperatorServer, StubOperatorServerType0 |
| `src/tests/47-01.RegionClose.spec.ts` | 68 | 単体テストへ（D-016） | スタブ: StubCatalogServer, StubCatalogServerRegionClose, StubOperatorServer 他 |
| `src/tests/48-01.IndAppendixPut.spec.ts` | 3 | 単体テストへ（D-016） | スタブ: StubOperatorServer |
| `src/tests/48-02.AppendixPut.spec.ts` | 3 | 単体テストへ（D-016） | スタブ: StubOperatorServer |
| `src/tests/49-01.CooperateUserPost.spec.ts` | 32 | 単体テストへ（D-016） | スタブ: StubOperatorServer, StubOperatorService |
| `src/tests/50-01.CooperateReleaseRequest.spec.ts` | 31 | 単体テストへ（D-016） | スタブ: StubOperatorServer, StubOperatorService |
| `src/tests/51-01.SearchUserPost.spec.ts` | 19 | 単体テストへ（D-016） | スタブ: StubOperatorServer, StubOperatorServerBookSearchUser |
| `src/tests/52-01.SearchCooperatePost.spec.ts` | 11 | 単体テストへ（D-016） | スタブ: StubOperatorServer |
| `src/tests/53-01.SettingsUpdatePost.spec.ts` | 9 | 単体テストへ（D-016） | スタブ: StubOperatorServer |
| `src/tests/54-01.SettingsTargetFindGet.spec.ts` | 9 | 単体テストへ（D-016） | スタブ: StubCatalogServer, StubCatalogServerSettingsTargetFind, StubOperatorServer |
| `src/tests/55-01.SettingsTargetGet.spec.ts` | 9 | 単体テストへ（D-016） | スタブ: StubCatalogServer, StubCatalogServerSettingsTargetFind, StubOperatorServer |
| `src/tests/56-01.SettingsNotificationCompletePost.spec.ts` | 10 | 単体テストへ（D-016） | スタブ: StubOperatorServer |
| `src/tests/64-00.PermissionAnalyzer.store.spec.ts` | 260 | 対象外（supertest なし） | なし |
| `src/tests/64-01.PermissionAnalyzer.share.spec.ts` | 273 | 対象外（supertest なし） | なし |
| `src/tests/64-02.PermissionAnalyzer.storeEvent.spec.ts` | 250 | 対象外（supertest なし） | なし |
| `src/tests/64-03.PermissionAnalyzer.error.spec.ts` | 50 | 対象外（supertest なし） | なし |
| `src/tests/65-00.catalogAccessor.spec.ts` | 28 | 対象外（supertest なし） | スタブ: StubCatalogServer, StubCatalogServerDataStore, StubCatalogServerShare |
| `src/tests/65-01.DataStorePermission.spec.ts` | 15 | 単体テストへ（D-016） | jest.mock 1 件; jest.spyOn 1 件 |
| `src/tests/65-02.DataSharePermission.spec.ts` | 28 | 単体テストへ（D-016） | スタブ: StubCTokenLedgerServer; jest.mock 1 件; jest.spyOn 4 件 |

### pxr-book-operate-service（単位：book-operate）

試験ファイル 28、試験 1677 件。単体テストへ 1677 件、E2E の候補（未確認）0 件

| 試験ファイル | 件数 | 分類 | 依存 |
| --- | --- | --- | --- |
| `src/tests/01-01.BookUserCreate.spec.ts` | 62 | 単体テストへ（D-017） | スタブ: StubBookManageServer, StubCatalogServer, StubNotificationServer 他; jest.spyOn 2 件 |
| `src/tests/01-02.BookUserdbInsertError.spec.ts` | 1 | 単体テストへ（D-017） | スタブ: StubBookManageServer, StubCatalogServer, StubOperatorServer; jest.mock 1 件 |
| `src/tests/01-03.BookUserdbInsertError.spec.ts` | 1 | 単体テストへ（D-017） | スタブ: StubBookManageServer, StubCatalogServer, StubOperatorServer; jest.mock 1 件 |
| `src/tests/02-01.BookUserList.spec.ts` | 52 | 単体テストへ（D-017） | スタブ: StubBookManageServer, StubCatalogServer, StubNotificationServer 他; jest.spyOn 1 件 |
| `src/tests/02-02.BookUserListdbError.spec.ts` | 1 | 単体テストへ（D-017） | スタブ: StubCatalogServer; jest.mock 1 件 |
| `src/tests/03-01.Event.outerBlockStoreOn.spec.ts` | 168 | 単体テストへ（D-017） | スタブ: StubCTokenServer, StubCatalogServer, StubOperatorServer; jest.mock 1 件 |
| `src/tests/03-02.Event.outerBlockStoreOff.spec.ts` | 179 | 単体テストへ（D-017） | スタブ: StubCTokenServer, StubCatalogServer, StubOperatorServer; jest.mock 1 件; jest.spyOn 1 件 |
| `src/tests/03-03.Event.dbError.spec.ts` | 5 | 単体テストへ（D-017） | スタブ: StubCatalogServer; jest.mock 1 件 |
| `src/tests/03-04.Event.dbInsertError.spec.ts` | 5 | 単体テストへ（D-017） | スタブ: StubCatalogServer; jest.mock 1 件 |
| `src/tests/04-01.Thing.outerBlockStoreOn.spec.ts` | 224 | 単体テストへ（D-017） | スタブ: StubCTokenServer, StubCatalogServer, StubOperatorServer; jest.mock 1 件 |
| `src/tests/04-02.Thing.outerBlockStoreOff.spec.ts` | 224 | 単体テストへ（D-017） | スタブ: StubCTokenServer, StubCatalogServer, StubOperatorServer; jest.mock 1 件 |
| `src/tests/04-03.Thing.dbError.spec.ts` | 5 | 単体テストへ（D-017） | スタブ: StubCatalogServer; jest.mock 1 件 |
| `src/tests/04-04.Thing.dbInsertError.spec.ts` | 5 | 単体テストへ（D-017） | スタブ: StubCatalogServer; jest.mock 1 件 |
| `src/tests/05-01.Book.spec.ts` | 49 | 単体テストへ（D-017） | スタブ: StubCatalogServer, StubOperatorServer |
| `src/tests/05-02.BookOutsideStore.spec.ts` | 2 | 単体テストへ（D-017） | スタブ: StubCatalogServer, StubOperatorServer, StubOutsideStoreServer; jest.mock 1 件 |
| `src/tests/06-01.IndAccessLog.spec.ts` | 48 | 単体テストへ（D-017） | スタブ: StubOperatorServer |
| `src/tests/06-02.IndAccessLog.spec.ts` | 48 | 単体テストへ（D-017） | スタブ: StubOperatorServer |
| `src/tests/07-01.Document.outerBlockStoreOn.spec.ts` | 174 | 単体テストへ（D-017） | スタブ: StubCTokenServer, StubCatalogServer, StubOperatorServer; jest.mock 1 件 |
| `src/tests/07-02.Document.outerBlockStoreOff.spec.ts` | 185 | 単体テストへ（D-017） | スタブ: StubCTokenServer, StubCatalogServer, StubOperatorServer; jest.mock 1 件; jest.spyOn 1 件 |
| `src/tests/07-03.Document.dbError.spec.ts` | 3 | 単体テストへ（D-017） | スタブ: StubCatalogServer; jest.mock 1 件 |
| `src/tests/07-04.Document.dbInsertError.spec.ts` | 3 | 単体テストへ（D-017） | スタブ: StubCatalogServer; jest.mock 1 件 |
| `src/tests/08-01.GetDataByTemporarilySharedCode.spec.ts` | 52 | 単体テストへ（D-017） | スタブ: StubBookManageServer, StubOperatorServer, StubOperatorServerType0 他 |
| `src/tests/09-01.PostShare.spec.ts` | 9 | 単体テストへ（D-017） | スタブ: StubBookManageServer, StubCatalogServer, StubOperatorServer 他 |
| `src/tests/10-01.StoreEventReceive.spec.ts` | 78 | 単体テストへ（D-017） | スタブ: StubCatalogServer, StubOperatorServer |
| `src/tests/11-01.ShareData.spec.ts` | 56 | 単体テストへ（D-017） | スタブ: StubBookManageServer, StubCatalogServer, StubOperatorServer 他 |
| `src/tests/11-02.ShareData.OutsideStore.spec.ts` | 3 | 単体テストへ（D-017） | スタブ: StubBookManageServer, StubCatalogServer, StubOperatorServer 他; jest.mock 1 件 |
| `src/tests/13-01.deleteUserStoreData.spec.ts` | 23 | 単体テストへ（D-017） | スタブ: StubCTokenServer, StubCatalogServer, StubOperatorServer 他 |
| `src/tests/15-01.BookUserCreateBatch.spec.ts` | 12 | 単体テストへ（D-017） | スタブ: StubBookManageServer, StubCatalogServer, StubNotificationServer 他; jest.spyOn 1 件 |
