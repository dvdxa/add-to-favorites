package repositories

import (
	"context"
	"fmt"
	"github.com/dvdxa/add-to-favorites/internal/domain"
	"github.com/dvdxa/add-to-favorites/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TerminalRepository struct {
	pgxpool *pgxpool.Pool
}

func NewTerminalRepository(pgxpool *pgxpool.Pool) *TerminalRepository {
	return &TerminalRepository{
		pgxpool: pgxpool,
	}
}

func (tr *TerminalRepository) AddToFavorites(terminalId int, userId int) error {
	log := logger.GetLogger()
	preCheckQuery := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM %s WHERE id = $1)", "terminals")
	preCheckIfTerminalFavorited := `SELECT COUNT(*) FROM favorite_terminals WHERE user_id = $1 AND $2 = ANY(terminal_id)`
	preCheckUserExists := `SELECT EXISTS (SELECT 1 FROM favorite_terminals WHERE user_id = $1)`
	insertCommand := `INSERT INTO favorite_terminals(user_id, terminal_id) VALUES ($1, ARRAY[$2::integer])`
	updateCommand := `UPDATE favorite_terminals SET terminal_id = array_append(terminal_id, $2) WHERE user_id = $1`
	tx, err := tr.pgxpool.Begin(context.Background())

	if err != nil {
		log.Error(err)
		return fmt.Errorf("failed to begin tx: %v", err)
	}
	var exists bool
	err = tx.QueryRow(context.Background(), preCheckQuery, terminalId).Scan(&exists)
	if err != nil {
		err = tx.Rollback(context.Background())
		if err != nil {
			log.Error(err)
			return fmt.Errorf("failed to rollback tx: %v", err)
		}
		return fmt.Errorf("error executing query: %v", err)
	}
	if !exists {
		return fmt.Errorf("terminal with ID %d doesnt exist in table %s", terminalId, "terminals")
	}

	var count int
	err = tx.QueryRow(context.Background(), preCheckIfTerminalFavorited, userId, terminalId).Scan(&count)
	if err != nil {
		log.Error(err)
		err = tx.Rollback(context.Background())
		if err != nil {
			return fmt.Errorf("failed to rollback tx: %v", err)
		}
		return err
	}
	if count > 0 {
		return fmt.Errorf("terminal with given ID %d is alredy favorited by user_id %d", terminalId, userId)
	}

	//check if user_id exists in favorite_terminals table
	var userExists bool
	err = tx.QueryRow(context.Background(), preCheckUserExists, userId).Scan(&userExists)
	if err != nil {
		log.Error(err)
		err = tx.Rollback(context.Background())
		if err != nil {
			log.Error(err)
			return fmt.Errorf("failed to rollback tx: %v", err)
		}
		return err
	}
	if userExists {
		_, err = tx.Exec(context.Background(), updateCommand, userId, terminalId)
		if err != nil {
			log.Errorf("failed to append element to terminal_id array: %v", err)
			err = tx.Rollback(context.Background())
			if err != nil {
				log.Errorf("failed to rollback: %v", err)
				return fmt.Errorf("failed to rollback tx: %v", err)
			}
			return err
		}
	} else {
		_, err = tx.Exec(context.Background(), insertCommand, userId, terminalId)
		if err != nil {
			log.Error("failed to insert a value to terminal_id array: %v", err)
			err = tx.Rollback(context.Background())
			if err != nil {
				return fmt.Errorf("failed to rollback tx: %v", err)
			}
			return err
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return fmt.Errorf("failed to commit tx: %v", err)
	}
	return nil
}

func (tr *TerminalRepository) GetFavoriteTerminalIds(userId int) ([]int, error) {
	query := `SELECT terminal_id FROM favorite_terminals WHERE user_id = $1;`
	var terminalIDs []int
	row := tr.pgxpool.QueryRow(context.Background(), query, userId)
	_ = row.Scan(&terminalIDs)
	//if err != nil {
	//	log.Errorf("failed to scan row: %v", err)
	//	return nil, err
	//}
	return terminalIDs, nil
}

func (tr *TerminalRepository) GetDefaultTerminalsList() ([]domain.Terminal, error) {
	query := `SELECT * FROM terminals`
	tx, err := tr.pgxpool.Begin(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to begin tx: %v", err)
	}
	rows, err := tx.Query(context.Background(), query)
	if err != nil {
		err = tx.Rollback(context.Background())
		if err != nil {
			return nil, fmt.Errorf("failed to rollback tx: %v", err)
		}
		return nil, err
	}
	defer rows.Close()

	terminals := make([]domain.Terminal, 0)
	for rows.Next() {
		var terminal domain.Terminal
		err = rows.Scan(&terminal.ID, &terminal.Name, &terminal.Status)
		if err != nil {
			err = tx.Rollback(context.Background())
			if err != nil {
				return nil, fmt.Errorf("failed to rollback tx: %v", err)
			}
			return nil, err
		}
		terminals = append(terminals, terminal)
	}
	return terminals, nil
}

func (tr *TerminalRepository) RemoveFromFavoriteTerminal(terminalID int, userId int) error {
	preCheckQuery := `SELECT COUNT(*) FROM favorite_terminals WHERE user_id = $1 AND $2 = ANY(terminal_id)`
	command := `UPDATE favorite_terminals SET terminal_id = array_remove(terminal_id, $2) WHERE user_id = $1`

	tx, err := tr.pgxpool.Begin(context.Background())
	if err != nil {
		return err
	}

	var count int
	err = tx.QueryRow(context.Background(), preCheckQuery, userId, terminalID).Scan(&count)
	if err != nil {
		err = tx.Rollback(context.Background())
		if err != nil {
			return fmt.Errorf("failed to rollback tx: %v", err)
		}
		return fmt.Errorf("error executing query: %v", err)
	}
	if count > 0 {
		_, err = tx.Exec(context.Background(), command, userId, terminalID)
		if err != nil {
			err = tx.Rollback(context.Background())
			if err != nil {
				return fmt.Errorf("failed to rollback tx: %v", err)
			}
			return err
		}
		err = tx.Commit(context.Background())
		if err != nil {
			return fmt.Errorf("failed to commit tx: %v", err)
		}
		return nil
	} else {
		return fmt.Errorf("terminal with %v ID has already removed", terminalID)
	}

}
