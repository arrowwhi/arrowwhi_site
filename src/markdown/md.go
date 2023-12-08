package markdown

import (
	"github.com/labstack/echo/v4"
	"github.com/russross/blackfriday/v2"
	"net/http"
)

func MarkdownHandler(c echo.Context) error {
	// Получите текст в формате Markdown из запроса
	markdownText := c.FormValue("markdown")

	// Преобразуйте Markdown в HTML с помощью библиотеки blackfriday
	html := blackfriday.Run([]byte(markdownText))

	// Отправьте полученный HTML в качестве ответа
	return c.HTML(http.StatusOK, string(html))
}
