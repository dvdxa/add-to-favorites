package repositories

import (
	"context"
	"errors"
	"fmt"
	"github.com/dvdxa/add-to-favorites/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pgxpool *pgxpool.Pool
}

func NewUserRepository(pgxpool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		pgxpool: pgxpool,
	}
}
func (ur *UserRepository) CreateUser(user domain.User) error {
	existingUserQuery := `SELECT COUNT(*) FROM users WHERE name = $1`
	var count int
	err := ur.pgxpool.QueryRow(context.Background(), existingUserQuery, user.Name).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("username already exists")
	}
	command := `INSERT INTO users (name, password) VALUES ($1, $2)`
	_, err = ur.pgxpool.Exec(context.Background(), command, user.Name, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) GetUser(username string) (domain.User, error) {
	var user domain.User
	query := `SELECT * FROM users WHERE name = $1`
	err := ur.pgxpool.QueryRow(context.Background(), query, username).Scan(&user.ID, &user.Name, &user.Password)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.User{}, errors.New("no user found with given name")
		}
		return domain.User{}, err
	}
	return user, nil
}
