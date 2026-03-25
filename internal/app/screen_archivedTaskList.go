package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type archivedTaskListKeyMap struct {
	Quit    key.Binding
	Back    key.Binding
	Delete  key.Binding
	Restore key.Binding
	Open    key.Binding
	Up      key.Binding
	Down    key.Binding
}

var archivedTaskListKeys = archivedTaskListKeyMap{
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete task"),
	),
	Restore: key.NewBinding(
		key.WithKeys("R"),
		key.WithHelp("R", "restore archived"),
	),
	Open: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "open task"),
	),
	Up: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("j", "down"),
	),
}

func (m rootModel) updateArchiveTaskKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if key.Matches(msg, archivedTaskListKeys.Quit) {
		return m, tea.Quit
	}

	if key.Matches(msg, archivedTaskListKeys.Back) {
		m.screen = m.archiveReturnScreen
		m.status = ""
		return m, nil
	}

	selected, hasTask := m.selectedTaskFromArchiveList()

	switch {
	case key.Matches(msg, archivedTaskListKeys.Open):
		if !hasTask {
			m.status = "No archived tasks selected."
			return m, nil
		}
		m.status = fmt.Sprintf("Archived task from folder: %s", selected.folderName)
		return m, nil
	case key.Matches(msg, archivedTaskListKeys.Delete):
		if !hasTask {
			m.status = "No archived tasks to delete."
			return m, nil
		}
		if !m.deleteArchivedTask(selected.folderID, selected.task.ID) {
			m.status = "Archived task no longer exists."
			return m, nil
		}
		m.persistArchiveOnly(fmt.Sprintf("Deleted archived task: %s", selected.task.Title))
		m.refreshArchiveTaskList("", "")
		return m, nil
	case key.Matches(msg, archivedTaskListKeys.Restore):
		if !hasTask {
			m.status = "No archived tasks to restore."
			return m, nil
		}
		if !m.restoreArchivedTask(selected.folderID, selected.task.ID) {
			m.status = "Archived task no longer exists."
			return m, nil
		}
		m.persistDataAndArchive(fmt.Sprintf("Restored task: %s", selected.task.Title))
		m.refreshArchiveTaskList("", "")
		m.refreshFolderList(selected.folderID)
		m.refreshUnifiedTaskList(selected.folderID, "")
		m.refreshTaskListForFolder(m.selectedFolderID, "")
		return m, nil
	}

	var cmd tea.Cmd
	m.archiveList, cmd = m.archiveList.Update(msg)
	return m, cmd
}

func (m rootModel) viewArchiveTasks() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Archive"))
	b.WriteString("\n")
	b.WriteString(mutedStyle.Render("Archived tasks from archive.json"))
	b.WriteString("\n\n")
	b.WriteString(m.archiveList.View())
	if s := renderStatus(m.status, m.err); s != "" {
		b.WriteString("\n\n")
		b.WriteString(s)
	}
	return b.String()
}

func (m *rootModel) refreshArchiveTaskList(selectedFolderID, selectedTaskID string) {
	items := make([]list.Item, 0)
	for _, f := range m.archiveData.Folders {
		for _, t := range f.Tasks {
			items = append(items, taskListItem{task: t, folderID: f.ID, folderName: f.Name})
		}
	}
	m.archiveList.SetItems(items)
	if len(items) == 0 {
		m.archiveList.Select(0)
		return
	}
	if strings.TrimSpace(selectedTaskID) != "" {
		for i, it := range items {
			ti := it.(taskListItem)
			if ti.task.ID == selectedTaskID && (selectedFolderID == "" || ti.folderID == selectedFolderID) {
				m.archiveList.Select(i)
				return
			}
		}
	}
	m.archiveList.Select(clampCursor(m.archiveList.Index(), len(items)))
}

func (m rootModel) selectedTaskFromArchiveList() (taskListItem, bool) {
	selected := m.archiveList.SelectedItem()
	if selected == nil {
		return taskListItem{}, false
	}
	it, ok := selected.(taskListItem)
	return it, ok
}
