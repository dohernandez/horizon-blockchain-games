package internal

import (
	"context"

	"github.com/dohernandez/horizon-blockchain-games/internal/entities"
)

//go:generate mockery --name=WarehouseProvider --outpkg=mocks --output=mocks --filename=warehouse_provider.go --with-expecter

// WarehouseProvider is the interface that provides the ability to save the flatten entity.
type WarehouseProvider interface {
	// Save saves the flatten entity.
	Save(ctx context.Context, flatten entities.Flatten) error
}

// Insert inserts the flatten entity into the target.
func Insert(ctx context.Context, target WarehouseProvider, input entities.Flatten) error {
	if err := target.Save(ctx, input); err != nil {
		return err
	}

	return nil
}
