package database

import (
	"time"
)

type Client struct {
	Login      string    `gorm:"type:varchar;primaryKey;unique"`
	Name       string    `gorm:"type:varchar"`
	Password   string    `gorm:"type:varchar;not null"`
	CreateDate time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func (dbe *dbEngine) SelectClientByLogin(login string) (*Client, error) {
	var client Client
	if err := dbe.db.Where("login = ?", login).First(&client).Error; err != nil {
		return nil, err
	}
	return &client, nil
}

func (dbe *dbEngine) UpdateClientPassword(login string, newPass string) error {
	return dbe.db.Model(&Client{}).Where("login = ?", login).Update("password", newPass).Error
}

func (dbe *dbEngine) AddClient(name, login, password string) (*Client, error) {
	newClient := &Client{
		Name:       name,
		Login:      login,
		Password:   password,
		CreateDate: time.Now(),
	}

	if err := dbe.db.Create(newClient).Error; err != nil {
		return nil, err
	}

	return newClient, nil
}

func (dbe *dbEngine) GetAllLogins() ([]string, error) {

	// Выборка всех значений Login из таблицы и сохранение их в массиве
	var logins []string
	if err := dbe.db.Model(&Client{}).Pluck("login", &logins).Error; err != nil {
		return nil, err
	}

	return logins, nil

}
