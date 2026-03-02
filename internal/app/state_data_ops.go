package app

import (
	"ayo/internal/storage"
)


// Function to save data into json
func (m *rootModel) persistData(okStatus string){
	normalizeData(&m.data)
	if err := storage.Save(m.storagePath, m.data); err != nil {
		m.err = err
		m.status = "Failed to save active tasks."
		return
	}
	m.err = nil
	m.status = okStatus
}

// Function to save data and archive into json
func (m *rootModel) persistDataAndArchive(okStatus string){
	normalizeData(&m.data)
	normalizeData(&m.archiveData)
	if err := storage.Save(m.storagePath, m.data); err != nil{
		m.err = err
		m.status = "Failed to save active tasks."
		return
	}
	if err := storage.Save(m.archivePath, m.archiveData); err != nil {
		m.err = err
		m.status = "Failed to save archive tasks."
		return
	}
	m.err = nil
	m.status = okStatus
}

// Function to only save archive into json. 
func (m *rootModel) persistArchiveOnly(okStatus string) {
	normalizeData(&m.archiveData)
	if err := storage.Save(m.archivePath, m.archiveData); err != nil {
		m.err = err
		m.status = "Failed to save archive tasks."
		return
	}
	m.err = nil
	m.status = okStatus
}

/*
func (m *rootModel) setTaskState(folderID, taskID, state string) bool {
	fi, ti := m.findTaskIndex(folderID, taskID)
	if fi < 0 || ti < 0{
		return false
	}
	t := &m.data.Folders[fi].Tasks[ti]
	t.State = state
	t.Completed = state == taskStateDone
	return true
}

func (m *rootModel) deleteTask(folderID, taskID string) bool{
	fi, ti := m.findTaskIndex(folderID, taskID)
	if fi < 0 || ti < 0{
		return false
	}
	tasks := m.data.Folders[fi].Tasks
	m.data.Folders[fi].Tasks = append(tasks[:ti], tasks[ti+1]...)
	return true
}

func (m *rootModel) deleteArchivedTask(folderID, taskID string) bool {
	fi, ti := m.findArchiveTaskIndex(folderID, taskId)
	if fi < 0 || ti < 0{
		return false
	}
	tasks := m.archiveData.Folders[fi].Tasks
	m.archiveData.Folders[fi].Tasks = append(tasks[:ti], tasks[ti+1:]...)
	return true
}
*/
//from line 70

// Find folder index by its ID
func (m rootModel) folderIndexByID(id string) int {
	for i, f := range m.data.Folders {
		if f.ID == id {
			return i
		}
	}
	return -1
}

// Find task index by its ID
func taskIndexByID(tasks []Task, id string) int {
	for i, t := range tasks {
		if t.ID == id {
			return i
		}
	}
	return -1
}

// Forces cursor to not go out of range
func clampCursor(cursor, n int) int {
	if n <= 0 {
		return 0
	}
	if cursor < 0 {
		return 0
	}
	if cursor >= n {
		return n - 1
	}
	return cursor
}
