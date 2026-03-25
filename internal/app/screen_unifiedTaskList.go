package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type unifiedTaskListKeyMap struct {
	Quit    key.Binding
	Back    key.Binding
	Add     key.Binding
	Edit    key.Binding
	Delete  key.Binding
	Toggle  key.Binding
	Mark    key.Binding
	ViewArc key.Binding
	Archive key.Binding
	Open    key.Binding
	Up      key.Binding
	Down    key.Binding
}

var unifiedTaskListKeys = unifiedTaskListKeyMap{
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Add: key.NewBinding( // #TODO need to figure out how to add a task from the unified screen
		key.WithKeys("a"),
		key.WithHelp("a", "add task"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit task"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete task"),
	),
	Toggle: key.NewBinding(
		key.WithKeys(" ", "x"),
		key.WithHelp("space", "toggle done"),
	),
	Mark: key.NewBinding(
		key.WithKeys("m"),
		key.WithHelp("m", "toggle mark"),
	),
	ViewArc: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "archive view"),
	),
	Archive: key.NewBinding(
		key.WithKeys("A"),
		key.WithHelp("A", "archive done"),
	),
	Open: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "open folder"),
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

func (m rootModel) updateUnifiedTaskKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if key.Matches(msg, unifiedTaskListKeys.Quit) {
		return m, tea.Quit
	}

	if key.Matches(msg, unifiedTaskListKeys.Back) {
		m.screen = screenFolderList
		m.status = ""
		return m, nil
	}

	selected, hasTask := m.selectedTaskFromUnifiedList()

	switch {
	case key.Matches(msg, unifiedTaskListKeys.Open):
		if !hasTask {
			m.status = "No tasks to open."
			return m, nil
		}
		m.selectedFolderID = selected.folderID
		m.refreshTaskListForFolder(selected.folderID, selected.task.ID)
		m.screen = screenFolderTasks
		m.status = ""
	case key.Matches(msg, unifiedTaskListKeys.Edit):
		if !hasTask {
			m.status = "No tasks to open."
			return m, nil
		}
		m.taskForm = newTaskFormModel(selected.folderID, selected.task, true)
		m.taskFormReturnScreen = screenUnifiedTasks
		m.screen = screenTaskForm
		return m, m.taskForm.Init()
	case key.Matches(msg, unifiedTaskListKeys.Toggle):
		if !hasTask {
			m.status = "No tasks to toggle."
			return m, nil
		}
		next := taskStateDone
		if selected.task.State == taskStateDone {
			next = taskStateTodo
		}
		if !m.setTaskState(selected.folderID, selected.task.ID, next) {
			m.status = "Task no longer exists."
			return m, nil
		}
		m.persistData(fmt.Sprintf("Updated task status: %s", selected.task.Title))
		m.refreshUnifiedTaskList(selected.folderID, selected.task.ID)
		m.refreshTaskListForFolder(m.selectedFolderID, "")
		return m, nil

	case key.Matches(msg, unifiedTaskListKeys.Mark):
		if !hasTask {
			m.status = "No tasks to toggle."
			return m, nil
		}

		next := taskStateMarked
		if selected.task.State == taskStateMarked {
			next = taskStateTodo
		}

		if !m.setTaskState(selected.folderID, selected.task.ID, next) {
			m.status = "Task no longer exists."
			return m, nil
		}
		m.persistData(fmt.Sprintf("Updated task status: %s", selected.task.Title))
		m.refreshUnifiedTaskList(selected.folderID, selected.task.ID)
		m.refreshTaskListForFolder(m.selectedFolderID, "")
		return m, nil
	case key.Matches(msg, unifiedTaskListKeys.Archive):
		moved := m.archiveCompletedTasks("")
		if moved == 0 {
			m.status = "No completed tasks to archive."
			return m, nil
		}
		m.persistDataAndArchive(fmt.Sprintf("Archived %d completed task(s).", moved))
		m.refreshUnifiedTaskList("", "")
		m.refreshTaskListForFolder(m.selectedFolderID, "")
		m.refreshFolderList("")
		m.refreshArchiveTaskList("", "")
		return m, nil
	case key.Matches(msg, unifiedTaskListKeys.ViewArc):
		m.refreshArchiveTaskList("", "")
		m.archiveReturnScreen = screenUnifiedTasks
		m.screen = screenArchiveTasks
		m.status = ""
		return m, nil

	}
	var cmd tea.Cmd
	m.unifiedList, cmd = m.unifiedList.Update(msg)
	return m, cmd
}

func (m rootModel) viewUnifiedTasks() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("All Tasks"))
	b.WriteString("\n")
	b.WriteString(mutedStyle.Render("Across all folders"))
	b.WriteString("\n\n")
	b.WriteString(m.unifiedList.View())
	b.WriteString("\n")
	if s := renderStatus(m.status, m.err); s != "" {
		b.WriteString("\n\n")
		b.WriteString(s)
	}
	return b.String()
}

func (m *rootModel) refreshUnifiedTaskList(selectedFolderID, selectedTaskID string) {
	items := make([]list.Item, 0)
	for _, f := range m.data.Folders {
		for _, t := range f.Tasks {
			items = append(items, taskListItem{task: t, folderID: f.ID, folderName: f.Name})
		}
	}

	m.taskList.SetItems(items)
	if len(items) == 0 {
		m.taskList.Select(0)
		return
	}
	if strings.TrimSpace(selectedTaskID) != "" {
		for i, it := range items {
			ti := it.(taskListItem)
			if ti.task.ID == selectedTaskID && (selectedFolderID == "" || ti.folderID == selectedFolderID) {
				m.unifiedList.Select(i)
				return
			}
		}
	}
	m.unifiedList.Select(clampCursor(m.unifiedList.Index(), len(items)))
}

func (m rootModel) selectedTaskFromUnifiedList() (taskListItem, bool) {
	selected := m.unifiedList.SelectedItem()
	if selected == nil {
		return taskListItem{}, false
	}

	it, ok := selected.(taskListItem)
	return it, ok
}
