package apperror

func NewPaymentAlreadyUploaded(err error) error {
	return NewAppError(
		CodeBadRequest,
		"payment proof already uploaded",
		err,
	)
}

func NewPaymentNotYetUploaded(err error) error {
	return NewAppError(
		CodeBadRequest,
		"payment proof not yet uploaded",
		err,
	)
}

func NewPaymentAlreadyConfirmed(err error) error {
	return NewAppError(
		CodeBadRequest,
		"payment already confirmed",
		err,
	)
}
