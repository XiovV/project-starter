package repository

type NotFoundErr struct {
	Message string
}

func (e NotFoundErr) Error() string {
	return e.Message
}
func notFoundErr(msg string) *NotFoundErr {
	return &NotFoundErr{msg}
}

type AlreadyExistsErr struct {
	Message string
}

func (e AlreadyExistsErr) Error() string {
	return e.Message
}
func alreadyExistsErr(msg string) *AlreadyExistsErr {
	return &AlreadyExistsErr{msg}
}
