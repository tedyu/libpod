package main

import (
	"os"

	"github.com/containers/common/pkg/auth"
	"github.com/containers/image/v5/types"
	"github.com/containers/libpod/cmd/podman/registry"
	"github.com/containers/libpod/pkg/domain/entities"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	logoutOptions = auth.LogoutOptions{}
	logoutCommand = &cobra.Command{
		Use:   "logout [flags] REGISTRY",
		Short: "Logout of a container registry",
		Long:  "Remove the cached username and password for the registry.",
		RunE:  logout,
		Args:  cobra.MaximumNArgs(1),
		Example: `podman logout quay.io
  podman logout --authfile dir/auth.json quay.io
  podman logout --all`,
	}
)

func init() {
	// Note that the local and the remote client behave the same: both
	// store credentials locally while the remote client will pass them
	// over the wire to the endpoint.
	registry.Commands = append(registry.Commands, registry.CliCommand{
		Mode:    []entities.EngineMode{entities.ABIMode, entities.TunnelMode},
		Command: logoutCommand,
	})
	flags := logoutCommand.Flags()

	// Flags from the auth package.
	flags.AddFlagSet(auth.GetLogoutFlags(&logoutOptions))
	logoutOptions.Stdin = os.Stdin
	logoutOptions.Stdout = os.Stdout
}

// Implementation of podman-logout.
func logout(cmd *cobra.Command, args []string) error {
	sysCtx := types.SystemContext{AuthFilePath: logoutOptions.AuthFile}

	registry := ""
	if len(args) > 0 {
		if logoutOptions.All {
			return errors.New("--all takes no arguments")
		}
		registry = args[0]
	}

	return auth.Logout(&sysCtx, &logoutOptions, registry)
}
