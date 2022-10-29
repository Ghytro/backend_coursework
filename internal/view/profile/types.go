package profile

type AnyProfileViewData struct {
	Username,
	FullName,
	CountryCode,
	CountryFullName,
	Bio string

	HasRecentPolls bool
	Polls          []Poll
}

type Poll struct {
	CreatedAt,
	Title,
	IsAnonymousStr string

	Options []string
}

type MyProfileViewData struct {
	Username,
	FullName,
	CountryCode,
	CountryFullName,
	Bio string

	HasRecentPolls bool
	Polls          []Poll
}
