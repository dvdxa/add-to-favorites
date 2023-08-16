package user_service

import (
	"errors"
	"github.com/dvdxa/add-to-favorites/internal/domain"
	"github.com/dvdxa/add-to-favorites/internal/repositories"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

var (
	ErrEmptyToken           = errors.New("empty token")
	ErrInvalidSigningMethod = errors.New("invalid signing method")
	ErrInvalidClaims        = errors.New("invalid claims")
	ErrTokenExpired         = errors.New("token expired")
	ErrInvalidToken         = errors.New("signature is invalid")
)

type UserService struct {
	userRepositoryPort repositories.UserRepositoryPort
}

func NewUserService(userRepositoryPort repositories.UserRepositoryPort) *UserService {
	return &UserService{
		userRepositoryPort: userRepositoryPort,
	}
}

func (us *UserService) CreateUser(user domain.User) error {
	passHash, _ := us.HashPassword(user.Password, 14)
	user.Password = string(passHash)
	err := us.userRepositoryPort.CreateUser(user)
	if err != nil {
		return err
	}
	return nil
}

func (us *UserService) GenerateToken(user domain.User) (tokenString string, err error) {
	var MySigningKey = []byte(os.Getenv("SECRET_KEY"))
	exUser, err := us.userRepositoryPort.GetUser(user.Name)
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(exUser.Password), []byte(user.Password))
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"authorized": true,
		"userId":     exUser.ID,
		"iss":        "jwtgo.io",
		"exp":        time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err = token.SignedString(MySigningKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (us *UserService) ParseToken(tokenStr string) (interface{}, error) {
	var MySigningKey = []byte(os.Getenv("SECRET_KEY"))
	if len(tokenStr) == 0 {
		return 0, ErrEmptyToken
	}
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSigningMethod
		}
		return MySigningKey, nil
	})
	if err != nil {
		return 0, err
	}
	if !token.Valid {
		return 0, ErrInvalidToken
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, ErrInvalidClaims
	}
	exp, ok := claims["exp"].(float64)
	if !ok {
		return 0, ErrInvalidClaims
	}
	expTime := time.Unix(int64(exp), 0)
	if time.Now().After(expTime) {
		return 0, ErrTokenExpired
	}
	userId := claims["userId"]
	return userId, nil
}

func (us *UserService) HashPassword(password string, cost int) ([]byte, error) {
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return nil, err
	}
	return passHash, nil
}
