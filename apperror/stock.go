package apperror

func NewStockNotEnough(err error) error {
	return NewAppError(
		CodeBadRequest,
		"stock is not enough",
		err,
	)
}

func NewTransferSameStock(err error) error {
	return NewAppError(
		CodeBadRequest,
		"cannot transfer to the same stock",
		err,
	)
}

func NewTransferDifferentProduct(err error) error {
	return NewAppError(
		CodeBadRequest,
		"source and target stock is not of the same product",
		err,
	)
}

func NewNotPending(err error) error {
	return NewAppError(
		CodeBadRequest,
		"request is already processed (not pending)",
		err,
	)
}
