package internal

import (
	"context"

	"github.com/dohernandez/horizon-blockchain-games/internal/entities"
)

//go:generate mockery --name=LoadProvider --outpkg=mocks --output=mocks --filename=load_provider.go --with-expecter

// LoadProvider is the interface that provides the ability to save the flatten entity.
type LoadProvider interface {
	// Save saves the flatten entity.
	Save(ctx context.Context, flatten entities.Flatten) error
}

// Insert inserts the flatten entity into the target.
func Insert(ctx context.Context, target LoadProvider, input entities.Flatten) error {
	if err := target.Save(ctx, input); err != nil {
		return err
	}

	return nil
}
