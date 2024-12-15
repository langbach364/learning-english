package main

type RandomWordParams struct {
	HasDictionaryDef    bool   `json:"hasDictionaryDef,omitempty"`    // Chỉ trả về các từ có trong từ điển
	IncludePartOfSpeech string `json:"includePartOfSpeech,omitempty"` // Chỉ bao gồm các từ loại được chỉ định (noun, verb, adjective,...) Ví dụ: "noun,verb,adjective"
	ExcludePartOfSpeech string `json:"excludePartOfSpeech,omitempty"` // Loại trừ các từ loại được chỉ định Ví dụ: "adverb,preposition"
	MinCorpusCount      int    `json:"minCorpusCount,omitempty"`      // Tần suất xuất hiện tối thiểu trong ngữ liệu
	MaxCorpusCount      int    `json:"maxCorpusCount,omitempty"`      // Tần suất xuất hiện tối đa trong ngữ liệu
	MinDictionaryCount  int    `json:"minDictionaryCount,omitempty"`  // Số lượng từ điển tối thiểu chứa từ này
	MaxDictionaryCount  int    `json:"maxDictionaryCount,omitempty"`  // Số lượng từ điển tối đa chứa từ này
	MinLength           int    `json:"minLength,omitempty"`           // Độ dài tối thiểu của từ
	MaxLength           int    `json:"maxLength,omitempty"`           // Độ dài tối đa của từ
	SortBy              string `json:"sortBy,omitempty"`              // Thuộc tính để sắp xếp kết quả Ví dụ: "alpha" (theo bảng chữ cái), "count" (theo tần suất)
	SortOrder           string `json:"sortOrder,omitempty"`           // Thứ tự sắp xếp Giá trị: "asc" (tăng dần) hoặc "desc" (giảm dần)
	Limit               int    `json:"limit,omitempty"`               // Số lượng từ tối đa trả về
}

type infoWord struct {
	Word       string `json:"word"`
	Frequency  int    `json:"frequency"`
	ErrorCount int    `json:"error_count"`
}

type TargetDate struct {
	Date string `json:"target_date"`
}

type Timege struct {
	Range string `json:"range"`
	Date  string `json:"date"`
}
