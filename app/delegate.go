package main

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func newItemDelegate(keys *delegateKeyMap, selected *selectInfo) list.DefaultDelegate {
	d := list.NewDefaultDelegate()
	d.Styles.SelectedTitle = d.Styles.SelectedTitle.Foreground(accent).BorderLeftForeground(accent)
	d.Styles.SelectedDesc = d.Styles.SelectedDesc.Foreground(accent).BorderLeftForeground(accent)

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		var title string
		var ep *episode

		if i, ok := m.SelectedItem().(episode); ok {
			title = i.Title()
			ep = &i
		} else {
			return nil
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.choose):
				(*selected).ep = ep
				return m.NewStatusMessage(statusMessageStyle("You chose " + title))
			}
		}

		return nil
	}

	help := []key.Binding{keys.choose}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}

	return d
}

type delegateKeyMap struct {
	choose key.Binding
}

func newDelegateKeyMap() *delegateKeyMap {
	return &delegateKeyMap{
		choose: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "choose"),
		),
	}
}
