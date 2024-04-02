package cmd

import (
	"fmt"
	"log"
	"os"
	
)

const INIT_ACTION = "init"

type (
	Init struct {
		name string
	}
)

func NewInit() *Init {
	return &Init{name: INIT_ACTION}
}

func (a *Init) Name() string {
	return a.name
}

func (a *Init) Exec() error {
	// Crear los directorios necesarios para un nuevo repositorio Git
	dirs := []string{".git", ".git/objects", ".git/refs"}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("creating directory %s: %w", dir, err)
		}
	}

	// Crear el archivo .git/HEAD
	headContents := []byte("ref: refs/heads/main\n")
	if err := os.WriteFile(".git/HEAD", headContents, 0644); err != nil {
		return fmt.Errorf("writing .git/HEAD: %w", err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %s", err.Error())
	}
	log.Println("Initialized empty Git repository in", cwd)
	return nil
}

func (a *Init) Help() string {
	return "init: Initialize a new, empty Git repository."
}
