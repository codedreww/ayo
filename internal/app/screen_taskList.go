package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type taskListKeyMap struct {
	Quit    key.Binding
	Back    key.Binding
	Add     key.Binding
	Edit    key.Binding
	Delete  key.Binding
	Toggle  key.Binding
	Mark    key.Binding
	Archive key.Binding
	Restore key.Binding
	Open    key.Binding
	Up      key.Binding
	Down    key.Binding
}

var taskListKeys = taskListKeyMap{
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Add: key.NewBinding(
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
	Archive: key.NewBinding(
		key.WithKeys("A"),
		key.WithHelp("A", "archive done"),
	),
	Restore: key.NewBinding(
		key.WithKeys("R"),
		key.WithHelp("R", "restore archived"),
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

func (m rootModel) updateFolderTaskKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if key.Matches(msg, taskListKeys.Quit) {
		return m, tea.Quit
	}

	if key.Matches(msg, taskListKeys.Back) {
		m.screen = screenFolderList
		m.status = ""
		return m, nil
	}

	selected, hasTask := m.selectedTaskFromTaskList()

	switch {
	case key.Matches(msg, taskListKeys.Add):
		if strings.TrimSpace(m.selectedFolderID) == "" {
			m.status = "No folder selected."
			return m, nil
		}

		m.taskForm = newTaskFormModel(m.selectedFolderID, Task{Priority: defaultPriority, State: taskStateTodo}, false)
		m.taskFormReturnScreen = screenFolderTasks
		m.screen = screenTaskForm
		return m, m.taskForm.Init()
	case key.Matches(msg, taskListKeys.Edit):
		if !hasTask {
			m.status = "No tasks to edit."
		}
		m.taskForm = newTaskFormModel(selected.folderID, selected.task, true)
		m.taskFormReturnScreen = screenFolderTasks
		m.screen = screenTaskForm
		return m, m.taskForm.Init()
	case key.Matches(msg, taskListKeys.Delete):
		if !hasTask {
			m.status = "No tasks to delete."
			return m, nil
		}
		if !m.deleteTask(selected.folderID, selected.task.ID) {
			m.status = "Task no longer exists"
			return m, nil
		}
		m.persistData(fmt.Sprintf("Deleted Task: %s", selected.task.Title))
		m.refreshTaskListForFolder(m.selectedFolderID, "")
		m.refreshUnifiedTaskList("", "")
		return m, nil
	case key.Matches(msg, taskListKeys.Toggle):
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
		m.refreshTaskListForFolder(m.selectedFolderID, selected.task.ID)
		m.refreshUnifiedTaskList(selected.folderID, selected.task.ID)
		return m, nil
	case key.Matches(msg, taskListKeys.Mark):
		if !hasTask {
			m.status = "No tasks to toggle."
			return m, nil
		}
		next := taskStateMarked
		if selected.task.State == taskStateDone {
			next = taskStateTodo
		}
		if !m.setTaskState(selected.folderID, selected.task.ID, next) {
			m.status = "Task no longer exists."
			return m, nil
		}
		m.persistData(fmt.Sprintf("Updated task mark: %s", selected.task.Title))
		m.refreshTaskListForFolder(m.selectedFolderID, selected.task.ID)
		m.refreshUnifiedTaskList(selected.folderID, selected.task.ID)
		return m, nil

	case key.Matches(msg, taskListKeys.Archive):
		moved := m.archiveCompletedTasks(m.selectedFolderID)
		if moved == 0 {
			m.status = "No completed tasks to archive in this folder."
			return m, nil
		}
		m.persistDataAndArchive(fmt.Sprintf("Archived %d completed task(s).", moved))
		m.refreshTaskListForFolder(m.selectedFolderID, "")
		m.refreshUnifiedTaskList("", "")
		m.refreshArchiveTaskList("", "")
		return m, nil
	}
	var cmd tea.Cmd
	m.taskList, cmd = m.taskList.Update(msg)
	return m, cmd
}

func (m rootModel) viewFolderTasks() string {
	var b strings.Builder
	folderName := "Tasks"
	if fi := m.folderIndexByID(m.selectedFolderID); fi >= 0 {
		folderName = "Tasks - " + m.data.Folders[fi].Name
	}
	b.WriteString(titleStyle.Render(folderName))
	b.WriteString("\n\n")
	b.WriteString(m.taskList.View())
	if s := renderStatus(m.status, m.err); s != "" {
		b.WriteString("\n\n")
		b.WriteString(s)
	}
	return b.String()

}

func (m *rootModel) refreshTaskListForFolder(folderID, selectedTaskID string) {
	fi := m.folderIndexByID(folderID)
	if fi < 0 {
		m.taskList.SetItems([]list.Item{})
		m.taskList.Select(0)
		return
	}

	items := make([]list.Item, 0, len(m.data.Folders[fi].Tasks))
	for _, t := range m.data.Folders[fi].Tasks {
		items = append(items, taskListItem{task: t, folderID: folderID})
	}
	m.taskList.SetItems(items)
	if len(items) == 0 {
		m.taskList.Select(0)
		return
	}
	if strings.TrimSpace(selectedTaskID) != "" {
		for i, it := range items {
			ti := it.(taskListItem)
			if ti.task.ID == selectedTaskID {
				m.taskList.Select(i)
				return
			}
		}
	}
	m.taskList.Select(clampCursor(m.taskList.Index(), len(items)))
}

func (m rootModel) selectedTaskFromTaskList() (taskListItem, bool) {
	selected := m.taskList.SelectedItem()
	if selected == nil {
		return taskListItem{}, false
	}
	it, ok := selected.(taskListItem)
	return it, ok
}
