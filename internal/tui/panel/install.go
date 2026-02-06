package panel

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/k1empka/dotfiles/internal/installer"
)

var (
	inHeader    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#c4a7e7"))
	inInstalled = lipgloss.NewStyle().Foreground(lipgloss.Color("#9ccfd8"))
	inMissing   = lipgloss.NewStyle().Foreground(lipgloss.Color("#f6c177"))
	inError     = lipgloss.NewStyle().Foreground(lipgloss.Color("#eb6f92"))
	inActive    = lipgloss.NewStyle().Foreground(lipgloss.Color("#e0def4")).Bold(true)
	inNormal    = lipgloss.NewStyle().Foreground(lipgloss.Color("#908caa"))
	inCursor    = lipgloss.NewStyle().Foreground(lipgloss.Color("#9ccfd8"))
	inMuted     = lipgloss.NewStyle().Foreground(lipgloss.Color("#6e6a86"))
	inVersion   = lipgloss.NewStyle().Foreground(lipgloss.Color("#6e6a86"))
	inConfirm   = lipgloss.NewStyle().Foreground(lipgloss.Color("#9ccfd8"))
	inDeny      = lipgloss.NewStyle().Foreground(lipgloss.Color("#eb6f92"))
)

// Install displays application install status and allows installation.
type Install struct {
	width, height int
	statuses      []installer.AppStatus
	selected      int
	spinner       spinner.Model
	loading       bool
	installing    int // index of app currently installing, -1 if none
	confirming    bool
	loaded        bool
}

// NewInstall creates a new Install panel.
func NewInstall() *Install {
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#c4a7e7"))
	return &Install{
		spinner:    sp,
		installing: -1,
	}
}

func (p *Install) Title() string { return "Install" }

func (p *Install) ShortHelp() string {
	if p.confirming {
		return "y: confirm install  n: cancel"
	}
	if p.installing >= 0 {
		return "installing..."
	}
	return "j/k: navigate  enter: install  r: refresh  q: quit  ?: help"
}

func (p *Install) Init() tea.Cmd { return nil }

func (p *Install) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case SetSizeMsg:
		p.width = msg.Width
		p.height = msg.Height

	case InstallStatusMsg:
		p.statuses = msg.Statuses
		p.loading = false
		p.loaded = true
		p.installing = -1

	case InstallResultMsg:
		for i, s := range p.statuses {
			if s.App.BinName == msg.App.BinName {
				if msg.Err != nil {
					p.statuses[i].Status = installer.StatusError
					p.statuses[i].Err = msg.Err
				} else {
					p.statuses[i].Status = installer.StatusInstalled
					p.statuses[i].Version = "(just installed)"
				}
				break
			}
		}
		p.installing = -1
		return p, func() tea.Msg { return CheckInstallMsg{} }

	case spinner.TickMsg:
		if p.loading || p.installing >= 0 {
			var cmd tea.Cmd
			p.spinner, cmd = p.spinner.Update(msg)
			return p, cmd
		}

	case tea.KeyMsg:
		if p.confirming {
			return p.updateConfirming(msg)
		}
		return p.updateNormal(msg)
	}
	return p, nil
}

func (p *Install) updateConfirming(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y":
		p.confirming = false
		if p.selected >= 0 && p.selected < len(p.statuses) {
			s := p.statuses[p.selected]
			if s.Status == installer.StatusMissing || s.Status == installer.StatusError {
				p.installing = p.selected
				p.statuses[p.selected].Status = installer.StatusInstalling
				return p, tea.Batch(
					p.spinner.Tick,
					func() tea.Msg { return InstallRequestMsg{App: s.App} },
				)
			}
		}
	case "n", "N", "esc":
		p.confirming = false
	}
	return p, nil
}

func (p *Install) updateNormal(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if p.installing >= 0 {
		return p, nil
	}
	switch msg.String() {
	case "j", "down":
		if p.selected < len(p.statuses)-1 {
			p.selected++
		}
	case "k", "up":
		if p.selected > 0 {
			p.selected--
		}
	case "enter":
		if p.selected >= 0 && p.selected < len(p.statuses) {
			s := p.statuses[p.selected]
			if s.Status == installer.StatusMissing || s.Status == installer.StatusError {
				p.confirming = true
			}
		}
	case "r":
		p.loading = true
		return p, tea.Batch(
			p.spinner.Tick,
			func() tea.Msg { return CheckInstallMsg{} },
		)
	}
	return p, nil
}

func (p *Install) View() string {
	var b strings.Builder

	b.WriteString(inHeader.Render("Application Status"))
	b.WriteString("\n\n")

	if !p.loaded {
		if p.loading {
			b.WriteString("  " + p.spinner.View() + " Checking applications...")
		} else {
			b.WriteString(inMuted.Render("  Loading..."))
		}
		return b.String()
	}

	for i, s := range p.statuses {
		cursor := "  "
		nameStyle := inNormal
		if i == p.selected {
			cursor = inCursor.Render("> ")
			nameStyle = inActive
		}

		var indicator string
		switch s.Status {
		case installer.StatusInstalled:
			indicator = inInstalled.Render("[ok]")
		case installer.StatusMissing:
			indicator = inMissing.Render("[--]")
		case installer.StatusInstalling:
			indicator = p.spinner.View()
		case installer.StatusError:
			indicator = inError.Render("[!!]")
		default:
			indicator = inMuted.Render("[??]")
		}

		line := fmt.Sprintf("%s%s %s", cursor, indicator, nameStyle.Render(s.App.Name))

		if s.Status == installer.StatusInstalled && s.Version != "" {
			line += "  " + inVersion.Render(s.Version)
		}
		if s.Status == installer.StatusError && s.Err != nil {
			line += "  " + inError.Render(s.Err.Error())
		}

		b.WriteString(line + "\n")
	}

	if p.confirming && p.selected >= 0 && p.selected < len(p.statuses) {
		app := p.statuses[p.selected].App
		b.WriteString("\n")
		fmt.Fprintf(&b, "  Install %s? ", inActive.Render(app.Name))
		b.WriteString(inMuted.Render("(may require elevated privileges)") + "\n")
		b.WriteString("  " + inConfirm.Render("[y]") + inMuted.Render(" yes  "))
		b.WriteString(inDeny.Render("[n]") + inMuted.Render(" no"))
	}

	return b.String()
}
