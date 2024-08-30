package server

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Model struct {
	db *sqlx.DB
}

func newModel(db *sqlx.DB) Model {
	return Model{
		db: db,
	}
}

func (m *Model) insertByte(b []byte, time int) error {
	_, err := m.db.Exec("INSERT INTO NEWS (NEWS, TIME) VALUES (?, ?)", b, time)
	if err != nil {
		return fmt.Errorf("inserting news byte into db: %w", err)
	}
	return nil
}

func (m *Model) getByte(time int) ([]byte, error) {
	var b []byte
	err := m.db.Get(&b, "SELECT NEWS FROM NEWS WHERE TIME = ?", time)
	if err != nil {
		return nil, fmt.Errorf("selecting news byte from db: %w", err)
	}
	return b, nil
}

func (m *Model) getTimes() ([]int, error) {
	var t []int
	err := m.db.Select(&t, "SELECT TIME FROM NEWS ORDER BY TIME DESC")
	if err != nil {
		return nil, fmt.Errorf("selecting news byte from db: %w", err)
	}
	return t, nil
}
