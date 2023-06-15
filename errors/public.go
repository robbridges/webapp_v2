package errors

func Public(err error, msg string) error {
	return PublicError{err, msg}
}

type PublicError struct {
	err error
	msg string
}

func (pe PublicError) Error() string {
	return pe.err.Error()
}

func (pe PublicError) Public() string {
	return pe.msg
}

func (pe PublicError) Unwrap() string {
	return pe.Error()
}
