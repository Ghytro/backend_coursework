package entity

import (
	"time"

	"github.com/go-pg/pg/v10/types"
	"github.com/google/uuid"
)

type PK uint

type baseEntity struct {
	CreatedAt time.Time      `pg:"created_at"`
	DeletedAt types.NullTime `pg:"deleted_at"`
}

type pkID struct {
	ID PK `pg:"id,pk"`
}

type User struct {
	tableName struct{} `pg:"users"`

	pkID
	baseEntity

	Username  string  `pg:"username" form:"username"`
	Password  string  `pg:"password" form:"password"`
	Bio       *string `pg:"bio" form:"bio"`
	AvatarUrl *string `pg:"avatar_url" form:"avatar_url"`
	Polls     []*Poll `pg:"rel-has-many"`
	Votes     []*Vote `pg:"rel-has-many"`
}

type Poll struct {
	tableName struct{} `pg:"polls"`

	pkID
	baseEntity

	CreatorID PK    `pg:"creator_id"`
	Creator   *User `pg:"rel:has-one"`

	Topic          string        `pg:"topic"`
	IsAnonymous    bool          `pg:"is_anonymous"`
	MultipleChoice bool          `pg:"multiple_choice"`
	RevoteAbility  bool          `pg:"revote_ability"`
	Options        []*PollOption `pg:"rel-has-many"`
}

type PollOption struct {
	tableName struct{} `pg:"poll_options"`

	pkID

	PollID PK    `pg:"poll_id"`
	Poll   *Poll `pg:"rel:has-one"`

	Votes []*Vote `pg:"rel-has-many"`

	Index     int       `pg:"index"`
	Option    string    `pg:"option"`
	UpdatedAt time.Time `pg:"updated_at"`
}

type Vote struct {
	tableName struct{} `pg:"votes"`

	baseEntity
	ID uuid.UUID `pg:"id"`

	UserID PK    `pg:"user_id"`
	User   *User `pg:"rel:has-one"`

	OptionID PK          `pg:"option_id"`
	Option   *PollOption `pg:"rel:has-one"`
}

type ErrResponse struct {
	StatusCode int
	Err        error
}

func (e *ErrResponse) Error() string {
	return e.Err.Error()
}

func (e *ErrResponse) Unwrap() error {
	return e.Err
}
