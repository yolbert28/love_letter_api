package store

import (
	"database/sql"

	"github.com/yolbert28/deliver_love_letter/internal/domain"
)

type letterRepository struct {
	db *sql.DB
}

func NewLetterRepository(db *sql.DB) domain.LetterRepository {
	return &letterRepository{db: db}
}

func (r *letterRepository) Create(letter domain.Letter) (domain.Letter, error) {

	err := r.db.QueryRow(`INSERT INTO letter(content, show_date) VALUES ($1, $2) RETURNING id, content, opened_count, COALESCE(last_opened::text, ''), show_date
	`, letter.Content, letter.Date).
		Scan(
			&letter.ID,
			&letter.Content,
			&letter.OpenedCount,
			&letter.LastOpened,
			&letter.Date,
		)

	if err != nil {
		return domain.Letter{}, err
	}

	return letter, nil
}

func (r *letterRepository) GetAll() ([]domain.Letter, error) {
	rows, err := r.db.Query(`
		SELECT id, content, opened_count, COALESCE(last_opened::text, ''), show_date 
		FROM letter 
		ORDER BY show_date DESC
	`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	letters := make([]domain.Letter, 0)
	for rows.Next() {
		var l domain.Letter
		if err := rows.Scan(&l.ID, &l.Content, &l.OpenedCount, &l.LastOpened, &l.Date); err != nil {
			return nil, err
		}
		letters = append(letters, l)
	}
	return letters, rows.Err()
}

func (r *letterRepository) GetByDate(date string) (*domain.Letter, error) {
	var letter domain.Letter

	err := r.db.QueryRow(`
		SELECT id, content, opened_count, COALESCE(last_opened::text, ''), show_date 
		FROM letter 
		WHERE show_date = $1
	`, date).Scan(
		&letter.ID,
		&letter.Content,
		&letter.OpenedCount,
		&letter.LastOpened,
		&letter.Date,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &letter, nil
}

func (r *letterRepository) IncrementTapCount(date string) error {

	_, err := r.db.Exec(`UPDATE letter SET opened_count = opened_count + 1, last_opened = NOW() WHERE show_date = $1`, date)

	return err
}

func (r *letterRepository) Update(letter domain.Letter) (*domain.Letter, error) {
	err := r.db.QueryRow(`
		UPDATE letter 
		SET content = $1, show_date = $2
		WHERE id = $3 
		RETURNING id, content, opened_count, COALESCE(last_opened::text, ''), show_date
	`, letter.Content, letter.Date, letter.ID).
		Scan(&letter.ID, &letter.Content, &letter.OpenedCount, &letter.LastOpened, &letter.Date)

	if err != nil {
		return nil, err
	}
	return &letter, nil
}

func (r *letterRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM letter WHERE id = $1", id)
	return err
}
