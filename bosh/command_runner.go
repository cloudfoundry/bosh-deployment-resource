package bosh

import (
	"io"

	boshcmd "github.com/cloudfoundry/bosh-cli/cmd"
)

//go:generate counterfeiter . Runner
type Runner interface {
	ExecuteWithWriter(commandOpts interface{}, writer io.Writer) error
	Execute(commandOpts interface{}) error
	ExecuteWithDefaultOverride(commandOpts interface{}, override func(interface{}) (interface{}, error), writer io.Writer) error
}

type CommandRunner struct {
	cliCoordinator CLICoordinator
}

func NewCommandRunner(cliCoordinator CLICoordinator) CommandRunner {
	return CommandRunner{
		cliCoordinator: cliCoordinator,
	}
}

func (c CommandRunner) ExecuteWithWriter(commandOpts interface{}, writer io.Writer) error {
	return c.ExecuteWithDefaultOverride(commandOpts, func(opts interface{}) (interface{}, error) { return opts, nil }, writer)
}

func (c CommandRunner) Execute(commandOpts interface{}) error {
	return c.ExecuteWithDefaultOverride(commandOpts, func(opts interface{}) (interface{}, error) { return opts, nil }, nil)
}

func (c CommandRunner) ExecuteWithDefaultOverride(commandOpts interface{}, override func(interface{}) (interface{}, error), writer io.Writer) error {
	deps := c.cliCoordinator.BasicDeps(writer)
	globalOpts := c.cliCoordinator.GlobalOpts()
	setDefaults(commandOpts)

	commandOpts, err := override(commandOpts)
	if err != nil {
		return err
	}

	cmd := boshcmd.NewCmd(globalOpts, commandOpts, deps)

	return cmd.Execute()
}
