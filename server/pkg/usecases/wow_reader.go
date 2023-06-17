package usecases

import (
	"context"
	"fmt"

	"websmee/word-of-wisdom/server/pkg/domain"
)

type WOWRepository interface {
	GetRandomWOW(ctx context.Context) (*domain.WOW, error)
}

type WOWReader struct {
	wowRepository WOWRepository
}

func NewWOWReader(wowRepository WOWRepository) *WOWReader {
	return &WOWReader{wowRepository}
}

func (r *WOWReader) ReadRandomWOW(ctx context.Context) (*domain.WOW, error) {
	wow, err := r.wowRepository.GetRandomWOW(ctx)
	if err != nil {
		return nil, fmt.Errorf("read random wow failed, err: %w", err)
	}

	return wow, nil
}
