package apperror

func NewUpdateCategoryParentRestrict() error {
	return NewAppError(
		CodeBadRequest,
		"can't update category with parent category level 2",
		nil,
	)
}

func NewCreateCategoryParentRestrict() error {
	return NewAppError(
		CodeBadRequest,
		"can't create category with parent category level 2",
		nil,
	)
}
