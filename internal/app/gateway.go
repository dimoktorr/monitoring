package app

import "context"

type Gateway interface {
	Pay(ctx context.Context, amount float64) (string, error)
}
