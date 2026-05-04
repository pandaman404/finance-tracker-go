package user

import (
	"errors"
	"sync"
)

type MockRepository struct {
	data []User
	mu   sync.Mutex
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		data: []User{},
	}
}

func (m *MockRepository) Create(user *User) error {
	m.mu.Lock() // mock en memoria -> tu controlas concurrencia
	defer m.mu.Unlock()

	// validar email único
	for _, u := range m.data {
		if u.Email == user.Email {
			return errors.New("email already exists")
		}
	}

	m.data = append(m.data, *user)
	return nil
}

func (m *MockRepository) FindAll() ([]User, error) {
	return m.data, nil
}

func (m *MockRepository) FindByEmail(email string) (*User, error) {
	for _, u := range m.data {
		if u.Email == email {
			return &u, nil
		}
	}
	return nil, nil
}
