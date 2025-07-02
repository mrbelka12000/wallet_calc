package usecase

type (
	ClientError struct {
		Message string
	}
)

func (e ClientError) Error() string {
	return e.Message
}

func newClientError(msg string) ClientError {
	return ClientError{Message: msg}
}
