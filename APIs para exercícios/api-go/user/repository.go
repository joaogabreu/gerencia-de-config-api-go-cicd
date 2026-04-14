package user

import (
	"database/sql"
	"errors"
)

type InMemoryUserRepository struct {
	users  []User
	nextID int
}

func NewUserRepository() UserRepository {
	return &InMemoryUserRepository{
		users:  []User{},
		nextID: 1,
	}
}

func (r *InMemoryUserRepository) GetAll() ([]User, error) {
	return r.users, nil
}

func (r *InMemoryUserRepository) GetByID(id int) (*User, error) {
	for _, u := range r.users {
		if u.ID == id {
			return &u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (r *InMemoryUserRepository) Create(u User) (*User, error) {
	u.ID = r.nextID
	r.nextID++
	r.users = append(r.users, u)
	return &u, nil
}

func (r *InMemoryUserRepository) Update(id int, u User) (*User, error) {
	for i, user := range r.users {
		if user.ID == id {
			u.ID = id
			r.users[i] = u
			return &u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (r *InMemoryUserRepository) Delete(id int) error {
	for i, u := range r.users {
		if u.ID == id {
			r.users = append(r.users[:i], r.users[i+1:]...)
			return nil
		}
	}
	return errors.New("user not found")
}

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) UserRepository {
	return &PostgresUserRepository{db: db}
}

func EnsureSchema(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE
		);
	`)

	return err
}

func (r *PostgresUserRepository) GetAll() ([]User, error) {
	rows, err := r.db.Query("SELECT id, name, email FROM users ORDER BY id")
	if err != nil {
		return nil, err
	}

	users := []User{}
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *PostgresUserRepository) GetByID(id int) (*User, error) {
	var u User
	err := r.db.QueryRow("SELECT id, name, email FROM users WHERE id = $1", id).Scan(&u.ID, &u.Name, &u.Email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *PostgresUserRepository) Create(u User) (*User, error) {
	err := r.db.QueryRow(
		"INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id",
		u.Name,
		u.Email,
	).Scan(&u.ID)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *PostgresUserRepository) Update(id int, u User) (*User, error) {
	result, err := r.db.Exec(
		"UPDATE users SET name = $1, email = $2 WHERE id = $3",
		u.Name,
		u.Email,
		id,
	)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, errors.New("user not found")
	}

	u.ID = id
	return &u, nil
}

func (r *PostgresUserRepository) Delete(id int) error {
	result, err := r.db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}