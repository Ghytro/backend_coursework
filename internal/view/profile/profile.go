package profile

import (
	"backend_coursework/internal/entity"
	"fmt"
	"strconv"

	"github.com/go-pg/pg/v10"
	"github.com/gofiber/fiber/v2"
)

type View struct {
	controller ProfileController
}

func NewView(c ProfileController) *View {
	return &View{
		controller: c,
	}
}

func (v *View) Routers(app fiber.Router, authHandler fiber.Handler, middlewares ...fiber.Handler) {
	r := fiber.New()
	r.Get("/:id", v.getProfile)
	app.Mount("/profile", r)
}

func (v *View) getProfile(c *fiber.Ctx) error {
	idParam, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return &entity.ErrResponse{
			StatusCode: fiber.StatusBadRequest,
			Err:        err,
		}
	}
	userID := entity.PK(idParam)
	user, err := v.controller.GetUser(c.Context(), userID)
	if err != nil {
		if err == pg.ErrNoRows {
			return &entity.ErrResponse{
				StatusCode: fiber.StatusNotFound,
				Err:        err,
			}
		}
		return &entity.ErrResponse{
			StatusCode: fiber.StatusBadRequest,
			Err:        err,
		}
	}
	return c.SendString(fmt.Sprintf("%v", user))
}
