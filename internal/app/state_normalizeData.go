package app

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"
)

// Function to generate new ID using prefix of folder_ or task_ with the ID using a random byte buffer or current time stamp if random byte fails. 
func newID(prefix string) string{
	buf := make([]byte, 4)
	if _, err := rand.Read(buf); err != nil {
		return prefix + "_" + time.Now().Format("20060102150405.000000000")
	}
	return prefix + "_" + hex.EncodeToString(buf)
}

// Function to normalize data before saving into json
func normalizeData(data *AppData){
	for fi := range data.Folders {
		f := &data.Folders[fi]
		if strings.TrimSpace(f.ID) == "" {
			f.ID = newID("folder")
		}
		if strings.TrimSpace(f.Color) == "" {
			f.Color = defaultFolderColor
		}
		if f.Tasks == nil {
			f.Tasks = []Task{}
		}
		for ti := range f.Tasks{
			t := &f.Tasks[ti]
			if strings.TrimSpace(t.ID) == ""{
				t.ID = newID("task")
			}
			if strings.TrimSpace(t.State) == "" {
				if t.Completed{
					t.State = taskStateDone
				} else{
					t.State = taskStateTodo
				}
			}
			switch t.State{
			case taskStateTodo, taskStateMarked, taskStateDone, taskStateArchived:
			default:
				t.State = taskStateTodo
			}
			t.Completed = t.State == taskStateDone
			if t.Priority < 1 || t.Priority > 3{
				t.Priority = defaultPriority
			}
			t.Tags = normalizeTags(t.Tags)
		}
	}
	if data.Folders == nil{
		data.Folders = []Folder{}
	}
}

// function to normalize tags before saving. 
func normalizeTags(tags []string) []string{
	if len(tags) == 0{
		return []string{}
	}
	seen := make(map[string]struct{}, len(tags))
	out := make([]string, 0, len(tags))
	for _, raw := range tags{
		tag := strings.TrimSpace(raw)
		if tag == ""{
			continue
		}
		tag = strings.ToLower(tag)
		if _, exists := seen[tag]; exists{
			continue
		}
		seen[tag] = struct{}{}
		out = append(out, tag)
	}
	return out
}