package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kimnguyenlong/ketoz/internal/entity"
)

type Response[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

type ListResponse[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Records []T    `json:"records"`
}

func responseError(c *fiber.Ctx, err error) error {
	b, ok := err.(*entity.Error)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(
			Response[*struct{}]{
				Code:    entity.CodeInternalError,
				Message: entity.MessageInternalError,
				Data:    nil,
			},
		)
	}

	return c.Status(b.HttpCode).JSON(
		Response[*struct{}]{
			Code:    b.Code,
			Message: b.Message,
			Data:    nil,
		},
	)
}

func responseNilData(c *fiber.Ctx, code int) error {
	return c.Status(code).JSON(Response[*struct{}]{
		Code:    entity.CodeSuccess,
		Message: entity.MessageSuccess,
		Data:    nil,
	})
}

func responseData[T any](c *fiber.Ctx, code int, data T) error {
	return c.Status(code).JSON(Response[T]{
		Code:    entity.CodeSuccess,
		Message: entity.MessageSuccess,
		Data:    data,
	})
}

func responseRecords[T any](c *fiber.Ctx, code int, records []T) error {
	return c.Status(code).JSON(ListResponse[T]{
		Code:    entity.CodeSuccess,
		Message: entity.MessageSuccess,
		Records: records,
	})
}
