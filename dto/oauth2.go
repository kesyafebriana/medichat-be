package dto

import "medichat-be/domain"

type OAuth2CallbackQuery struct {
	Code  string `form:"code"`
	State string `form:"state"`
}

func (q *OAuth2CallbackQuery) ToOpts() domain.OAuth2CallbackOpts {
	return domain.OAuth2CallbackOpts{
		Code:  q.Code,
		State: q.State,
	}
}
