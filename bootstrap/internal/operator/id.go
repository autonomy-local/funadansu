package operator

// OperatorID は、オペレーターの内部の主キーです（旧 id、bigint）。JSON では integer。
type OperatorID int64

// PxrID は、業務上の識別子です（旧 pxr_id、varchar(255)）。JSON では string。
type PxrID string

// Operator は、この単位が返す値です。見本では id と pxr_id だけを持ちます。
type Operator struct {
	ID    OperatorID
	PxrID PxrID
}
