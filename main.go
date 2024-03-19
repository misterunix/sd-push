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

//go:embed t.tpl
var t string

//go:embed t1.tpl
var t1 string

//go:embed models.txt
var basemodels string

var timedir string
var tss string
var cmd *exec.Cmd

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
}

var width, height int

func firstpass(prompt, nprompt, model string, r int, thesteps int) {
	os.Remove("mm.py")
	os.Remove("mn.py")
	fmt.Println("firstpass")
	sd := Stable{}
	sd.RandomNumber = r
	sd.SmallImage = timedir + "/" + tss + "-" + "small.jpg"
	sd.LargeImage = timedir + "/" + tss + "-" + "large.jpg"
	sd.Prompt = prompt
	sd.NPrompt = nprompt
	sd.Model = model
	sd.Steps = thesteps
	sd.Width = width
	sd.Height = height

	jsonString, _ := json.Marshal(sd)
	os.WriteFile(timedir+"/"+tss+".json", jsonString, os.ModePerm)

	// Create a new template and parse the letter into it.
	passOne := "t.tpl"
	// tmpl, err := template.New(passOne).ParseFiles(passOne)
	tmpl, err := template.New(passOne).Parse(t)
	CheckFatal(err)

	small, err := os.OpenFile("mm.py", os.O_CREATE|os.O_WRONLY, 0644)
	CheckFatal(err)

	defer small.Close()

	err = tmpl.Execute(small, sd)
	if err != nil {
		fmt.Fprintln(os.Stderr, "s1 Execute:", err)
		return
	}

	cmd = exec.Command("./installer_files/env/bin/python", "mm.py")
	err = cmd.Run()
	CheckFatal(err)

	// err = cmd.Wait()
	// if err != nil {
	// 	fmt.Fprintln(os.Stderr, "s1 Wait:", err)
	// 	return
	// }

}

func secondpass(prompt, nprompt, model string, r int, thesteps int) {
	fmt.Println("secondpass")
	sd := Stable{}
	sd.RandomNumber = r
	sd.SmallImage = timedir + "/" + tss + "-" + "small.jpg"
	sd.LargeImage = timedir + "/" + tss + "-" + "large.jpg"
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
	CheckFatal(err)

	small, err := os.OpenFile("mn.py", os.O_CREATE|os.O_WRONLY, 0644)
	CheckFatal(err)
	defer small.Close()

	err = tmpl.Execute(small, sd)
	if err != nil {
		fmt.Fprintln(os.Stderr, "s2 Execute:", err)
		return
	}

	cmd = exec.Command("./installer_files/env/bin/python", "mn.py")
	err = cmd.Run()
	CheckFatal(err)
	// err = cmd.Wait()
	// if err != nil {
	// 	fmt.Fprintln(os.Stderr, "s2 Wait:", err)
	// 	return
	// }
	os.Remove(sd.SmallImage)

}

func runAllModels(prompt, nprompt string, theseed int, thesteps int) {

	m, err := os.ReadFile("models.txt")
	CheckFatal(err)
	mo := string(m)
	models := strings.Split(mo, "\n")
	var r int
	if theseed != 0 {
		r = theseed
	} else {
		r = time.Now().Nanosecond()
	}
	for index, model := range models {

		model = strings.TrimSpace(model)
		if strings.HasPrefix(model, "#") {
			continue
		}
		tss = fmt.Sprintf("%d", time.Now().Unix())

		fmt.Println("index:", index, "tss:", tss, "model:", model)

		firstpass(prompt, nprompt, model, r, thesteps)

		secondpass(prompt, nprompt, model, r, thesteps)

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

	flag.StringVar(&prompt, "prompt", "prompt", "prompt")
	flag.StringVar(&nprompt, "nprompt", "nprompt", "nprompt")
	flag.StringVar(&modelcli, "model", "", "model")
	flag.IntVar(&theseed, "seed", 0, "seed")
	flag.IntVar(&thesteps, "steps", 16, "steps")
	flag.IntVar(&width, "width", 512, "width")
	flag.IntVar(&height, "height", 768, "height")
	flag.IntVar(&count, "count", 1, "count")

	//flag.IntVar(&r, "r", 0, "random number")
	flag.Parse()

	os.Remove("mm.py")
	os.Remove("mn.py")

	timedir = "/mnt/nfs_clientshare/stable/" + time.Now().Format("2006-01-02-15-04-05")
	err := os.MkdirAll(timedir, 0777)
	CheckFatal(err)

	if modelcli == "" {
		for i := 0; i < count; i++ {
			runAllModels(prompt, nprompt, theseed, thesteps)
		}
		os.Exit(0)
	}

	//m, err := os.ReadFile("models.txt")
	//CheckFatal(err)

	//mo := string(m)
	models := strings.Split(basemodels, "\n")

	found := false
	for _, model := range models {
		model = strings.TrimSpace(model)
		fmt.Println("model:", model, "modelcli:", modelcli)
		if strings.HasPrefix(model, "#") {
			continue
		}
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
