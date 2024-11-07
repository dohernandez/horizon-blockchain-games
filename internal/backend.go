package internal

import (
	"github.com/dohernandez/horizon-blockchain-games/internal/conversor"
	"github.com/dohernandez/horizon-blockchain-games/internal/source"
	"github.com/dohernandez/horizon-blockchain-games/internal/storage"
	"github.com/dohernandez/horizon-blockchain-games/internal/target"
)

// Config holds the configuration for the backend to create the providers dependencies.
type Config struct {
	Dir    string
	File   string
	IsTest bool
}

// Backend is the main struct that holds the providers dependencies.
type Backend struct {
	cfg Config

	extractProvider ExtractProvider
	conversor       Conversor
	loadProvider    LoadProvider
	stepProvider    StepProvider
}

// NewBackend creates a new backend with the given configuration.
func NewBackend(cfg Config) *Backend {
	b := Backend{
		cfg: cfg,
	}

	if cfg.IsTest {
		b.prepareBackendsForTest()

		return &b
	}

	return &b
}

// prepareBackendsForTest prepares the backend for testing purposes.
//
// It sets the providers dependencies mainly to use files, hardcoded values, and print the output.
func (b *Backend) prepareBackendsForTest() {
	b.extractProvider = source.NewFile(b.cfg.Dir, b.cfg.File)

	b.stepProvider = storage.NewFile(b.cfg.Dir)

	b.conversor = conversor.NewHardcoded()

	b.loadProvider = &target.Print{}
}

// ExtractProvider returns the extract provider.
func (b *Backend) ExtractProvider() ExtractProvider {
	return b.extractProvider
}

// Conversor returns the conversor.
func (b *Backend) Conversor() Conversor {
	return b.conversor
}

// LoadProvider returns the load provider.
func (b *Backend) LoadProvider() LoadProvider {
	return b.loadProvider
}

// StepProvider returns the step provider.
func (b *Backend) StepProvider() StepProvider {
	return b.stepProvider
}
