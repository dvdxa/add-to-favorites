package user_service

import (
	"errors"
	repoMock "github.com/dvdxa/add-to-favorites/internal/mocks"
	"github.com/golang-jwt/jwt"
	"github.com/golang/mock/gomock"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func TestParseToken(t *testing.T) {
	ctl := gomock.NewController(t)
	repo := repoMock.NewMockUserRepositoryPort(ctl)
	service := NewUserService(repo)
	err := godotenv.Load("../../.env")
	if err != nil {
		require.Error(t, err)
	}
	var MySigningKey = []byte(os.Getenv("SECRET_KEY"))

	expUserID := float64(1)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"authorized": true,
		"userId":     expUserID,
		"iss":        "jwtgo.io",
		"exp":        time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(MySigningKey)
	if err != nil {
		require.Error(t, err)
	}

	userID, err := service.ParseToken(tokenString)
	require.NoError(t, err)
	require.Equal(t, expUserID, userID)
}

type InvalidClaims struct {
	authorized bool
	userId     int
	iss        string
	exp        int64
}

func (c InvalidClaims) Valid() error {
	return jwt.NewValidationError("invalid claims", jwt.ValidationErrorClaimsInvalid)
}

func TestParseTokenErr(t *testing.T) {
	ctl := gomock.NewController(t)
	repo := repoMock.NewMockUserRepositoryPort(ctl)
	service := NewUserService(repo)
	err := godotenv.Load("../../.env")
	if err != nil {
		require.Error(t, err)
	}
	var MySigningKey = []byte(os.Getenv("SECRET_KEY"))
	mockUserID := 1

	tokenStr1 := ""

	token2 := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"authorized": true,
		"userId":     mockUserID,
		"iss":        "jwtgo.io",
		"exp":        time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenStr2, err := token2.SignedString(MySigningKey)
	if err != nil {
		require.Error(t, err)
	}

	token3 := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"authorized": true,
		"userId":     mockUserID,
		"iss":        "jwtgo.io",
		"exp":        time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenStr3, err := token3.SignedString([]byte("fake_key"))
	if err != nil {
		require.Error(t, err)
	}

	token4 := jwt.NewWithClaims(jwt.SigningMethodHS256, InvalidClaims{
		authorized: true,
		userId:     mockUserID,
		iss:        "jwtgo.io",
		exp:        time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenStr4, err := token4.SignedString(MySigningKey)
	if err != nil {
		require.Error(t, err)
	}
	token5 := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"authorized": true,
		"userId":     mockUserID,
		"iss":        "jwtgo.io",
		"exp":        time.Now().Add(time.Hour * 0).Unix(),
	})
	tokenStr5, err := token5.SignedString(MySigningKey)
	if err != nil {
		require.Error(t, err)
	}

	cases := []struct {
		name   string
		token  string
		expErr error
	}{
		{
			name:   "empty token",
			token:  tokenStr1,
			expErr: errors.New("empty token"),
		},
		{
			name:   "invalid signing method",
			token:  tokenStr2,
			expErr: errors.New("invalid signing method"),
		},
		{
			name:   "invalid token",
			token:  tokenStr3,
			expErr: errors.New("signature is invalid"),
		},
		{
			name:   "invalid claims",
			token:  tokenStr4,
			expErr: errors.New("invalid claims"),
		},
		{
			name:   "expired token",
			token:  tokenStr5,
			expErr: errors.New("token expired"),
		},
	}
	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			_, err = service.ParseToken(tCase.token)
			if tCase.name == "invalid signing method" {
				require.NoError(t, err)
				return
			}
			require.Error(t, err)
			require.EqualError(t, err, tCase.expErr.Error())
		})
	}
}
