package terminal_service

import (
	"github.com/dvdxa/add-to-favorites/internal/domain"
	"github.com/dvdxa/add-to-favorites/internal/repositories"
	"sort"
)

type TerminalService struct {
	terminalRepositoryPort repositories.TerminalRepositoryPort
}

func NewTerminalService(terminalRepositoryPort repositories.TerminalRepositoryPort) *TerminalService {
	return &TerminalService{
		terminalRepositoryPort: terminalRepositoryPort,
	}
}

func (ts *TerminalService) AddToFavorite(terminalId int, userId int) error {
	return ts.terminalRepositoryPort.AddToFavorites(terminalId, userId)
}

func (ts *TerminalService) SortTerminals(userTerminalIDs []int) ([]domain.FakeTerminal, error) {
	idIndexMap := make(map[int]int)
	for i, id := range userTerminalIDs {
		idIndexMap[id] = i
	}
	terminals, err := ts.terminalRepositoryPort.GetDefaultTerminalsList()
	if err != nil {
		return nil, err
	}
	joinTerminals := make([]domain.FakeTerminal, len(terminals))
	for i, val := range terminals {
		joinTerminals[i] = ConvertToFakeTerminal(val)
		joinTerminals[i].IsFavorite = false
		for _, userTerminalID := range userTerminalIDs {
			if userTerminalID == val.ID {
				joinTerminals[i].IsFavorite = true
			}
		}
	}
	sort.Slice(joinTerminals, func(i, j int) bool {
		_, iExists := idIndexMap[joinTerminals[i].ID]
		_, jExists := idIndexMap[joinTerminals[j].ID]

		if iExists == jExists {
			return false
		}
		return iExists
	})
	return joinTerminals, nil
}

func (ts *TerminalService) GetFavoriteTerminalIds(userId int) ([]int, error) {
	return ts.terminalRepositoryPort.GetFavoriteTerminalIds(userId)
}

func (ts *TerminalService) RemoveFromFavoriteTerminal(terminalID int, userId int) error {
	return ts.terminalRepositoryPort.RemoveFromFavoriteTerminal(terminalID, userId)
}

func ConvertToFakeTerminal(terminal domain.Terminal) domain.FakeTerminal {
	fakeTerminal := domain.FakeTerminal{
		ID:     terminal.ID,
		Name:   terminal.Name,
		Status: terminal.Status,
	}
	return fakeTerminal
}
