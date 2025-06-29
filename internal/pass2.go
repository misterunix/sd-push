package main 

import (
	"sd-push/internal/common"
)


func secondpass(sd common.Stable)
	prompt, nprompt, model string, r int, thesteps int) error {
	fmt.Println("secondpass")
	sd := common.Stable{}
	sd.RandomNumber = r
	tmpS := strconv.Itoa(r)
	sd.SmallImage = path.Join(common.UserDir, common.ImageDirectory, tmpS+"-small.png")
	sd.LargeImage = path.Join(common.UserDir, common.ImageDirectory, tmpS+"-large.png")

	// sd.SmallImage = common.UserDir + "/" + common.Tss + "-" + "small.png"
	// sd.LargeImage = common.UserDir + "/" + common.Tss + "-" + "large.png"
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