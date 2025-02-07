package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"os"
)


var DB *sql.DB

func InitDB() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbname)

	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Ошибка подключения к бд", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Не удалось подключиться к бд", err)
	}

	createUser := `CREATE TABLE IF NOT EXISTS users (
    	id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    	name VARCHAR(255) NOT NULL
	);`

	_, err = DB.Exec(createUser)
	if err != nil {
		log.Fatal("Не удалось создать таблицу users", err)
	}

	createNotes := `CREATE TABLE IF NOT EXISTS notes (
		id INT AUTO_INCREMENT,
		user_id INT NOT NULL,
		content TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(user_id) REFERENCES users(id)
	);`

	_, err = DB.Exec(createNotes)
	if err != nil {
		log.Fatal("Не удалось создать таблицу notes", err)
	}

}
