-- 見本のクエリ。旧スキーマ pxr_operator の表を読むだけで、書き込みはしない。
-- 機密の列（hpassword など）と個人情報（name、mobile_phone、mail）は、見本でも選ばない。

-- name: FindOperatorByID :one
SELECT id, pxr_id
FROM pxr_operator.operator
WHERE id = $1;
