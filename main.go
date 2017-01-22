package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/rasecoiac03/clauda/pkg/config"
)

var fileFolder = "files/"

func main() {
	downloadFiles()
}

func downloadFiles() {
	files := []string{
		"lotomania_file_enpoint",
		"lotofacil_file_enpoint",
	}
	executeCmd("files removed", exec.Command("rm", "-rf", fileFolder))
	executeCmd("files folder exists", exec.Command("mkdir", "-p", fileFolder))
	for _, file := range files {
		filePath := config.GetConfig(file)
		filePathParts := strings.Split(filePath, "/")
		fileName := filePathParts[len(filePathParts)-1]
		executeCmd(fmt.Sprintf("file downloaded: %s", fileName), exec.Command("wget", filePath))
		executeCmd(fmt.Sprintf("file unzipped: %s", fileName), exec.Command("unzip", "-u", fileName, "-d", fileFolder))
		executeCmd(fmt.Sprintf("file zip removed: %s", fileName), exec.Command("rm", "-f", fileName))
		executeCmd("cleaning", exec.Command("/bin/sh", "-c", "rm -f "+fileFolder+"*.GIF"))
	}
}

func executeCmd(message string, cmd *exec.Cmd) {
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(message)
}
