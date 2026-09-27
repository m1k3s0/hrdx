package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestFuzzyMatch(t *testing.T) {
	cases := []struct {
		haystack, needle string
		want             bool
	}{
		{"api › zot 1", "", true},
		{"api › zot 1", "apzot", true},
		{"api › zot 1", "z1", true},
		{"api › zot 1", "xyz", false},
		{"web › shell 2", "wsh2", true},
	}
	for _, c := range cases {
		if got := fuzzyMatch(c.haystack, c.needle); got != c.want {
			t.Fatalf("fuzzyMatch(%q, %q) = %v, want %v", c.haystack, c.needle, got, c.want)
		}
	}
}

func TestFindFiltersAndJumps(t *testing.T) {
	model := newTestModel("/tmp/api", "/tmp/web")
	target := model.spaces[1].tab().panes[0]
	model.paneAttention[target.id] = true

	updated, _ := model.updateKey(tea.KeyMsg{Type: tea.KeyCtrlB})
	model = updated.(Model)
	updated, _ = model.updateKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	model = updated.(Model)
	if model.mode != modeFind {
		t.Fatalf("mode = %d, want modeFind after ctrl+b /", model.mode)
	}
	if got := len(model.findCandidates()); got != 2 {
		t.Fatalf("candidates = %d, want 2", got)
	}

	updated, _ = model.updateKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("web")})
	model = updated.(Model)
	candidates := model.findCandidates()
	if len(candidates) != 1 || candidates[0].spaceIndex != 1 {
		t.Fatalf("filtered candidates = %+v, want only web", candidates)
	}

	updated, _ = model.updateKey(tea.KeyMsg{Type: tea.KeyEnter})
	model = updated.(Model)
	if model.mode != modeFind && model.selected != 1 {
		t.Fatalf("selected = %d, want 1 after jump", model.selected)
	}
	if model.mode != modeTerminal {
		t.Fatalf("mode = %d, want modeTerminal after jump", model.mode)
	}
	if model.paneAttention[target.id] {
		t.Fatal("find jump did not clear focused attention")
	}
}

func TestFindCustomNavigationKeys(t *testing.T) {
	model := newTestModel("/tmp/api", "/tmp/web")
	model.navKeys = buildNavigationKeys(map[string]string{
		"navigate-up": "home", "navigate-down": "end",
	})
	updated, _ := model.openFind()
	model = updated.(Model)
	model.findIndex = 1

	updated, _ = model.updateKey(tea.KeyMsg{Type: tea.KeyHome})
	model = updated.(Model)
	if model.findIndex != 0 {
		t.Fatalf("findIndex after home = %d, want 0", model.findIndex)
	}
	updated, _ = model.updateKey(tea.KeyMsg{Type: tea.KeyEnd})
	model = updated.(Model)
	if model.findIndex != 1 {
		t.Fatalf("findIndex after end = %d, want 1", model.findIndex)
	}
}

func TestFindEscCloses(t *testing.T) {
	model := newTestModel("/tmp/api")
	updated, _ := model.openFind()
	model = updated.(Model)
	updated, _ = model.updateKey(tea.KeyMsg{Type: tea.KeyEsc})
	model = updated.(Model)
	if model.mode != modeTerminal {
		t.Fatalf("mode = %d, want modeTerminal after esc", model.mode)
	}
}
