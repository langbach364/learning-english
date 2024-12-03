package main

import "time"

type FieldsWord struct {
	HasDictionaryDef    bool   // Chỉ trả về các từ có trong từ điển
	IncludePartOfSpeech string // Chỉ bao gồm các từ loại được chỉ định (noun, verb, adjective,...) Ví dụ: "noun,verb,adjective"
	ExcludePartOfSpeech string // Loại trừ các từ loại được chỉ định Ví dụ: "adverb,preposition"
	MinCorpusCount      int    // Tần suất xuất hiện tối thiểu trong ngữ liệu
	MaxCorpusCount      int    // Tần suất xuất hiện tối đa trong ngữ liệu
	MinDictionaryCount  int    // Số lượng từ điển tối thiểu chứa từ này
	MaxDictionaryCount  int    // Số lượng từ điển tối đa chứa từ này
	MinLength           int    // Độ dài tối thiểu của từ
	MaxLength           int    // Độ dài tối đa của từ
	SortBy              string // Thuộc tính để sắp xếp kết quả Ví dụ: "alpha" (theo bảng chữ cái), "count" (theo tần suất)
	SortOrder           string // Thứ tự sắp xếp Giá trị: "asc" (tăng dần) hoặc "desc" (giảm dần)
	Limit               int    // Số lượng từ tối đa trả về
}

type InfoWord struct {
	time time.Time
	word string
}
