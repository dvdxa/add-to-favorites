package user_handler

import (
	"bytes"
	"encoding/json"
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

func TestValidateUser(t *testing.T) {
	log := logger.GetLogger()
	ctl := gomock.NewController(t)
	repo := repoMock.NewMockUserRepositoryPort(ctl)
	service := user_service.NewUserService(repo)
	h := NewUserHandler(*log, service)
	pass := "1234567"
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(pass), 14)
	if err != nil {
		require.Error(t, err)
	}
	expUser := domain.User{
		ID:       1,
		Name:     "Bountyhunter",
		Password: string(hashBytes),
	}
	user := domain.User{
		Name:     "Bountyhunter",
		Password: "1234567",
	}
	repo.EXPECT().GetUser(expUser.Name).Return(expUser, nil).Times(1)
	token, err := service.GenerateToken(user)
	if err != nil {
		require.NoError(t, err)
	}
	router := gin.Default()
	router.GET("/terminals", h.ValidateUser)
	w := httptest.NewRecorder()
	reqUser := domain.User{
		Name:     "Bountyhunter",
		Password: "1234567",
	}
	jsonData, err := json.Marshal(reqUser)
	if err != nil {
		require.Error(t, err)
	}
	req := httptest.NewRequest(http.MethodGet, "/terminals", bytes.NewBuffer(jsonData))
	req.Header.Set("token", token)

	router.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	expected := ""
	require.Equal(t, expected, string(data))
}

func TestValidateUserServiceErr(t *testing.T) {
	log := logger.GetLogger()
	ctl := gomock.NewController(t)
	repo := repoMock.NewMockUserRepositoryPort(ctl)
	service := user_service.NewUserService(repo)
	h := NewUserHandler(*log, service)
	token := ""
	router := gin.Default()
	router.GET("/terminals", h.ValidateUser)
	w := httptest.NewRecorder()
	reqUser := domain.User{
		Name:     "Bountyhunter",
		Password: "1234567",
	}
	jsonData, err := json.Marshal(reqUser)
	if err != nil {
		require.Error(t, err)
	}
	req := httptest.NewRequest(http.MethodGet, "/terminals", bytes.NewBuffer(jsonData))
	req.Header.Set("token", token)

	router.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	expected := "{\"err\":\"empty token\"}"
	require.Equal(t, expected, string(data))
}
