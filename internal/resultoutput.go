package internal

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"willbehn/ht/models"
)

func ResultOutputShort(results []models.CmdEvent) {
	fmt.Println("\033[2m" + strings.Repeat("-", 72) + "\033[0m")

	for _, ev := range results {
		t := time.Unix(ev.TS, 0).Local()
		since := TimeSince(t)

		fmt.Printf("\033[2m %-4d \033[0m \033[32m%-8s\033[0m  \033[2m%-20s\033[0m  \033[1m%s\033[0m\n",
			ev.Id,
			since,
			prettyDir(ev.Dir, 24),
			ev.Cmd,
		)
	}
}

func ResultOutputLong(results []models.CmdEvent) {
	fmt.Println("\033[2m" + strings.Repeat("-", 80) + "\033[0m")

	for _, ev := range results {
		t := time.Unix(ev.TS, 0).Local()
		absT := t.Format("2006-01-02 15:04:05")

		fmt.Printf("\033[2m %-4d \033[0m \033[32m%-19s\033[0m  \033[2m%-7s\033[0m  \033[2m%-20s\033[0m  \033[1m%s\033[0m\n",
			ev.Id,
			absT,
			ev.Shell,
			prettyDir(ev.Dir, 24),
			ev.Cmd,
		)
	}
}

func prettyDir(dir string, width int) string {
	if dir == "" {
		return fmt.Sprintf("%-*s", width, "-")
	}

	var out string

	clean := filepath.Clean(dir)
	split := strings.Split(clean, string(filepath.Separator))

	if len(split) > 3 {
		keep := split[len(split)-2:]
		out = "…" + string(filepath.Separator) + filepath.Join(keep...)
	} else {
		out = dir
	}

	if len(out) > width {
		out = out[:width]
	}

	return fmt.Sprintf("%-*s", width, out)
}
