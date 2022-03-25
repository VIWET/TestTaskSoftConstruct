package repository

import (
	"database/sql"

	"github.com/VIWET/TestTaskSoftConstruct/internal/domain"
)

type GameRepository interface {
	GetAllGames() ([]*domain.Game, error)
	GetGame(id int) (*domain.Game, error)
}

type gameRepository struct {
	db *sql.DB
}

func NewGameRepository(db *sql.DB) GameRepository {
	return &gameRepository{
		db: db,
	}
}

func (r *gameRepository) GetAllGames() ([]*domain.Game, error) {
	rows, err := r.db.Query("SELECT id, title FROM games")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []*domain.Game

	for rows.Next() {
		game := &domain.Game{}
		err := rows.Scan(&game.ID, &game.Title)
		if err != nil {
			return games, err
		}

		games = append(games, game)
	}

	return games, nil
}

func (r *gameRepository) GetGame(id int) (*domain.Game, error) {
	game := &domain.Game{}
	err := r.db.QueryRow("SELECT id, title FROM games WHERE id=?", id).Scan(&game.ID, &game.Title)
	if err != nil {
		return nil, err
	}

	return game, nil
}
