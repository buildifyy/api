package models

type Dropdown struct {
	Label  string `json:"label"`
	Value  string `json:"value"`
	Symbol string `json:"symbol"`
}

type ParentTemplateDropdown struct {
	Label        string `json:"label"`
	Value        string `json:"value"`
	RootTemplate string `json:"rootTemplate"`
}
