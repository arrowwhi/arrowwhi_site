package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"html/template"
	"io"
	"net/http"
	"site/auth"
	"site/database"
	"time"
)

// TemplateRenderer представляет собой структуру, реализующую интерфейс echo.Renderer
type TemplateRenderer struct {
	templates *template.Template
}

// Render выполняет рендеринг шаблона и возвращает результат
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func MainHandler(c echo.Context) error {
	var username = ""
	cookie, err := c.Cookie("token")
	if err == nil {
		username, err = auth.VerifyAndExtractUsername(cookie.Value)
	}

	// Данные для передачи в шаблон
	data := map[string]interface{}{
		"title":     "Главная",
		"LoginInfo": username,
	}

	err = renderBase(c, "home.page.tmpl", data)
	if err != nil {
		log.Error(err.Error())
	}
	return err
}

func ChatTemplateHandler(c echo.Context) error {
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
	logins, err := database.Get().GetAllLogins()
	if err != nil {
		log.Error(err.Error())
	}
	data := map[string]interface{}{
		"title":     "Чат",
		"LoginInfo": username,
		"chatId":    "ws://" + c.Request().Host + "/chat_s",
		"logins":    logins,
	}

	err = renderBase(c, "chat.page.tmpl", data)
	if err != nil {
		log.Error(err.Error())
	}
	return err
}

func LoginHandler(c echo.Context) error {

	// Данные для передачи в шаблон
	data := map[string]interface{}{
		"title": "Авторизация",
	}

	err := renderBase(c, "login.page.tmpl", data)
	if err != nil {
		log.Error(err.Error())
	}
	return err
}

func RegHandler(c echo.Context) error {

	// Данные для передачи в шаблон
	data := map[string]interface{}{
		"title": "Регистрация",
	}

	err := renderBase(c, "reg.page.tmpl", data)
	if err != nil {
		log.Error(err.Error())
	}
	return err
}

func ChessHandler(c echo.Context) error {

	// Данные для передачи в шаблон
	data := map[string]interface{}{
		"title": "Шахматы",
	}

	err := renderBase(c, "chess.page.tmpl", data)
	if err != nil {
		log.Error(err.Error())
	}
	return err
}

func LogoutHandler(c echo.Context) error {
	cookie := &http.Cookie{
		Name:    "token",
		Value:   "",
		Expires: time.Now().Add(-1 * time.Hour), // Устанавливаем истекшее время жизни
	}
	c.SetCookie(cookie)
	return c.Redirect(http.StatusFound, "/")
}

func renderBase(c echo.Context, page string, data map[string]interface{}) error {
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseFiles("ui/html/header.partial.tmpl",
			"ui/html/footer.partial.tmpl",
			"ui/html/base.layout.tmpl",
			"ui/html/"+page)),
	}

	// Рендерим шаблон и отправляем результат клиенту
	err := renderer.Render(c.Response().Writer, page, data)
	if err != nil {
		return err
	}
	return err

}
