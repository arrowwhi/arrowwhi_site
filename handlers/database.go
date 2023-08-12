package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"site/auth"
	"site/database"
)

func GetMessagesHistory(c echo.Context) error {

	var username = ""
	cookie, err := c.Cookie("token")
	if err == nil {
		username, err = auth.VerifyAndExtractUsername(cookie.Value)
		if err != nil {
			log.Error(err.Error())
			return c.Redirect(http.StatusFound, "/logout")
		}
	} else {
		return c.Redirect(http.StatusFound, "/login")
	}

	input := new(struct {
		Username string `json:"username"`
		LastId   int    `json:"lastId"`
		Count    int    `json:"count"`
	})
	if err := c.Bind(input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}
	msgs := database.Get().SelectMessages(input.Username, username, input.Count, input.LastId)
	// Здесь вы можете использовать message.Username и message.LastID
	// для выполнения нужных действий в вашем приложении
	// Возвращаем ответ
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":   "success",
		"messages": msgs,
	})
}

func TakeFeedback(c echo.Context) error {
	input := new(database.Feedback)
	if err := c.Bind(input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}
	if err := database.Get().AddFeedback(input); err != nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"error": err.Error(),
		})
	}
	return c.JSON(http.StatusCreated, map[string]string{
		"status": "success",
	})
}
