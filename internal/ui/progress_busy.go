package ui

import (
	"image/color"
	"math/rand/v2"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type ProgressBusy struct {
	width int
	color color.Color

	pattern []rune
}

type ProgressBusyTickMsg struct{}

func NewProgressBusy(width int, clr color.Color) ProgressBusy {
	return ProgressBusy{
		width:   width,
		color:   clr,
		pattern: generateBraillePattern(width),
	}
}

func (p ProgressBusy) Init() tea.Cmd {
	return p.tick()
}

func (p ProgressBusy) Update(msg tea.Msg) (ProgressBusy, tea.Cmd) {
	switch msg.(type) {
	case ProgressBusyTickMsg:
		p.pattern = generateBraillePattern(p.width)
		return p, p.tick()
	}
	return p, nil
}

func (p ProgressBusy) View() string {
	if p.width <= 0 {
		return ""
	}

	return lipgloss.NewStyle().Foreground(p.color).Render(string(p.pattern))
}

// Private

func (p ProgressBusy) tick() tea.Cmd {
	return tea.Tick(50*time.Millisecond, func(time.Time) tea.Msg {
		return ProgressBusyTickMsg{}
	})
}

// Helpers

func generateBraillePattern(width int) []rune {
	pattern := make([]rune, width)
	for i := range pattern {
		// Braille patterns: U+2800 to U+28FF (256 patterns)
		pattern[i] = rune(0x2800 + rand.IntN(256))
	}
	return pattern
}
