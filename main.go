package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/joho/godotenv"
	"github.com/rivo/tview"
)

type appState struct {
	app             *tview.Application
	list            *tview.List
	input           *tview.InputField
	status          *tview.TextView
	entries         []string
	filtered        []string
	cursor          int
	fileMode        bool
	searchMode      bool
	statusClearTime time.Time
}

func filterEntries(entries []string, query string) []string {
	if query == "" {
		return entries
	}
	var out []string
	q := strings.ToLower(query)
	for _, e := range entries {
		if strings.Contains(strings.ToLower(e), q) {
			out = append(out, e)
		}
	}
	return out
}

func newApp(entries []string, fileMode bool) *appState {

	a := &appState{
		app:      tview.NewApplication(),
		list:     tview.NewList().ShowSecondaryText(false),
		input:    tview.NewInputField().SetLabel("Search: "),
		status:   tview.NewTextView().SetDynamicColors(true),
		entries:  entries,
		filtered: entries,
		fileMode: fileMode,
	}

	// configure input field
	a.input.
		SetFieldWidth(0).
		SetChangedFunc(func(text string) {
			a.filtered = filterEntries(a.entries, text)
			a.refreshList()
		}).
		SetDoneFunc(func(key tcell.Key) {
			switch key {
			case tcell.KeyEnter:
				a.exitSearch()
			case tcell.KeyEsc:
				a.input.SetText("")
				a.filtered = a.entries
				a.exitSearch()
			}
		})

	a.refreshList()

	return a
}

func (a *appState) refreshList() {
	a.list.Clear()
	for _, e := range a.filtered {
		a.list.AddItem(e, "", 0, nil)
	}
	if a.cursor >= len(a.filtered) {
		a.cursor = len(a.filtered) - 1
	}
	if a.cursor < 0 {
		a.cursor = 0
	}
	a.list.SetCurrentItem(a.cursor)
}

func (a *appState) setStatus(msg string) {
	a.status.SetText(fmt.Sprintf("[yellow]%s", msg))
	a.statusClearTime = time.Now().Add(2 * time.Second)
}

func (a *appState) exitSearch() {
	a.searchMode = false
	a.app.SetRoot(a.layout(), true).SetFocus(a.list)
}

func (a *appState) layout() tview.Primitive {
	var top tview.Primitive = tview.NewBox()
	if a.searchMode {
		top = a.input
	}
	inst := "[y]ank  [s]earch  [q]uit"
	footer := tview.NewFlex().
		AddItem(a.status, 0, 1, false).
		AddItem(tview.NewTextView().SetText(inst), len(inst)+1, 0, false)

	return tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(top, 1, 0, a.searchMode).
		AddItem(a.list, 0, 1, !a.searchMode).
		AddItem(footer, 1, 0, false)
}

func (a *appState) run() error {
	root := a.layout()
	a.app.SetRoot(root, true)

	a.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// auto-clear status if expired
		if !a.statusClearTime.IsZero() && time.Now().After(a.statusClearTime) {
			a.status.SetText("")
		}

		// In search mode, let input field handle keys
		if a.searchMode {
			return event
		}

		r := event.Rune()
		k := event.Key()

		// Commands
		switch {
		case r == 'q' || k == tcell.KeyCtrlC:
			a.app.Stop()
		case r == 's':
			a.searchMode = true
			a.input.SetText("")
			a.filtered = a.entries
			a.refreshList()
			a.app.SetRoot(a.layout(), true).SetFocus(a.input)
		case r == 'y':
			if len(a.filtered) > 0 {
				entry := a.filtered[a.cursor]
				if err := clipboard.WriteAll(entry); err != nil {
					a.setStatus("Error copying to clipboard")
				} else {
					a.setStatus("Copied to clipboard")
				}
			}
		}

		// Navigation (arrow keys + j/k)
		if k == tcell.KeyUp || r == 'k' {
			if a.cursor > 0 {
				a.cursor--
				a.list.SetCurrentItem(a.cursor)
			}
		} else if k == tcell.KeyDown || r == 'j' {
			if a.cursor < len(a.filtered)-1 {
				a.cursor++
				a.list.SetCurrentItem(a.cursor)
			}
		}

		return event
	})

	return a.app.Run()
}

func main() {
	fileMode := false
	var entries []string
	if len(os.Args) > 1 {
		fileMode = true
		envMap, err := godotenv.Read(os.Args[1])
		if err != nil {
			log.Fatalf("Error reading .env file: %v", err)
		}
		for k, v := range envMap {
			entries = append(entries, fmt.Sprintf("%s=%s", k, v))
		}
	} else {
		entries = os.Environ()
	}

	app := newApp(entries, fileMode)
	if err := app.run(); err != nil {
		log.Fatalf("Error running application: %v", err)
	}
}
