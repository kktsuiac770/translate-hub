package models

type Project struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	SourceLang string `json:"source_lang"`
	TargetLang string `json:"target_lang"`
}
