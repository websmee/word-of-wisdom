package infrastructure

import (
	"context"
	"math/rand"

	"websmee/word-of-wisdom/server/pkg/domain"
)

type WOWRepository struct {
	quotes []string
}

func NewWOWRepository(quotes []string) *WOWRepository {
	return &WOWRepository{quotes}
}

func (r *WOWRepository) GetRandomWOW(_ context.Context) (*domain.WOW, error) {
	if len(r.quotes) == 0 {
		return &domain.WOW{}, nil
	}

	return &domain.WOW{Quote: r.quotes[rand.Intn(len(r.quotes))]}, nil
}
