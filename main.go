package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/rasecoiac03/clauda/pkg/config"
)

var fileFormat = "files/%s"

func main() {
	downloadFiles()
}

func downloadFiles() {
	files := []string{
		"lotomania_file_enpoint",
		"lotofacil_file_enpoint",
	}
	executeCmd("files folder exists", exec.Command("mkdir", "-p", fmt.Sprintf(fileFormat, "")))
	executeCmd("files removed", exec.Command("rm", "-f", fmt.Sprintf(fileFormat, "*")))
	for _, file := range files {
		filePath := config.GetConfig(file)
		filePathParts := strings.Split(filePath, "/")
		fileName := filePathParts[len(filePathParts)-1]
		executeCmd(fmt.Sprintf("file downloaded: %s", fileName), exec.Command("wget", filePath))
		executeCmd(fmt.Sprintf("file moved to folder: %s", fileName), exec.Command("mv", fileName, fmt.Sprintf(fileFormat, fileName)))
	}
}

func executeCmd(message string, cmd *exec.Cmd) {
	defer log.Println(message)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
