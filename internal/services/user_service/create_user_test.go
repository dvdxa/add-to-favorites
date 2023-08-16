package user_service

import (
	"errors"
	"github.com/dvdxa/add-to-favorites/internal/domain"
	repoMock "github.com/dvdxa/add-to-favorites/internal/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestCreateUser(t *testing.T) {
	ctl := gomock.NewController(t)
	repo := repoMock.NewMockUserRepositoryPort(ctl)
	service := NewUserService(repo)

	password := "12345yusuf"
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		require.Error(t, err)
	}

	mockUser := domain.User{
		Name:     "Yusuf",
		Password: string(passwordHash),
	}

	expUser := domain.User{
		Name:     "Yusuf",
		Password: "12345yusuf",
	}

	repo.EXPECT().CreateUser(gomock.AssignableToTypeOf(mockUser)).Return(nil).Times(1)
	err = service.CreateUser(expUser)
	require.NoError(t, err)
}

func TestCreateUserPassErr(t *testing.T) {
	ctl := gomock.NewController(t)
	repo := repoMock.NewMockUserRepositoryPort(ctl)
	service := NewUserService(repo)

	pass := "123456Yusuf"
	bytesHash, err := service.HashPassword(pass, 14)
	if err != nil {
		require.Error(t, err)
	}

	expUser := domain.User{
		Name:     "Yusuf",
		Password: string(bytesHash),
	}

	mockUser := domain.User{
		Name:     "Yusuf",
		Password: "123456Yusuf",
	}
	expErr := errors.New("DB is down")
	repo.EXPECT().CreateUser(gomock.AssignableToTypeOf(expUser)).Return(expErr).Times(1)
	err = service.CreateUser(mockUser)
	require.EqualError(t, err, expErr.Error())
}
