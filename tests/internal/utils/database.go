package utils

import (
	"database/sql"
	"fmt"
)

func AssertTestDatabase(host, user, password, dbname, port, sslmode string) error {
	connStr := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=America/Sao_Paulo",
		host,
		user,
		password,
		dbname,
		port,
		sslmode,
	)

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return err
	}

	defer db.Close()

	db.Exec(fmt.Sprintf("CREATE DATABASE \"%s_test\"", dbname))

	return err
}
