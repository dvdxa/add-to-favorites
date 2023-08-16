package terminal_service

import (
	"errors"
	"fmt"
	"github.com/dvdxa/add-to-favorites/internal/domain"
	repoMock "github.com/dvdxa/add-to-favorites/internal/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAddToFavorite(t *testing.T) {
	ctl := gomock.NewController(t)
	repo := repoMock.NewMockTerminalRepositoryPort(ctl)
	service := NewTerminalService(repo)
	cases := []struct {
		name       string
		terminalId int
		userId     int
		mockErr    error
	}{
		{
			name:       "invalid_terminal_id",
			terminalId: 4321,
			userId:     4,
			mockErr:    fmt.Errorf("terminal with ID %d doesnt exist in table %s", 4321, "terminals"),
		},
		{
			name:       "terminal_favorited",
			terminalId: 4,
			userId:     4,
			mockErr:    fmt.Errorf("terminal with given ID %d is alredy favorited by user_id %d", 4, 4),
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			repo.EXPECT().AddToFavorites(tCase.terminalId, tCase.userId).Return(tCase.mockErr).Times(1)
			err := service.AddToFavorite(tCase.terminalId, tCase.userId)
			require.Error(t, err)
			require.EqualError(t, err, tCase.mockErr.Error())
		})
	}
}

func TestSortTerminals(t *testing.T) {
	ctl := gomock.NewController(t)
	repo := repoMock.NewMockTerminalRepositoryPort(ctl)
	service := NewTerminalService(repo)

	userTerminalIDs := []int{1, 2, 4}
	mockResp := []domain.Terminal{
		{
			ID:     1,
			Name:   "terminal1",
			Status: "active",
		},
		{
			ID:     2,
			Name:   "terminal2",
			Status: "active",
		},
		{
			ID:     3,
			Name:   "terminal3",
			Status: "active",
		},
		{
			ID:     4,
			Name:   "terminal4",
			Status: "active",
		},
	}

	expTerminals := []domain.FakeTerminal{
		{
			ID:         1,
			Name:       "terminal1",
			Status:     "active",
			IsFavorite: true,
		},
		{
			ID:         2,
			Name:       "terminal2",
			Status:     "active",
			IsFavorite: true,
		},
		{
			ID:         4,
			Name:       "terminal4",
			Status:     "active",
			IsFavorite: true,
		},
		{
			ID:         3,
			Name:       "terminal3",
			Status:     "active",
			IsFavorite: false,
		},
	}
	repo.EXPECT().GetDefaultTerminalsList().Return(mockResp, nil).Times(1)
	terminals, err := service.SortTerminals(userTerminalIDs)
	require.NoError(t, err)
	require.Equal(t, expTerminals, terminals)
}

func TestSortTerminalsRepoErr(t *testing.T) {
	ctl := gomock.NewController(t)
	repo := repoMock.NewMockTerminalRepositoryPort(ctl)
	service := NewTerminalService(repo)

	expErr := errors.New("DB is down")
	repo.EXPECT().GetDefaultTerminalsList().Return(nil, expErr).Times(1)
	userTerminalIDs := []int{1, 2, 4}
	_, err := service.SortTerminals(userTerminalIDs)
	require.Equal(t, expErr, err)
}

func TestGetFavoriteTerminalIds(t *testing.T) {
	ctl := gomock.NewController(t)
	repo := repoMock.NewMockTerminalRepositoryPort(ctl)
	service := NewTerminalService(repo)

	userId := 1
	favoriteTerminalIDs := []int{1, 2, 3}

	repo.EXPECT().GetFavoriteTerminalIds(userId).Return(favoriteTerminalIDs, nil).Times(1)
	ids, err := service.GetFavoriteTerminalIds(userId)
	require.NoError(t, err)
	require.Equal(t, favoriteTerminalIDs, ids)
}

func TestRemoveFromFavoriteTerminal(t *testing.T) {
	ctl := gomock.NewController(t)
	repo := repoMock.NewMockTerminalRepositoryPort(ctl)
	service := NewTerminalService(repo)

	terminalID := 4
	userID := 1

	repo.EXPECT().RemoveFromFavoriteTerminal(terminalID, userID).Return(nil).Times(1)
	err := service.RemoveFromFavoriteTerminal(terminalID, userID)
	require.NoError(t, err)
}
