package service

import (
	"log"
	"reminder/internal/database/mysql"
)



func CreateNote(userId int64, content string) error {
	query := `INSERT INTO notes (user_id, content) VALUES (?, ?)`
	_, err := mysql.DB.Exec(query, userId, content)
	if err != nil {
		log.Fatal("Ошибка при выполнении запроса на создание заметки:", err)
		return err
	}
	return nil
}

func GetNotes(userId int64) ([]string, error) {
	query := `SELECT content FROM notes WHERE user_id = ?`
	rows, err := mysql.DB.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []string
	for rows.Next() {
		var content string
		if err := rows.Scan(&content); err != nil {
			return nil, err
		}
		notes = append(notes, content)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notes, nil
}
