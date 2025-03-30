package cli

import (
	"context"
	"fmt"

	"github.com/NordSecurity/nordvpn-linux/daemon/pb"
	"github.com/NordSecurity/nordvpn-linux/internal"
	"github.com/NordSecurity/nordvpn-linux/nstrings"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

// SetNotifyUsageText is shown next to notify command by nordvpn set --help
const SetNotifyUsageText = "Enables or disables notifications"

func (c *cmd) SetNotify(ctx *cli.Context) error {
	if ctx.NArg() != 1 {
		return formatError(argsCountError(ctx))
	}

	flag, err := nstrings.BoolFromString(ctx.Args().First())
	if err != nil {
		return formatError(argsParseError(ctx))
	}

	daemonResp, err := c.client.SetNotify(context.Background(), &pb.SetNotifyRequest{
		Notify: flag,
	})
	if err != nil {
		return formatError(err)
	}

	printMessage := func() {}
	defer func() {
		printMessage()
	}()

	messageNothingToSet := func() {
		color.Yellow(fmt.Sprintf(SetNotifyNothingToSet, nstrings.GetBoolLabel(flag)))
	}
	messageSuccess := func() {
		color.Green(fmt.Sprintf(SetNotifySuccess, nstrings.GetBoolLabel(flag)))
	}

	switch daemonResp.Type {
	case internal.CodeConfigError:
		return formatError(ErrConfig)
	case internal.CodeNothingToDo:
		printMessage = messageNothingToSet
	case internal.CodeSuccess:
		printMessage = messageSuccess
	}

	return nil
}
