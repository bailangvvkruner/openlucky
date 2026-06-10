package modules

type State string

const (
	StateReady State = "ready"
	StateStub  State = "stub"
)

type Module struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Route       string   `json:"route"`
	Description string   `json:"description"`
	Category    string   `json:"category"`
	Phase       string   `json:"phase"`
	State       State    `json:"state"`
	Enabled     bool     `json:"enabled"`
	Stub        bool     `json:"stub"`
	Implemented bool     `json:"implemented"`
	Endpoints   []string `json:"endpoints"`
}

type Snapshot struct {
	Modules []Module `json:"modules"`
}
