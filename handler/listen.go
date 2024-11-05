package main

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"time"

	htgotts "github.com/hegedustibor/htgo-tts"
)

type MpvPlayer struct{}

func (app *MpvPlayer) Play(fileName string) error {
	cmd := exec.Command("mpv", fileName)
	return cmd.Run()
}

func delete_file(fileName string) {
	err := os.Remove(fileName)
	if err != nil {
		log.Fatal(err)
	}
}

func listen_word(text, name,language_code string) bool {
    speech := htgotts.Speech{
        Folder:   "audio",
        Language: language_code,
        Handler:  &MpvPlayer{},
    }

    fileName := (time.Now()).String() + name
    filePath, err := speech.CreateSpeechFile(text, fileName)
    if err != nil {
        log.Fatal(err)
        return false
    }

    err = speech.Handler.Play(filePath)
    if err != nil {
        log.Fatal(err)
        return false
    }

    delete_file(filePath)
    return true
}

func get_data(name, language_code string) {
    file, err := read_file("../middleware/listen.txt")
    check_err(err)

    defer file.Close()

    scanner := bufio.NewScanner(file)
    var data string

    for scanner.Scan() {
        data = scanner.Text()
    }
    
    listen_word(data, name, language_code)
}