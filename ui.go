package main

import (
	"github.com/rivo/tview"
)

var app *tview.Application
var list1 *tview.List


func getmodels() []string {
	fp := "models/stable-diffusion/"
	
}

func BuildUI() {
app = tview.NewApplication()

list1 := tview.NewList().
	AddItem("All Models","",'a',nil).
	AddItem("First Pass","",'f',nil).


}