package exception

import (
	commonentity "PBD_backend_go/commonentity"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

func PanicLogging(err interface{}) {
	if err != nil {
		panic(err)
	}
}

func ErrorHandler(ctx *fiber.Ctx, err error) error {

	_, validationError := err.(ValidationError)
	if validationError {
		data := err.Error()
		var messages []map[string]interface{}

		json.Unmarshal([]byte(data), &messages)
		return ctx.Status(fiber.StatusBadRequest).JSON(commonentity.GeneralResponse{
			Code:    400,
			Message: "Bad Request",
			Data:    err.Error(),
		})
	}

	_, notFoundError := err.(NotFoundError)
	if notFoundError {
		return ctx.Status(fiber.StatusNotFound).JSON(commonentity.GeneralResponse{
			Code:    404,
			Message: "Not Found",
			Data:    err.Error(),
		})
	}

	_, unauthorizedError := err.(UnauthorizedError)
	if unauthorizedError {
		return ctx.Status(fiber.StatusUnauthorized).JSON(commonentity.GeneralResponse{
			Code:    401,
			Message: "Unauthorized",
			Data:    err.Error(),
		})
	}

	_, accessDenialError := err.(AccessDenialError)
	if accessDenialError {
		return ctx.Status(fiber.StatusForbidden).JSON(commonentity.GeneralResponse{
			Code:    403,
			Message: "Access Denial",
			Data:    err.Error(),
		})
	}

	return ctx.Status(fiber.StatusInternalServerError).JSON(commonentity.GeneralResponse{
		Code:    500,
		Message: "General Error",
		Data:    err.Error(),
	})
}

type ValidationError struct {
	Message string
}

func (validationError ValidationError) Error() string {
	return validationError.Message
}

type UnauthorizedError struct {
	Message string
}

func (unauthorizedError UnauthorizedError) Error() string {
	return unauthorizedError.Message
}

type NotFoundError struct {
	Message string
}

func (notFoundError NotFoundError) Error() string {
	return notFoundError.Message
}

type AccessDenialError struct {
	Message string
}

func (accessDenialError AccessDenialError) Error() string {
	return accessDenialError.Message
}
