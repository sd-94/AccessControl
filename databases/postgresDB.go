package databases

import (
	"database/sql"
	"fmt"
	"log"
	"prac/types"

	_ "github.com/lib/pq"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func GetDB(cfg PostgresConfig) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, types.ConnectionError.Wrap(err, "failed to open connection to database")
	}

	if err = db.Ping(); err != nil {
		return nil, types.ConnectionError.Wrap(err, "failed to ping database")
	}

	if err = createAccTable(db); err != nil {
		return nil, types.TableCreationError.Wrap(err, "failed to create accounts table")
	}

	if err = ensureAdminAccount(db); err != nil {
		return nil, types.AdminAccountError.Wrap(err, "failed to ensure SuperUser account exists")
	}

	return db, nil
}

func createAccTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS accounts (
			acc_id SERIAL PRIMARY KEY,
			first_name TEXT,
			last_name TEXT,
			email TEXT,
			password TEXT,
			rights TEXT
		);`

	_, err := db.Exec(query)
	if err != nil {
		return types.TableCreationError.Wrap(err, "error creating accounts table")
	}
	return nil
}

func ensureAdminAccount(db *sql.DB) error {
	email := "superuser@example.com"
	password := "admin123"
	firstName := "SuperUser"
	lastName := ""
	rights := "Superuser"

	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM accounts WHERE email=$1);`
	err := db.QueryRow(query, email).Scan(&exists)
	if err != nil {
		return types.AdminAccountError.Wrap(err, "error checking if admin account exists")
	}

	if !exists {
		insertQuery := `
			INSERT INTO accounts (first_name, last_name, email, password, rights)
			VALUES ($1, $2, $3, $4, $5);
		`
		_, err = db.Exec(insertQuery, firstName, lastName, email, password, rights)
		if err != nil {
			return types.AdminAccountError.Wrap(err, "error creating admin account")
		}
		log.Println("Superuser account created.")
	}

	return nil
}
