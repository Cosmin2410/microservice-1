package db

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/microservice/server/domain"
)

func SetupDBConn(dsn string) (*sql.DB, error) {
	DBConn, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	fmt.Println("Connected!")

	pingErr := DBConn.Ping()
	if pingErr != nil {
		return nil, pingErr
	}

	return DBConn, nil
}

func Find(db *sql.DB) ([]domain.User, error) {
	query := `SELECT id, name, email FROM users;`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func Insert(db *sql.DB, user domain.User) (int, error) {
	query := `INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id;`
	var id int
	err := db.QueryRow(query, user.Name, user.Email).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error inserting user: %v", err)
	}
	return id, nil
}
