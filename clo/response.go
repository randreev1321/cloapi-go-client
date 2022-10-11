package clo

type Response struct {
	Code  int    `json:"code,omitempty"`
	Title string `json:"title,omitempty"`
}

type Link struct {
	Rel  string `json:"rel,omitempty"`
	Href string `json:"href,omitempty"`
}
