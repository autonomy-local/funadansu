-- 見本のクエリ。旧スキーマ pxr_operator の表を読むだけで、書き込みはしない。
-- 機密の列（hpassword など）と個人情報（name、mobile_phone、mail）は、見本でも選ばない。

-- name: FindOperatorByID :one
SELECT id, pxr_id
FROM pxr_operator.operator
WHERE id = $1;

-- テスト用のデータの投入と削除。本番の経路では使わない（テストの後始末のため）。
-- name: InsertOperatorForTest :one
INSERT INTO pxr_operator.operator (
    type, login_id, hpassword, pxr_id, user_information, name, mobile_phone, mail,
    auth, attributes, lock_flg, user_id, region_catalog_code, app_catalog_code,
    wf_catalog_code, client_id, created_by, updated_by, unique_check_login_id
) VALUES (
    0, $1, 'test', $2, '{}', 'test', '', '', '{}', '{}', false, $3, 0, 0, 0,
    'test', 'test', 'test', $4
) RETURNING id;

-- name: DeleteOperatorForTest :exec
DELETE FROM pxr_operator.operator WHERE id = $1;
