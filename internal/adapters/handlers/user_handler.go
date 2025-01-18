package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/go-feedback-backend/internal/core/domain"
	"github.com/yourusername/go-feedback-backend/internal/ports"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserHandler struct {
	userService ports.UserService
}

func NewUserHandler(userService ports.UserService) *UserHandler {
	return &UserHandler{userService: userService}

