package user

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo      Repository
	jwtSecret string
}

func NewService(repo Repository, jwtSecret string) *Service {
	return &Service{repo: repo, jwtSecret: jwtSecret}
}

func (s *Service) CreateUser(req CreateUserRequest) (*UserResponse, error) {
	existing, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrEmailExists
	}

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &User{
		ID:           uuid.New(),
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err = s.repo.Create(user); err != nil {
		return nil, err
	}

	return toResponse(user), nil
}

func (s *Service) GetUsers() ([]UserResponse, error) {
	users, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	responses := make([]UserResponse, len(users))
	for i, u := range users {
		responses[i] = *toResponse(&u)
	}
	return responses, nil
}

func (s *Service) GetUserByID(id uuid.UUID) (*UserResponse, error) {

	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}
	return toResponse(user), nil
}

func (s *Service) GetUserByEmail(email string) (*UserResponse, error) {
	user, err := s.repo.FindByEmail(email)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, nil
	}

	return toResponse(user), nil
}

func (s *Service) UpdateUser(id uuid.UUID, req UpdateUserRequest) (*UserResponse, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}

	if req.Email != "" && req.Email != user.Email {
		existing, err := s.repo.FindByEmail(req.Email)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return nil, ErrEmailExists
		}
		user.Email = req.Email
	}

	if req.Name != "" {
		user.Name = req.Name
	}

	user.UpdatedAt = time.Now()

	if err = s.repo.Update(user); err != nil {
		return nil, err
	}

	return toResponse(user), nil
}

func (s *Service) DeleteUser(id uuid.UUID) (bool, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return false, err
	}
	if user == nil {
		return false, nil
	}

	if err = s.repo.Delete(id); err != nil {
		return false, err
	}

	return true, nil
}

func (s *Service) Login(req LoginRequest) (*TokenResponse, error) {
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	claims := jwt.MapClaims{
		"sub": user.ID.String(),
		"exp": time.Now().Add(24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, err
	}

	return &TokenResponse{Token: signed}, nil
}

func toResponse(u *User) *UserResponse {
	return &UserResponse{
		ID:        u.ID.String(),
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
		UpdatedAt: u.UpdatedAt.Format(time.RFC3339),
	}
}

func hashPassword(p string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
