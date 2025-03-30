package service

import (
	"context"
	"github.com/Ex6linz/BookSwap/backend/internal/book/domain"
)

// BookRepository interfejs definiujący metody dostępu do danych
type BookRepository interface {
	GetAll(ctx context.Context) ([]domain.Book, error)
	// Inne metody
}

type BookService struct {
	repo BookRepository
}

func NewBookService(repo BookRepository) *BookService {
	return &BookService{
		repo: repo,
	}
}

func (s *BookService) GetAllBooks(ctx context.Context) ([]domain.Book, error) {
	return s.repo.GetAll(ctx)
}
