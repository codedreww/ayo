package app

import "github.com/charmbracelet/bubbles/list"

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
	archiveLlist list.Model

	status string
	err error

}

func InitialModel() tea.Model{
	return rootModel{screen: screenFolderList}
}	


func (m rootModel) Init() tea.Cmd{
	return nil
}

func (m rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd){

}

func (m rootModel) View() string{
	switch m.screen{

	}
}