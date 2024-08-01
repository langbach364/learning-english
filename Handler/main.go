package main

import "fmt"


func scan_verbs_and_compare() {
	verbs_file, err := open_file("./Database/Data-txt/Verbs/verbs.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer verbs_file.Close()

	text_file, err := open_file("./Database//Data-txt/Text.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer text_file.Close()

	text_words := read_text(text_file)

	lack_word_file, err := create_lack("lack-word.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer lack_word_file.Close()

	scan_verbs(verbs_file, text_words, lack_word_file)

	fmt.Println("Quá trình so sánh và ghi file hoàn tất.")
}

func main () {
	scan_verbs_and_compare()
}