package user_service

import (
	"errors"
	"github.com/dvdxa/add-to-favorites/internal/domain"
	repoMock "github.com/dvdxa/add-to-favorites/internal/mocks"
	"github.com/golang/mock/gomock"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

type Token interface {
	SignedString(key []byte) (string, error)
}

func TestGenerateToken(t *testing.T) {
	ctl := gomock.NewController(t)
	repo := repoMock.NewMockUserRepositoryPort(ctl)
	service := NewUserService(repo)
	err := godotenv.Load("../../.env")
	if err != nil {
		require.Error(t, err)
	}

	user := domain.User{
		Name:     "Yahya",
		Password: "12345Yahya",
	}

	expPass := "12345Yahya"
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(expPass), 14)
	if err != nil {
		require.Error(t, err)
	}
	expUser := domain.User{
		ID:       1,
		Name:     "Yahya",
		Password: string(passwordHash),
	}

	repo.EXPECT().GetUser(user.Name).Return(expUser, nil).Times(1)
	_, err = service.GenerateToken(user)
	require.NoError(t, err)
}

func TestGenerateTokenErr(t *testing.T) {
	ctl := gomock.NewController(t)
	repo := repoMock.NewMockUserRepositoryPort(ctl)
	service := NewUserService(repo)

	user := domain.User{
		Name:     "Yahya",
		Password: "123",
	}
	expPass := "12345Yahya"
	passwordHash2, err := bcrypt.GenerateFromPassword([]byte(expPass), 14)
	if err != nil {
		require.Error(t, err)
	}
	expUser := domain.User{
		ID:       1,
		Name:     "Yahya",
		Password: string(passwordHash2),
	}
	expErr := errors.New("crypto/bcrypt: hashedPassword is not the hash of the given password")
	repo.EXPECT().GetUser(user.Name).Return(expUser, nil).Times(1)
	_, err = service.GenerateToken(user)
	require.Error(t, err)
	require.EqualError(t, err, expErr.Error())
}

func TestGenerateTokenRepoErr(t *testing.T) {
	ctl := gomock.NewController(t)
	repo := repoMock.NewMockUserRepositoryPort(ctl)
	service := NewUserService(repo)

	user := domain.User{
		Name:     "Yahya",
		Password: "123",
	}
	expErr := errors.New("DB is down")
	repo.EXPECT().GetUser(user.Name).Return(domain.User{}, expErr).Times(1)
	_, err := service.GenerateToken(user)
	require.Error(t, err)
	require.EqualError(t, err, expErr.Error())
}
