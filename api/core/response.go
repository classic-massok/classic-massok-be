package core

import (
	"encoding/json"

	"github.com/labstack/echo/v4"
)

func JSON(c echo.Context, code int, data interface{}, errors ...error) error {
	response := c.Response()
	header := c.Response().Header()
	header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response.Status = code

	return json.NewEncoder(response).Encode(&resp{stringSliceErrors(errors...), data})
}

type resp struct {
	Errors interface{} `json:"errors,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}

func stringSliceErrors(errors ...error) []string {
	var ss []string
	for _, err := range errors {
		if err == nil {
			continue
		}

		ss = append(ss, err.Error())
	}

	if len(ss) == 0 {
		return nil
	}

	return ss
}
