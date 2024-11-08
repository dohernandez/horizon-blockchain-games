package warehouse

import (
	"context"
	"fmt"

	"github.com/dohernandez/horizon-blockchain-games/internal/entities"
)

// PrintType is the type of the Print warehouse.
const PrintType = "print"

// Print is a target to print the flatten entity.
type Print struct{}

// Save saves the flatten entity.
func (p *Print) Save(_ context.Context, f entities.Flatten) error {
	fmt.Printf("Save: %v\n", f)

	return nil
}
