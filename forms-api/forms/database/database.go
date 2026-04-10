package database

import (
	"database/sql"
	"fmt"

	"forms/config"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Connect(cfg *config.Config) error {
	var err error
	DB, err = sql.Open("mysql", cfg.GetDSN())
	if err != nil {
		return fmt.Errorf("error abriendo conexión a la base de datos: %w", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("error conectando a la base de datos: %w", err)
	}

	// Configurar zona horaria de Panamá para todas las conexiones
	_, err = DB.Exec("SET time_zone = '-05:00'")
	if err != nil {
		return fmt.Errorf("error configurando zona horaria: %w", err)
	}

	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)

	return nil
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

