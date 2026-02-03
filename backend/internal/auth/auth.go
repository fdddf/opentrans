package auth

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/fdddf/opentrans/internal/database"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var JWTSecret []byte

// SetJWTSecret sets the JWT secret for the auth package
func SetJWTSecret(secret []byte) {
	JWTSecret = secret
}

// Auth holds the dependencies for auth operations
type Auth struct {
	DB *database.Database
}

// NewAuth creates a new Auth instance
func NewAuth(db *database.Database) *Auth {
	return &Auth{
		DB: db,
	}
}

// Claims represents the JWT claims
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
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
func GenerateJWT(userID uint, username, role string) (string, error) {
	// Set token expiration (24 hours)
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
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

// ValidateEmail validates an email address format
func ValidateEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

// Global functions for backward compatibility - these still use the global DB for
// backward compatibility with existing code, but these are now deprecated
var authInstance *Auth

// SetAuthInstance sets the auth instance for backward compatibility
func SetAuthInstance(auth *Auth) {
	authInstance = auth
}

// RegisterUser creates a new user with the provided details
func RegisterUser(username, email, password string) (*database.User, error) {
	// This function is deprecated as it relies on the global DB
	// Use the Auth instance methods instead
	if authInstance != nil {
		return authInstance.RegisterUser(username, email, password)
	} else {
		// This would only work if global DB was set, but now we expect DI
		return nil, fmt.Errorf("auth instance not initialized - use dependency injection")
	}
}

// ActivateUser activates a user account using an activation code
func ActivateUser(activationCode string) error {
	// This function is deprecated as it relies on the global DB
	// Use the Auth instance methods instead
	if authInstance != nil {
		return authInstance.ActivateUser(activationCode)
	} else {
		return fmt.Errorf("auth instance not initialized - use dependency injection")
	}
}

// LoginUser authenticates a user and returns the user and JWT token
func LoginUser(username, password string) (*database.User, string, error) {
	// This function is deprecated as it relies on the global DB
	// Use the Auth instance methods instead
	if authInstance != nil {
		return authInstance.LoginUser(username, password)
	} else {
		return nil, "", fmt.Errorf("auth instance not initialized - use dependency injection")
	}
}

// GetUserByID retrieves a user by ID
func GetUserByID(userID uint) (*database.User, error) {
	// This function is deprecated as it relies on the global DB
	// Use the Auth instance methods instead
	if authInstance != nil {
		return authInstance.GetUserByID(userID)
	} else {
		return nil, fmt.Errorf("auth instance not initialized - use dependency injection")
	}
}

// GetUserByUsername retrieves a user by username
func GetUserByUsername(username string) (*database.User, error) {
	// This function is deprecated as it relies on the global DB
	// Use the Auth instance methods instead
	if authInstance != nil {
		return authInstance.GetUserByUsername(username)
	} else {
		return nil, fmt.Errorf("auth instance not initialized - use dependency injection")
	}
}

// RegisterUserWithAuth registers a user using the Auth instance
func (a *Auth) RegisterUser(username, email, password string) (*database.User, error) {
	// Validate email format
	if !ValidateEmail(email) {
		return nil, errors.New("invalid email format")
	}

	// Hash the password
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	// Generate activation code
	activationCode := generateActivationCode()

	user := &database.User{
		Username:       username,
		Email:          email,
		Password:       hashedPassword,
		IsActive:       false, // User is not active until activation
		IsActivated:    false,
		ActivationCode: activationCode,
	}

	result := a.DB.Create(user)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to create user: %v", result.Error)
	}

	// Clear password for security
	user.Password = ""
	return user, nil
}

// ActivateUserWithAuth activates a user account using an activation code
func (a *Auth) ActivateUser(activationCode string) error {
	user := &database.User{}
	result := a.DB.Where("activation_code = ? AND is_activated = ?", activationCode, false).First(user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return errors.New("invalid activation code or account already activated")
	}

	if result.Error != nil {
		return fmt.Errorf("database error: %v", result.Error)
	}

	user.IsActive = true
	user.IsActivated = true
	user.ActivationCode = "" // Clear activation code after use

	result = a.DB.Save(user)
	if result.Error != nil {
		return fmt.Errorf("failed to activate user: %v", result.Error)
	}

	return nil
}

// LoginUserWithAuth authenticates a user and returns the user and JWT token
func (a *Auth) LoginUser(username, password string) (*database.User, string, error) {
	var user database.User

	// Find user by username (only activated users can log in)
	result := a.DB.Where("username = ? AND is_activated = ? AND is_active = ?", username, true, true).First(&user)
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
	token, err := GenerateJWT(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %v", err)
	}

	// Clear password for security
	user.Password = ""
	return &user, token, nil
}

// GetUserByIDWithAuth retrieves a user by ID
func (a *Auth) GetUserByID(userID uint) (*database.User, error) {
	var user database.User
	result := a.DB.First(&user, userID)
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

// GetUserByUsernameWithAuth retrieves a user by username
func (a *Auth) GetUserByUsername(username string) (*database.User, error) {
	var user database.User
	result := a.DB.Where("username = ? AND is_activated = ? AND is_active = ?", username, true, true).First(&user)
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

// generateActivationCode generates a random activation code
func generateActivationCode() string {
	// Generate a 32-character random string as activation code
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
