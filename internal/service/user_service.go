package service

import (
	"log"
	"reminder/internal/database/mysql"
)

func AddUser(userId int64, name string) error {
	query := `INSERT IGNORE INTO users (id, name) VALUES (?, ?)`
	_, err := mysql.DB.Exec(query, userId, name)
	if err != nil {
		log.Fatal("Ошибка добавления пользователя", err)
	}
	return err
}
