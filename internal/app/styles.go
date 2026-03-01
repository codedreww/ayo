package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("69"))
	mutedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))

	selectedRowStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("45"))
	statusStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("76"))
	errorStyle       = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("196"))

	doneTaskStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("242")).Strikethrough(true)

	priorityHighStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("230")).Background(lipgloss.Color("160")).Padding(0, 1)
	priorityMediumStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("0")).Background(lipgloss.Color("226")).Padding(0, 1)
	priorityLowStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("0")).Background(lipgloss.Color("114")).Padding(0, 1)

	dueStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("230")).Background(lipgloss.Color("62")).Padding(0, 1)
	tagStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("230")).Background(lipgloss.Color("238")).Padding(0, 1)
	folderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("0")).Background(lipgloss.Color("153")).Padding(0, 1)

	todoStateStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("0")).Background(lipgloss.Color("250")).Padding(0, 1)
	markedStateStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("0")).Background(lipgloss.Color("220")).Padding(0, 1)
	doneStateStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("230")).Background(lipgloss.Color("34")).Padding(0, 1)
	archivedStateStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("230")).Background(lipgloss.Color("240")).Padding(0, 1)
)

func priorityLabel(p int) string{
	switch p{
		case 3: return "High"
		case 2: return "Med"
		default: return "Low"
	}
}

func priorityChip(p int) string {
	label := priorityLabel(p)
	switch p {
	case 3:
		return priorityHighStyle.Render(label)
	case 2:
		return priorityMediumStyle.Render(label)
	default:
		return priorityLowStyle.Render(label)
	}
}

func dueChip(dueDate string) string {
	if strings.TrimSpace(dueDate) == "" {
		return ""
	}
	return dueStyle.Render("due " + dueDate)
}

func tagChip(tag string) string {
	tag = strings.TrimSpace(tag)
	if tag == "" {
		return ""
	}
	bg := tagColor(tag)
	style := tagStyle.Background(bg)
	return style.Render("#" + tag)
}

func folderChip(name string) string {
	return folderStyle.Render(name)
}

func taskStateTag(state string) string {
	switch state {
	case taskStateMarked:
		return markedStateStyle.Render("MARKED")
	case taskStateDone:
		return doneStateStyle.Render("DONE")
	case taskStateArchived:
		return archivedStateStyle.Render("ARCHIVED")
	default:
		return todoStateStyle.Render("TODO")
	}
}

func folderColorToLipglossColor(name string) lipgloss.Color {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "red":
		return lipgloss.Color("160")
	case "orange":
		return lipgloss.Color("208")
	case "yellow":
		return lipgloss.Color("226")
	case "green":
		return lipgloss.Color("34")
	case "purple":
		return lipgloss.Color("99")
	default:
		return lipgloss.Color("69") // blue fallback
	}
}

func tagColor(tag string) lipgloss.Color {
	colors := []lipgloss.Color{
		lipgloss.Color("31"),  // blue
		lipgloss.Color("35"),  // magenta
		lipgloss.Color("64"),  // teal
		lipgloss.Color("166"), // orange
		lipgloss.Color("99"),  // purple
		lipgloss.Color("28"),  // green
	}
	sum := 0
	for _, r := range strings.ToLower(tag) {
		sum += int(r)
	}
	return colors[sum%len(colors)]
}


func rowWithRightMeta(prefix, left string, rightParts []string, width int) string {
	rightParts = compactStrings(rightParts)
	line := prefix + left
	if len(rightParts) == 0 {
		return line
	}
	right := strings.Join(rightParts, " ")
	if width <= 0 {
		return line + "  " + right
	}
	leftWidth := lipgloss.Width(line)
	rightWidth := lipgloss.Width(right)
	gap := width - leftWidth - rightWidth
	if gap < 2 {
		gap = 2
	}
	return line + strings.Repeat(" ", gap) + right
}

func renderStatus(status string, err error) string {
	if err != nil {
		return errorStyle.Render(fmt.Sprintf("Error: %v", err))
	}
	if strings.TrimSpace(status) == "" {
		return ""
	}
	return statusStyle.Render(status)
}

func compactStrings(in []string) []string {
	if len(in) == 0 {
		return nil
	}
	out := make([]string, 0, len(in))
	for _, s := range in {
		if strings.TrimSpace(s) != "" {
			out = append(out, s)
		}
	}
	return out
}
