package config

type TaskResult struct {
	ID               int
	Directory        string
	StartTime        string
	EndTime          string
	TotalRuntime     string
	FilesAdded       string
	FilesDeleted     string
	MagicStringCount string
	Status           string
}

type Task struct {
	TaskID      int    `json:"task_id"`
	Directory   string `json:"directory"`
	Interval    int    `json:"interval"`
	MagicString string `json:"magic_string"`
}
