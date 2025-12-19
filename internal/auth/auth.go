package auth

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/fdddf/xcstrings-translator/internal/database"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var JWTSecret = []byte(os.Getenv("JWT_SECRET"))

// Claims represents the JWT claims
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// HashPassword hashes a password with bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPassword compares a password with its hash
func CheckPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// GenerateJWT generates a JWT token for a user
func GenerateJWT(userID uint, username string) (string, error) {
	// Set token expiration (24 hours)
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecret)
}

// ParseJWT parses and validates a JWT token
func ParseJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return JWTSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// RegisterUser creates a new user with the provided details
func RegisterUser(username, email, password string) (*database.User, error) {
	// Hash the password
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	user := &database.User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
		IsActive: true,
	}

	result := database.DB.Create(user)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to create user: %v", result.Error)
	}

	// Clear password for security
	user.Password = ""
	return user, nil
}

// LoginUser authenticates a user and returns the user and JWT token
func LoginUser(username, password string) (*database.User, string, error) {
	var user database.User

	// Find user by username
	result := database.DB.Where("username = ? AND is_active = ?", username, true).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, "", errors.New("invalid username or password")
	}

	if result.Error != nil {
		return nil, "", fmt.Errorf("database error: %v", result.Error)
	}

	// Check password
	if err := CheckPassword(password, user.Password); err != nil {
		return nil, "", errors.New("invalid username or password")
	}

	// Generate JWT token
	token, err := GenerateJWT(user.ID, user.Username)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %v", err)
	}

	// Clear password for security
	user.Password = ""
	return &user, token, nil
}

// GetUserByID retrieves a user by ID
func GetUserByID(userID uint) (*database.User, error) {
	var user database.User
	result := database.DB.First(&user, userID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}

	if result.Error != nil {
		return nil, fmt.Errorf("database error: %v", result.Error)
	}

	// Clear password for security
	user.Password = ""
	return &user, nil
}

// GetUserByUsername retrieves a user by username
func GetUserByUsername(username string) (*database.User, error) {
	var user database.User
	result := database.DB.Where("username = ? AND is_active = ?", username, true).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}

	if result.Error != nil {
		return nil, fmt.Errorf("database error: %v", result.Error)
	}

	// Clear password for security
	user.Password = ""
	return &user, nil
}
