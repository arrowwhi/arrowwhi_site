package database

import (
	"github.com/labstack/gommon/log"
	"strings"
	"time"
)

var queueMessages = make(chan Message, 100)

type Message struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	From       string    `json:"user_from" gorm:"not null;column:user_from"`
	To         string    `json:"user_to" gorm:"not null;column:user_to"`
	Message    string    `json:"message" gorm:"not null"`
	IsRead     bool      `json:"is_read" gorm:"default:false"`
	CreateDate time.Time `json:"create_date" gorm:"default:current_timestamp"`
}

func (dbe *DbEngine) AddMessage(message, from, to string) Message {
	msg := Message{
		From:    strings.TrimSpace(from),
		To:      strings.TrimSpace(to),
		Message: strings.TrimSpace(message),
	}
	if err := dbe.db.Create(&msg).Error; err != nil {
		log.Error(err.Error())
	}
	return msg
}

//func (dbe *DbEngine) processQueue() {
//	for {
//		msg := <-queueMessages
//		if err := dbe.db.Create(&msg).Error; err != nil {
//			log.Error(err.Error())
//		}
//	}
//}

// SelectMessages Метод для выборки сообщений из базы данных
func (dbe *DbEngine) SelectMessages(from, to string, count, firstID int) []Message {
	var messages []Message
	from = strings.TrimSpace(from)
	to = strings.TrimSpace(to)
	query := dbe.db.Where("(\"user_from\" = ? AND \"user_to\" = ?) OR (\"user_from\" = ? AND \"user_to\" = ?)", from, to, to, from)
	if firstID != 0 {
		query = query.Where("id < ?", firstID)
	}
	query = query.Order("id DESC")
	if count != 0 {
		query = query.Limit(count)
	}
	query = query.Find(&messages)

	if query.Error != nil {
		log.Error(query.Error)
	}
	return messages
}

//func (dbe *DbEngine) FindUsersWithMessages(user string) ([]string, error) {
//	var users []string
//
//	// Подзапросы для получения имен пользователей с сообщениями
//	subQueryFrom := dbe.db.Model(&Message{}).Select("DISTINCT `from`").Where("to = ?", user)
//	subQueryTo := dbe.db.Model(&Message{}).Select("DISTINCT `to`").Where("from = ?", user)
//
//	// Объединение подзапросов с помощью OR
//	query := dbe.db.Model(&Message{}).Select("DISTINCT `from`").Where("? IN ?", user, subQueryFrom).
//		Or(dbe.db.Model(&Message{}).Select("DISTINCT `from`").Where("? IN ?", user, subQueryTo)).
//		Pluck("from", &users)
//
//	if query.Error != nil {
//		return nil, query.Error
//	}
//	return users, nil
//}

func (dbe *DbEngine) MakeMessagesRead(ids []int) ([]int, error) {
	// Обновляем статус сообщений
	updateQuery := dbe.db.Model(&Message{}).
		Where("id IN (?)", ids).
		Update("is_read", true)

	if updateQuery.Error != nil {
		return nil, updateQuery.Error
	}

	// Получаем обновленные идентификаторы
	var updatedIDs []int
	selectQuery := dbe.db.Model(&Message{}).
		Where("id IN (?) AND is_read = ?", ids, true).
		Pluck("id", &updatedIDs)

	if selectQuery.Error != nil {
		return nil, selectQuery.Error
	}

	return updatedIDs, nil
}

func (dbe *DbEngine) SelectMessageById(id int) (*Message, error) {
	var message Message
	if err := dbe.db.Where("id = ?", id).First(&message).Error; err != nil {
		return nil, err
	}
	return &message, nil
}
