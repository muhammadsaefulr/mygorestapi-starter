package response

import (
	response_dto_auth "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/auth/response"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model/user"
)

type Common struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type SuccessWithCommonData[T any] struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Results []T    `json:"data"`
}

type SuccessWithDetail[T any] struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

type SuccessWithUser struct {
	Code    int        `json:"code"`
	Status  string     `json:"status"`
	Message string     `json:"message"`
	User    model.User `json:"user"`
}

type SuccessWithTokens struct {
	Code    int                      `json:"code"`
	Status  string                   `json:"status"`
	Message string                   `json:"message"`
	User_id string                   `json:"user_id"`
	Tokens  response_dto_auth.Tokens `json:"tokens"`
}

type SuccessWithPaginate[T any] struct {
	Code         int    `json:"code"`
	Status       string `json:"status"`
	Message      string `json:"message"`
	Results      []T    `json:"data"`
	Page         int    `json:"page"`
	Limit        int    `json:"limit"`
	TotalPages   int64  `json:"total_pages"`
	TotalResults int64  `json:"total_results"`
}

type ErrorDetails struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors"`
}
