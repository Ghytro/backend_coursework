package polls

type NewPollRequest struct {
	Topic          string   `form:"topic"`
	Options        []string `form:"options"`
	IsAnonymous    string   `form:"is_anonymous"`
	MultipleChoice string   `form:"multiple_choice"`
	CantRevote     string   `form:"cant_revote"`
}

type GetPollViewData struct {
	Topic,
	UserID,
	Username string
	Options []Option
}

type Option struct {
	Option, VotesNumber string
}
