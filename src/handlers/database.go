package handlers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"io"
	"net/http"
	"os"
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
	fmt.Println(msgs)
	user, err := database.Get().SelectClientByLogin(input.Username)
	if err != nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"status": "error",
			"error":  err.Error(),
		})
	}
	if user.ProfilePhoto == "" {
		user.ProfilePhoto = "/profiles/default.jpg"
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
		"name":   user.FirstName + " " + user.LastName,
		"photo":  user.ProfilePhoto,

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

func GetLogins(c echo.Context) error {
	logins, err := database.Get().GetLogins()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
		"logins": logins,
	})
}

func ChangeProfilePhoto(c echo.Context) error {
	var username = ""
	cookie, err := c.Cookie("token")
	if err == nil {
		username, err = auth.VerifyAndExtractUsername(cookie.Value)
		if err != nil {
			log.Error(err.Error())
			return err
		}
	} else {
		return err
	}

	// Получаем файл из запроса
	file, err := c.FormFile("image")
	if err != nil {
		return c.String(http.StatusBadRequest, "Error reading the file")
	}

	// Открываем файл на диске для записи
	src, err := file.Open()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error opening the file")
	}
	defer src.Close()

	path := "ui/static/profiles/"

	// Создаем файл на диске для сохранения изображения
	dst, err := os.Create(path + username + ".jpg")
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error creating the file")
	}
	defer dst.Close()

	// Копируем данные из файла запроса в файл на диске
	_, err = io.Copy(dst, src)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error copying file data")
	}

	err = database.Get().ChangeProfilePhoto(username, dst.Name())
	if err != nil {
		return err
	}
	if err := database.Get().ChangeProfilePhoto(username, dst.Name()); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
		"path":   dst.Name(),
	})
}

func TakeUserLogins(c echo.Context) error {
	var username = ""
	cookie, err := c.Cookie("token")
	if err == nil {
		username, err = auth.VerifyAndExtractUsername(cookie.Value)
		if err != nil {
			log.Error(err.Error())
			return err
		}
	} else {
		return err
	}
	logins := database.Get().GetLoginsToLine(username)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
		"path":   logins,
	})
}

//func MessagesMakeRead(c echo.Context) error {
//
//	input := new(struct {
//		Ids []int `json:"ids"`
//	})
//
//	if err := c.Bind(input); err != nil {
//		return c.JSON(http.StatusBadRequest, map[string]string{
//			"status": "error",
//			"error":  "Invalid request body",
//		})
//	}
//
//	ids, err := database.Get().MakeMessagesRead(input.Ids)
//	if err != nil {
//		return c.JSON(http.StatusServiceUnavailable, map[string]string{
//			"status": "error",
//			"error":  err.Error(),
//		})
//	}
//
//	//TODO переделать жeсткий костыль
//	for _, v := range ids {
//		msg, err := database.Get().SelectMessageById(v)
//		if err != nil {
//			log.Print(err.Error())
//			continue
//		}
//		chat.SendRead(msg.From, msg.ID)
//		chat.SendRead(msg.To, msg.ID)
//
//	}
//
//	return c.JSON(http.StatusOK, map[string]interface{}{
//		"status": "success",
//		"count":  len(input.Ids),
//		"ids":    ids,
//	})
//}
