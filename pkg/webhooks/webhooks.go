package webhooks

type Webhook struct {
	FormResponse FormResponse `json:"form_response"`
}

type FormResponse struct {
	Answers []Answers `json:"answers"`
}

type Answers struct {
	Type   string  `json:"type"`
	Number float32 `json:"number"`
	Text   string  `json:"text"`
	Date   string  `json:"date"`
	Choice Choice  `json:"choice"`
	PicURL string  `json:"file_url"`
	Field  Field   `json:"field"`
}

type Choice struct {
	Label string `json:"label"`
}

type Field struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Ref  string `json:"ref"`
}
