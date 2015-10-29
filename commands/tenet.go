package commands

import (
	"fmt"

	"strings"

	"github.com/codegangsta/cli"
	"github.com/lingo-reviews/lingo/commands/review"
	"github.com/lingo-reviews/lingo/tenet"
)

// TenetCMD is a fallthrough CMD which treats command as the tenet name and
// passes through any arguments to the tenet.
func TenetCMD(ctx *cli.Context, command string) {
	var commandIsTenet bool
	var c tenet.Config
	for _, c = range listTenets(ctx) {
		if command == c.Name {
			commandIsTenet = true
			break
		}
	}
	if !commandIsTenet {
		fmt.Println("command not found")
		return
	}

	// TODO(matt) read about bash completion on
	// https://github.com/codegangsta/cli. Is there a nice way that we could
	// add bash completion for tenet names (as they'll be long and clumsy).
	t, err := tenet.New(ctx, c)
	if err != nil {
		oserrf("command or tenet not found: %s", err.Error())
		return
	}

	// Initialise the tenet driver
	err = t.InitDriver()
	if err != nil {
		oserrf(err.Error())
		return
	}

	method := "Help"
	args := ctx.Args()
	if len(args[1:]) > 0 {
		method = methodFromArg(args[1])
	}

	switch method {
	case "Review":
		reviewResult, err := t.Review(args[2:]...)
		if err != nil {
			oserrf("error running method %q, %s", method, err.Error())
			return
		}
		for _, i := range reviewResult.Issues {
			fmt.Println(review.FormatPlainText(i))
		}
		for _, e := range reviewResult.Errs {
			fmt.Println(e)
		}
	case "Help":
		text, err := t.Help()
		if err != nil {
			oserrf("error running method %q, %s", method, err.Error())
			return
		}
		fmt.Println(text)
	case "Version":
		text, err := t.Version()
		if err != nil {
			oserrf("error running method %q, %s", method, err.Error())
			return
		}
		fmt.Println(text)
	case "Debug":
		fmt.Println(t.Debug(args...))
	default:
		oserrf("tenet does not have method %q", method)
		return
	}
}

func methodFromArg(arg string) (method string) {
	method = strings.Title(arg)
	return strings.Replace(method, "-", "", -1)
}
