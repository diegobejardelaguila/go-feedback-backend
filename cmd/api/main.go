package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/diegobejardelaguila/go-feedback-backend/internal/adapters/db"
	"github.com/diegobejardelaguila/go-feedback-backend/internal/adapters/handlers"
	"github.com/diegobejardelaguila/go-feedback-backend/internal/core/services"
	"github.com/diegobejardelaguila/go-feedback-backend/internal/ports"
)

func main() {
	// Initialize MongoDB connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI environment variable is not set")
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	// Initialize repositories
	userRepo := db.NewMongoUserRepository(client.Database("feedback_app"))
	feedbackRepo := db.NewMongoFeedbackRepository(client.Database("feedback_app"))
	voteRepo := db.NewMongoVoteRepository(client.Database("feedback_app"))

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	// Initialize services
	userService := services.NewUserService(userRepo, jwtSecret)
	feedbackService := services.NewFeedbackService(feedbackRepo)
	voteService := services.NewVoteService(voteRepo)
	analyticsService := services.NewAnalyticsService(feedbackRepo, voteRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	feedbackHandler := handlers.NewFeedbackHandler(feedbackService)
	voteHandler := handlers.NewVoteHandler(voteService)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService)

	// Set up Gin router
	router := gin.Default()

	// Define API routes
	api := router.Group("/api")
	{
		api.POST("/register", userHandler.Register)
		api.POST("/login", userHandler.Login)

		authorized := api.Group("/")
		authorized.Use(authMiddleware(userService))
		{
			authorized.GET("/user", userHandler.GetUser)
			authorized.PUT("/user", userHandler.UpdateUser)

			authorized.POST("/feedback", feedbackHandler.CreateFeedback)
			authorized.GET("/feedback", feedbackHandler.ListFeedback)
			authorized.GET("/feedback/:id", feedbackHandler.GetFeedback)
			authorized.PUT("/feedback/:id", feedbackHandler.UpdateFeedback)
			authorized.DELETE("/feedback/:id", feedbackHandler.DeleteFeedback)

			authorized.POST("/vote", voteHandler.CreateVote)
			authorized.DELETE("/vote/:id", voteHandler.DeleteVote)

			authorized.GET("/analytics", analyticsHandler.GetAnalytics)
		}
	}

	// Start server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}

func authMiddleware(userService ports.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		token := bearerToken[1]
		claims, err := validateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		user, err := userService.GetUser(c, primitive.ObjectIDFromHex(userID))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

func validateToken(token string) (jwt.MapClaims, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

