package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"go-parser/internal/dto"
)

// BookRepo is a struct that provides
// all functions to execute SQL queries
// related to 'Book' entity
type BookRepo struct {
	*sql.DB
}

func NewBookRepo(db *sql.DB) *BookRepo {
	return &BookRepo{db}
}

// InsertBooks creates a new books' records in the database
func (br *BookRepo) InsertBooks(ctx context.Context, books []*dto.Book) error {
	tx, err := br.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO books(title, author, price)
		VALUES($1, $2, $3) 
	`
	for _, book := range books {
		result, err := tx.ExecContext(ctx, query, book.Title, book.Author, book.Price)
		if err != nil {
			return err
		}
		rows, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rows != 1 {
			err := fmt.Errorf("expected to affect 1 row, affected %d", rows)
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}
