package user_handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/dvdxa/add-to-favorites/internal/domain"
	repoMock "github.com/dvdxa/add-to-favorites/internal/mocks"
	"github.com/dvdxa/add-to-favorites/internal/services/user_service"
	"github.com/dvdxa/add-to-favorites/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSignUp(t *testing.T) {
	log := logger.GetLogger()
	ctl := gomock.NewController(t)
	repo := repoMock.NewMockUserRepositoryPort(ctl)
	pass := "12345Khalid"
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(pass), 14)
	if err != nil {
		require.Error(t, err)
	}
	user := domain.User{
		Name:     "Khalid",
		Password: string(hashBytes),
	}
	repo.EXPECT().CreateUser(gomock.AssignableToTypeOf(user)).Return(nil).Times(1)
	service := user_service.NewUserService(repo)
	h := NewUserHandler(*log, service)

	router := gin.Default()
	router.POST("/user/sign-up", h.SignUp)
	w := httptest.NewRecorder()
	reqUser := domain.User{
		Name:     "Khalid",
		Password: "12345Khalid",
	}
	jsonData, err := json.Marshal(reqUser)
	if err != nil {
		require.Error(t, err)
	}
	req := httptest.NewRequest(http.MethodPost, "/user/sign-up", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "Application/Json")

	router.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	expected := "\"user created\""
	require.Equal(t, expected, string(data))
}

func TestSignUpBadJSON(t *testing.T) {
	log := logger.GetLogger()
	ctl := gomock.NewController(t)
	repo := repoMock.NewMockUserRepositoryPort(ctl)
	service := user_service.NewUserService(repo)
	h := NewUserHandler(*log, service)

	router := gin.Default()
	router.POST("/user/sign-up", h.SignUp)
	w := httptest.NewRecorder()
	invalidUser := domain.User{
		Name: "Tiger",
	}
	jsonData, err := json.Marshal(invalidUser)
	if err != nil {
		require.Error(t, err)
	}
	req := httptest.NewRequest(http.MethodPost, "/user/sign-up", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "Application/Json")

	router.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	expected := "{\"err\":\"Key: 'User.Password' Error:Field validation for 'Password' failed on the 'required' tag\"}"
	require.Equal(t, expected, string(data))
}

func TestSignUpValidateRequestErr(t *testing.T) {
	log := logger.GetLogger()
	ctl := gomock.NewController(t)
	repo := repoMock.NewMockUserRepositoryPort(ctl)
	service := user_service.NewUserService(repo)
	h := NewUserHandler(*log, service)

	router := gin.Default()
	router.POST("/user/sign-up", h.SignUp)
	w := httptest.NewRecorder()
	invalidUser := domain.User{
		Name:     "Tig",
		Password: "123",
	}
	jsonData, err := json.Marshal(invalidUser)
	if err != nil {
		require.Error(t, err)
	}
	req := httptest.NewRequest(http.MethodPost, "/user/sign-up", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "Application/Json")

	router.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	expected := "{\"err\":\"username or password must be at least 5 characters\"}"
	require.Equal(t, expected, string(data))
}

func TestSignUpServiceErr(t *testing.T) {
	log := logger.GetLogger()
	ctl := gomock.NewController(t)
	repo := repoMock.NewMockUserRepositoryPort(ctl)
	repoErr := errors.New("can't create user DB is down")
	user := domain.User{
		Name:     "Timbersaw",
		Password: "12345Tiger",
	}
	repo.EXPECT().CreateUser(gomock.AssignableToTypeOf(user)).Return(repoErr).Times(1)
	service := user_service.NewUserService(repo)
	h := NewUserHandler(*log, service)
	router := gin.Default()
	router.POST("/user/sign-up", h.SignUp)
	w := httptest.NewRecorder()

	jsonData, err := json.Marshal(user)
	if err != nil {
		require.Error(t, err)
	}
	req := httptest.NewRequest(http.MethodPost, "/user/sign-up", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "Application/Json")

	router.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	expected := "{\"err\":\"can't create user DB is down\"}"
	require.Equal(t, expected, string(data))
}
