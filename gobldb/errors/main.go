package errors

const (
	ErrCodeMarshal   = "MashalDataFailed"
	ErrCodeUnMarshal = "UnMashalDataFailed"
	ErrCodeSave      = "WriteDataFailed"
	ErrCodeGet       = "GetDataFailed"
	ErrCodeDelete    = "DeleteDataFailed"
	ErrCodeNotFound  = "EntityNotFound"
	ErrFilterOptions = "InvalidFilterOptions"

	ErrDBDriver = "InvalidDBDriver"
)
