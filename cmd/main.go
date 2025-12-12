package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"pomogoro/internal/interface/cli"
	"pomogoro/internal/service"
)

type Command struct {
	Name  string
	Usage string
	Flag  *flag.FlagSet
	Run   func(cmd *Command, args []string)
}

func main() {
	cmdRun := &Command{
		Name:  "run",
		Usage: "run pomodoro timer",
		Flag:  flag.NewFlagSet("run", flag.ExitOnError),
	}

	cmdHelp := &Command{
		Name:  "help",
		Usage: "show help for commands",
		Flag:  flag.NewFlagSet("help", flag.ExitOnError),
	}

	cmdVersion := &Command{
		Name:  "version",
		Usage: "print version",
		Flag:  flag.NewFlagSet("version", flag.ExitOnError),
	}

	var workMin int
	var breakMin int
	var cycles int
	var confirmPhase bool

	cmdRun.Flag.IntVar(&workMin, "w", 25, "work duration in minutes")
	cmdRun.Flag.IntVar(&breakMin, "b", 5, "break duration in minutes")
	cmdRun.Flag.IntVar(&cycles, "c", 4, "number of cycles")
	cmdRun.Flag.BoolVar(&confirmPhase, "confirm-phase", false, "require pressing any key to continue after phase ends")

	commandInfos := []cli.CommandInfo{
		{Name: cmdRun.Name, Usage: cmdRun.Usage, Flag: cmdRun.Flag},
		{Name: cmdHelp.Name, Usage: cmdHelp.Usage, Flag: cmdHelp.Flag},
		{Name: cmdVersion.Name, Usage: cmdVersion.Usage, Flag: cmdVersion.Flag},
	}

	cmdRun.Run = func(cmd *Command, args []string) {
		if err := cmd.Flag.Parse(args); err != nil {
			fmt.Println("Error parsing flats:", err)
		}

		if workMin <= 0 || breakMin <= 0 || cycles <= 0 {
			fmt.Fprintln(os.Stderr, "work, break and cycles must be positive integers")
			os.Exit(2)
		}

		timerCLI := cli.NewTimerCLI(workMin, breakMin, confirmPhase)
		timerUC := service.NewTimerUseCase(
			time.Duration(workMin)*time.Minute,
			time.Duration(breakMin)*time.Minute,
			cycles,
		)

		timerUC.Start(timerCLI.OnTick, timerCLI.OnPhaseEnd)

		fmt.Println("All cycles completed. Good job!")
	}

	cmdHelp.Run = func(cmd *Command, args []string) {
		if len(args) > 0 {
			name := args[0]
			for _, ci := range commandInfos {
				if ci.Name == name {
					cli.PrintHelp([]cli.CommandInfo{ci})
					return
				}
			}
			fmt.Fprintf(os.Stderr, "Unknown command: %s\n", name)
			os.Exit(2)
		}
		cli.PrintHelp(commandInfos)
	}

	cmdVersion.Run = func(cmd *Command, args []string) {
		fmt.Println("pomogoro version 0.1.0")
	}

	if len(os.Args) < 2 {
		cli.PrintHelp(commandInfos)
		os.Exit(0)
	}

	name := os.Args[1]
	switch name {
	case cmdRun.Name:
		cmdRun.Run(cmdRun, os.Args[2:])
	case cmdHelp.Name:
		cmdHelp.Run(cmdHelp, os.Args[2:])
	case cmdVersion.Name:
		cmdVersion.Run(cmdVersion, os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %q\n", name)
		cli.PrintHelp(commandInfos)
		os.Exit(2)
	}
}
