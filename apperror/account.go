package apperror

func NewEmailAlreadyVerified(err error) error {
	return NewAppError(
		CodeBadRequest,
		"email already verified",
		err,
	)
}

func NewEmailNotVerified(err error) error {
	return NewAppError(
		CodeBadRequest,
		"email is not verified",
		err,
	)
}
