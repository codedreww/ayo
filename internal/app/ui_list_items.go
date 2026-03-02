package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

// This fiel defines bubbles/list item adapters and delegate setup

// for folder and tasks, including metadata shown under each row

type folderListItem struct{
	folder Folder
}

func (i folderListItem) FilterValue() string{
	return strings.TrimSpace(i.folder.Name + " " + i.folder.Description)
}

func(i folderListItem) Title() string{
	return i.folder.Name
}

func (i folderListItem) Description() string{
	desc := strings.TrimSpace(i.folder.Description)
	if desc == "" {
		desc = "No Description"
	}
	return fmt.Sprintf("%s | %d tasks", desc, len(i.folder.Tasks))

}

type taskListItem struct{
	task Task
	folderID string
	folderName string
}

func (i taskListItem) FilterValue() string{
	return strings.TrimSpace(i.task.Title + " " + i.task.Description + " " + strings.Join(i.task.Tags, " ") + " " + i.folderName)
}

func (i taskListItem) Title() string{
	state := i.task.State
	if strings.TrimSpace(state) == "" {
		state = taskStateTodo
	}
	return fmt.Sprintf("%s %s", taskStateTag(state), i.task.Title)
}

func (i taskListItem) Description() string{
	desc := strings.TrimSpace(i.task.Description)
	if desc == ""{
		desc = "(no description)"
	}

	chips := []string{priorityChip(i.task.Priority)}
	if due := dueChip(i.task.DueDate); due != ""{
		chips = append(chips, due)
	}
	for _, tag := range i.task.Tags{
		chips = append(chips, tagChip(tag))
	}
	if strings.TrimSpace(i.folderName) != ""{
		chips = append(chips, folderChip(i.folderName))
	}

	chipLine := strings.Join(chips, " ")
	if chipLine == ""{
		return desc
	}
	return desc + "\n" + chipLine
}

func newList(width, height int, delegate list.DefaultDelegate) list.Model{
	l := list.New([]list.Item{}, delegate, width, height)
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(true)
	l.SetShowTitle(true)
	l.Styles.NoItems = mutedStyle
	return l
}

func newFolderDelegate() list.DefaultDelegate {
	d := list.NewDefaultDelegate()
	d.SetHeight(2)
	d.SetSpacing(1)
	d.Styles.NormalTitle = d.Styles.NormalTitle.PaddingLeft(1)
	d.Styles.NormalDesc = d.Styles.NormalDesc.PaddingLeft(1)
	d.Styles.SelectedTitle = d.Styles.SelectedTitle.PaddingLeft(1)
	d.Styles.SelectedDesc = d.Styles.SelectedDesc.PaddingLeft(1)

	d.ShortHelpFunc = func() []key.Binding {
		return [] key.Binding{
			folderListKeys.Add,
			folderListKeys.Edit,
			folderListKeys.Delete,
			folderListKeys.Open,
			folderListKeys.Unified,
			folderListKeys.ViewArc,
			folderListKeys.Archive,

		}
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{
			{folderListKeys.Add, folderListKeys.Edit, folderListKeys.Delete, folderListKeys.Open}, 
			{folderListKeys.Unified,folderListKeys.ViewArc,folderListKeys.Archive},
		}
	}


	return d
}

func newTaskDelegate() list.DefaultDelegate {
	d := list.NewDefaultDelegate()
	d.SetHeight(3)
	d.SetSpacing(1)
	d.Styles.NormalTitle = d.Styles.NormalTitle.PaddingLeft(1)
	d.Styles.NormalDesc = d.Styles.NormalDesc.PaddingLeft(1)
	d.Styles.SelectedTitle = d.Styles.SelectedTitle.PaddingLeft(1)
	d.Styles.SelectedDesc = d.Styles.SelectedDesc.PaddingLeft(1)

	d.ShortHelpFunc = func() []key.Binding {
		return [] key.Binding{
			//TODO

		}
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{
			//TODO
		}
	}

	return d
}

func withSelectedColor(d list.DefaultDelegate, color lipgloss.Color) list.DefaultDelegate {
	d.Styles.SelectedTitle = d.Styles.SelectedTitle.
		BorderForeground(color).
		Foreground(color)
	d.Styles.SelectedDesc = d.Styles.SelectedDesc.
		BorderForeground(color).
		Foreground(color)
	return d
}
