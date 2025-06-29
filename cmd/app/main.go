package main

import (
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand/v2"
	"os"
	"path"
	"sd-push/internal/common"
	"sd-push/internal/pass1"
	"sd-push/internal/pass2"
	"strings"
	"time"
)

//go:embed templates/t.tpl
var t string

//go:embed templates/t1.tpl
var t1 string

//go:embed templates/samplers.txt
var SamplersFile string

func NewStable() *common.Stable {
	sd := common.Stable{}
	sd.RandomNumber = 0 // Random number for image generation
	sd.SmallImage = ""  // Path and Name of the Small image
	sd.LargeImage = ""  // Path and Name of the Large image
	sd.JsonDesc = ""    // Path and Name of the JSON description file
	sd.Prompt = ""      // Prompt for image generation
	sd.NPrompt = ""     // Negative prompt for image generation
	sd.Model = ""       // Model name for image generation
	sd.Steps = 15       // Default steps for image generation
	sd.Width = 512      // Default width for image generation
	sd.Height = 512     // Default height for image generation
	sd.Sampler = ""     // Default sampler for image generation
	sd.Seed = 0         // Random seed for image generation
	return &sd
}

// get all the models from the directory
func LoadModels(sd *common.Stable) error {
	fmt.Println("Loading Models")

	// get a list of models from a directory
	dir, err := os.ReadDir(sd.ModelsLocation)
	if err != nil {
		return err
	}

	for _, file := range dir {
		if file.IsDir() {
			continue
		}
		if strings.HasSuffix(file.Name(), ".txt") {
			continue
		}
		if strings.HasSuffix(file.Name(), ".safetensors") {
			fmt.Println("Found model: ", file.Name())
			sd.Models = append(sd.Models, file.Name())
		}
	}
	//fmt.Println("basemodels:", len(basemodels))
	if len(sd.Models) == 0 {
		return fmt.Errorf("no models found")
	}

	fmt.Println()
	fmt.Println()

	return nil
}

func runAllModels(sd *common.Stable) {

	// m, err := os.ReadFile("models.txt")
	// CheckFatal(err)
	// mo := string(m)
	// models := strings.Split(mo, "\n")

	err := LoadModels(sd)
	if err != nil {
		fmt.Println("LoadModels:", err)
		os.Exit(1)
	}

	if sd.Seed == 0 {
		sd.Seed = time.Now().Nanosecond()
	}

	for index, model := range common.BaseModels {

		fmt.Println("index:", index, "model:", model)

		model = strings.TrimSpace(model)
		// if strings.HasPrefix(model, "#") {
		// 	continue
		// }

		common.Tss = fmt.Sprintf("%d", time.Now().Unix())

		fmt.Println("index:", index, "tss:", common.Tss, "model:", model)

		err := pass1.FirstPass(sd) // Set the current run parameter
		//err := firstpass(prompt, nprompt, model, r, thesteps)
		if err != nil {
			fmt.Fprintln(os.Stderr, "firstpass:", err)
			continue
		}

		err = pass2.SecondPass(sd)
		//err = secondpass(prompt, nprompt, model, r, thesteps)
		if err != nil {
			fmt.Fprintln(os.Stderr, "secondpass:", err)
			continue
		}

	}
}

// Main function
func main() {

	var prompt string
	var nprompt string
	var createstartjson bool

	flag.BoolVar(&createstartjson, "c", false, "create a blank sd-push-run.json")
	flag.Parse()

	if createstartjson {
		sd := NewStable()

		sd.RandomNumber = 0
		sd.Prompt = prompt
		sd.NPrompt = nprompt
		sd.ModelsLocation = "models/stable-diffusion"
		sd.Model = ""
		sd.LoraLocation = "models/lora"
		sd.Steps = 50
		sd.Width = 512
		sd.Height = 512
		sd.ScaleUp = true
		sd.RemoveSmall = true
		sd.Sampler = ""
		sd.Count = 1

		jsonString, _ := json.MarshalIndent(sd, " ", "  ")
		err := os.WriteFile("sd-push-run.json", jsonString, os.ModePerm)
		if err != nil {
			fmt.Fprintln(os.Stderr, "WriteFile:", err)
			os.Exit(1)
		}

		os.Exit(0)

	}

	// struct for the program
	sd := NewStable()

	//sd.Width = 512
	//sd.Height = 512

	//common.CurrentRun = common.Startup{}

	uh, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "UserHomeDir:", err)
		os.Exit(1)
	}
	sd.UserHome = uh
	sd.ImageHome = path.Join(uh, "images")

	// set the current working directory
	sd.CurrentHome, err = os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Getwd:", err)
		os.Exit(1)
	}

	sd.Smallpy = path.Join(sd.CurrentHome, "sd-push-small.py")
	sd.Largepy = path.Join(sd.CurrentHome, "sd-push-large.py")

	// check if sd-push-run exists
	configFile := path.Join(sd.UserHome, "sd-push-run.json")
	_, err = os.Stat(configFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, configFile+" does not exist")
		os.Exit(1)
	}
	data, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ReadFile:", err)
		os.Exit(1)
	}

	err = json.Unmarshal(data, &sd)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unmarshal:", err)
		os.Exit(1)
	}

	fmt.Println(sd)

	os.Remove(sd.Smallpy)
	os.Remove(sd.Largepy)

	// os.Remove("mm.py")
	// os.Remove("mn.py")

	// check if the models directory exists

	//var userdir string

	//webserverDir := "/images/"
	//serverPath := path.Join(common.UserDir, webserverDir)
	//fmt.Println(common.UserDir, webserverDir, serverPath)
	// common.UserDir = serverPath // quick fix - remove later
	// err = os.MkdirAll(serverPath, 0755)
	// if err != nil {
	// 	fmt.Fprintln(os.Stderr, "MkdirAll:", err)
	// 	os.Exit(1)
	// }

	//timedir = "/mnt/nfs_clientshare/stable/" + time.Now().Format("2006-01-02-15-04-05")
	// err = os.MkdirAll(timedir, 0777)
	// if err != nil {
	// 	fmt.Fprintln(os.Stderr, "MkdirAll:", err)
	// 	os.Exit(1)
	// }

	fmt.Println("Current folder:", sd.CurrentHome)

	err = LoadModels(sd)
	if err != nil {
		fmt.Println("LoadModels:", err)
		os.Exit(1)
	}

	s := strings.Split(SamplersFile, "\n")
	common.Samplers = s
	/*
		for _, v := range s {
			v = strings.TrimSpace(v)
			if v == "" {
				continue
			}
			samplers = append(samplers, v)
		}
	*/

	// if modelcli == "" {
	// 	fmt.Println("Running all models")
	// 	for i := 0; i < count; i++ {
	// 		runAllModels(prompt, nprompt, theseed, thesteps)
	// 	}
	// 	os.Exit(0)
	// }

	//m, err := os.ReadFile("models.txt")
	//CheckFatal(err)

	//mo := string(m)
	//models := strings.Split(basemodels, "\n")

	//found := false
	/*
		for _, model := range common.BaseModels {
			model = strings.TrimSpace(model)
			fmt.Println("model: ", model, " modelcli: ", modelcli)
			// if strings.HasPrefix(model, "#") {
			// 	continue
			// }

			// if model == modelcli {
			// 	found = true
			// 	break
			// }
		}
	*/

	// if !found {
	// 	fmt.Println("model not found")
	// 	os.Exit(1)
	// }

	sd.Prompt = strings.ReplaceAll(sd.Prompt, "\n", " ")
	sd.NPrompt = strings.ReplaceAll(sd.NPrompt, "\n", " ")

	sd.Prompt = strings.ReplaceAll(sd.Prompt, "\n", " ")
	sd.NPrompt = strings.ReplaceAll(sd.NPrompt, "\n", " ")

	sd.Prompt = strings.ReplaceAll(sd.Prompt, "'", "\\'")
	sd.NPrompt = strings.ReplaceAll(sd.NPrompt, "'", "\\'")

	sd.Prompt = strings.ReplaceAll(sd.Prompt, "\"", "\\\"")
	sd.NPrompt = strings.ReplaceAll(sd.NPrompt, "\"", "\\\"")

	if common.CurrentRun.Sampler == "" {
		common.CurrentRun.Sampler = common.Samplers[rand.IntN(len(common.Samplers))]
	}

	//for _, common.RunThisSampler = range common.Samplers {

	for _, modeltorun := range common.BaseModels {

		common.CurrentRun.Model = modeltorun
		fmt.Println("Running Model:", common.CurrentRun.Model)

		RunModels()

	}

}

func RunModels() {

	totalstart := time.Now()

	var r int

	for i := 0; i < common.CurrentRun.Count; i++ {

		fmt.Println("count:", i)

		if common.CurrentRun.Seed != 0 {
			r = common.CurrentRun.Seed
		} else {
			r = time.Now().Nanosecond()
		}

		//r = time.Now().Nanosecond()
		common.Tss = fmt.Sprintf("%d", time.Now().Unix())
		fmt.Println("i:", i, "tss:", common.Tss, "modelcli:", common.CurrentRun.Model)
		start := time.Now()
		firstpass(common.CurrentRun.Prompt, common.CurrentRun.NPrompt, common.CurrentRun.Model, r, common.CurrentRun.Steps)
		secondpass(common.CurrentRun.NPrompt, common.CurrentRun.NPrompt, common.CurrentRun.Model, r, common.CurrentRun.Steps)
		fmt.Println("time:", time.Since(start).Minutes())
	}
	fmt.Println("total time:", time.Since(totalstart).Minutes())
}
