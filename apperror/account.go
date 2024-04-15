package apperror

func NewEmailAlreadyVerified(err error) error {
	return NewAppError(
		CodeBadRequest,
		"email already verified",
		err,
	)
}
