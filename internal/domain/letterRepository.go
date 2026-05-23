package domain

type LetterRepository interface {
	Create(letter Letter) (Letter, error)
	GetAll() ([]Letter, error)
	GetByDate(date string) (*Letter, error)
	IncrementTapCount(date string) error
	Update(letter Letter) (*Letter, error)
	Delete(id int) error
}
