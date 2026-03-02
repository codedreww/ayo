package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type folderListKeyMap struct {
	Quit    key.Binding
	Add     key.Binding
	Edit    key.Binding
	Delete  key.Binding
	Open    key.Binding
	Unified key.Binding
	ViewArc key.Binding
	Archive key.Binding
	Up      key.Binding
	Down    key.Binding
}

var folderListKeys = folderListKeyMap{
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Add: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "add folder"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit folder"),
	),
	Delete: key.NewBinding(
		key.WithKeys("D"),
		key.WithHelp("D", "delete folder"),
	),
	Open: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "open tasks"),
	),
	Unified: key.NewBinding(
		key.WithKeys("u"),
		key.WithHelp("u", "all tasks"),
	),
	ViewArc: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "archive view"),
	),
	Archive: key.NewBinding(
		key.WithKeys("A"),
		key.WithHelp("A", "archive done"),
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

func (m rootModel) updateFolderListKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, folderListKeys.Quit):
		return m, tea.Quit
	case key.Matches(msg, folderListKeys.Add):
		m.folderForm = newFolderFormModel(Folder{Color: ""}, false)
		m.screen = screenFolderForm
		return m, m.folderForm.Init()
	case key.Matches(msg, folderListKeys.Edit):
		selected, ok := m.selectedFolderItem()
		if !ok {
			m.status = "No folders to edit."
			return m, nil
		}
		m.folderForm = newFolderFormModel(selected.folder, true)
		m.screen = screenFolderForm
		return m, m.folderForm.Init()
	case key.Matches(msg, folderListKeys.Delete):
		selected, ok := m.selectedFolderItem()
		if !ok {
			m.status = "No folders to delete."
			return m, nil
		}
		fi := m.folderIndexByID(selected.folder.ID)
		if fi < 0 {
			m.status = "Folder no longer exists"
			return m, nil
		}
		m.data.Folders = append(m.data.Folders[:fi], m.data.Folders[fi+1:]...)
		m.persistData(fmt.Sprintf("Deleted Folder: %s", selected.folder.Name))
		m.refreshFolderList("")
		m.refreshUnifiedTaskList("", "")
		return m, nil
	case key.Matches(msg, folderListKeys.Open):
		selected, ok := m.selectedFolderItem()
		if !ok {
			m.status = "No folders to open."
			return m, nil
		}
		m.selectedFolderID = selected.folder.ID
		m.refreshTaskListForFolder(selected.folder.ID, "")
		m.screen = screenTasksList
		m.status = ""
		return m, nil
	case key.Matches(msg, folderListKeys.Unified):
		m.refreshUnifiedTaskList("", "")
		m.screen = screenUnifiedTasks
		m.status = ""
		return m, nil
	case key.Matches(msg, folderListKeys.ViewArc):
		m.refreshArchiveTaskList("", "")
		m.screen = screenArchiveTasks
		m.status = ""
		return m, nil
	case key.Matches(msg, folderListKeys.Archive):
		moved := m.archiveCompletedTasks("")
		if moved == 0 {
			m.status = "No completed tasks to archive."
			return m, nil
		}
		m.persistDataAndArchive(fmt.Sprintf("Archived %d completed task(s).", moved))
		m.refreshFolderList("")
		m.refreshUnifiedTaskList("", "")
		m.refreshArchiveTaskList("", "")
		return m, nil
	}
	var cmd tea.Cmd
	m.folderList, cmd = m.folderList.Update(msg)
	m.syncFolderHoverStyle()
	return m, cmd
}

func (m rootModel) viewFolders() string {
	var b strings.Builder
	b.WriteString(m.folderList.View())

	if s := renderStatus(m.status, m.err); s != "" {
		b.WriteString("\n\n")
		b.WriteString(s)
	}
	return b.String()
}

func (m *rootModel) refreshFolderList(selectedFolderID string) {
	items := make([]list.Item, 0, len(m.data.Folders))
	for _, f := range m.data.Folders {
		items = append(items, folderListItem{folder: f})
	}
	m.folderList.SetItems(items)
	if len(items) == 0 {
		m.folderList.Select(0)
		m.syncFolderHoverStyle()
		return
	}
	if strings.TrimSpace(selectedFolderID) != "" {
		for i, it := range items {
			fi := it.(folderListItem)
			if fi.folder.ID == selectedFolderID {
				m.folderList.Select(i)
				m.syncFolderHoverStyle()
				return
			}
		}
	}
	m.folderList.Select(clampCursor(m.folderList.Index(), len(items)))
	m.syncFolderHoverStyle()
}
func (m *rootModel) syncFolderHoverStyle() {
	selected := m.folderList.SelectedItem()
	if selected == nil {
		m.folderDelegate = newFolderDelegate()
		m.folderList.SetDelegate(m.folderDelegate)
		return
	}

	item, ok := selected.(folderListItem)
	if !ok {
		m.folderDelegate = newFolderDelegate()
		m.folderList.SetDelegate(m.folderDelegate)
		return
	}
	color := folderColorToLipglossColor(item.folder.Color)
	m.folderDelegate = withSelectedColor(newFolderDelegate(), color)
	m.folderList.SetDelegate(m.folderDelegate)
}

func (m rootModel) selectedFolderItem() (folderListItem, bool) {
	selected := m.folderList.SelectedItem()
	if selected == nil {
		return folderListItem(), false
	}
	it, ok := selected.(folderListItem)
	return it, ok
}
