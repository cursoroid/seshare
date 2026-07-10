package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/cursoroid/seshare/internal/contacts"
	"github.com/cursoroid/seshare/internal/session"
	"github.com/cursoroid/seshare/internal/transport"
)

// cmdTUI browses sessions across all projects and sends a chosen one to a
// paired contact. Selection happens in the TUI; the actual croc transfer runs
// after the UI exits so its progress output isn't fighting the screen.
func cmdTUI(args []string) error {
	entries, err := session.ListAll()
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		return fmt.Errorf("no sessions found under ~/.claude/projects")
	}
	names, err := contacts.ListNames()
	if err != nil {
		return err
	}
	if len(names) == 0 {
		return fmt.Errorf("no contacts yet — run: seshare pair <name>")
	}

	app := tview.NewApplication()
	pages := tview.NewPages()

	list := tview.NewList().ShowSecondaryText(false)
	list.SetBorder(true).SetTitle(" sessions — enter to send, q to quit ")
	preview := tview.NewTextView().SetDynamicColors(true).SetWrap(true)
	preview.SetBorder(true).SetTitle(" preview ")

	showPreview := func(i int) {
		if i < 0 || i >= len(entries) {
			return
		}
		m, err := session.Preview(entries[i].Path)
		if err != nil {
			preview.SetText("[red]" + err.Error())
			return
		}
		preview.SetText(formatPreview(entries[i], m))
	}

	for _, e := range entries {
		list.AddItem(fmt.Sprintf("%-9s  %s", relTime(e.ModTime), e.Project), "", 0, nil)
	}
	list.SetChangedFunc(func(i int, _, _ string, _ rune) { showPreview(i) })
	showPreview(0)

	var chosenPath, chosenContact string
	list.SetSelectedFunc(func(i int, _, _ string, _ rune) {
		cl := tview.NewList().ShowSecondaryText(false)
		cl.SetBorder(true).SetTitle(" send to — esc to cancel ")
		for _, n := range names {
			n := n
			cl.AddItem(n, "", 0, func() {
				chosenPath, chosenContact = entries[i].Path, n
				app.Stop()
			})
		}
		cl.SetDoneFunc(func() { pages.SwitchToPage("sessions") })
		pages.AddAndSwitchToPage("contacts", cl, true)
	})

	flex := tview.NewFlex().
		AddItem(list, 0, 2, true).
		AddItem(preview, 0, 3, false)
	pages.AddPage("sessions", flex, true, true)

	app.SetInputCapture(func(ev *tcell.EventKey) *tcell.EventKey {
		if name, _ := pages.GetFrontPage(); name == "sessions" && ev.Rune() == 'q' {
			app.Stop()
			return nil
		}
		return ev
	})

	if err := app.SetRoot(pages, true).Run(); err != nil {
		return err
	}

	if chosenPath == "" || chosenContact == "" {
		return nil // user quit without choosing
	}
	return sendFile(chosenPath, chosenContact)
}

// sendFile gzips a transcript and sends it to a paired contact over croc.
func sendFile(path, contactName string) error {
	code, err := contacts.Get(contactName)
	if err != nil {
		return err
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	gz, err := gzipToTemp(raw, "seshare-"+strings.TrimSuffix(filepath.Base(path), ".jsonl")+".jsonl.gz")
	if err != nil {
		return err
	}
	defer os.Remove(gz)
	fmt.Printf("sending %s to %q…\n", filepath.Base(path), contactName)
	return transport.Send(gz, code)
}

func formatPreview(e session.Entry, m session.Meta) string {
	title := m.Title
	if title == "" {
		title = "(untitled)"
	}
	b := &strings.Builder{}
	fmt.Fprintf(b, "[yellow]%s[-]\n\n", title)
	fmt.Fprintf(b, "[gray]dir:[-]      %s\n", orDash(m.Cwd))
	fmt.Fprintf(b, "[gray]when:[-]     %s\n", e.ModTime.Format("2006-01-02 15:04"))
	fmt.Fprintf(b, "[gray]messages:[-] %d\n\n", m.MsgCount)
	if m.FirstPrompt != "" {
		fmt.Fprintf(b, "[gray]first:[-] %s\n\n", truncate(m.FirstPrompt, 300))
	}
	if m.LastPrompt != "" {
		fmt.Fprintf(b, "[gray]last:[-]  %s\n", truncate(m.LastPrompt, 300))
	}
	return b.String()
}

func orDash(s string) string {
	if s == "" {
		return "—"
	}
	return s
}

func truncate(s string, n int) string {
	s = strings.ReplaceAll(strings.TrimSpace(s), "\n", " ")
	if len(s) > n {
		return s[:n] + "…"
	}
	return s
}

func relTime(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	default:
		return fmt.Sprintf("%dd ago", int(d.Hours())/24)
	}
}
