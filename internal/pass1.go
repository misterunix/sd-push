package pass1

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path"
	"sd-push/internal/common"
	"strconv"
)

func FirstPass(sd *common.Stable) error {
	os.Remove("mm.py")
	os.Remove("mn.py")
	fmt.Println("firstpass")

	common.CurrentRun.RandomNumber = r

	//sd.Seed = time.Now().Nanosecond() // Use the current nanosecond as the seed

	tmpS := strconv.Itoa(sd.Seed)

	sd.SmallImage = path.Join(common.UserDir, common.ImageDirectory, tmpS+"-small.png")
	sd.LargeImage = path.Join(common.UserDir, common.ImageDirectory, tmpS+"-large.png")

	common.CurrentRun.SmallImage = path.Join(common.UserDir, common.ImageDirectory, tmpS+"-small.png")
	common.CurrentRun.LargeImage = path.Join(common.UserDir, common.ImageDirectory, tmpS+"-large.png")

	// sd.SmallImage = common.UserDir + "/" + common.ImageDirectory + "-" + "small.png"
	// sd.LargeImage = common.UserDir + "/" + common.ImageDirectory + "-" + "large.png"
	sd.Prompt = prompt
	sd.NPrompt = nprompt
	sd.Model = model
	sd.Steps = thesteps
	sd.Width = common.Width
	sd.Height = common.Height
	sd.Sampler = "dpmpp_2s_a"
	jsonString, _ := json.MarshalIndent(sd, " ", "  ")

	tmpS2 := path.Join(common.UserDir, common.ImageDirectory, tmpS+".json")
	os.WriteFile(tmpS2, jsonString, os.ModePerm)

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
