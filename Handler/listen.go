package main

import (
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
