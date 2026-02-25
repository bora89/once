package command

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/basecamp/once/internal/docker"
)

type StopCommand struct {
	cmd *cobra.Command
}

func NewStopCommand(root *RootCommand) *StopCommand {
	s := &StopCommand{}
	s.cmd = &cobra.Command{
		Use:   "stop <app>",
		Short: "Stop an application",
		Args:  cobra.ExactArgs(1),
		RunE:  WithNamespace(s.run),
	}
	return s
}

func (s *StopCommand) Command() *cobra.Command {
	return s.cmd
}

// Private

func (s *StopCommand) run(ns *docker.Namespace, cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	appName := args[0]

	app := ns.Application(appName)
	if app == nil {
		return fmt.Errorf("application %q not found", appName)
	}

	if err := app.Stop(ctx); err != nil {
		return fmt.Errorf("stopping application: %w", err)
	}

	fmt.Printf("Stopped %s\n", appName)
	return nil
}
