package main

import "easyssh/cmd"

func main() {
	err := cmd.Execute()
	if err != nil {
		return
	}
}
