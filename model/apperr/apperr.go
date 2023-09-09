package apperr

type AppErr struct {
	StatusCode int
	Message    string
	Log        string
}

// Error returns error message.
// AppErr satisfies error interface.
func (e AppErr) Error() string {
	return e.Message
}
