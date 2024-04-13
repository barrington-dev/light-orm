package server

import (
	"github.com/labstack/echo/v4"
	"light-orm/internal/services"
	"net/http"
	"strconv"
)

func (s *Server) getUserHandler(c echo.Context) error {
	var (
		statusCode int
		errMessage error
	)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	user, err := services.NewUserService(s.db).GetUser(c.Request().Context(), int64(id))
	if err != nil {
		if err.Error() == "user not found" {
			errMessage = err
			statusCode = http.StatusNotFound
		} else {
			errMessage = err
			statusCode = http.StatusInternalServerError
		}

		return c.JSON(statusCode, map[string]*JSONResponseError[string]{
			"error": NewJSONResponseError(statusCode, errMessage.Error()),
		})
	}

	return c.JSON(http.StatusOK, NewJSONResponseSuccess(user))
}
