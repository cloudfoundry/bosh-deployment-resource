package bosh

import (
	boshcmd "github.com/cloudfoundry/bosh-cli/cmd"
)

type Runner interface {
	Execute(commandOpts interface{}) error
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
	deps := c.cliCoordinator.StreamingBasicDeps()
	globalOpts := c.cliCoordinator.GlobalOpts()
	setDefaults(commandOpts)

	cmd := boshcmd.NewCmd(globalOpts, commandOpts, deps)

	return cmd.Execute()
}
