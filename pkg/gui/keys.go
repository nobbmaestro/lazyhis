package gui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/nobbmaestro/lazyhis/pkg/config"
)

type keyMap struct {
	ActionAcceptSelected  key.Binding
	ActionPrefillSelected key.Binding
	ActionNextFilter      key.Binding
	ActionPrevFilter      key.Binding
	ActionJumpDown        key.Binding
	ActionJumpUp          key.Binding
	ActionMoveDown        key.Binding
	ActionMoveUp          key.Binding
	ActionQuit            key.Binding
	ActionShowHelp        key.Binding
}

func createKeyMap(keys config.GuiKeyMap) keyMap {
	return keyMap{
		ActionMoveUp: key.NewBinding(
			key.WithKeys(keys.MoveUp...),
			key.WithHelp(prettyKey(keys.MoveUp), "move up"),
		),
		ActionMoveDown: key.NewBinding(
			key.WithKeys(keys.MoveDown...),
			key.WithHelp(prettyKey(keys.MoveDown), "move down"),
		),
		ActionJumpUp: key.NewBinding(
			key.WithKeys(keys.JumpUp...),
			key.WithHelp(prettyKey(keys.JumpUp), "jump up"),
		),
		ActionJumpDown: key.NewBinding(
			key.WithKeys(keys.JumpDown...),
			key.WithHelp(prettyKey(keys.JumpDown), "jump down"),
		),
		ActionAcceptSelected: key.NewBinding(
			key.WithKeys(keys.AcceptSelected...),
			key.WithHelp(prettyKey(keys.AcceptSelected), "accept"),
		),
		ActionPrefillSelected: key.NewBinding(
			key.WithKeys(keys.PrefillSelected...),
			key.WithHelp(prettyKey(keys.PrefillSelected), "prefill"),
		),
		ActionQuit: key.NewBinding(
			key.WithKeys(keys.Quit...),
			key.WithHelp(prettyKey(keys.Quit), "quit"),
		),
		ActionShowHelp: key.NewBinding(
			key.WithKeys(keys.ShowHelp...),
			key.WithHelp(prettyKey(keys.ShowHelp), "toggle help"),
		),
		ActionNextFilter: key.NewBinding(
			key.WithKeys(keys.NextFilter...),
			key.WithHelp(prettyKey(keys.NextFilter), "next filter"),
		),
		ActionPrevFilter: key.NewBinding(
			key.WithKeys(keys.PrevFilter...),
			key.WithHelp(prettyKey(keys.PrevFilter), "prev filter"),
		),
	}
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.ActionShowHelp,
		k.ActionAcceptSelected,
		k.ActionNextFilter,
		k.ActionQuit,
	}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.ActionAcceptSelected, k.ActionPrefillSelected},
		{k.ActionNextFilter, k.ActionPrevFilter},
		{k.ActionMoveUp, k.ActionMoveDown},
		{k.ActionJumpDown, k.ActionJumpUp},
		{k.ActionShowHelp, k.ActionQuit},
	}
}

func prettyKey(keys []string) string {
	if len(keys) == 0 {
		return ""
	}

	key := keys[0]

	replacements := map[string]string{
		"ctrl":  "⌃",
		"shift": "⇧",
		"enter": "↵",
		"tab":   "⇥",
	}

	parts := strings.Split(key, "+")
	for i, part := range parts {
		if sym, ok := replacements[part]; ok {
			parts[i] = sym
		}
	}

	return fmt.Sprintf("%3s", strings.Join(parts, ""))
}
