package entities

type OperationStatusType string

const (
	OperationStatusTypePending   OperationStatusType = "pending"
	OperationStatusTypeActive    OperationStatusType = "active"
	OperationStatusTypeCancelled OperationStatusType = "closed"
)
