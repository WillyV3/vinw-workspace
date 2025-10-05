package main

import (
	"os/exec"
)

type DependencyStatus struct {
	Name      string
	Available bool
	Required  bool
}

func checkDependencies(terminal, agent string) []DependencyStatus {
	deps := []DependencyStatus{
		{Name: "vinw", Available: commandExists("vinw"), Required: true},
		{Name: "tmux", Available: commandExists("tmux"), Required: true},
	}

	if terminal == "nextui" {
		deps = append(deps, DependencyStatus{
			Name:      "nextui",
			Available: commandExists("nextui"),
			Required:  true,
		})
	}

	if agent != "none" && agent != "" {
		deps = append(deps, DependencyStatus{
			Name:      agent,
			Available: commandExists(agent),
			Required:  true,
		})
	}

	return deps
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func allDependenciesAvailable(deps []DependencyStatus) bool {
	for _, dep := range deps {
		if dep.Required && !dep.Available {
			return false
		}
	}
	return true
}
