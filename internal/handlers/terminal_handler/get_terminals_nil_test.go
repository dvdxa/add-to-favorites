package terminal_handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/dvdxa/add-to-favorites/internal/domain"
	repoMock "github.com/dvdxa/add-to-favorites/internal/mocks"
	"github.com/dvdxa/add-to-favorites/internal/services/terminal_service"
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

// test case: terminal_id = 0 & is_favorite=nil, OK
func TestGetTerminalWithFavorites(t *testing.T) {
	log := logger.GetLogger()
	ctl := gomock.NewController(t)
	userRepo := repoMock.NewMockUserRepositoryPort(ctl)
	terminalRepo := repoMock.NewMockTerminalRepositoryPort(ctl)
	userService := user_service.NewUserService(userRepo)
	terminalService := terminal_service.NewTerminalService(terminalRepo)
	h := NewTerminalHandler(*log, terminalService)

	pass := "123456Khalid"
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(pass), 14)
	userID := 1
	favoriteTerminalIDS := []int{1, 2, 3}
	expUser := domain.User{
		ID:       1,
		Name:     "Khalid",
		Password: string(hashBytes),
	}

	user := domain.User{
		Name:     "Khalid",
		Password: "123456Khalid",
	}
	defaultTerminals := []domain.Terminal{
		{
			ID:     1,
			Name:   "terminal1",
			Status: "active",
		},
		{
			ID:     2,
			Name:   "terminal1",
			Status: "active",
		},
		{
			ID:     3,
			Name:   "terminal1",
			Status: "active",
		},
		{
			ID:     4,
			Name:   "terminal1",
			Status: "active",
		},
	}

	userRepo.EXPECT().GetUser(expUser.Name).Return(expUser, nil).Times(1)
	terminalRepo.EXPECT().GetFavoriteTerminalIds(userID).Return(favoriteTerminalIDS, nil)
	terminalRepo.EXPECT().GetDefaultTerminalsList().Return(defaultTerminals, nil)
	token, err := userService.GenerateToken(user)
	require.NoError(t, err)

	body := Request{
		TerminalID: 0,
		IsFavorite: "nil",
	}

	router := gin.Default()
	router.GET("/terminals", func(c *gin.Context) {
		c.Set("userId", float64(userID))
		h.GetTerminalsWithFavorites(c)
	})
	w := httptest.NewRecorder()

	jsonData, err := json.Marshal(body)
	if err != nil {
		require.Error(t, err)
	}
	req := httptest.NewRequest(http.MethodGet, "/terminals", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "Application/Json")
	req.Header.Set("token", token)

	router.ServeHTTP(w, req)
	resp := w.Result()
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	expected :=
		"[{\"id\":1,\"name\":\"terminal1\",\"status\":\"active\",\"is_favorite\":true}," +
			"{\"id\":2,\"name\":\"terminal1\",\"status\":\"active\",\"is_favorite\":true}," +
			"{\"id\":3,\"name\":\"terminal1\",\"status\":\"active\",\"is_favorite\":true},{" +
			"\"id\":4,\"name\":\"terminal1\",\"status\":\"active\",\"is_favorite\":false}]"
	require.Equal(t, expected, string(data))
}

// test case: terminal_id = 0 & is_favorite=nil, UserId Err
func TestGetTerminalWithFavoritesUserIdErr(t *testing.T) {
	log := logger.GetLogger()
	ctl := gomock.NewController(t)
	userRepo := repoMock.NewMockUserRepositoryPort(ctl)
	terminalRepo := repoMock.NewMockTerminalRepositoryPort(ctl)
	userService := user_service.NewUserService(userRepo)
	terminalService := terminal_service.NewTerminalService(terminalRepo)
	h := NewTerminalHandler(*log, terminalService)

	pass := "123456Khalid"
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(pass), 14)
	expUser := domain.User{
		ID:       1,
		Name:     "Khalid",
		Password: string(hashBytes),
	}

	user := domain.User{
		Name:     "Khalid",
		Password: "123456Khalid",
	}

	body := Request{
		TerminalID: 0,
		IsFavorite: "nil",
	}
	userRepo.EXPECT().GetUser(expUser.Name).Return(expUser, nil)
	token, err := userService.GenerateToken(user)
	if err != nil {
		require.Error(t, err)
	}
	router := gin.Default()
	router.GET("/terminals", func(c *gin.Context) {
		h.GetTerminalsWithFavorites(c)
	})
	w := httptest.NewRecorder()

	jsonData, err := json.Marshal(body)
	if err != nil {
		require.Error(t, err)
	}
	req := httptest.NewRequest(http.MethodGet, "/terminals", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "Application/Json")
	req.Header.Set("token", token)

	router.ServeHTTP(w, req)
	resp := w.Result()
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	expected := "{\"err\":\"failed to get user_id\"}"
	require.Equal(t, expected, string(data))
}

// test case: terminal_id = 0 & is_favorite=nil: invalid body
func TestGetTerminalWithFavoritesBadJSON(t *testing.T) {
	log := logger.GetLogger()
	ctl := gomock.NewController(t)
	userRepo := repoMock.NewMockUserRepositoryPort(ctl)
	terminalRepo := repoMock.NewMockTerminalRepositoryPort(ctl)
	userService := user_service.NewUserService(userRepo)
	terminalService := terminal_service.NewTerminalService(terminalRepo)
	h := NewTerminalHandler(*log, terminalService)
	userID := 1
	pass := "123456Khalid"
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(pass), 14)
	expUser := domain.User{
		ID:       1,
		Name:     "Khalid",
		Password: string(hashBytes),
	}

	user := domain.User{
		Name:     "Khalid",
		Password: "123456Khalid",
	}

	body := Request{
		TerminalID: 0,
	}
	userRepo.EXPECT().GetUser(expUser.Name).Return(expUser, nil)
	token, err := userService.GenerateToken(user)
	if err != nil {
		require.Error(t, err)
	}
	router := gin.Default()
	router.GET("/terminals", func(c *gin.Context) {
		c.Set("userId", float64(userID))
		h.GetTerminalsWithFavorites(c)
	})
	w := httptest.NewRecorder()

	jsonData, err := json.Marshal(body)
	if err != nil {
		require.Error(t, err)
	}
	req := httptest.NewRequest(http.MethodGet, "/terminals", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "Application/Json")
	req.Header.Set("token", token)

	router.ServeHTTP(w, req)
	resp := w.Result()
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	expected := "{\"err\":\"Key: 'Request.IsFavorite' Error:Field validation for 'IsFavorite' failed on the 'required' tag\"}"
	require.Equal(t, expected, string(data))
}

// test case: terminal_id = 0 & is_favorite=nil: GetFavoriteTerminalIds error in service layer
func TestGetTerminalWithFavoritesServiceErr(t *testing.T) {
	log := logger.GetLogger()
	ctl := gomock.NewController(t)
	userRepo := repoMock.NewMockUserRepositoryPort(ctl)
	terminalRepo := repoMock.NewMockTerminalRepositoryPort(ctl)
	userService := user_service.NewUserService(userRepo)
	terminalService := terminal_service.NewTerminalService(terminalRepo)
	h := NewTerminalHandler(*log, terminalService)

	pass := "123456Khalid"
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(pass), 14)
	userID := 1
	expUser := domain.User{
		ID:       1,
		Name:     "Khalid",
		Password: string(hashBytes),
	}

	user := domain.User{
		Name:     "Khalid",
		Password: "123456Khalid",
	}

	repoErr := errors.New("DB is down")
	userRepo.EXPECT().GetUser(expUser.Name).Return(expUser, nil).Times(1)
	terminalRepo.EXPECT().GetFavoriteTerminalIds(userID).Return(nil, repoErr)
	token, err := userService.GenerateToken(user)
	require.NoError(t, err)

	body := Request{
		TerminalID: 0,
		IsFavorite: "nil",
	}

	router := gin.Default()
	router.GET("/terminals", func(c *gin.Context) {
		c.Set("userId", float64(userID))
		h.GetTerminalsWithFavorites(c)
	})
	w := httptest.NewRecorder()

	jsonData, err := json.Marshal(body)
	if err != nil {
		require.Error(t, err)
	}
	req := httptest.NewRequest(http.MethodGet, "/terminals", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "Application/Json")
	req.Header.Set("token", token)

	router.ServeHTTP(w, req)
	resp := w.Result()
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	expected := "{\"failed to get user terminal ids\":\"DB is down\"}"
	require.Equal(t, expected, string(data))
}

// test case: terminal_id = 0 & is_favorite=nil SortErr
func TestGetTerminalWithFavoritesSortErr(t *testing.T) {
	log := logger.GetLogger()
	ctl := gomock.NewController(t)
	userRepo := repoMock.NewMockUserRepositoryPort(ctl)
	terminalRepo := repoMock.NewMockTerminalRepositoryPort(ctl)
	userService := user_service.NewUserService(userRepo)
	terminalService := terminal_service.NewTerminalService(terminalRepo)
	h := NewTerminalHandler(*log, terminalService)

	pass := "123456Khalid"
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(pass), 14)
	favoriteTerminalIDs := []int{1, 2, 3}
	userID := 1
	expUser := domain.User{
		ID:       1,
		Name:     "Khalid",
		Password: string(hashBytes),
	}

	user := domain.User{
		Name:     "Khalid",
		Password: "123456Khalid",
	}

	repoErr := errors.New("DB is down")
	userRepo.EXPECT().GetUser(expUser.Name).Return(expUser, nil).Times(1)
	terminalRepo.EXPECT().GetFavoriteTerminalIds(userID).Return(favoriteTerminalIDs, nil)
	terminalRepo.EXPECT().GetDefaultTerminalsList().Return(nil, repoErr)
	token, err := userService.GenerateToken(user)
	require.NoError(t, err)

	body := Request{
		TerminalID: 0,
		IsFavorite: "nil",
	}

	router := gin.Default()
	router.GET("/terminals", func(c *gin.Context) {
		c.Set("userId", float64(userID))

		h.GetTerminalsWithFavorites(c)
	})
	w := httptest.NewRecorder()

	jsonData, err := json.Marshal(body)
	if err != nil {
		require.Error(t, err)
	}
	req := httptest.NewRequest(http.MethodGet, "/terminals", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "Application/Json")
	req.Header.Set("token", token)

	router.ServeHTTP(w, req)
	resp := w.Result()
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	expected := "{\"failed to sort terminals\":\"DB is down\"}"
	require.Equal(t, expected, string(data))
}
