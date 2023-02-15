package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gophermart/internal/storage"
	"log"
	"time"

	"gophermart/internal/config"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type UserDB struct {
	db *sql.DB
}

func NewUserDB(cfg *config.Config) *UserDB {

	ctx := context.Background()

	db, err := sql.Open("pgx", cfg.DB)
	if err != nil {
		log.Println("Не возожно подключиться к бд: ", err)
	}

	err = createAllTablesWithContext(ctx, db)
	if err != nil {
		log.Println(err)
	}

	return &UserDB{
		db: db,
	}
}

func createAllTablesWithContext(ctx context.Context, db *sql.DB) error {

	childCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	queries := []string{
		"CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, firstname VARCHAR(20), lastname VARCHAR(20), login VARCHAR(20), passwd VARCHAR(20), token VARCHAR(50), balance BIGINT);",
		"CREATE TABLE IF NOT EXISTS orders (id SERIAL PRIMARY KEY, order_title VARCHAR(20), user_token VARCHAR(20), balls INT);",
	}

	for _, query := range queries {
		r, err := db.ExecContext(childCtx, query)
		if err != nil {
			return errors.New(fmt.Sprintf("не удалось создать необходимые таблицы в базе данных. \nОшибка: %s\nОтвет базы данных: %s", err, r))
		}
	}

	return nil
}

func (db *UserDB) InsertNewUserWithContext(ctx context.Context, user *storage.User) error {
	childCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if db.db == nil {
		return errors.New("отсутствует открытая база данных")
	}

	query := `INSERT INTO users (id, firstname, lastname, login, passwd)
							VALUES(@id, @firstname, @lastname, @login, @passwd);`

	r, err := db.db.ExecContext(childCtx, query,
		sql.Named("id", user.Firstname),
		sql.Named("firstname", user.Firstname),
		sql.Named("lastname", user.Lastname),
		sql.Named("login", user.Login),
		sql.Named("passwd", user.Passwd),
	)
	if err != nil {
		errors.New(fmt.Sprintf("не удалось отправить данные в базу данных.\n Ошибка: %s\nОтвет базы данных: %s", err, r))
	}

	return nil
}
