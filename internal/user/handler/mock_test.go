package handler

import (
	"github.com/google/uuid"
	"github.com/qw-trading/platform/internal/models"
)

type mockUserRepo struct {
	createFn      func(user *models.User) error
	getByIDFn     func(id uuid.UUID) (*models.User, error)
	getByEmailFn  func(email string) (*models.User, error)
	emailExistsFn func(email string) (bool, error)
}

func (m *mockUserRepo) Create(user *models.User) error {
	if m.createFn != nil {
		return m.createFn(user)
	}
	return nil
}

func (m *mockUserRepo) GetByID(id uuid.UUID) (*models.User, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(id)
	}
	return nil, nil
}

func (m *mockUserRepo) GetByEmail(email string) (*models.User, error) {
	if m.getByEmailFn != nil {
		return m.getByEmailFn(email)
	}
	return nil, nil
}

func (m *mockUserRepo) EmailExists(email string) (bool, error) {
	if m.emailExistsFn != nil {
		return m.emailExistsFn(email)
	}
	return false, nil
}
