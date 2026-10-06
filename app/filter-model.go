package main

import (
	"fmt"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type filterModel struct {
	height          int
	selectedFilters filters
	pointer         int
	viewport        viewport.Model
	ready           bool
	tab             int
	help            help.Model
	keys            keyMap
}

func (f filterModel) content() string {
	filterView := ""
	switch filterTabs[f.tab] {
	case "Show":
		for i, s := range showList {
			content := " "
			if slices.Contains(*f.selectedFilters.shows, s.Abbr) {
				content = "x"
			}
			str := fmt.Sprintf(" [%s] %s", content, s.FullName)
			if i == f.pointer {
				filterView += selectedStyle.Render(">"+str) + "\n"
			} else {
				filterView += " " + str + "\n"
			}
		}
	case "Species":
		for i, s := range speciesList {
			content := " "
			if slices.Contains(*f.selectedFilters.species, s) {
				content = "x"
			}
			str := fmt.Sprintf(" [%s] %s", content, s)
			if i == f.pointer {
				filterView += selectedStyle.Render(">"+str) + "\n"
			} else {
				filterView += " " + str + "\n"
			}
		}
	case "Other":
		for i, e := range eventList {
			content := " "
			if slices.Contains(*f.selectedFilters.events, e.Short) {
				content = "x"
			}
			str := fmt.Sprintf(" [%s] %s", content, e.FullName)
			if i == f.pointer {
				filterView += selectedStyle.Render(">"+str) + "\n"
			} else {
				filterView += " " + str + "\n"
			}
		}
	}
	return filterStyle.Render(filterView)
}

type filterKeyMap struct {
	up             key.Binding
	down           key.Binding
	left           key.Binding
	right          key.Binding
	toggleSelected key.Binding
}

var filterKeys = filterKeyMap{
	down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "prev page"),
	),
	right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "next page"),
	),
	toggleSelected: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "choose"),
	),
}

func (f filterModel) Update(msg tea.Msg) (filterModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		f.height = msg.Height
		footerHeight := lipgloss.Height(f.footerView())
		headerHeight := lipgloss.Height(f.headerView()) + 1
		verticalMarginHeight := headerHeight + footerHeight

		if !f.ready {
			f.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			f.viewport.YPosition = headerHeight
			f.ready = true
		} else {
			f.viewport.Width = msg.Width
			f.viewport.Height = msg.Height - verticalMarginHeight
			f.viewport.YPosition = headerHeight
		}

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, filterKeys.up):
			newVal := max(f.pointer-1, 0)
			f.pointer = newVal

		case key.Matches(msg, filterKeys.down):
			var max int
			switch filterTabs[f.tab] {
			case "Show":
				max = len(showList)
			case "Species":
				max = len(speciesList)
			case "Other":
				max = len(eventList)
			}
			newVal := min(f.pointer+1, max-1)
			f.pointer = newVal

		case key.Matches(msg, filterKeys.left):
			newVal := max(f.tab-1, 0)
			if f.tab != newVal {
				f.tab = newVal
				f.pointer = 0
			}

		case key.Matches(msg, filterKeys.right):
			newVal := min(f.tab+1, len(filterTabs)-1)
			if f.tab != newVal {
				f.tab = newVal
				f.pointer = 0
			}

		case key.Matches(msg, filterKeys.toggleSelected):
			switch filterTabs[f.tab] {
			case "Show":
				if f.pointer < len(showList) {
					selected := showList[f.pointer]
					f.selectedFilters.toggleShow(selected)
				}
			case "Species":
				if f.pointer < len(speciesList) {
					selected := speciesList[f.pointer]
					f.selectedFilters.toggleSpecies(selected)
				}
			case "Other":
				if f.pointer < len(eventList) {
					selected := eventList[f.pointer]
					f.selectedFilters.toggleEvent(selected)
				}
			}

		case key.Matches(msg, f.keys.toggleHelpMenu):
			f.help.ShowAll = !f.help.ShowAll

			headerHeight := lipgloss.Height(f.headerView())
			footerHeight := lipgloss.Height(f.footerView())
			verticalMarginHeight := headerHeight + footerHeight + 1

			f.viewport.Height = f.height - verticalMarginHeight
		}
	}

	// Handle keyboard and mouse events in the viewport
	newViewport, cmd := f.viewport.Update(msg)
	f.viewport = newViewport

	f.viewport.SetContent(f.content())

	return f, cmd
}

func (f filterModel) headerView() string {
	header := []string{}
	for i, tab := range filterTabs {
		var tabStr = tab
		if f.tab == i {
			tabStr = titleStyle.Render(tabStr)
		} else {
			tabStr = " " + tabStr + " "
		}
		header = append(header, tabStr)
	}
	tabs := strings.Join(header, "|")
	full := fmt.Sprintf("%s   %s", descStyle.Render("Filters"), tabs)
	return lipgloss.NewStyle().Padding(0, 0, 1).Render(full)
}

func (f filterModel) footerView() string {
	return lipgloss.NewStyle().Padding(1, 0, 0).Render(f.help.View(f))
}

func (f filterModel) ShortHelp() []key.Binding {
	return []key.Binding{filterKeys.down, filterKeys.up, filterKeys.toggleSelected, f.keys.toggleHelpMenu}
}

func (f filterModel) FullHelp() [][]key.Binding {
	return [][]key.Binding{f.ShortHelp(), {filterKeys.right, filterKeys.left, f.keys.randomGenerator, f.keys.goBack}}
}

func (f filterModel) View() string {
	return fmt.Sprintf("%s\n%s\n%s", f.headerView(), f.viewport.View(), f.footerView())
}

func newFilters(keys keyMap) filterModel {
	f := filterModel{}
	f.selectedFilters = filters{&[]string{}, &[]string{}, &[]string{}}
	f.ready = false
	f.keys = keys
	f.help = help.New()
	return f
}
