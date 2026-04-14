package user

import (
	"errors"
	"testing"
)

type userRepositoryMock struct {
	users []User
	err   error
}

func (m *userRepositoryMock) GetAll() ([]User, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.users, nil
}

func (m *userRepositoryMock) GetByID(id int) (*User, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, u := range m.users {
		if u.ID == id {
			copy := u
			return &copy, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *userRepositoryMock) Create(u User) (*User, error) {
	if m.err != nil {
		return nil, m.err
	}
	u.ID = len(m.users) + 1
	m.users = append(m.users, u)
	return &u, nil
}

func (m *userRepositoryMock) Update(id int, u User) (*User, error) {
	if m.err != nil {
		return nil, m.err
	}
	for i := range m.users {
		if m.users[i].ID == id {
			u.ID = id
			m.users[i] = u
			return &u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *userRepositoryMock) Delete(id int) error {
	if m.err != nil {
		return m.err
	}
	for i := range m.users {
		if m.users[i].ID == id {
			m.users = append(m.users[:i], m.users[i+1:]...)
			return nil
		}
	}
	return errors.New("user not found")
}

func TestCreateAndListUsers(t *testing.T) {
	repo := &userRepositoryMock{}
	service := NewUserService(repo)

	created, err := service.CreateUser(User{Name: "Jon", Email: "jon@example.com"})
	if err != nil {
		t.Fatalf("unexpected error creating user: %v", err)
	}
	if created.ID == 0 {
		t.Fatalf("expected a generated ID")
	}

	users, err := service.ListUsers()
	if err != nil {
		t.Fatalf("unexpected error listing users: %v", err)
	}
	if len(users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(users))
	}
}

func TestDeleteUser(t *testing.T) {
	repo := &userRepositoryMock{users: []User{{ID: 1, Name: "Jon", Email: "jon@example.com"}}}
	service := NewUserService(repo)

	if err := service.DeleteUser(1); err != nil {
		t.Fatalf("unexpected error deleting user: %v", err)
	}

	_, err := service.GetUser(1)
	if err == nil {
		t.Fatalf("expected user to be deleted")
	}
}
