package cli

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/tabwriter"
	"time"
)

type CommandInfo struct {
	Name  string
	Usage string
	Flag  *flag.FlagSet
}

type TimerCLI struct {
	Phase     string
	Remain    time.Duration
	WorkTime  int
	BreakTime int
}

func NewTimerCLI(workTime, breakTime int) *TimerCLI {
	return &TimerCLI{
		Phase:     "Work",
		Remain:    time.Duration(workTime) * time.Minute,
		WorkTime:  workTime,
		BreakTime: breakTime,
	}
}

func (t *TimerCLI) OnTick(remain time.Duration, phase string) {
	t.Remain = remain
	t.Phase = phase
	t.Render()
}

func (t *TimerCLI) OnPhaseEnd(phase string) {
	t.Phase = phase
	t.Render()
	fmt.Printf("\n--- %s phase ended! ---\n", t.Phase)
	t.playNotificationSound()
}

func (t *TimerCLI) Render() {
	fmt.Print("\033[H\033[2J")

	var phaseColor string
	if t.Phase == "Work" {
		phaseColor = "\033[38;5;207m"
	} else {
		phaseColor = "\033[38;5;39m"
	}

	minutes := int(t.Remain.Minutes())
	seconds := int(t.Remain.Seconds()) % 60

	bar := t.renderProgressBar()

	fmt.Printf("%sPhase: %s\033[0m\n", phaseColor, t.Phase)
	fmt.Printf("[%s] %02d:%02d\n", bar, minutes, seconds)
	fmt.Printf("Press Ctrl+C to stop. Work=%dm Break=%dm\n", t.WorkTime, t.BreakTime)
}

func (t *TimerCLI) renderProgressBar() string {
	total := 1
	if t.Phase == "Work" {
		total = t.WorkTime * 60
	} else {
		total = t.BreakTime * 60
	}

	if total <= 0 {
		total = 1
	}
	elapsed := total - int(t.Remain.Seconds())
	if elapsed < 0 {
		elapsed = 0
	}
	progress := float64(elapsed) / float64(total)
	if progress < 0 {
		progress = 0
	}
	if progress > 1 {
		progress = 1
	}

	width := 30
	filled := int(progress * float64(width))

	filledStr := strings.Repeat("#", filled)
	emptyStr := strings.Repeat(" ", width-filled)

	color := "\033[38;5;154m"
	if progress > 0.66 {
		color = "\033[38;5;208m"
	} else if progress > 0.9 {
		color = "\033[38;5;196m"
	}

	bar := fmt.Sprintf("%s%s\033[0m%s", color, filledStr, emptyStr)
	return bar
}

func (t *TimerCLI) playNotificationSound() {
	url := "https://actions.google.com/sounds/v1/alarms/alarm_clock.ogg"
	cmd := exec.Command("ffplay", "-nodisp", "-autoexit", "-loglevel", "quiet", url)
	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "unable to play sound via ffplay: %v\n", err)
	}
}

func PrintHelp(cmds []CommandInfo) {
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintln(w, "Usage: pomogoro <command> [flags]\n")
	fmt.Fprintln(w, "Available commands:")
	for _, c := range cmds {
		fmt.Fprintf(w, "  %s\t%s\n", c.Name, c.Usage)
		if c.Flag != nil {
			c.Flag.VisitAll(func(f *flag.Flag) {
				fmt.Fprintf(w, "    -%s\t%s (default %v)\n", f.Name, f.Usage, f.DefValue)
			})
		}
	}
	fmt.Fprintln(w, "\nUse 'pomogoro help <command>' for details.")
	w.Flush()
}
