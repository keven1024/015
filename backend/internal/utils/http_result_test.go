package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHTTPSuccessHandler(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	data := map[string]interface{}{"result": "success"}
	err := HTTPSuccessHandler(c, data)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	expected := map[string]interface{}{
		"code":    float64(http.StatusOK),
		"message": "success",
		"data":    data,
	}
	assert.Equal(t, expected, response)
}

func TestHTTPErrorHandler(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := HTTPErrorHandler(c, assert.AnError)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	expected := map[string]interface{}{
		"code":    float64(http.StatusBadRequest),
		"message": assert.AnError.Error(),
		"data":    map[string]interface{}{},
	}
	assert.Equal(t, expected, response)
}
