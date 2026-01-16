package ui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

// Chart renders a histogram-style chart using braille characters.
// Each data point is one character wide, and the height scales dynamically
// so the maximum value fills the available height.
type Chart struct {
	Width  int
	Height int
	Data   []float64
	Color  lipgloss.Style
	Title  string
}

// braille bit patterns for left and right columns
// Each column has 4 dots, allowing 2 data points per character.
// Left column dots (bottom to top): 7, 3, 2, 1
// Right column dots (bottom to top): 8, 6, 5, 4
var (
	leftDots  = [4]rune{0x40, 0x04, 0x02, 0x01} // dots 7, 3, 2, 1
	rightDots = [4]rune{0x80, 0x20, 0x10, 0x08} // dots 8, 6, 5, 4
)

func NewChart(width, height int, data []float64) Chart {
	return Chart{
		Width:  width,
		Height: height,
		Data:   data,
		Color:  lipgloss.NewStyle().Foreground(Colors.Secondary),
	}
}

func (c Chart) View() string {
	if len(c.Data) == 0 || c.Width == 0 || c.Height == 0 {
		return ""
	}

	// Account for border (2 chars width, 2 lines height)
	innerWidth := c.Width - 2

	maxVal := c.maxValue()
	displayMax := maxVal
	if maxVal == 0 {
		maxVal = 1
	}

	// Format labels and calculate label width
	maxLabel := formatChartValue(displayMax)
	labelWidth := max(len(maxLabel), 1)
	chartWidth := innerWidth - labelWidth - 1 // -1 for space between label and chart

	if chartWidth <= 0 {
		return ""
	}

	// Each character row represents 4 vertical dots
	dotsHeight := c.Height * 4

	// Calculate the height in dots for each data point
	heights := make([]int, len(c.Data))
	for i, v := range c.Data {
		heights[i] = int((v / maxVal) * float64(dotsHeight))
		if v > 0 && heights[i] == 0 {
			heights[i] = 1 // ensure non-zero values show at least 1 dot
		}
	}

	var content []string

	// Title line (centered over inner width)
	if c.Title != "" {
		titleLine := lipgloss.NewStyle().Width(innerWidth).Align(lipgloss.Center).Render(c.Title)
		content = append(content, titleLine)
	}

	// Build the chart row by row, from top to bottom
	// Use rightmost data points if we have more data than chart columns
	dataOffset := max(0, len(heights)-chartWidth*2)

	labelStyle := lipgloss.NewStyle().Width(labelWidth).Align(lipgloss.Left)
	for row := range c.Height {
		var sb strings.Builder
		rowBottomDot := (c.Height - 1 - row) * 4
		rowTopDot := rowBottomDot + 4

		for col := range chartWidth {
			dataIdxLeft := dataOffset + col*2
			dataIdxRight := dataOffset + col*2 + 1

			var char rune = 0x2800 // braille base character

			// Left column (first data point)
			if dataIdxLeft < len(heights) {
				char |= brailleColumn(heights[dataIdxLeft], rowBottomDot, rowTopDot, leftDots)
			}

			// Right column (second data point)
			if dataIdxRight < len(heights) {
				char |= brailleColumn(heights[dataIdxRight], rowBottomDot, rowTopDot, rightDots)
			}

			sb.WriteRune(char)
		}

		// Add label on first and last row
		var label string
		switch row {
		case 0:
			label = labelStyle.Render(maxLabel)
		case c.Height - 1:
			label = labelStyle.Render("0")
		default:
			label = labelStyle.Render("")
		}

		chartRow := c.Color.Render(sb.String())
		content = append(content, label+" "+chartRow)
	}

	// Wrap in rounded border
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#6272a4"))

	return borderStyle.Render(strings.Join(content, "\n"))
}

// brailleColumn returns the braille bits for a single column based on height
func brailleColumn(h, rowBottom, rowTop int, dots [4]rune) rune {
	if h <= rowBottom {
		return 0
	}

	var bits rune
	dotsToFill := min(h-rowBottom, 4)
	for i := range dotsToFill {
		bits |= dots[i]
	}
	return bits
}

func (c Chart) maxValue() float64 {
	if len(c.Data) == 0 {
		return 0
	}
	max := c.Data[0]
	for _, v := range c.Data[1:] {
		if v > max {
			max = v
		}
	}
	return max
}

// Helpers

func formatChartValue(v float64) string {
	if v >= 1_000_000 {
		return fmt.Sprintf("%.1fM", v/1_000_000)
	}
	if v >= 1_000 {
		return fmt.Sprintf("%.1fK", v/1_000)
	}
	return fmt.Sprintf("%.0f", v)
}
