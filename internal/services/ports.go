package services

import (
	"github.com/dvdxa/add-to-favorites/internal/domain"
	"github.com/dvdxa/add-to-favorites/internal/repositories"
	"github.com/dvdxa/add-to-favorites/internal/services/terminal_service"
	"github.com/dvdxa/add-to-favorites/internal/services/user_service"
)

type UserServicePort interface {
	CreateUser(user domain.User) error
	GenerateToken(user domain.User) (tokenString string, err error)
	ParseToken(tokenStr string) (interface{}, error)
}

type TerminalServicePort interface {
	AddToFavorite(terminalId int, userId int) error
	SortTerminals(userTerminalIDs []int) ([]domain.FakeTerminal, error)
	GetFavoriteTerminalIds(userId int) ([]int, error)
	RemoveFromFavoriteTerminal(terminalID int, userId int) error
}

type ServicePort struct {
	UserServicePort
	TerminalServicePort
}

func NewServicePort(repo *repositories.RepositoryPort) *ServicePort {
	return &ServicePort{
		UserServicePort:     user_service.NewUserService(repo.UserRepositoryPort),
		TerminalServicePort: terminal_service.NewTerminalService(repo.TerminalRepositoryPort),
	}
}
