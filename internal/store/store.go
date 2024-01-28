package store

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Store struct {
	config            *Config
	db                *sql.DB
	messageRepository *MessageRepository
}

// Инициалзация БД
func New(config *Config) *Store {
	return &Store{
		config: config,
	}
}

// Соединение с БД
func (s *Store) Open() error {
	Db, err := sql.Open("postgres", s.config.Database_URL)
	if err != nil {
		return err
	}
	s.db = Db
	return nil
}

// Закрытие БД
func (s *Store) Close() {
	s.db.Close()
}

// Метод для взаимодествия с БД из стороннего пакета
func (s *Store) Message() *MessageRepository {
	if s.messageRepository != nil {
		return s.messageRepository
	}
	s.messageRepository = &MessageRepository{
		store: s,
	}
	return s.messageRepository
}
