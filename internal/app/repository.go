package app

import (
	"context"
	"github.com/dimoktorr/monitoring/internal/domain"
)

type Repository interface {
	GetProduct(ctx context.Context, id int) (*domain.Product, error)
}
