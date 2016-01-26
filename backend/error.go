package kahinah

// DBError represents an error occuring with the database backend.
type DBError struct {
	WrappedError error
}

func (e *DBError) Error() string {
	return e.WrappedError.Error()
}
