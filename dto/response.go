package dto

import (
	"medichat-be/apperror"
	"medichat-be/constants"
)

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func ResponseOk(data any) Response {
	return Response{
		Message: constants.MessageOK,
		Data:    data,
	}
}

func ResponseSeeOther(data any) Response {
	return Response{
		Message: constants.MessageSeeOther,
		Data:    data,
	}
}

func ResponseCreated(data any) Response {
	return Response{
		Message: constants.MessageCreated,
		Data:    data,
	}
}

func ResponseError(appErr *apperror.AppError) Response {
	return Response{
		Message: appErr.Message,
		Data:    nil,
	}
}
