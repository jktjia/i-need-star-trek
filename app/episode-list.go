package main

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type episodelist struct {
	list         list.Model
	keys         *keyMap
	delegateKeys *delegateKeyMap
	selected     *selectInfo
}

func (e episodelist) Update(msg tea.Msg) (episodelist, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		e.list.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, e.keys.goBack):
			e.selected.ep = nil
			return e, nil
		}
	}

	newListModel, cmd := e.list.Update(msg)
	e.list = newListModel
	cmds = append(cmds, cmd)

	return e, tea.Batch(cmds...)
}

func (e *episodelist) SetEpisodes(episodes []episode) tea.Cmd {
	items := []list.Item{}
	for _, ep := range episodes {
		items = append(items, ep)
	}

	cmd := e.list.SetItems(items)
	e.list.ResetSelected()
	return cmd
}

func (e episodelist) View() string {
	return e.list.View()
}

func newEpisodeList(episodes []episode, keyMap *keyMap, selected *selectInfo) episodelist {
	var (
		delegateKeys = newDelegateKeyMap()
	)

	items := []list.Item{}
	for _, e := range episodes {
		items = append(items, e)
	}

	// Setup list
	delegate := newItemDelegate(delegateKeys, selected)
	episodeList := list.New(items, delegate, 0, 0)
	episodeList.Title = "Star Trek Episodes"
	episodeList.Styles.Title = titleStyle
	episodeList.SetFilteringEnabled(false)
	episodeList.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{keyMap.showFilters, keyMap.randomGenerator}
	}
	episodeList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{keyMap.showFilters, keyMap.randomGenerator}
	}

	return episodelist{
		episodeList,
		keyMap,
		delegateKeys,
		selected,
	}
}
