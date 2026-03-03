package app

// taskAdd.go defines the huh form model for creating and editing tasks

import (
	"errors"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/mattn/go-runewidth"
)

type taskFormVars struct {
	confirm     bool
	title       string
	description string
	priority    int
	dueDate     string
	tags        string
}

type taskFormModel struct {
	form     *huh.Form
	vars     *taskFormVars
	edit     bool
	folderID string
	taskID   string
	state    string
}

func newTaskFormModel(folderID string, task Task, edit bool) taskFormModel {
	state := strings.TrimSpace(task.State)
	if state == "" {
		if task.Completed {
			state = taskStateDone
		} else {
			state = taskStateTodo
		}
	}

	priority := task.Priority
	if priority < 1 || priority > 3 {
		priority = defaultPriority
	}

	v := &taskFormVars{
		confirm:     true,
		title:       task.Title,
		description: task.Description,
		priority:    priority,
		dueDate:     task.DueDate,
		tags:        strings.Join(task.Tags, ", "),
	}

	confirmQ := "Create task?"
	if edit {
		confirmQ = "Save task changes?"
	}

	tm := taskFormModel{
		vars:     v,
		edit:     edit,
		folderID: folderID,
		taskID:   task.ID,
		state:    state,
	}

	tm.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Title").
				Value(&tm.vars.title).
				Validate(func(s string) error {
					if len(strings.TrimSpace(s)) == 0 {
						return errors.New("Title cannot be empty.")
					}
					if runewidth.StringWidth(s) > 120 {
						return errors.New("Title is too long (max 120 chars)")
					}
					return nil
				}),

			huh.NewText().
				Title("Description -- Optional").
				Value(&tm.vars.description),

			huh.NewSelect[int]().
				Title("Priority").
				Options(
					huh.NewOption("1 (Low)", 1),
					huh.NewOption("2 (Medium)", 2),
					huh.NewOption("3 (High)", 3),
				).
				Value(&tm.vars.priority),
			huh.NewInput().
				Title("Due Date  (Optional, YYYY-MM-DD)").
				Validate(validateDueDate).
				Value(&tm.vars.dueDate),
			huh.NewInput().
				Title("Tags -- Optional, comma separated").
				Value(&tm.vars.tags),
			huh.NewConfirm().
				Title(confirmQ).
				Affirmative("Yes").
				Negative("No").
				Value(&tm.vars.confirm),
		),
	).WithWidth(78)
	return tm
}

func validateDueDate(s string) error {
	const layout = "2006-01-02"
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}

	parsedDate, err := time.Parse(layout, s)
	if err != nil {
		return fmt.Errorf("Invalid due date format, use YYYY-MM-DD")
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	inputDate := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, now.Location())

	if inputDate.Before(today) {
		return fmt.Errorf("Please input a valid due date")
	}
	return nil
}

func parseTagsInput(s string) []string {
	parts := strings.Split(s, ",")
	return normalizeTags(parts)
}

func (m taskFormModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m taskFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}
	return m, cmd
}

func (m taskFormModel) View() string {
	return m.form.View()
}

func (m rootModel) updateTaskForm(msg tea.Msg) (tea.Model, tea.Cmd) {
	updated, cmd := m.taskForm.Update(msg)
	if tm, ok := updated.(taskFormModel); ok {
		m.taskForm = tm
	}

	if !m.taskForm.isDone() {
		return m, cmd
	}

	m.screen = m.taskFormReturnScreen
	if m.taskForm.isCanceled() {
		action := "Add"
		if m.taskForm.edit {
			action = "Edit"
		}
		m.status = action + " task cancelled."
		return m, nil
	}

	fi := m.folderIndexByID(m.taskForm.folderID)
	if fi < 0 {
		m.screen = screenFolderList
		m.status = "Target folder no longer exists."
		return m, nil
	}

	folder := &m.data.Folders[fi]
	result := m.taskForm.result()
	if result.Priority < 1 || result.Priority > 3 {
		result.Priority = defaultPriority
	}

	if m.taskForm.edit {
		ti := taskIndexByID(folder.Tasks, result.ID)
		if ti < 0 {
			m.status = "Task no longer exists."
			return m, nil
		}
		folder.Tasks[ti] = result
		m.persistData(fmt.Sprintf("Updated task: %s", result.Title))
		m.refreshTaskListForFolder(folder.ID, result.ID)
		m.refreshUnifiedTaskList(folder.ID, result.ID)
		return m, nil
	}

	result.ID = newID("task")
	if strings.TrimSpace(result.State) == "" {
		result.State = taskStateTodo
	}
	result.Completed = result.State == taskStateDone
	folder.Tasks = append(folder.Tasks, result)
	m.persistData(fmt.Sprintf("Added task: %s", result.Title))
	m.refreshTaskListForFolder(folder.ID, result.ID)
	m.refreshUnifiedTaskList(folder.ID, result.ID)

	return m, nil
}

func (m taskFormModel) isDone() bool {
	return m.form != nil && (m.form.State == huh.StateCompleted || m.form.State == huh.StateAborted)
}

func (m taskFormModel) isCanceled() bool {
	if m.form == nil {
		return true
	}
	if m.form.State == huh.StateAborted {
		return true
	}
	return !m.vars.confirm
}

func (m taskFormModel) result() Task {
	taskState := m.state
	completed := taskState == taskStateDone

	return Task{
		ID:          m.taskID,
		Title:       strings.TrimSpace(m.vars.title),
		Description: strings.TrimSpace(m.vars.description),
		Priority:    m.vars.priority,
		DueDate:     strings.TrimSpace(m.vars.dueDate),
		Tags:        parseTagsInput(m.vars.tags),
		State:       taskState,
		Completed:   completed,
	}
}
