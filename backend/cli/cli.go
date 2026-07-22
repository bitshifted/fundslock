// Copyright 2026 Bitshift ED
// SPDX-License-Identifier: MPL-2.0

package cli

import (
	"bitshifted/fundslock-be/srv"
	"fmt"

	"github.com/alecthomas/kong"
)

type CLI struct {
	Version     VersionCmd `cmd:"" name:"version" help:"Display version information"`
	Start       StartCmd   `cmd:"" name:"start" help:"Start the server"`
	EnableDebug bool       `help:"Enable debug logging"`
}

type VersionCmd struct {
}

type StartCmd struct {
}

func (vc *VersionCmd) Run(ctx *kong.Context) error {
	fmt.Printf("Version: %s\nBuild number: %s\nCommit ID: %s\n",
		ProgramVersion.Version, ProgramVersion.BuildNumber, ProgramVersion.CommitID)
	return nil
}

func (sc *StartCmd) Run(ctx *kong.Context) error {
	return srv.Start()
}
