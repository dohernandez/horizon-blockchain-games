package internal

import (
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"io"
	"sync"

	"golang.org/x/sync/errgroup"

	"github.com/dohernandez/horizon-blockchain-games/internal/entities"
)

// chanCap is the capacity of the channels used as default.
const chanCap = 100

// Step is the type to represent the step of the pipeline.
type Step string

// String returns the string representation of the step.
func (s Step) String() string {
	return string(s)
}

const (
	// extractionStep is the extraction step.
	extractionStep Step = "extraction"
	// calculationStep is the calculation step.
	calculationStep Step = "calculation"
)

//go:generate mockery --name=StepProvider --outpkg=mocks --output=mocks --filename=step_provider.go --with-expecter

// StepProvider is the interface that provides the ability to load and save the step data.
type StepProvider interface {
	LoadStep(ctx context.Context, step string) ([]byte, error)
	SaveStep(ctx context.Context, step string, data []byte) error
}

//go:generate mockery --name=PipelineBackend --outpkg=mocks --output=mocks --filename=pipeline_backend.go --with-expecter

// PipelineBackend is the interface that provides the providers dependencies for the pipeline.
type PipelineBackend interface {
	ExtractProvider() ExtractProvider
	Conversor() Conversor
	LoadProvider() WarehouseProvider

	StepProvider() StepProvider
}

// PipelineConfig holds the configuration for the pipeline.
type PipelineConfig struct {
	// Workers is the number of workers to step's pipeline.
	// If it is 0, it will be set to 1.
	// It is used for the calculation steps.
	Workers int

	// ExtractStepEnabled enable the extraction step.
	ExtractStepEnabled bool
	// CalculateStepEnabled enable the calculation step.
	CalculateStepEnabled bool
	// InsertStepEnabled enable the insertion step.
	InsertStepEnabled bool
}

// Pipeline is the struct that holds the pipeline configuration and the backend dependencies.
type Pipeline struct {
	b PipelineBackend

	cfg PipelineConfig
}

// NewPipeline creates a new pipeline with the given backend dependencies and configuration.
func NewPipeline(b PipelineBackend, cfg PipelineConfig) *Pipeline {
	if cfg.Workers == 0 {
		cfg.Workers = 1
	}

	return &Pipeline{b: b, cfg: cfg}
}

// Run runs the pipeline.
//
// It runs the extraction, calculation and insertion steps in parallel using the number of workers defined in the configuration.
// Each step is optional and can be enabled or disabled.
//
// The extraction step is responsible for loading the data from the provider, normalize it and send the normalized transactions
// to the calculation step.
// The calculation step is responsible for calculating the total volume of the transactions in USD
// and send the flatten entities to the insertion step.
// The insertion step is responsible for saving the flatten entities into the target.
func (p *Pipeline) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	var (
		transactions chan entities.Transaction
		flattens     chan entities.Flatten
	)

	if p.cfg.ExtractStepEnabled {
		transactions = p.runExtraction(ctx, g)
	}

	if p.cfg.CalculateStepEnabled {
		flattens = p.runCalculation(ctx, g, transactions)
	}

	if p.cfg.InsertStepEnabled {
		p.runInsertion(ctx, g, flattens)
	}

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

// runExtraction runs the extraction step.
func (p *Pipeline) runExtraction(ctx context.Context, g *errgroup.Group) chan entities.Transaction {
	transactions := make(chan entities.Transaction, chanCap)

	g.Go(func() error {
		defer close(transactions)

		err := Extract(ctx, p.b.ExtractProvider(), transactions)
		if err != nil {
			return err
		}

		return nil
	})

	// Since CalculateStepEnabled is not enable, there is a need to save the step data.
	if !p.cfg.CalculateStepEnabled {
		p.saveExtractionStepData(ctx, g, transactions)
	}

	return transactions
}

// saveExtractionStepData saves the extraction step data.
//
// The data is saved when the calculation step is not enabled.
func (p *Pipeline) saveExtractionStepData(ctx context.Context, g *errgroup.Group, transactions <-chan entities.Transaction) {
	data := make(chan encoder, chanCap)

	p.saveStepData(ctx, g, extractionStep, data)

	g.Go(func() error {
		defer close(data)

		for t := range transactions {
			if ctx.Err() != nil {
				return ctx.Err()
			}

			data <- t
		}

		return nil
	})
}

// encoder is the interface that provides the ability to encode the data.
type encoder interface {
	Encode() []string
}

// saveStepData saves the step data.
func (p *Pipeline) saveStepData(ctx context.Context, g *errgroup.Group, step Step, data <-chan encoder) {
	g.Go(func() error {
		// Create a bytes.Buffer to store the transaction CSV data.
		var buf bytes.Buffer
		writer := csv.NewWriter(&buf)

	loop:
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case d, ok := <-data:
				if !ok {
					break loop
				}

				err := writer.Write(d.Encode())
				if err != nil {
					return err
				}
			}
		}

		writer.Flush()

		if err := writer.Error(); err != nil {
			return err
		}

		return p.b.StepProvider().SaveStep(ctx, step.String(), buf.Bytes())
	})
}

// runCalculation runs the calculation step.
//
// It runs the calculation step in parallel using the number of workers defined in the configuration.
// Along with the calculation step, it runs the flatten step in parallel to flatten the results into a map,
// and then send the flatten entities to the insertion step.
func (p *Pipeline) runCalculation(ctx context.Context, g *errgroup.Group, transactions chan entities.Transaction) chan entities.Flatten {
	// Since ExtractStepEnabled is not enable, there is a need to load the step data.
	if !p.cfg.ExtractStepEnabled {
		transactions = p.loadExtractionStepData(ctx, g)
	}

	var (
		fs       = make(chan entities.Flatten, chanCap)
		fsClosed bool

		fsSm sync.Mutex
	)

	for range p.cfg.Workers {
		g.Go(func() error {
			defer func() {
				fsSm.Lock()
				if !fsClosed {
					close(fs)

					fsClosed = true
				}
				fsSm.Unlock()
			}()

			err := Calculate(ctx, p.b.Conversor(), transactions, fs)
			if err != nil {
				return err
			}

			return nil
		})
	}

	var (
		mfs   = make(map[string]*entities.Flatten)
		mfsSm sync.Mutex

		flattens = make(chan entities.Flatten, chanCap)
	)

	// Flatten the results (flatten the map).
	g.Go(func() error {
		defer func() {
			close(flattens)
		}()

	loop:
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case f, ok := <-fs:
				if !ok {
					break loop
				}

				mfsSm.Lock()

				date := f.Date

				if _, ok := mfs[date]; !ok {
					mfs[date] = &entities.Flatten{
						Date:      date,
						ProjectID: f.ProjectID,
					}
				}

				mfs[date].NumTxs++

				mfs[date].TotalVolume += f.TotalVolume

				mfsSm.Unlock()
			}
		}

		for _, f := range mfs {
			flattens <- *f
		}

		return nil
	})

	// Since InsertStepEnabled is not enable, there is a need to save the step data.
	if !p.cfg.InsertStepEnabled {
		p.saveCalculationStepData(ctx, g, flattens)
	}

	return flattens
}

// loadExtractionStepData loads the extraction step data.
//
// It loads the extraction step data when the extraction step is not enabled.
func (p *Pipeline) loadExtractionStepData(ctx context.Context, g *errgroup.Group) chan entities.Transaction {
	data := p.loadDataStep(ctx, g, extractionStep, func(v []string) (any, error) {
		var t entities.Transaction

		err := t.Decode(v)
		if err != nil {
			return nil, err
		}

		return t, nil
	})

	transactions := make(chan entities.Transaction, chanCap)

	g.Go(func() error {
		defer close(transactions)

		for d := range data {
			if ctx.Err() != nil {
				return ctx.Err()
			}

			transactions <- d.(entities.Transaction)
		}

		return nil
	})

	return transactions
}

// decoderFunc is the function type to decode the data.
type decoderFunc func([]string) (any, error)

// loadDataStep loads the data step.
func (p *Pipeline) loadDataStep(ctx context.Context, g *errgroup.Group, step Step, decoder decoderFunc) chan any {
	data := make(chan any, chanCap)

	g.Go(func() error {
		defer close(data)

		dataLoaded, err := p.b.StepProvider().LoadStep(ctx, step.String())
		if err != nil {
			return err
		}

		reader := csv.NewReader(bytes.NewReader(dataLoaded))

		for {
			if ctx.Err() != nil {
				return ctx.Err()
			}

			record, err := reader.Read()
			if errors.Is(err, io.EOF) {
				break // End.
			}

			if err != nil {
				return err
			}

			d, err := decoder(record)
			if err != nil {
				return err
			}

			data <- d
		}

		return nil
	})

	return data
}

// saveCalculationStepData saves the calculation step data.
//
// The data is saved when the insertion step is not enabled.
func (p *Pipeline) saveCalculationStepData(ctx context.Context, g *errgroup.Group, flattens <-chan entities.Flatten) {
	data := make(chan encoder, chanCap)

	p.saveStepData(ctx, g, calculationStep, data)

	g.Go(func() error {
		defer close(data)

		for f := range flattens {
			data <- f
		}

		return nil
	})
}

// runInsertion runs the insertion step.
func (p *Pipeline) runInsertion(ctx context.Context, g *errgroup.Group, flattens chan entities.Flatten) {
	// Since ExtractStepEnabled is not enable, there is a need to load the step data.
	if !p.cfg.CalculateStepEnabled {
		flattens = p.loadCalculationStepData(ctx, g)
	}

	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case f, ok := <-flattens:
				if !ok {
					return nil
				}

				err := Insert(ctx, p.b.LoadProvider(), f)
				if err != nil {
					return err
				}
			}
		}
	})
}

// loadCalculationStepData loads the calculation step data.
//
// It loads the calculation step data when the calculation step is not enabled.
func (p *Pipeline) loadCalculationStepData(ctx context.Context, g *errgroup.Group) chan entities.Flatten {
	data := p.loadDataStep(ctx, g, calculationStep, func(v []string) (any, error) {
		var f entities.Flatten

		err := f.Decode(v)
		if err != nil {
			return nil, err
		}

		return f, nil
	})

	flattens := make(chan entities.Flatten, chanCap)

	g.Go(func() error {
		defer close(flattens)

		for d := range data {
			if ctx.Err() != nil {
				return ctx.Err()
			}

			flattens <- d.(entities.Flatten)
		}

		return nil
	})

	return flattens
}
