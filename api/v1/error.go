package v1

import (
	"net/http"

	"github.com/bellamy-labs/auth-api/api/v1/handlers"
	"github.com/labstack/echo/v4"
)

func customHTTPErrorHandler(err error, c echo.Context) {
	e := handlers.ErrorResponse{
		Status:  http.StatusInternalServerError,
		Message: http.StatusText(http.StatusInternalServerError),
	}

	if httpError, ok := err.(*handlers.ErrorResponse); ok {
		e.Status = httpError.Status
		e.Message = httpError.Message
	} else {
		c.Logger().Error(err)
	}

	c.JSON(e.Status, e.Message)
}
