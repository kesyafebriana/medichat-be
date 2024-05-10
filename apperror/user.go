package apperror

func NewUserLocationCannotDeleteMain(err error) error {
	return NewAppError(
		CodeBadRequest,
		"cannot delete or unactivate main location",
		err,
	)
}

func NewUserLocationShouldHaveActive(err error) error {
	return NewAppError(
		CodeBadRequest,
		"should have at least one active location",
		err,
	)
}

func NewUserLocationIsNotActive(err error) error {
	return NewAppError(
		CodeBadRequest,
		"location is not active",
		err,
	)
}
