package main

import (
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand/v2"
	"os"
	"os/exec"
	"path"
	"sd-push/internal/common"
	"strings"
	"text/template"
	"time"
)

//go:embed templates/t.tpl
var t string

//go:embed templates/t1.tpl
var t1 string

//go:embed templates/samplers.txt
var SamplersFile string

func firstpass(prompt, nprompt, model string, r int, thesteps int) error {
	os.Remove("mm.py")
	os.Remove("mn.py")
	fmt.Println("firstpass")
	sd := common.Stable{}
	sd.RandomNumber = r
	sd.SmallImage = common.UserDir + "/" + common.Tss + "-" + "small.png"
	sd.LargeImage = common.UserDir + "/" + common.Tss + "-" + "large.png"
	sd.Prompt = prompt
	sd.NPrompt = nprompt
	sd.Model = model
	sd.Steps = thesteps
	sd.Width = common.Width
	sd.Height = common.Height
	sd.Sampler = "dpmpp_2s_a"
	jsonString, _ := json.MarshalIndent(sd, " ", "  ")

	os.WriteFile(common.UserDir+"/"+common.Tss+".json", jsonString, os.ModePerm)

	// Create a new template and parse the letter into it.
	passOne := "t.tpl"
	// tmpl, err := template.New(passOne).ParseFiles(passOne)
	tmpl, err := template.New(passOne).Parse(t)
	if err != nil {
		fmt.Fprintln(os.Stderr, "s1 Parse:", err)
		return err
	}
	//CheckFatal(err)

	small, err := os.OpenFile("mm.py", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, "s1 OpenFile:", err)
		return err
	}

	//CheckFatal(err)

	defer small.Close()

	err = tmpl.Execute(small, sd)
	if err != nil {
		fmt.Fprintln(os.Stderr, "s1 Execute:", err)
		return err
	}

	common.Cmd = exec.Command("./installer_files/env/bin/python", "mm.py")
	err = common.Cmd.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "s1 Run:", err)
		return err
	}

	//CheckFatal(err)

	// err = cmd.Wait()
	// if err != nil {
	// 	fmt.Fprintln(os.Stderr, "s1 Wait:", err)
	// 	return
	// }
	return nil
}

func secondpass(prompt, nprompt, model string, r int, thesteps int) error {
	fmt.Println("secondpass")
	sd := common.Stable{}
	sd.RandomNumber = r
	sd.SmallImage = common.UserDir + "/" + common.Tss + "-" + "small.png"
	sd.LargeImage = common.UserDir + "/" + common.Tss + "-" + "large.png"
	sd.Prompt = prompt
	sd.NPrompt = nprompt
	sd.Model = model
	sd.Steps = thesteps
	sd.Width = common.Width
	sd.Height = common.Height

	// Create a new template and parse the letter into it.
	passOne := "t1.tpl"
	//tmpl, err := template.New(passOne).ParseFiles(passOne)
	tmpl, err := template.New(passOne).Parse(t1)
	if err != nil {
		fmt.Fprintln(os.Stderr, "s2 Parse:", err)
		return err
	}

	small, err := os.OpenFile("mn.py", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, "s2 OpenFile:", err)
		return err
	}

	defer small.Close()

	err = tmpl.Execute(small, sd)
	if err != nil {
		fmt.Fprintln(os.Stderr, "s2 Execute:", err)
		return err
	}

	common.Cmd = exec.Command("./installer_files/env/bin/python", "mn.py")
	err = common.Cmd.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "s2 Run:", err)
		return err
	}
	// err = cmd.Wait()
	// if err != nil {
	// 	fmt.Fprintln(os.Stderr, "s2 Wait:", err)
	// 	return
	// }

	os.Remove(sd.SmallImage)

	return nil
}

// get all the models from the directory
func LoadModels() error {
	fmt.Println("Loading Models")

	// get a list of models from a directory
	dir, err := os.ReadDir("models/stable-diffusion")
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
			fmt.Println("Model Name(): ", file.Name())
			common.BaseModels = append(common.BaseModels, file.Name())
		}
	}
	//fmt.Println("basemodels:", len(basemodels))
	if len(common.BaseModels) == 0 {
		return fmt.Errorf("no models found")
	}
	return nil
}

func runAllModels(prompt, nprompt string, theseed int, thesteps int) {

	// m, err := os.ReadFile("models.txt")
	// CheckFatal(err)
	// mo := string(m)
	// models := strings.Split(mo, "\n")

	err := LoadModels()
	if err != nil {
		fmt.Println("LoadModels:", err)
		os.Exit(1)
	}

	var r int
	if theseed != 0 {
		r = theseed
	} else {
		r = time.Now().Nanosecond()
	}

	for index, model := range common.BaseModels {

		fmt.Println("index:", index, "model:", model)

		model = strings.TrimSpace(model)
		// if strings.HasPrefix(model, "#") {
		// 	continue
		// }
		common.Tss = fmt.Sprintf("%d", time.Now().Unix())

		fmt.Println("index:", index, "tss:", common.Tss, "model:", model)

		err := firstpass(prompt, nprompt, model, r, thesteps)
		if err != nil {
			fmt.Fprintln(os.Stderr, "firstpass:", err)
			continue
		}

		err = secondpass(prompt, nprompt, model, r, thesteps)
		if err != nil {
			fmt.Fprintln(os.Stderr, "secondpass:", err)
			continue
		}

	}
}

// Main function
func main() {

	//getmodels()
	//os.Exit(0)
	var prompt string
	var nprompt string
	//var modelcli string
	var r int
	//var count int
	var theseed int
	//var thesteps int
	var createstartjson bool

	// flag.StringVar(&prompt, "prompt", "prompt", "prompt")
	// flag.StringVar(&nprompt, "nprompt", "nprompt", "nprompt")
	// flag.StringVar(&modelcli, "model", "", "model")
	// flag.IntVar(&theseed, "seed", 0, "seed")
	// flag.IntVar(&thesteps, "steps", 16, "steps")
	// flag.IntVar(&width, "width", 512, "width")
	// flag.IntVar(&height, "height", 512, "height")
	// flag.IntVar(&count, "count", 1, "count")
	flag.BoolVar(&createstartjson, "cj", false, "create sd-push-run.json")
	//flag.IntVar(&r, "r", 0, "random number")
	flag.Parse()

	if createstartjson {
		sd := common.Startup{}
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

	// check if sd-push-run exists
	_, err := os.Stat("sd-push-run.json")
	if err != nil {
		fmt.Fprintln(os.Stderr, "sd-push-run.json does not exist")
		os.Exit(1)
	}

	common.Width = 512
	common.Height = 512

	params := common.Startup{}
	data, err := os.ReadFile("sd-push-run.json")
	if err != nil {
		fmt.Fprintln(os.Stderr, "ReadFile:", err)
		os.Exit(1)
	}
	err = json.Unmarshal(data, &params)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unmarshal:", err)
		os.Exit(1)
	}

	fmt.Println(params)

	os.Remove("mm.py")
	os.Remove("mn.py")

	// check if the models directory exists

	//var userdir string
	common.UserDir, err = os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "UserHomeDir:", err)
		os.Exit(1)
	}
	webserverDir := "/images/"

	serverPath := path.Join(common.UserDir, webserverDir)

	fmt.Println(common.UserDir, webserverDir, serverPath)

	common.UserDir = serverPath // quick fix - remove later

	err = os.MkdirAll(serverPath, 0755)
	if err != nil {
		fmt.Fprintln(os.Stderr, "MkdirAll:", err)
		os.Exit(1)
	}

	//timedir = "/mnt/nfs_clientshare/stable/" + time.Now().Format("2006-01-02-15-04-05")
	// err = os.MkdirAll(timedir, 0777)
	// if err != nil {
	// 	fmt.Fprintln(os.Stderr, "MkdirAll:", err)
	// 	os.Exit(1)
	// }

	fmt.Println("Current folder:", common.UserDir)

	err = LoadModels()
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

	params.Prompt = strings.ReplaceAll(params.Prompt, "\n", " ")
	params.NPrompt = strings.ReplaceAll(params.NPrompt, "\n", " ")

	params.Prompt = strings.ReplaceAll(params.Prompt, "'", "\\'")
	params.NPrompt = strings.ReplaceAll(params.NPrompt, "'", "\\'")

	params.Prompt = strings.ReplaceAll(params.Prompt, "\"", "\\\"")
	params.NPrompt = strings.ReplaceAll(params.NPrompt, "\"", "\\\"")

	if params.Sampler == "" {
		params.Sampler = common.Samplers[rand.IntN(len(common.Samplers))]
	}

	//for _, common.RunThisSampler = range common.Samplers {

	for _, modeltorun := range common.BaseModels {

		fmt.Println("modeltorun:", modeltorun)

		totalstart := time.Now()

		for i := 0; i < params.Count; i++ {

			fmt.Println("count:", i)

			if theseed != 0 {
				r = theseed
			} else {
				r = time.Now().Nanosecond()
			}

			//r = time.Now().Nanosecond()
			common.Tss = fmt.Sprintf("%d", time.Now().Unix())
			fmt.Println("i:", i, "tss:", common.Tss, "modelcli:", modeltorun)
			start := time.Now()
			firstpass(params.Prompt, params.NPrompt, modeltorun, r, params.Steps)
			secondpass(params.NPrompt, params.NPrompt, modeltorun, r, params.Steps)
			fmt.Println("time:", time.Since(start).Minutes())
		}
		fmt.Println("total time:", time.Since(totalstart).Minutes())
	}

}
