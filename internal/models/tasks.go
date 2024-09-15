package models

type TasksResponse struct {
	Tasks       []any `json:"tasks"`
	SubSections []struct {
		Title string `json:"title"`
		Tasks []Task `json:"tasks"`
	} `json:"subSections"`
}

type Task struct {
	ID                 string `json:"id"`
	Kind               string `json:"kind"`
	Type               string `json:"type"`
	Status             string `json:"status"`
	ValidationType     string `json:"validationType"`
	IconFileKey        string `json:"iconFileKey"`
	BannerFileKey      any    `json:"bannerFileKey"`
	Title              string `json:"title"`
	ProductName        any    `json:"productName"`
	Description        any    `json:"description"`
	Reward             string `json:"reward"`
	SocialSubscription struct {
		OpenInTelegram bool   `json:"openInTelegram"`
		URL            string `json:"url"`
	} `json:"socialSubscription"`
	IsHidden             bool `json:"isHidden"`
	IsDisclaimerRequired bool `json:"isDisclaimerRequired"`
}
