package app

import (
	"errors"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/mattn/go-runewidth"
)

type folderFormVars struct {
	confirm bool
	title string
	description string
	color string
}

type folderFormModel struct{
	form *huh.Form
	vars *folderFormVars
	edit bool
	folderID string
}

func newFolderFormModel(folder Folder, edit bool) folderFormModel{
	color := strings.TrimSpace (folder.Color)
	if color == "" {
		color = defaultFolderColor
	}

	v := &folderFormVars{
		confirm: true,
		title: folder.Name,
		description: folder.Description,
		color: color,
	}

	confirmQ := "Create folder?"
	if edit {
		confirmQ = "Save folder changes?"
	}

	fm := folderFormModel{
		vars: v,
		edit: edit,
		folderID: folder.ID,
	}

	fm.form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Color").
				Options(huh.NewOptions("green", "orange", "purple", "blue", "red", "yellow")...).
				Value(&fm.vars.color),
		
			huh.NewInput().
				Title("Folder Name").
				Value(&fm.vars.title).
				Validate(func(s string) error {
					if len(strings.TrimSpace(s)) == 0{
						return errors.New("name cannot be empty")
					}
					if runewidth.StringWidth(s) > 32 {
						return errors.New("name is too long (Max 32 characters)")
					}
					return nil
				}),

			huh.NewText().
				Title("Description -- Optional").
				Value(&fm.vars.description).

			huh.NewConfirm().
				Title(confirmQ).
				Affermative("Yes").
				Negative("No").
				Value(&fm.vars.confirm),
		),
	). WithWidth(70)
	return fm
}

func (m folderFormModel) Init() tea.Cmd{
	return m.form.Init()
}

func (m folderFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd){
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}
	return m, cmd
}

func (m folderFormModel) View() string{
	var b strings.Builder

	b.WriteString(m.form.View())
	return b.String()
}

func (m rootModel) updateFolderForm(msg tea.Msg) (tea.Model, tea.Cmd){
	updated, cmd := m.folderForm.Update(msg)
	if fm, ok := updated.(folderFormModel); ok {
		m.folderForm = fm
	}

	if !m.folderForm.isDone(){
		return m, cmd
	}

	m.screen = screenFolderList
	if m.folderForm.isCancelled(){
		if m.folderForm.edit {
			m.status = "Edit folder cancelled."
		} else {
			m.status = "Add folder cancelled."
		}
		return m, nil
	}

	result := m.folderForm.result()
	if m.folderForm.edit{
		fi := m.folderIndexByID(result.ID)
		if fi < 0 {
			m.status = "Folder no longer exists."
			return m, nil
		}
		result.Tasks = m.data.Folders[fi].Tasks
		m.data.Folders[fi] = result
		m.persistData(fmt.Sprintf("Updated folder %s", result.Name))
		m.refreshFolderList(result.ID)
		m.refreshUnifiedTaskList("", "")
		return m, nil
		}

		result.ID = newID("folder")
		if result.Color == "" {
			result.Color = defaultFolderColor
		}
		result.Tasks = []Task{}
		m.data.Folders = append(m.data.Folders, result)
		m.persistData(fmt.Sprintf("Added folder %s", result.Name))
		m.refreshFolderList(result.ID)
		m.refreshUnifiedTaskList("", "")
		return m, nil
}

func (m folderFormModel) isDone() bool{
	return m.form != nil && m.form.state == huh.StateCompleted || m.form.state == huh.StateAborted}

func (m folderFormModel) isCancelled() bool {
	if m.form == nil {
		return true
	}
	if m.form.State == huh.StateAborted {
		return true
	}
	return !m.vars.confirm
}

func (m folderFormModel) result() Folder{
	return Folder{
		ID: m.folderID,
		Name: strings.TrimSpace(m.vars.title),
		Description: strings.TrimSpace(m.vars.description),
		Color: strings.TrimSpace(m.vars.color),
	}
}