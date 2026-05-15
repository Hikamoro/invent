package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

var (
	DB_HOST     = os.Getenv("DB_HOST")
	DB_PORT     = os.Getenv("DB_PORT")
	DB_USER     = os.Getenv("DB_USER")
	DB_PASSWORD = os.Getenv("DB_PASSWORD")
	DB_NAME     = os.Getenv("DB_NAME")
)

func InitDB() {
	// подключаемся к postgres чтобы создать БД
	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable",
		DB_HOST, DB_PORT, DB_USER, DB_PASSWORD,
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// создаем БД если нет
	_, err = db.Exec("CREATE DATABASE " + DB_NAME)
	if err != nil {
		fmt.Println("База возможно уже существует")
	}

	fmt.Println("Проверка базы данных завершена")

	// подключаемся уже к нужной БД
	psqlInfo = fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME,
	)

	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Подключение к БД успешно")

	createTables()
}

func createTables() {
	// таблица пользователей
	goods := `
	CREATE TABLE IF NOT EXISTS goods(
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) UNIQUE NOT NULL,
		unit VARCHAR(255) NOT NULL,
		quantity INT NOT NULL,
		category_id INT NOT NULL,
		description TEXT
	);
	`
	categories := `
	CREATE TABLE IF NOT EXISTS categories(
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) UNIQUE NOT NULL
	);
	`
	_, err := DB.Exec(goods)
	if err != nil {
		log.Fatal("Ошибка создания таблицы Goods:", err)
	}

	_, err = DB.Exec(categories)
	if err != nil {
		log.Fatal("Ошибка создания таблицы Categories:", err)
	}
	fmt.Println("Таблицы успешно проверены/созданы")
}
