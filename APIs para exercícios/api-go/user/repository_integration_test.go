package user

import (
	"database/sql"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

func TestPostgresRepositoryCRUD(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL not set")
	}

	db := openDBWithRetry(t, dsn)
	defer func() {
		_ = db.Close()
	}()

	if err := EnsureSchema(db); err != nil {
		t.Fatalf("failed to ensure schema: %v", err)
	}

	if _, err := db.Exec("TRUNCATE TABLE users RESTART IDENTITY"); err != nil {
		t.Fatalf("failed to truncate users table: %v", err)
	}

	repo := NewPostgresUserRepository(db)

	created, err := repo.Create(User{Name: "Ada", Email: "ada@example.com"})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	fetched, err := repo.GetByID(created.ID)
	if err != nil {
		t.Fatalf("get by id failed: %v", err)
	}
	if fetched.Email != "ada@example.com" {
		t.Fatalf("unexpected email: %s", fetched.Email)
	}

	updated, err := repo.Update(created.ID, User{Name: "Ada Lovelace", Email: "ada.lovelace@example.com"})
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}
	if updated.Name != "Ada Lovelace" {
		t.Fatalf("unexpected updated name: %s", updated.Name)
	}

	users, err := repo.GetAll()
	if err != nil {
		t.Fatalf("get all failed: %v", err)
	}
	if len(users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(users))
	}

	if err := repo.Delete(created.ID); err != nil {
		t.Fatalf("delete failed: %v", err)
	}
}

func openDBWithRetry(t *testing.T, dsn string) *sql.DB {
	t.Helper()

	var lastErr error
	for attempt := 1; attempt <= 20; attempt++ {
		db, err := sql.Open("postgres", dsn)
		if err != nil {
			lastErr = err
			time.Sleep(1 * time.Second)
			continue
		}

		if err := db.Ping(); err != nil {
			lastErr = err
			_ = db.Close()
			time.Sleep(1 * time.Second)
			continue
		}

		return db
	}

	t.Fatalf("database not available after retries: %v", lastErr)
	return nil
}
