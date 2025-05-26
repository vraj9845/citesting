package citesting

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func waitForMySQL(dsn string) error {
	var db *sql.DB
	var err error

	for i := 0; i < 10; i++ {
		db, err = sql.Open("mysql", dsn)
		if err == nil {
			err = db.Ping()
		}
		if err == nil {
			db.Close()
			return nil
		}
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("mysql not ready: %w", err)
}

func TestMySQLBasicFunctionality(t *testing.T) {
	dsn := "testuser:testpass@tcp(localhost:3306)/testdb"

	if err := waitForMySQL(dsn); err != nil {
		t.Fatalf("MySQL did not become ready in time: %v", err)
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("Failed to open DB: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (id INT PRIMARY KEY AUTO_INCREMENT, name VARCHAR(50))`)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	_, err = db.Exec(`INSERT INTO users (name) VALUES ('Alice')`)
	if err != nil {
		t.Fatalf("Failed to insert user: %v", err)
	}

	var name string
	err = db.QueryRow(`SELECT name FROM users WHERE name='Alice'`).Scan(&name)
	if err != nil {
		t.Fatalf("Failed to query user: %v", err)
	}

	if name != "Alice" {
		t.Errorf("Expected 'Alice', got '%s'", name)
	}
}
