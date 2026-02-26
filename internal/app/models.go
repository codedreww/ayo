package app

type Task struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description,omitempty"`
	Priority    int      `json:"priority"`
	DueDate     string   `json:"due_date,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	State       string   `json:"state,omitempty"`
	Completed   bool     `json:"completed"`
}

type Folder struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Color       string `json:"color,omitempty"`
	Tasks       []Task `json:"tasks,omitempty"`
}

type AppData struct {
	Folders []Folder `json:"folders"`
}