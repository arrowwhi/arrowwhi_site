package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"net/http"
	"site/database"
	"time"
)

// Ключ для подписи JWT токена
var jwtKey []byte

// Claims Вспомогательная структура для передачи данных в JWT токен
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func SetJwtKey(key string) {
	jwtKey = []byte(key)
}

// Функция для создания JWT токена
func createToken(username string) (string, error) {
	// Устанавливаем время истечения токена через 1 неделю
	expirationTime := time.Now().Add(24 * 7 * time.Hour)

	// Создаем структуру Claims с данными для токена
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Создаем JWT токен с указанным ключом подписи
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// VerifyAndExtractUsername Функция для проверки валидности JWT токена и извлечения атрибута "username"
func VerifyAndExtractUsername(tokenString string) (string, error) {
	// Парсим токен с использованием ключа подписи
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Убедитесь, что метод подписи токена совпадает с тем, который использовался при создании токена
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неверный метод подписи токена")
		}
		return jwtKey, nil
	})

	if err != nil {
		return "", err
	}

	// Проверяем валидность токена
	if !token.Valid {
		return "", fmt.Errorf("неверный токен")
	}

	// Извлекаем атрибут "username" из структуры Claims
	if claims, ok := token.Claims.(*Claims); ok {
		return claims.Username, nil
	} else {
		return "", fmt.Errorf("не удалось извлечь атрибут username из токена")
	}
}

// TakeAuthHandler Обработчик для авторизации
func TakeAuthHandler(c echo.Context) error {

	username := c.FormValue("username")
	pass := c.FormValue("password")

	client, err := database.Get().SelectClientByLogin(username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":       "error",
			"error_reason": "database error",
			"error":        err.Error(),
		})
	}
	if client == nil {
		return c.JSON(http.StatusForbidden, map[string]interface{}{
			"status":       "error",
			"error-reason": "invalid login",
		})
	} else if !ComparePassword(client.Password, pass) {
		return c.JSON(http.StatusForbidden, map[string]interface{}{
			"status":       "error",
			"error-reason": "invalid password",
		})
	}
	// Создаем JWT токен
	tokenString, err := createToken(username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":       "error",
			"error-reason": "error when creating token",
		})
	}
	// Устанавливаем куки с токеном в браузер
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = tokenString
	cookie.Expires = time.Now().Add(24 * 7 * time.Hour)
	c.SetCookie(cookie)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
	})
}

// TakeRegHandler регистрация пользователя
func TakeRegHandler(c echo.Context) error {
	fName := c.FormValue("fName")
	fmt.Println(fName)
	lName := c.FormValue("lName")
	fmt.Println(lName)
	login := c.FormValue("login")
	pass := c.FormValue("password")
	// TODO validate fields
	pass, err := HashPassword(pass)
	if err != nil {
		return c.String(http.StatusServiceUnavailable, err.Error())
	}
	_, err = database.Get().AddClient(fName, lName, login, pass)
	if err != nil {
		return c.String(http.StatusServiceUnavailable, err.Error())
	}

	// Создаем JWT токен
	tokenString, err := createToken(login)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create token")
	}

	// Устанавливаем куки с токеном в браузер
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = tokenString
	cookie.Expires = time.Now().Add(24 * 7 * time.Hour)
	c.SetCookie(cookie)

	return c.Redirect(http.StatusSeeOther, "/")
}
