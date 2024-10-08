package main

type WordDefinition struct {
	ID               string        `json:"id"`               // ID định danh duy nhất cho định nghĩa
	PartOfSpeech     string        `json:"partOfSpeech"`     // Từ loại (ví dụ: danh từ, động từ, tính từ)
	AttributionText  string        `json:"attributionText"`  // Văn bản ghi nhận nguồn của định nghĩa
	SourceDictionary string        `json:"sourceDictionary"` // Tên của từ điển cung cấp định nghĩa này
	Define           string        `json:"text"`             // Nội dung chính của định nghĩa
	Sequence         string        `json:"sequence"`         // Số thứ tự của định nghĩa trong danh sách
	Score            float64       `json:"score"`            // Điểm số liên quan đến mức độ phổ biến hoặc liên quan
	Word             string        `json:"word"`             // Từ được định nghĩa
	AttributionURL   string        `json:"attributionUrl"`   // URL dẫn đến nguồn của định nghĩa
	WordnikURL       string        `json:"wordnikUrl"`       // URL dẫn đến trang Wordnik cho từ này
	Citations        []interface{} `json:"citations"`        // Danh sách các trích dẫn minh họa cách sử dụng từ
	ExampleUses      []interface{} `json:"exampleUses"`      // Danh sách các ví dụ về cách sử dụng từ
	Labels           []interface{} `json:"labels"`           // Các nhãn phân loại hoặc mô tả bổ sung cho từ
	Notes            []interface{} `json:"notes"`            // Các ghi chú bổ sung về từ hoặc định nghĩa
	RelatedWords     []interface{} `json:"relatedWords"`     // Danh sách các từ có liên quan
	TextProns        []interface{} `json:"textProns"`        // Danh sách các cách phát âm của từ được biểu diễn bằng văn bản
}

type AudioWord struct {
	CommentCount    string `json:"commentCount"`    // Số lượng bình luận về file âm thanh
	CreatedBy       string `json:"createdBy"`       // Người tạo file âm thanh
	CreatedAt       string `json:"createdAt"`       // Thời gian tạo file âm thanh
	ID              string `json:"id"`              // ID định danh duy nhất cho file âm thanh
	Word            string `json:"word"`            // Từ được phát âm trong file âm thanh
	Duration        string `json:"duration"`        // Thời lượng của file âm thanh
	AudioType       string `json:"audioType"`       // Cách phát âm
	AttributionText string `json:"attributionText"` // Văn bản ghi nhận nguồn của file âm thanh
	AttributionUrl  string `json:"attributionUrl"`  // URL dẫn đến nguồn của file âm thanh
	FileUrl         string `json:"fileUrl"`         // URL để tải xuống file âm thanh
}
