package errors

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	ErrNotFound = echo.NewHTTPError(
		http.StatusNotFound,
		"ERR_NOT_FOUND",
	)
	ErrBearerTokenInvalid = echo.NewHTTPError(
		http.StatusUnauthorized,
		"ERR_BEARER_TOKEN_INVALID",
	)
	ErrInsufficientPermission = echo.NewHTTPError(
		http.StatusUnauthorized,
		"ERR_INSUFFICIENT_PERMISSION",
	)
	ErrBadRequest    = echo.NewHTTPError(http.StatusBadRequest, "ERR_BAD_REQUEST")
	ErrDuplicated    = echo.NewHTTPError(http.StatusBadRequest, "ERR_DUPLICATED")
	ErrInternalError = echo.NewHTTPError(http.StatusInternalServerError, "ERR_INTERNAL_SERVER")
)
