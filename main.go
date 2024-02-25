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
}

func firstpass(prompt, nprompt, model string, r int) {
	sd := Stable{}
	sd.RandomNumber = r
	sd.SmallImage = timedir + "/" + tss + "-" + "small.jpg"
	sd.LargeImage = timedir + "/" + tss + "-" + "large.jpg"
	sd.Prompt = prompt
	sd.NPrompt = nprompt
	sd.Model = model

	jsonString, _ := json.Marshal(sd)
	os.WriteFile(timedir+"/"+tss+".json", jsonString, os.ModePerm)

	// Create a new template and parse the letter into it.
	passOne := "t.tpl"
	// tmpl, err := template.New(passOne).ParseFiles(passOne)
	tmpl, err := template.New(passOne).Parse(t)
	if err != nil {
		panic(err)
	}

	small, err := os.OpenFile("mm.py", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer small.Close()

	err = tmpl.Execute(small, sd)
	if err != nil {
		fmt.Fprintln(os.Stderr, "s1 Execute:", err)
		return
	}

	cmd = exec.Command("./installer_files/env/bin/python", "mm.py")
	err = cmd.Run()
	if err != nil {
		os.Exit(1)
		fmt.Fprintln(os.Stderr, "s1 Start:", err)
		return
	}
	// err = cmd.Wait()
	// if err != nil {
	// 	fmt.Fprintln(os.Stderr, "s1 Wait:", err)
	// 	return
	// }

}

func secondpass(prompt, nprompt, model string, r int) {
	sd := Stable{}
	sd.RandomNumber = r
	sd.SmallImage = timedir + "/" + tss + "-" + "small.jpg"
	sd.LargeImage = timedir + "/" + tss + "-" + "large.jpg"
	sd.Prompt = prompt
	sd.NPrompt = nprompt
	sd.Model = model

	// Create a new template and parse the letter into it.
	passOne := "t1.tpl"
	//tmpl, err := template.New(passOne).ParseFiles(passOne)
	tmpl, err := template.New(passOne).Parse(t1)
	if err != nil {
		panic(err)
	}

	small, err := os.OpenFile("mn.py", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer small.Close()

	err = tmpl.Execute(small, sd)
	if err != nil {
		fmt.Fprintln(os.Stderr, "s2 Execute:", err)
		return
	}

	cmd = exec.Command("./installer_files/env/bin/python", "mn.py")
	err = cmd.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "s2 Start:", err)
		return
	}
	// err = cmd.Wait()
	// if err != nil {
	// 	fmt.Fprintln(os.Stderr, "s2 Wait:", err)
	// 	return
	// }

}

func main() {

	var prompt string
	var nprompt string
	//var model string
	var r int

	flag.StringVar(&prompt, "prompt", "prompt", "prompt")
	flag.StringVar(&nprompt, "nprompt", "nprompt", "nprompt")
	//flag.StringVar(&model, "model", "model", "model")
	//flag.IntVar(&r, "r", 0, "random number")
	flag.Parse()

	timedir = "/mnt/nfs_clientshare/stable/" + time.Now().Format("2006-01-02-15-04-05")
	err := os.MkdirAll(timedir, 0777)
	if err != nil {
		panic(err)
	}

	m, err := os.ReadFile("models.txt")
	if err != nil {
		panic(err)
	}
	mo := string(m)
	models := strings.Split(mo, "\n")
	r = time.Now().Nanosecond()

	for index, model := range models {
		strings.TrimSpace(model)
		if strings.HasPrefix(model, "#") {
			continue
		}
		tss = fmt.Sprintf("%d", time.Now().Unix())

		//tss = time.Now().Format("2006-01-02-15-04-05")
		fmt.Println("index:", index, "tss:", tss, "model:", model)

		firstpass(prompt, nprompt, model, r)

		// err = cmd.Wait()
		// if err != nil {
		// 	fmt.Fprintln(os.Stderr, "Wait:", err)
		// 	return
		// }

		secondpass(prompt, nprompt, model, r)

		// err = cmd.Wait()
		// if err != nil {
		// 	fmt.Fprintln(os.Stderr, "Wait:", err)
		// 	return
		// }
	}
}
