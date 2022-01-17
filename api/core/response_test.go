package core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func Test_JSON(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	data := struct {
		Name string `json:"name"`
	}{"test"}
	err := fmt.Errorf("this is an error")

	require.NoError(t, JSON(c, http.StatusBadRequest, data, err))
	require.Equal(t, http.StatusBadRequest, res.Result().StatusCode)

	var testRes resp
	expectedRes := resp{
		Data: map[string]interface{}{
			"name": "test",
		},
		Errors: []interface{}{
			"this is an error",
		},
	}

	require.NoError(t, json.NewDecoder(res.Body).Decode(&testRes))
	require.EqualValues(t, expectedRes, testRes)
}

func Test_JSON_error(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	err := fmt.Errorf("some error")
	enc := &json.Encoder{}
	monkey.PatchInstanceMethod(reflect.TypeOf(enc), "Encode", func(_ *json.Encoder, _ interface{}) error {
		return err
	})
	defer monkey.Unpatch(enc)

	require.EqualError(t, JSON(c, http.StatusBadRequest, nil, nil), err.Error())
}
