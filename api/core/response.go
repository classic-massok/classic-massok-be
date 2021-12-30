package core

import (
	"encoding/json"

	"github.com/labstack/echo/v4"
)

func JSON(c echo.Context, code int, data interface{}, errors ...interface{}) error {
	response := c.Response()
	enc := json.NewEncoder(response)
	header := c.Response().Header()
	header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response.Status = code

	return enc.Encode(&resp{errors, data})
}

type resp struct {
	Errors interface{} `json:"errors,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}
