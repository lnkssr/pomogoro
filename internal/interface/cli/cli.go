package cli

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/tabwriter"
	"time"
	"unicode/utf8"

	"golang.org/x/term"
)

type CommandInfo struct {
	Name  string
	Usage string
	Flag  *flag.FlagSet
}

type TimerCLI struct {
	Phase        string
	Remain       time.Duration
	WorkTime     int
	BreakTime    int
	ConfirmPhase bool
}

func NewTimerCLI(workTime, breakTime int, confirmPhase bool) *TimerCLI {
	return &TimerCLI{
		Phase:        "Work",
		Remain:       time.Duration(workTime) * time.Minute,
		WorkTime:     workTime,
		BreakTime:    breakTime,
		ConfirmPhase: confirmPhase,
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
	t.playNotificationSound()

	if t.ConfirmPhase {
		prompt := fmt.Sprintf("--- %s phase ended. Press any key to continue ---", t.Phase)
		t.centerPromptAndWait(prompt)
	} else {
		prompt := fmt.Sprintf("--- %s phase ended ---", t.Phase)
		t.centerPrompt(prompt)
		time.Sleep(700 * time.Millisecond)
	}
}

func (t *TimerCLI) Render() {
	var lines []string

	var phaseColorStart, phaseColorEnd string
	if t.Phase == "Work" {
		phaseColorStart = "\033[38;5;207m"
	} else {
		phaseColorStart = "\033[38;5;39m"
	}
	phaseColorEnd = "\033[0m"

	minutes := int(t.Remain.Minutes())
	seconds := int(t.Remain.Seconds()) % 60
	timerLine := fmt.Sprintf("%02d:%02d", minutes, seconds)

	bar := t.renderProgressBar()

	lines = append(lines, fmt.Sprintf("%sPhase: %s%s", phaseColorStart, t.Phase, phaseColorEnd))
	lines = append(lines, fmt.Sprintf("[%s]", bar))
	lines = append(lines, timerLine)
	lines = append(lines, fmt.Sprintf("Work=%dm Break=%dm", t.WorkTime, t.BreakTime))

	t.centerOutput(lines)
}

func (t *TimerCLI) renderProgressBar() string {
	var total int = 1
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
	if filled < 0 {
		filled = 0
	}
	if filled > width {
		filled = width
	}

	filledStr := strings.Repeat("#", filled)
	emptyStr := strings.Repeat(" ", width-filled)

	color := "\033[38;5;154m"
	if progress > 0.9 {
		color = "\033[38;5;196m"
	} else if progress > 0.66 {
		color = "\033[38;5;208m"
	}

	return fmt.Sprintf("%s%s\033[0m%s", color, filledStr, emptyStr)
}

func (t *TimerCLI) playNotificationSound() {
	url := "https://actions.google.com/sounds/v1/alarms/alarm_clock.ogg"
	cmd := exec.Command("ffplay", "-nodisp", "-autoexit", "-loglevel", "quiet", url)
	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "unable to play sound via ffplay: %v\n", err)
	}
}

func (t *TimerCLI) centerOutput(lines []string) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || width <= 0 || height <= 0 {
		fmt.Print("\033[H\033[2J")
		for _, ln := range lines {
			fmt.Println(ln)
		}
		fmt.Printf("\033[999B")
		return
	}

	maxLineWidth := 0
	for _, ln := range lines {
		l := visibleWidth(ln)
		if l > maxLineWidth {
			maxLineWidth = l
		}
	}

	totalLines := len(lines)
	topPad := (height - totalLines) / 2
	if topPad < 0 {
		topPad = 0
	}

	fmt.Print("\033[H\033[2J")

	for i := 0; i < topPad; i++ {
		fmt.Println()
	}

	for _, ln := range lines {
		lineWidth := visibleWidth(ln)
		leftPad := (width - lineWidth) / 2
		if leftPad < 0 {
			leftPad = 0
		}

		fmt.Print(strings.Repeat(" ", leftPad))
		fmt.Println(ln)
	}

	fmt.Printf("\033[%d;1H", height)
}

func (t *TimerCLI) centerPrompt(msg string) {
	t.centerOutput([]string{msg})
}

func (t *TimerCLI) centerPromptAndWait(msg string) {
	t.centerPrompt(msg)
	t.readAnyKey()
}

func (t *TimerCLI) readAnyKey() {
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		var tmp string
		if _, err := fmt.Scanln(&tmp); err != nil {
			fmt.Println("Error Scanln: ", err)
		}
		return
	}
	defer func() { _ = term.Restore(fd, oldState) }()

	buf := make([]byte, 1)
	_, err = os.Stdin.Read(buf)
	if err != nil {
		return
	}
}

func visibleWidth(s string) int {
	return utf8.RuneCountInString(stripAnsi(s))
}

func stripAnsi(s string) string {
	var b strings.Builder
	inEsc := false
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if !inEsc {
			if ch == 0x1b { // ESC
				inEsc = true
				continue
			}
			b.WriteByte(ch)
		} else {
			if (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') {
				inEsc = false
			}
		}
	}
	return b.String()
}

func PrintHelp(cmds []CommandInfo) {
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
	if _, err := fmt.Fprintln(w, "Usage: pomogoro <command> [flags]"); err != nil {
		return
	}
	if _, err := fmt.Fprintln(w, "Available commands:"); err != nil {
		return
	}
	for _, c := range cmds {
		if _, err := fmt.Fprintf(w, "  %s\t%s\n", c.Name, c.Usage); err != nil {
			return
		}
		if c.Flag != nil {
			c.Flag.VisitAll(func(f *flag.Flag) {
				if _, err := fmt.Fprintf(w, "    -%s\t%s (default %v)\n", f.Name, f.Usage, f.DefValue); err != nil {
					return
				}
			})
		}
	}
	if _, err := fmt.Fprintln(w, "\nUse 'pomogoro help <command>' for details."); err != nil {
		return
	}
	if err := w.Flush(); err != nil {
		return
	}
}
