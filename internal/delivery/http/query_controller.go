package http

import (
	"database-remote-commander/internal/model"
	"database-remote-commander/internal/usecase"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type QueryController struct {
	Log     *logrus.Logger
	UseCase *usecase.RemoteQueryUseCase
}

func NewQueryController(useCase *usecase.RemoteQueryUseCase, logger *logrus.Logger) *QueryController {
	return &QueryController{
		UseCase: useCase,
		Log:     logger,
	}
}

func (c *QueryController) ExecQuery(ctx *fiber.Ctx) error {
	request := new(model.QueryRequest)

	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}

	result, err := c.UseCase.ExecQuery(ctx.UserContext(), request)

	if err != nil {
		errorMessage := fmt.Sprintf("Failed to execute query: %+v", err)
		c.Log.Warn(errorMessage)
		return errors.New(errorMessage)
	}

	return ctx.JSON(fiber.Map{"command": request.Query, "result": result})
}
