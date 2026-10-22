package interactive

import (
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
)

type selectInfo struct{ ep *episode }

func newItemDelegate() list.DefaultDelegate {
	d := list.NewDefaultDelegate()
	d.Styles.SelectedTitle = d.Styles.SelectedTitle.Foreground(accent).BorderLeftForeground(accent)
	d.Styles.SelectedDesc = d.Styles.SelectedDesc.Foreground(accent).BorderLeftForeground(accent)

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
