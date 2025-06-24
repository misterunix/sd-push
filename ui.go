package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/rivo/tview"
)

var app *tview.Application
var list1 *tview.List
var models []string

func getmodels() []string {
	fp := "models/stable-diffusion/"



	
entries, err := os.ReadDir(fp)
	if err != nil {
		log.Fatalf("Error reading directory: %v", err)
	}

	fmt.Printf("Files and directories in %s:\n", dirPath)
	for _, entry := range entries {
		fmt.Println
		if strings.HasSuffix(entry.Name(), ".safetensors") {
			models = append(models, entry.Name())
		}
	}

}

func BuildUI() {
app = tview.NewApplication()

list1 := tview.NewList().
	AddItem("All Models","",'a',nil).
	AddItem("First Pass","",'f',nil).


}