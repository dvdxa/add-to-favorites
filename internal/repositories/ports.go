package repositories

import (
	"github.com/dvdxa/add-to-favorites/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepositoryPort interface {
	CreateUser(user domain.User) error
	GetUser(username string) (domain.User, error)
}

type TerminalRepositoryPort interface {
	AddToFavorites(terminalId int, userId int) error
	GetFavoriteTerminalIds(userId int) ([]int, error)
	GetDefaultTerminalsList() ([]domain.Terminal, error)
	RemoveFromFavoriteTerminal(terminalID int, userId int) error
}

type RepositoryPort struct {
	UserRepositoryPort
	TerminalRepositoryPort
}

func NewRepositoryPort(pgx *pgxpool.Pool) *RepositoryPort {
	return &RepositoryPort{
		UserRepositoryPort:     NewUserRepository(pgx),
		TerminalRepositoryPort: NewTerminalRepository(pgx),
	}
}
