package profile

type AnyProfileViewData struct {
	UserName,
	FullName,
	CountryCode,
	CountryFullName,
	Bio string

	Polls []Poll
}

type Poll struct {
	CreatedAt,
	Title,
	IsAnonymousStr string

	Options []string
}

type MyProfileViewData struct {
	UserName string
}
