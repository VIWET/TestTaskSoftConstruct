package repository

import (
	"database/sql"

	"github.com/VIWET/TestTaskSoftConstruct/internal/domain"
)

type playerRepository struct {
	db *sql.DB
}

func NewPlayerRepository(db *sql.DB) *playerRepository {
	return &playerRepository{
		db: db,
	}
}

func (r *playerRepository) GetAllPlayers() ([]*domain.Player, error) {
	// row, err := r.db.QueryRow()
	return nil, nil
}
