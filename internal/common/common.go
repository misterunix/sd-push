package common

import "os/exec"

type Startup struct {
	RandomNumber   int
	Prompt         string
	NPrompt        string
	ModelsLocation string
	Model          string
	LoraLocation   string
	Steps          int
	Width          int
	Height         int
	ScaleUp        bool
	RemoveSmall    bool
	Sampler        string
	Count          int
}

type Stable struct {
	RandomNumber int
	SmallImage   string
	LargeImage   string
	Prompt       string
	NPrompt      string
	Model        string
	Steps        int
	Width        int
	Height       int
	Sampler      string
}

var Samplers []string   // List of samplers
var BaseModels []string // List of models in the directory

// var timedir string
var Tss string
var Cmd *exec.Cmd
var Width, Height int
var UserDir string
