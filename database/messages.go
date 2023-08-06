package database

import (
	"github.com/labstack/gommon/log"
	"time"
)

var queueMessages = make(chan Message, 100)

type Message struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	From       string    `json:"from" gorm:"not null"`
	To         string    `json:"to" gorm:"not null"`
	Message    string    `json:"message" gorm:"not null"`
	IsRead     bool      `json:"is_read" gorm:"default:false"`
	CreateDate time.Time `json:"create_date" gorm:"default:current_timestamp"`
}

func (dbe *DbEngine) AddMessage(message, from, to string) Message {
	msg := Message{
		From:    from,
		To:      to,
		Message: message,
	}
	queueMessages <- msg
	return msg
}

func (dbe *DbEngine) processQueue() {
	for {
		msg := <-queueMessages // Получаем сообщение из очереди
		if err := dbe.db.Create(&msg).Error; err != nil {
			log.Error(err.Error())
		}
	}
}

// SelectMessages Метод для выборки сообщений из базы данных
func (dbe *DbEngine) SelectMessages(from, to string, count, firstID int) []Message {
	var messages []Message

	query := dbe.db.
		Where("(`from` = ? AND `to` = ?) OR (`from` = ? AND `to` = ?)", from, to, to, from).
		Where("id < ?", firstID).
		Order("id DESC").
		Limit(count).
		Find(&messages)

	if query.Error != nil {
		log.Error(query.Error)
	}

	return messages
}

func (dbe *DbEngine) FindUsersWithMessages(user string) ([]string, error) {
	var users []string

	// Подзапросы для получения имен пользователей с сообщениями
	subQueryFrom := dbe.db.Model(&Message{}).Select("DISTINCT `from`").Where("to = ?", user)
	subQueryTo := dbe.db.Model(&Message{}).Select("DISTINCT `to`").Where("from = ?", user)

	// Объединение подзапросов с помощью OR
	query := dbe.db.Model(&Message{}).Select("DISTINCT `from`").Where("? IN ?", user, subQueryFrom).
		Or(dbe.db.Model(&Message{}).Select("DISTINCT `from`").Where("? IN ?", user, subQueryTo)).
		Pluck("from", &users)

	if query.Error != nil {
		return nil, query.Error
	}

	return users, nil
}
