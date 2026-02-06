package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/k1empka/dotfiles/internal/chezmoi"
	"github.com/k1empka/dotfiles/internal/installer"
	"github.com/k1empka/dotfiles/internal/tui/component"
	"github.com/k1empka/dotfiles/internal/tui/panel"
)

// App is the root model that orchestrates tabs and panels.
type App struct {
	client    *chezmoi.Client
	installer *installer.Installer
	tabs      tabBar
	panels    []panel.Panel
	help      component.Help
	width     int
	height    int
	err       error
}

// NewApp creates a new App with all panels.
func NewApp(client *chezmoi.Client, inst *installer.Installer) App {
	panels := []panel.Panel{
		panel.NewOverview(),
		panel.NewShell(),
		panel.NewGit(),
		panel.NewThemes(),
		panel.NewNeovim(),
		panel.NewAlacritty(),
		panel.NewChezmoi(),
		panel.NewInstall(),
	}
	titles := make([]string, len(panels))
	for i, p := range panels {
		titles[i] = p.Title()
	}
	return App{
		client:    client,
		installer: inst,
		tabs:      newTabBar(titles),
		panels:    panels,
		help:      component.NewHelp(),
	}
}

func (a App) Init() tea.Cmd {
	return tea.Batch(
		a.loadConfig(),
		a.loadShellConfigs(),
		a.loadGitConfig(),
		a.loadAlacrittyConfig(),
		a.loadNvimFiles(),
		a.loadThemeColors("rose-pine"),
		a.loadThemeColors("catppuccin"),
		a.loadThemeColors("tokyo-night"),
		a.checkInstallStatus(),
	)
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Help overlay intercepts keys when visible.
	if a.help.Visible() {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			var cmd tea.Cmd
			a.help, cmd = a.help.Update(msg)
			return a, cmd
		case tea.WindowSizeMsg:
			a.width = msg.Width
			a.height = msg.Height
			a.help.SetSize(msg.Width, msg.Height)
			return a, nil
		}
		return a, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, globalKeys.Quit):
			return a, tea.Quit
		case key.Matches(msg, globalKeys.Help):
			a.help.Toggle()
			return a, nil
		case key.Matches(msg, globalKeys.NextTab):
			a.tabs.next()
			return a, nil
		case key.Matches(msg, globalKeys.PrevTab):
			a.tabs.prev()
			return a, nil
		case key.Matches(msg, globalKeys.Tab1):
			a.tabs.setActive(0)
			return a, nil
		case key.Matches(msg, globalKeys.Tab2):
			a.tabs.setActive(1)
			return a, nil
		case key.Matches(msg, globalKeys.Tab3):
			a.tabs.setActive(2)
			return a, nil
		case key.Matches(msg, globalKeys.Tab4):
			a.tabs.setActive(3)
			return a, nil
		case key.Matches(msg, globalKeys.Tab5):
			a.tabs.setActive(4)
			return a, nil
		case key.Matches(msg, globalKeys.Tab6):
			a.tabs.setActive(5)
			return a, nil
		case key.Matches(msg, globalKeys.Tab7):
			a.tabs.setActive(6)
			return a, nil
		case key.Matches(msg, globalKeys.Tab8):
			a.tabs.setActive(7)
			return a, nil
		}

	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		a.help.SetSize(msg.Width, msg.Height)
		contentHeight := max(0, a.height-3)
		sizeMsg := panel.SetSizeMsg{Width: a.width, Height: contentHeight}
		for i, p := range a.panels {
			updated, cmd := p.Update(sizeMsg)
			a.panels[i] = updated.(panel.Panel)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
		return a, tea.Batch(cmds...)

	case ConfigLoadedMsg:
		if msg.Err != nil {
			a.err = msg.Err
			return a, nil
		}
		cfgMsg := panel.ConfigDataMsg{
			Name:    msg.Name,
			Email:   msg.Email,
			Theme:   msg.Theme,
			OS:      msg.OS,
			Version: msg.Version,
		}
		for i, p := range a.panels {
			updated, cmd := p.Update(cfgMsg)
			a.panels[i] = updated.(panel.Panel)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
		return a, tea.Batch(cmds...)

	// Panel -> App messages that require chezmoi operations.
	case panel.EditSubmitMsg:
		return a, a.writeConfigAndReload(msg.Values)

	case panel.ThemeApplyMsg:
		return a, a.writeConfigAndReload(map[string]string{"theme": msg.Theme})

	case panel.RunChezmoiMsg:
		return a, a.runChezmoiCmd(msg.Mode)

	case panel.NvimFileSelectMsg:
		return a, a.loadNvimFileContent(msg.Path)

	// Panel -> App messages that require installer operations.
	case panel.InstallRequestMsg:
		return a, a.installApp(msg.App)

	case panel.CheckInstallMsg:
		return a, a.checkInstallStatus()

	// Results from async operations — broadcast to all panels.
	case panel.ShellContentMsg, panel.GitContentMsg, panel.AlacrittyContentMsg,
		panel.NvimFilesMsg, panel.NvimContentMsg, panel.ThemeColorsMsg,
		panel.ChezmoiOutputMsg, panel.ConfigUpdatedMsg,
		panel.InstallStatusMsg, panel.InstallResultMsg:
		for i, p := range a.panels {
			updated, cmd := p.Update(msg)
			a.panels[i] = updated.(panel.Panel)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
		return a, tea.Batch(cmds...)

	case ErrorMsg:
		a.err = msg.Err
		return a, nil
	}

	// Delegate to active panel.
	active := a.tabs.active
	if active >= 0 && active < len(a.panels) {
		updated, cmd := a.panels[active].Update(msg)
		a.panels[active] = updated.(panel.Panel)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	return a, tea.Batch(cmds...)
}

func (a App) View() string {
	if a.width == 0 {
		return "Loading..."
	}

	// Help overlay takes over the entire screen.
	if a.help.Visible() {
		return a.help.View()
	}

	var b strings.Builder

	b.WriteString(a.tabs.View(a.width))
	b.WriteString("\n")

	active := a.tabs.active
	contentHeight := max(0, a.height-3)
	content := ""
	if active >= 0 && active < len(a.panels) {
		content = a.panels[active].View()
	}
	contentStyle := lipgloss.NewStyle().
		Width(a.width).
		Height(contentHeight)
	b.WriteString(contentStyle.Render(content))

	// Status bar.
	help := ""
	if active >= 0 && active < len(a.panels) {
		help = a.panels[active].ShortHelp()
	}
	errStr := ""
	if a.err != nil {
		errStr = errorStyle.Render(fmt.Sprintf(" %s ", a.err.Error()))
	}
	gap := max(0, a.width-lipgloss.Width(help)-lipgloss.Width(errStr))
	statusLine := statusBarStyle.Width(a.width).Render(
		help + strings.Repeat(" ", gap) + errStr,
	)
	b.WriteString(statusLine)

	return b.String()
}

// --- Async commands ---

func (a App) loadConfig() tea.Cmd {
	return func() tea.Msg {
		if a.client == nil {
			return ConfigLoadedMsg{Err: fmt.Errorf("chezmoi client not initialized")}
		}
		data, err := a.client.Data()
		if err != nil {
			return ConfigLoadedMsg{Err: err}
		}
		version, _ := a.client.Version()
		return ConfigLoadedMsg{
			Name:    data.Name,
			Email:   data.Email,
			Theme:   data.Theme,
			OS:      data.OS,
			Version: version,
		}
	}
}

func (a App) loadShellConfigs() tea.Cmd {
	return func() tea.Msg {
		if a.client == nil {
			return nil
		}
		configs, _ := a.client.ShellConfigs()
		return tea.BatchMsg{
			func() tea.Msg { return panel.ShellContentMsg{Index: 0, Content: configs[0]} },
			func() tea.Msg { return panel.ShellContentMsg{Index: 1, Content: configs[1]} },
			func() tea.Msg { return panel.ShellContentMsg{Index: 2, Content: configs[2]} },
		}
	}
}

func (a App) loadGitConfig() tea.Cmd {
	return func() tea.Msg {
		if a.client == nil {
			return nil
		}
		content, err := a.client.GitConfig()
		if err != nil {
			return panel.GitContentMsg{Content: fmt.Sprintf("Error: %v", err)}
		}
		return panel.GitContentMsg{Content: content}
	}
}

func (a App) loadAlacrittyConfig() tea.Cmd {
	return func() tea.Msg {
		if a.client == nil {
			return nil
		}
		content, err := a.client.AlacrittyConfig()
		if err != nil {
			return panel.AlacrittyContentMsg{Content: fmt.Sprintf("Error: %v", err)}
		}
		return panel.AlacrittyContentMsg{Content: content}
	}
}

func (a App) loadNvimFiles() tea.Cmd {
	return func() tea.Msg {
		if a.client == nil {
			return nil
		}
		files, err := a.client.NvimFiles()
		if err != nil {
			return panel.NvimFilesMsg{Files: []string{}}
		}
		return panel.NvimFilesMsg{Files: files}
	}
}

func (a App) loadNvimFileContent(path string) tea.Cmd {
	return func() tea.Msg {
		if a.client == nil {
			return nil
		}
		content, err := a.client.NvimFileContent(path)
		if err != nil {
			return panel.NvimContentMsg{Path: path, Content: fmt.Sprintf("Error: %v", err)}
		}
		return panel.NvimContentMsg{Path: path, Content: content}
	}
}

func (a App) loadThemeColors(theme string) tea.Cmd {
	return func() tea.Msg {
		if a.client == nil {
			return nil
		}
		colors, err := a.client.ThemeColors(theme)
		if err != nil {
			return panel.ThemeColorsMsg{Theme: theme, Colors: nil}
		}
		return panel.ThemeColorsMsg{Theme: theme, Colors: colors}
	}
}

func (a App) writeConfigAndReload(values map[string]string) tea.Cmd {
	return func() tea.Msg {
		if a.client == nil {
			return panel.ConfigUpdatedMsg{Err: fmt.Errorf("chezmoi client not initialized")}
		}
		if err := a.client.WriteConfigAndApply(values); err != nil {
			return panel.ConfigUpdatedMsg{Err: err}
		}
		return panel.ConfigUpdatedMsg{Err: nil}
	}
}

func (a App) checkInstallStatus() tea.Cmd {
	return func() tea.Msg {
		if a.installer == nil {
			return panel.InstallStatusMsg{Statuses: nil}
		}
		statuses := a.installer.CheckAll()
		return panel.InstallStatusMsg{Statuses: statuses}
	}
}

func (a App) installApp(app installer.App) tea.Cmd {
	return func() tea.Msg {
		if a.installer == nil {
			return panel.InstallResultMsg{App: app, Err: fmt.Errorf("installer not initialized")}
		}
		err := a.installer.Install(app)
		return panel.InstallResultMsg{App: app, Err: err}
	}
}

func (a App) runChezmoiCmd(mode string) tea.Cmd {
	return func() tea.Msg {
		if a.client == nil {
			return panel.ChezmoiOutputMsg{Mode: mode, Err: fmt.Errorf("chezmoi client not initialized")}
		}
		switch mode {
		case "status":
			output, err := a.client.Status()
			return panel.ChezmoiOutputMsg{Mode: mode, Output: output, Err: err}
		case "diff":
			output, err := a.client.Diff()
			return panel.ChezmoiOutputMsg{Mode: mode, Output: output, Err: err}
		case "apply":
			err := a.client.Apply()
			if err != nil {
				return panel.ChezmoiOutputMsg{Mode: mode, Err: err}
			}
			return panel.ChezmoiOutputMsg{Mode: mode, Output: "Apply completed successfully."}
		default:
			return panel.ChezmoiOutputMsg{Mode: mode, Err: fmt.Errorf("unknown mode: %s", mode)}
		}
	}
}
