package database

import (
	"fmt"
	"strings"
	"time"
)

type Client struct {
	Login        string    `gorm:"type:varchar;primaryKey;unique"`
	FirstName    string    `gorm:"type:varchar"`
	LastName     string    `gorm:"type:varchar"`
	Password     string    `gorm:"type:varchar;not null"`
	CreateDate   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	ProfilePhoto string    `gorm:"type:varchar;default null"`
}

func (dbe *DbEngine) SelectClientByLogin(login string) (*Client, error) {
	var client Client
	login = strings.TrimSpace(login)
	if err := dbe.db.Where("login = ?", login).First(&client).Error; err != nil {
		return nil, err
	}
	return &client, nil
}

func (dbe *DbEngine) ChangeProfilePhoto(login, path string) error {
	return dbe.db.Model(&Client{}).Where("login = ?", login).Update("profile_photo", path[10:]).Error
}

func (dbe *DbEngine) UpdateClientPassword(login string, newPass string) error {
	return dbe.db.Model(&Client{}).Where("login = ?", login).Update("password", newPass).Error
}

func (dbe *DbEngine) AddClient(fname, lname, login, password string) (*Client, error) {
	newClient := &Client{
		FirstName:  fname,
		LastName:   lname,
		Login:      login,
		Password:   password,
		CreateDate: time.Now(),
	}

	if err := dbe.db.Create(newClient).Error; err != nil {
		return nil, err
	}
	return newClient, nil
}

func (dbe *DbEngine) GetLogins() ([]string, error) {
	// Выборка всех значений Login из таблицы и сохранение их в массиве
	var logins []string
	if err := dbe.db.Model(&Client{}).Pluck("login", &logins).Error; err != nil {
		return nil, err
	}
	return logins, nil
}

func (dbe *DbEngine) GetLoginsToLine(login string) []map[string]interface{} {
	query := fmt.Sprintf("with unread as ("+
		"select user_from as \"user\", count(*) as unread from messages "+
		"where \"user_to\" = '%s' and is_read = false "+
		"group by \"user_from\" "+
		") "+
		"select case when user_to = '%s' then 'to' else 'from' end as rotation, "+
		"       case when user_to = '%s' then user_from else user_to end as \"user\", "+
		"       message, "+
		"       COALESCE(ur.unread, 0) AS unread, "+
		"       create_date "+
		"from (select usr, MAX(id) as max_id "+
		"      from (select id, user_from as usr "+
		"            from messages "+
		"            where user_to = '%s' "+
		"            union all "+
		"            select id, user_to as usr "+
		"            from messages "+
		"            where user_from = '%s') as mgs "+
		"      group by usr) as uniq_usrs "+
		"         join messages m on m.id = uniq_usrs.max_id "+
		"        left join unread ur on ur.\"user\" = uniq_usrs.usr;", login, login, login, login, login)
	var results []map[string]interface{}
	dbe.db.Raw(query).Scan(&results)
	for _, v := range results {
		fmt.Println(v)
	}
	return results
}
