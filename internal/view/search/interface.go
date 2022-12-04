package search

import (
	"backend_coursework/internal/common"
	"backend_coursework/internal/entity"
	"backend_coursework/internal/usecase/search"
	"context"
)

type UseCase interface {
	Search(ctx context.Context, query string, page *search.PageData) (*search.SearchResult, error)
	SearchUser(ctx context.Context, searchParams *search.UserSearchParams, page *common.PageData) ([]*entity.User, error)
	SearchPoll(ctx context.Context, searchParams *search.PollSearchParams, page *common.PageData) ([]*entity.Poll, error)
}
