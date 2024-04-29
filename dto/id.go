package dto

type IDPathRequest struct {
	ID int64 `uri:"id" binding:"required"`
}
