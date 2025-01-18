package services

import (
	"context"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/yourusername/go-feedback-backend/internal/core/domain"
	"github.com/yourusername/go-feedback-backend/internal/ports"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepo  ports.UserRepository
	jwtSecret []byte
}

func NewUserService(userRepo ports.UserRepository, jwtSecret string) ports.UserService {
	return &userService{
		userRepo:  userRepo,
		jwtSecret: []byte(jwtSecret),
	}
}

func (s *userService) Register(ctx context.Context, user *domain.User) error {
	// Check if user already exists
	existingUser, _ := s.userRepo.GetByEmail(ctx, user.Email)
	if existingUser != nil {
		return errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	return s.userRepo.Create(ctx, user)
}

func (s *userService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	return s.GenerateToken(user)
}

func (s *userService) GenerateToken(user *domain.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.Hex(),
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString(s.jwtSecret)
}

func (s *userService) GetUser(ctx context.Context, id primitive.ObjectID) (*domain.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *userService) UpdateUser(ctx context.Context, user *domain.User) error {
	user.UpdatedAt = time.Now()
	return s.userRepo.Update(ctx, user)
}

func (s *userService) GetUserByToken(token string) (*domain.User, error) {
	// This is a placeholder and should be replaced with actual JWT validation and user retrieval
	return &domain.User{ID: primitive.NewObjectID(), Email: "dummy@example.com"}, nil
}

