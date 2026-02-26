package app

import (
	"ayo/internal/storage"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	screenFolderList screen = iota
	screenFolderForm
	screenTasksList
	screenTaskForm
	screenUnifiedTasks
	screenArchiveTasks
)

type rootModel struct {
	storagePath string
	archivePath string
	
	data AppData
	archiveData AppData

	screen screen

	width int
	height int

	selectedFolderID string

	folderForm folderFormModel
	taskForm taskFormModel
	taskFormReturnScreen screen

	folderList list.Model
	taskList list.Model
	unifiedList list.Model
	archiveList list.Model

	status string
	err error

}

func InitialModel() tea.Model{
	m := rootModel{
		storagePath: dataFilePath(),
		archivePath: archiveFilePath(),
		screen: screenFolderList,
		status: "Welcome to Ayo! Use arrow keys to navigate, Enter to select, Press a to add a folder, and 'q' to quit.",
	}

	if err := storage.Load(m.storagePath, &m.data); err != nil {
		m.err = err
		m.status = "Failed to load active tasks file."
	}
	if err := storage.Load(m.archivePath, &m.archiveData); err != nil {
		m.err = err
		m.status = "Failed to load archive file."
	}

	if len(m.data.Folders) > 0 {
		m.status = "Loaded folders from disk."
	}
	return m
}	


func (m rootModel) Init() tea.Cmd{
	return nil
}

func (m rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd){
	if m.screen == screenFolderForm {
		return m.updateFolderForm(msg)
	}
	if m.screen == screenTaskForm {
		return m.updateTaskForm(msg)
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.resizeLists()
		return m, nil
	case tea.KeyMsg:
		switch m.screen {
		case screenFolderList:
			return m.updateFolderListKeys(msg)
		case screenTasksList:
			return m.updateTasksListKeys(msg)
		case screenUnifiedTasks:
			return m.updateUnifiedTaskKeys(msg)
		case screenArchiveTasks:
			return m.updateArchiveTaskKeys(msg)
		}
	}

	return m, nil
}

func (m rootModel) View() string{
	switch m.screen {
	case screenFolderForm:
		return m.folderForm.View()
	case screenTaskForm:
		return m.taskForm.View()
	case screenTasksList:
		return m.viewTasksList()
	case screenUnifiedTasks:
		return m.viewUnifiedTasks()
	case screenArchiveTasks:
		return m.viewArchiveTasks()
	default:
		return m.viewFolders()
	}
}

func (m *rootModel) resizeLists() {
	w := m.width - 2
	if w < 30 {
		w = 30
	}
	h := m.height - 8
	if h < 8 {
		h = 8
	}
	m.folderList.SetSize(w, h)
	m.taskList.SetSize(w, h)
	m.unifiedList.SetSize(w, h)
	m.archiveList.SetSize(w, h)
}


func dataFilePath() string {
	if p := strings.TrimSpace(os.Getenv("AYO_DATA_PATH")); p != "" {
		return p
	}
	return filepath.Join("data", "todos.json")
}

func archiveFilePath() string {
	if p := strings.TrimSpace(os.Getenv("AYO_ARCHIVE_PATH")); p != "" {
		return p
	}
	return filepath.Join("data", "archive.json")
}

