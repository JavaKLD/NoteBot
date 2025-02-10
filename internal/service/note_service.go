package service

import (
	"log"
	"reminder/internal/database/mysql"
)



func CreateNote(userId int64, content string) error {
	var noteCount int
	queryCount := `SELECT COUNT(*) FROM notes WHERE user_id = ?`
	err := mysql.DB.QueryRow(queryCount, userId).Scan(&noteCount)
	if err != nil {
		log.Fatal("Ошибка при получении количества заметок:", err)
		return err
	}

	// Генерация уникального note_id
	noteId := int64(noteCount + 1)

	// Вставка заметки с уникальным note_id
	query := `INSERT INTO notes (user_id, note_id, content) VALUES (?, ?, ?)`
	_, err = mysql.DB.Exec(query, userId, noteId, content)
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

func DeleteNote(userId int64, noteId int64) error {
	query:= `DELETE FROM notes WHERE user_id = ? AND note_id = ?`

	res, err := mysql.DB.Exec(query, userId, noteId)
	if err != nil {
		log.Fatal("Ошибка при удалении заметки", err)
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatal("Ошибка при проверке удаленный строк", err)
		return err
	}

	if rowsAffected == 0 {
		log.Fatal("Строки не были удалены", err)
		return err
	}

	return nil
}
