package node

import (
	"context"
	"flowctrl/uuid"
	"fmt"
)

type NodeInfo struct {
	Name        string
	Description string
	Version     string
}

type Node interface {
	Process(ctx context.Context)
	Cleanup()
	Info() NodeInfo
}

type Runner struct {
	uuid       uuid.Value
	name       string
	node       Node
	inputs     []Port
	outputs    []Port
	connectors []*Connector
	ctx        context.Context
}

func NewRunner(ctx context.Context, node Node) *Runner {
	inputs := Ports(node, DirectionInput)
	outputs := Ports(node, DirectionOutput)
	connectors := Connectors(inputs)
	info := node.Info()
	id := uuid.New()
	name := fmt.Sprintf("%s_%s_%d", info.Name, info.Version, id)
	return &Runner{id, name, node, inputs, outputs, connectors, ctx}
}

func (runner *Runner) Name() string {
	return runner.name
}

func (runner *Runner) UUID() uuid.Value {
	return runner.uuid
}

func (runner *Runner) Process() error {
	// trigger all connectors to input ports
	err := runner.processConnectors()
	if err != nil {
		return err
	}

	// run processing node
	runner.node.Process(runner.ctx)
	return nil
}

func (runner *Runner) Reset() {
	runner.resetConnectors()
}

func (runner *Runner) Outputs() []Port {
	return runner.outputs
}

func (runner *Runner) Inputs() []Port {
	return runner.inputs
}

func (runner *Runner) Connectors() []*Connector {
	return runner.connectors
}

func (runner *Runner) processConnectors() error {
	connectorsCount := len(runner.connectors)
	if connectorsCount == 0 {
		return nil
	}

	for i := 0; i < connectorsCount; i++ {
		conn := runner.connectors[i]
		err := conn.Trigger(runner.ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (runner *Runner) resetConnectors() {
	connectorsCount := len(runner.connectors)
	if connectorsCount == 0 {
		return
	}
	for i := 0; i < connectorsCount; i++ {
		runner.connectors[i].Reset()
	}
}