package main

import (
	"github.com/nextmv-io/nextroute"
)

// customConstraint is a struct that allows to implement a custom constraint.
type customConstraint struct {
	stopsToClients map[string]string
}

type solutionData struct {
	clientsVisited map[string]bool
}

func (d *solutionData) Copy() nextroute.Copier {
	clientsVisited := make(map[string]bool)
	for k, v := range d.clientsVisited {
		clientsVisited[k] = v
	}
	return &solutionData{
		clientsVisited: clientsVisited,
	}
}

// UpdateConstraintSolutionData is called when a stop is added to the solution
func (c *customConstraint) UpdateConstraintSolutionData(solution nextroute.Solution) (nextroute.Copier, error) {
	data := &solutionData{
		clientsVisited: make(map[string]bool),
	}

	for _, vehicle := range solution.Vehicles() {
		if vehicle.IsEmpty() {
			continue
		}

		stop := vehicle.Last().ConstraintData(c).(*solutionData)
		for k, v := range stop.clientsVisited {
			data.clientsVisited[k] = v
		}
	}

	return data, nil
}

func (c *customConstraint) UpdateConstraintStopData(solutionStop nextroute.SolutionStop) (nextroute.Copier, error) {
	if solutionStop.IsFirst() {
		return &solutionData{
			clientsVisited: make(map[string]bool),
		}, nil
	}

	currentData, ok := solutionStop.Solution().ConstraintData(c).(*solutionData)
	if !ok {
		currentData = &solutionData{
			clientsVisited: make(map[string]bool),
		}
	}

	if solutionStop.IsLast() {
		return currentData, nil
	}

	currentData.clientsVisited[c.stopsToClients[solutionStop.ModelStop().ID()]] = true
	return currentData, nil
}

// EstimateIsViolated returns true if the constraint is violated. If the
// constraint is violated, the solver needs a hint to determine if further
// moves should be generated for the vehicle.
func (c *customConstraint) EstimateIsViolated(move nextroute.Move) (isViolated bool, stopPositionsHint nextroute.StopPositionsHint) {
	constraintData := move.Solution().ConstraintData(c).(*solutionData)
	for _, stopPosition := range move.StopPositions() {
		stopId := stopPosition.Stop().ModelStop().ID()
		stopClientId := c.stopsToClients[stopId]
		if constraintData.clientsVisited[stopClientId] {
			return true, nextroute.NoPositionsHint()
		}
	}
	// If the constraint is not violated, the solver does not need a hint.
	return false, nextroute.NoPositionsHint()
}

// String returns the name of the constraint.
func (c *customConstraint) String() string {
	return "clients_visited_constraint"
}

// IsTemporal returns true if the constraint should be checked after all initial
// stops have been planned. It returns false if the constraint should be checked
// after each of the initial stops has been planned.
func (c *customConstraint) IsTemporal() bool {
	return false
}
