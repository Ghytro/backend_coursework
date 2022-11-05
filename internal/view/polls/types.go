package polls

type NewPollRequest struct {
	Topic          string   `form:"topic"`
	Options        []string `form:"options"`
	IsAnonymous    string   `form:"is_anonymous"`
	MultipleChoice string   `form:"multiple_choice"`
	CantRevote     string   `form:"cant_revote"`
}

type GetPollViewData struct {
	PollID,
	Topic,
	UserID,
	Username string
	MultipleChoice,
	RevoteAbility,
	IsAnonymous,
	CurrentUserVoted bool
	Options          []Option
	CurrentUserVotes []bool
}

type Option struct {
	Option, VotesNumber string
}
