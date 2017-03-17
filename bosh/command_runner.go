package bosh

import (
	boshcmd "github.com/cloudfoundry/bosh-cli/cmd"
)

type Runner interface {
	Execute(commandOpts interface{}) error
	ExecuteWithDefaultOverride(commandOpts interface{}, override func(interface{}) (interface{}, error)) error
}

type CommandRunner struct {
	cliCoordinator CLICoordinator
}

func NewCommandRunner(cliCoordinator CLICoordinator) CommandRunner {
	return CommandRunner{
		cliCoordinator: cliCoordinator,
	}
}

func (c CommandRunner) Execute(commandOpts interface{}) error {
	return c.ExecuteWithDefaultOverride(commandOpts, func(opts interface{}) (interface{}, error) { return opts, nil })
}

func (c CommandRunner) ExecuteWithDefaultOverride(commandOpts interface{}, override func(interface{}) (interface{}, error)) error {
	deps := c.cliCoordinator.StreamingBasicDeps()
	globalOpts := c.cliCoordinator.GlobalOpts()
	setDefaults(commandOpts)

	commandOpts, err := override(commandOpts)
	if err != nil {
		return err
	}

	cmd := boshcmd.NewCmd(globalOpts, commandOpts, deps)

	return cmd.Execute()
}
