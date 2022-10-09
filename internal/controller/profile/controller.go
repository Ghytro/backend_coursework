package profile

import (
	"backend_coursework/internal/entity"
	"context"
)

type Controller struct {
	model ProfileModel
}

func NewController(model ProfileModel) *Controller {
	return &Controller{
		model: model,
	}
}

func (c *Controller) CreateUser(ctx context.Context, user *entity.User) error {
	return c.model.CreateUser(ctx, user)
}

func (c *Controller) GetUser(ctx context.Context, userID entity.PK) (*entity.User, error) {
	return c.model.GetUser(ctx, userID)
}

func (c *Controller) UpdateUser(ctx context.Context, user *entity.User) error {
	return c.model.UpdateUser(ctx, user)
}

func (c *Controller) DeleteUser(ctx context.Context, userID entity.PK) error {
	return c.model.DeleteUser(ctx, userID)
}
