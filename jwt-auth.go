package pancham

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var secretKey = []byte("pancham")

type CustomClaims struct {
	UserID   primitive.ObjectID `json:"userId"`
	Username string             `json:"username"`
	Email    string             `json:"email"`
	Name     string             `json:"name"`
	Mobile   string             `json:"mobile"`
	jwt.StandardClaims
}

func CreateToken(userID primitive.ObjectID, username string, name string, email string, mobile string) (string, error) {
	// Set up the claims
	claims := CustomClaims{
		UserID:   userID,
		Username: username,
		Name:     name,
		Email:    email,
		Mobile:   mobile,
		StandardClaims: jwt.StandardClaims{
			// Token expiration time (e.g., 24 hours)
			ExpiresAt: time.Now().Add(time.Hour * 240000).Unix(),
		},
	}

	// Create the token with the claims and the secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ValidateToken(tokenString string) (*CustomClaims, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	// Check if the token is valid
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
