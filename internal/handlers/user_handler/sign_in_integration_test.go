package user_handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/dvdxa/add-to-favorites/internal/domain"
	"github.com/dvdxa/add-to-favorites/internal/handlers/user_handler"
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

func TestSignIn(t *testing.T) {
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
	repo.EXPECT().GetUser(user.Name).Return(user, nil).Times(1)
	service := user_service.NewUserService(repo)
	h := user_handler.NewUserHandler(*log, service)

	router := gin.Default()
	router.POST("/user/sign-in", h.SignIn)
	w := httptest.NewRecorder()
	reqUser := domain.User{
		Name:     "Khalid",
		Password: "12345Khalid",
	}
	jsonData, err := json.Marshal(reqUser)
	if err != nil {
		require.Error(t, err)
	}
	req := httptest.NewRequest(http.MethodPost, "/user/sign-in", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "Application/Json")

	router.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	expected := "{\"message\":\"access token in header\"}"
	require.Equal(t, expected, string(data))
}

func TestSignInBadJSON(t *testing.T) {
	log := logger.GetLogger()
	ctl := gomock.NewController(t)
	repo := repoMock.NewMockUserRepositoryPort(ctl)
	service := user_service.NewUserService(repo)
	h := user_handler.NewUserHandler(*log, service)

	router := gin.Default()
	router.POST("/user/sign-in", h.SignIn)
	w := httptest.NewRecorder()

	reqUser := domain.User{
		Name: "Khalid",
	}
	jsonData, err := json.Marshal(reqUser)
	if err != nil {
		require.Error(t, err)
	}

	req := httptest.NewRequest(http.MethodPost, "/user/sign-in", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "Application/Json")

	router.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	expected := "{\"err\":\"Key: 'User.Password' Error:Field validation for 'Password' failed on the 'required' tag\"}"
	require.Equal(t, expected, string(data))
}
func TestSignInValidateRequestErr(t *testing.T) {
	log := logger.GetLogger()
	ctl := gomock.NewController(t)
	repo := repoMock.NewMockUserRepositoryPort(ctl)
	service := user_service.NewUserService(repo)
	h := user_handler.NewUserHandler(*log, service)

	router := gin.Default()
	router.POST("/user/sign-in", h.SignIn)
	w := httptest.NewRecorder()
	invalidUser := domain.User{
		Name:     "Tig",
		Password: "123",
	}
	jsonData, err := json.Marshal(invalidUser)
	if err != nil {
		require.Error(t, err)
	}
	req := httptest.NewRequest(http.MethodPost, "/user/sign-in", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "Application/Json")

	router.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	expected := "{\"err\":\"username or password must be at least 5 characters\"}"
	require.Equal(t, expected, string(data))
}

func TestSignInServiceErr(t *testing.T) {
	log := logger.GetLogger()
	ctl := gomock.NewController(t)
	repo := repoMock.NewMockUserRepositoryPort(ctl)
	repoErr := errors.New("can't create user DB is down")
	pass := "12345Tiger"
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(pass), 14)
	if err != nil {
		require.NoError(t, err)
	}
	user := domain.User{
		Name:     "Timbersaw",
		Password: string(hashBytes),
	}
	repo.EXPECT().GetUser(user.Name).Return(domain.User{}, repoErr).Times(1)
	service := user_service.NewUserService(repo)
	h := user_handler.NewUserHandler(*log, service)
	router := gin.Default()
	router.POST("/user/sign-in", h.SignIn)
	w := httptest.NewRecorder()

	reqUser := domain.User{
		Name:     "Timbersaw",
		Password: "12345Tiger",
	}
	jsonData, err := json.Marshal(reqUser)
	if err != nil {
		require.Error(t, err)
	}
	req := httptest.NewRequest(http.MethodPost, "/user/sign-in", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "Application/Json")

	router.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	expected := "{\"err\":\"can't create user DB is down\"}"
	require.Equal(t, expected, string(data))
}
