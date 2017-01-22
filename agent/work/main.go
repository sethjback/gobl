package work

const (
	StateErrors   = "errors"
	StateSkipped  = "skipped"
	StateComplete = "complete"

	ErrorFileHash       = "FileHashFailed"
	ErrorSaveEngines    = "SaveEngineErrors"
	ErrorRestoreEngines = "RestoreEngineErrors"
	ErrorModifications  = "ModificationErrors"
	ErrorFileOps        = "FileOperationFaild"
	ErrorSave           = "SaveFailed"
	ErrorRestore        = "RestoreFailed"
)
