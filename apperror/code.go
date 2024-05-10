package apperror

const (
	CodeInternal = iota + 1
	CodeCanceled
	CodeBadRequest
	CodeValidationFailed
	CodeConstraintViolation
	CodeNotFound
	CodeAlreadyExists
	CodeUnauthorized
	CodeInvalidToken
	CodeForbidden
)
