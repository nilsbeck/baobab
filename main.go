// Package main holds the implementation for the app.
package main

import (
	"context"
	"log"
	"strings"

	"github.com/nextmv-io/nextroute"
	"github.com/nextmv-io/nextroute/check"
	"github.com/nextmv-io/nextroute/factory"
	"github.com/nextmv-io/nextroute/schema"
	"github.com/nextmv-io/sdk/run"
	runSchema "github.com/nextmv-io/sdk/run/schema"
)

func main() {
	runner := run.CLI(solver)
	err := runner.Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

type options struct {
	Model  factory.Options                `json:"model,omitempty"`
	Solve  nextroute.ParallelSolveOptions `json:"solve,omitempty"`
	Format nextroute.FormatOptions        `json:"format,omitempty"`
	Check  check.Options                  `json:"check,omitempty"`
}

// createMap creates a map of stop IDs to client IDs from the input stops
func createMap(input schema.Input) map[string]string {
	stopToClient := make(map[string]string)
	for _, stop := range input.Stops {
		// Extract client ID from stop ID (format: "s{number}-client{number}")
		stopToClient[stop.ID] = stop.ID[strings.Index(stop.ID, "-")+1:]
	}
	return stopToClient
}

func solver(
	ctx context.Context,
	input schema.Input,
	options options,
) (runSchema.Output, error) {
	// Create the stop to client mapping
	stopToClient := createMap(input)

	model, err := factory.NewModel(input, options.Model)
	if err != nil {
		return runSchema.Output{}, err
	}

	constraint := &customConstraint{
		stopsToClients: stopToClient,
	}

	err = model.AddConstraint(constraint)
	if err != nil {
		return runSchema.Output{}, err
	}

	solver, err := nextroute.NewParallelSolver(model)
	if err != nil {
		return runSchema.Output{}, err
	}

	solutions, err := solver.Solve(ctx, options.Solve)
	if err != nil {
		return runSchema.Output{}, err
	}

	last, err := solutions.Last()
	if err != nil {
		return runSchema.Output{}, err
	}

	output, err := check.Format(
		ctx,
		options,
		options.Check,
		solver,
		last,
	)
	if err != nil {
		return runSchema.Output{}, err
	}
	output.Statistics.Result.Custom = factory.DefaultCustomResultStatistics(last)

	return output, nil
}
