package app

import (
	"ayo/internal/storage"
	"strings"
)

// Function to save data into json
func (m *rootModel) persistData(okStatus string) {
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
func (m *rootModel) persistDataAndArchive(okStatus string) {
	normalizeData(&m.data)
	normalizeData(&m.archiveData)
	if err := storage.Save(m.storagePath, m.data); err != nil {
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

func (m *rootModel) setTaskState(folderID, taskID, state string) bool {
	fi, ti := m.findTaskIndex(folderID, taskID)
	if fi < 0 || ti < 0 {
		return false
	}
	t := &m.data.Folders[fi].Tasks[ti]
	t.State = state
	t.Completed = state == taskStateDone
	return true
}

func (m *rootModel) deleteTask(folderID, taskID string) bool {
	fi, ti := m.findTaskIndex(folderID, taskID)
	if fi < 0 || ti < 0 {
		return false
	}
	tasks := m.data.Folders[fi].Tasks
	m.data.Folders[fi].Tasks = append(tasks[:ti], tasks[ti+1:]...)
	return true
}

func (m *rootModel) deleteArchivedTask(folderID, taskID string) bool {
	fi, ti := m.findArchiveTaskIndex(folderID, taskID)
	if fi < 0 || ti < 0 {
		return false
	}
	tasks := m.archiveData.Folders[fi].Tasks
	m.archiveData.Folders[fi].Tasks = append(tasks[:ti], tasks[ti+1:]...)
	return true
}

//from line 70

func (m *rootModel) restoreArchivedTask(folderID, taskID string) bool {
	afi, ati := m.findArchiveTaskIndex(folderID, taskID)
	if afi < 0 || ati < 0 {
		return false
	}

	archFolder := m.archiveData.Folders[afi]
	task := m.archiveData.Folders[afi].Tasks[ati]
	task.State = taskStateTodo
	task.Completed = false

	tasks := m.archiveData.Folders[afi].Tasks
	m.archiveData.Folders[afi].Tasks = append(tasks[:ati], tasks[ati+1:]...)

	fi := m.folderIndexByID(folderID)
	if fi < 0 {
		m.data.Folders = append(m.data.Folders, Folder{
			ID:          archFolder.ID,
			Name:        archFolder.Name,
			Description: archFolder.Description,
			Color:       archFolder.Color,
			Tasks:       []Task{},
		})
		fi = len(m.data.Folders) - 1
	}
	m.data.Folders[fi].Tasks = append(m.data.Folders[fi].Tasks, task)
	return true
}

func (m rootModel) findTaskIndex(folderID, taskID string) (int, int) {
	fi := m.folderIndexByID(folderID)
	if fi < 0 {
		return -1, -1
	}
	for ti, t := range m.data.Folders[fi].Tasks {
		if t.ID == taskID {
			return fi, ti
		}
	}
	return -1, -1
}

func (m rootModel) findArchiveTaskIndex(folderID, taskID string) (int, int) {
	for fi := range m.archiveData.Folders {
		if m.archiveData.Folders[fi].ID != folderID {
			continue
		}
		for ti, t := range m.archiveData.Folders[fi].Tasks {
			if t.ID == taskID {
				return fi, ti
			}
		}
	}
	return -1, -1
}

func (m *rootModel) archiveCompletedTasks(folderID string) int {
	moved := 0
	for fi := range m.data.Folders {
		folder := &m.data.Folders[fi]
		if strings.TrimSpace(folderID) != "" && folder.ID != folderID {
			continue
		}

		remaining := make([]Task, 0, len(folder.Tasks))
		for _, t := range folder.Tasks {
			if t.State == taskStateDone || t.Completed {
				archived := t
				archived.State = taskStateArchived
				archived.Completed = true
				m.appendArchiveTask(*folder, archived)
				moved++
				continue
			}
			remaining = append(remaining, t)
		}
		folder.Tasks = remaining
	}
	return moved
}

func (m *rootModel) appendArchiveTask(folder Folder, task Task) {
	idx := -1
	for i, f := range m.archiveData.Folders {
		if f.ID == folder.ID {
			idx = i
			break
		}
	}
	if idx < 0 {
		m.archiveData.Folders = append(m.archiveData.Folders, Folder{
			ID:          folder.ID,
			Name:        folder.Name,
			Description: folder.Description,
			Color:       folder.Color,
			Tasks:       []Task{},
		})
		idx = len(m.archiveData.Folders) - 1
	}
	m.archiveData.Folders[idx].Tasks = append(m.archiveData.Folders[idx].Tasks, task)
}

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
