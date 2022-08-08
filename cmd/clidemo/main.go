package main

import (
	"fmt"

	"github.com/alecthomas/kong"
)

var CLI struct {
	Ver     string `default:"${version}"`
	Logging struct {
		Level string `enum:"debug,info,warn,error" default:"info" env:"Y" help:"Paths to remove."`
		Type  string `enum:"json,console" default:"console" env:"X" help:"Paths to remove."`
	} `embed:""` //`embed:"" prefix:"logging." envprefix:"XXX_"`

	Logging1 struct {
		Level string `enum:"debug,info,warn,error" default:"info" env:"Y" help:"Paths to remove."`
		Type  string `enum:"json,console" default:"console" env:"X" help:"Paths to remove."`
		Type1 string `kong:"-"` // Ignore the field
	} `embed:"" prefix:"logging." envprefix:"XXX_"`
	Debug bool `help:"Enable debug mode." env:"X"`
	Rm    struct {
		Force     bool `help:"Force removal." env:"X"`
		Recursive bool `help:"Recursively remove files."`

		Paths []string `arg:"" name:"path" help:"Paths to remove." type:"path"`
	} `cmd:"" help:"Remove files."`

	Ls struct {
		Paths []string `arg:"" optional:"" name:"path" help:"Paths to list." type:"path"`
	} `cmd:"" help:"List paths."`

	Help     helpCmd   `cmd:"" help:"Show help."`
	Question helpCmd   `cmd:"" hidden:"" name:"?" help:"Show help."`
	Status   statusCmd `cmd:"" help:"Show server status."`
}

func main() {
	ctx := kong.Parse(
		&CLI,
		kong.Name("test"),
		//kong.NoDefaultHelp(),
		kong.Description("An app demonstrating HelpProviders"),
		kong.UsageOnError(),
		kong.Vars{
			"version": "0.0.1",
		},
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: true,
		}),
	)
	ctx.FatalIfErrorf(ctx.Run())
	//assert.Must(ctx.PrintUsage(true))
	//switch ctx.Command() {
	//case "rm <path>":
	//case "ls":
	//default:
	//}
}

type statusCmd struct {
	Verbose bool `short:"v" help:"Show verbose status information."`
}

func (s *statusCmd) Run(ctx *kong.Context) error {
	//ctx.Printf("OK")
	fmt.Println("OK")
	return nil
}

type helpCmd struct {
	Command []string `arg:"" optional:"" help:"Show help on command."`
}

// Run shows help.
func (h *helpCmd) Run(realCtx *kong.Context) error {
	ctx, err := kong.Trace(realCtx.Kong, h.Command)
	if err != nil {
		return err
	}
	if ctx.Error != nil {
		return ctx.Error
	}
	err = ctx.PrintUsage(false)
	if err != nil {
		return err
	}
	fmt.Fprintln(realCtx.Stdout)
	return nil
}
