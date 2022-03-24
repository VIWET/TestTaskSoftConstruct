package repository

import (
	"database/sql"

	"github.com/VIWET/TestTaskSoftConstruct/internal/domain"
)

type PlayerRepository interface {
	GetAllPlayers() ([]*domain.Player, error)
	SetInGameStatus(id int, status int) error
}

type playerRepository struct {
	db *sql.DB
}

func NewPlayerRepository(db *sql.DB) PlayerRepository {
	return &playerRepository{
		db: db,
	}
}

func (r *playerRepository) GetAllPlayers() ([]*domain.Player, error) {
	rows, err := r.db.Query("SELECT id, name FROM players WHERE in_game=FALSE")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var players []*domain.Player

	for rows.Next() {
		p := &domain.Player{}
		err := rows.Scan(&p.ID, &p.Name)
		if err != nil {
			return players, err
		}

		players = append(players, p)
	}

	return players, nil
}

func (r *playerRepository) SetInGameStatus(id int, status int) error {
	res, err := r.db.Exec("UPDATE players SET in_game=? WHERE id=?", status, id)
	if err != nil {
		return err
	}

	c, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if c == 0 {
		// TODO
		return nil
	}
	return nil
}
