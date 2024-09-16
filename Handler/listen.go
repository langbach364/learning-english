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

// func check_word(text string) bool {
//     words := strings.Fields(text)
//         return len(words) == 1
// }

// func get_audio_word(text string) ([]AudioWord, error) {
//     apiKey := load_API_key("API_wordnik")
//     url := fmt.Sprintf("https://api.wordnik.com/v4/word.json/%s/audio?useCanonical=false&limit=50&api_key=%s", text, apiKey)

//     resp, err := http.Get(url)
// 	if err != nil {
// 		return nil, fmt.Errorf("lỗi khi gửi yêu cầu: %v", err)
// 	}
// 	defer resp.Body.Close()
//     var audioWord []AudioWord
    
//     decoder := json.NewDecoder(resp.Body)

// 	if err := decoder.Decode(&audioWord); err != nil {
// 		return nil, fmt.Errorf("lỗi khi giải mã JSON: %v", err)
// 	}

// 	return audioWord, nil
// }

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
