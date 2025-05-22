package main

import (
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"time"
)

//go:embed templates/t.tpl
var t string

//go:embed templates/t1.tpl
var t1 string

//go:embed samplers.txt
var samplersfile string

func firstpass(prompt, nprompt, model string, r int, thesteps int) error {
	os.Remove("mm.py")
	os.Remove("mn.py")
	fmt.Println("firstpass")
	sd := Stable{}
	sd.RandomNumber = r
	sd.SmallImage = userdir + "/" + tss + "-" + "small.jpg"
	sd.LargeImage = userdir + "/" + tss + "-" + "large.jpg"
	sd.Prompt = prompt
	sd.NPrompt = nprompt
	sd.Model = model
	sd.Steps = thesteps
	sd.Width = width
	sd.Height = height
	sd.Sampler = "dpmpp_2s_a"
	jsonString, _ := json.MarshalIndent(sd, " ", "  ")

	os.WriteFile(userdir+"/"+tss+".json", jsonString, os.ModePerm)

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

	cmd = exec.Command("./installer_files/env/bin/python", "mm.py")
	err = cmd.Run()
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
	sd := Stable{}
	sd.RandomNumber = r
	sd.SmallImage = userdir + "/" + tss + "-" + "small.jpg"
	sd.LargeImage = userdir + "/" + tss + "-" + "large.jpg"
	sd.Prompt = prompt
	sd.NPrompt = nprompt
	sd.Model = model
	sd.Steps = thesteps
	sd.Width = width
	sd.Height = height

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

	cmd = exec.Command("./installer_files/env/bin/python", "mn.py")
	err = cmd.Run()
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
	fmt.Println("LoadModels")
	// get a list of models from a directory
	dir, err := os.ReadDir("models/stable-diffusion")
	if err != nil {
		return err
	}

	//CheckFatal(err)
	for _, file := range dir {
		if file.IsDir() {
			continue
		}
		if strings.HasSuffix(file.Name(), ".txt") {
			continue
		}
		fmt.Println("file.Name(): ", file.Name())
		basemodels = append(basemodels, file.Name())
	}
	//fmt.Println("basemodels:", len(basemodels))
	if len(basemodels) == 0 {
		return fmt.Errorf("no models found")
	}
	return nil
}

func runAllModels(prompt, nprompt string, theseed int, thesteps int) {

	// m, err := os.ReadFile("models.txt")
	// CheckFatal(err)
	// mo := string(m)
	// models := strings.Split(mo, "\n")

	// LoadModels()

	var r int
	if theseed != 0 {
		r = theseed
	} else {
		r = time.Now().Nanosecond()
	}
	for index, model := range basemodels {

		model = strings.TrimSpace(model)
		// if strings.HasPrefix(model, "#") {
		// 	continue
		// }
		tss = fmt.Sprintf("%d", time.Now().Unix())

		fmt.Println("index:", index, "tss:", tss, "model:", model)

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

func main() {

	var prompt string
	var nprompt string
	var modelcli string
	var r int
	var count int
	var theseed int
	var thesteps int
	var createstartjson bool
	var samplers []string

	flag.StringVar(&prompt, "prompt", "prompt", "prompt")
	flag.StringVar(&nprompt, "nprompt", "nprompt", "nprompt")
	flag.StringVar(&modelcli, "model", "", "model")
	flag.IntVar(&theseed, "seed", 0, "seed")
	flag.IntVar(&thesteps, "steps", 16, "steps")
	flag.IntVar(&width, "width", 512, "width")
	flag.IntVar(&height, "height", 768, "height")
	flag.IntVar(&count, "count", 1, "count")
	flag.BoolVar(&createstartjson, "cj", false, "create start json")
	//flag.IntVar(&r, "r", 0, "random number")
	flag.Parse()

	if createstartjson {
		sd := Startup{}
		sd.RandomNumber = 0
		sd.Prompt = prompt
		sd.NPrompt = nprompt
		sd.ModelsLocation = "models/stable-diffusion"
		sd.Model = ""
		sd.LoraLocation = "models/lora"
		sd.Steps = 16
		sd.Width = 768
		sd.Height = 512
		sd.ScaleUp = true
		sd.RemoveSmall = true
		sd.Sampler = ""
		jsonString, _ := json.MarshalIndent(sd, " ", "  ")
		err := os.WriteFile("start.json", jsonString, os.ModePerm)
		if err != nil {
			fmt.Fprintln(os.Stderr, "WriteFile:", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	// check if start.json exists
	_, err := os.Stat("start.json")
	if err != nil {
		fmt.Fprintln(os.Stderr, "start.json does not exist")
		os.Exit(1)
	}

	params := Startup{}
	data, err := os.ReadFile("start.json")
	if err != nil {
		fmt.Fprintln(os.Stderr, "ReadFile:", err)
		os.Exit(1)
	}
	err = json.Unmarshal(data, &params)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unmarshal:", err)
		os.Exit(1)
	}

	os.Remove("mm.py")
	os.Remove("mn.py")

	// check if the models directory exists

	//var userdir string
	userdir, err = os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "UserHomeDir:", err)
		os.Exit(1)
	}
	userdir += "/webserver/"

	err = os.MkdirAll(userdir, 0755)
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

	fmt.Println("Current folder:", userdir)

	err = LoadModels()
	if err != nil {
		fmt.Println("LoadModels:", err)
		os.Exit(1)
	}

	s := strings.Split(samplersfile, "\n")
	for _, v := range s {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		samplers = append(samplers, v)
	}

	if modelcli == "" {
		fmt.Println("Running all models")
		for i := 0; i < count; i++ {
			runAllModels(prompt, nprompt, theseed, thesteps)
		}
		os.Exit(0)
	}

	//m, err := os.ReadFile("models.txt")
	//CheckFatal(err)

	//mo := string(m)
	//models := strings.Split(basemodels, "\n")

	found := false
	for _, model := range basemodels {
		model = strings.TrimSpace(model)
		fmt.Println("model: ", model, " modelcli: ", modelcli)
		// if strings.HasPrefix(model, "#") {
		// 	continue
		// }
		if model == modelcli {
			found = true
			break
		}
	}

	if !found {
		fmt.Println("model not found")
		os.Exit(1)
	}

	totalstart := time.Now()
	for i := 0; i < count; i++ {
		if theseed != 0 {
			r = theseed
		} else {
			r = time.Now().Nanosecond()
		}
		//r = time.Now().Nanosecond()
		tss = fmt.Sprintf("%d", time.Now().Unix())
		fmt.Println("i:", i, "tss:", tss, "modelcli:", modelcli)
		start := time.Now()
		firstpass(prompt, nprompt, modelcli, r, thesteps)
		secondpass(prompt, nprompt, modelcli, r, thesteps)
		fmt.Println("time:", time.Since(start).Minutes())
	}
	fmt.Println("total time:", time.Since(totalstart).Minutes())
}
