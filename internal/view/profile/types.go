package profile

type View struct {
	Username     string
	CreatedPolls uint
	RecentVotes  []Vote
}

type Vote struct {
	PollName    string
	PollOptions string
}
