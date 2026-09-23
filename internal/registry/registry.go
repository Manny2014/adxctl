package registry

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	registryDir  = ".adxctl"
	registryFile = "registry.json"
)

// Process represents a running poller or worker.
type Process struct {
	ID      string `json:"id"`
	Type    string `json:"type"` // "poller" or "worker"
	PID     int    `json:"pid"`
	Subject string `json:"subject,omitempty"`
	Runtime string `json:"runtime,omitempty"`
	Status  string `json:"status"`
}

// Registry manages the list of running processes.
type Registry struct {
	path string
}

// New creates a new registry.
func New() (*Registry, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	registryPath := filepath.Join(home, registryDir, registryFile)
	if err := os.MkdirAll(filepath.Dir(registryPath), 0755); err != nil {
		return nil, err
	}
	return &Registry{path: registryPath}, nil
}

// Add adds a process to the registry.
func (r *Registry) Add(p Process) error {
	processes, err := r.List()
	if err != nil {
		return err
	}
	processes = append(processes, p)
	return r.save(processes)
}

// Remove removes a process from the registry.
func (r *Registry) Remove(id string) error {
	processes, err := r.List()
	if err != nil {
		return err
	}
	var newProcesses []Process
	for _, p := range processes {
		if p.ID != id {
			newProcesses = append(newProcesses, p)
		}
	}
	return r.save(newProcesses)
}

// List lists all processes in the registry.
func (r *Registry) List() ([]Process, error) {
	data, err := os.ReadFile(r.path)
	if err != nil {
		if os.IsNotExist(err) {
			return []Process{}, nil
		}
		return nil, err
	}
	var processes []Process
	if err := json.Unmarshal(data, &processes); err != nil {
		return nil, err
	}
	return processes, nil
}

func (r *Registry) save(processes []Process) error {
	data, err := json.MarshalIndent(processes, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(r.path, data, 0644)
}
