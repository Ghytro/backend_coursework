package profile

import (
	"backend_coursework/internal/entity"
	"errors"
	"strconv"

	"github.com/go-pg/pg/v10"
	"github.com/gofiber/fiber/v2"
)

type View struct {
	controller ProfileController
}

type AnyProfileViewData struct {
	UserName string
}

type MyProfileViewData struct {
	UserName string
}

func NewView(c ProfileController) *View {
	return &View{
		controller: c,
	}
}

func (v *View) Routers(app fiber.Router, authHandler fiber.Handler, middlewares ...fiber.Handler) {
	r := fiber.New()
	for _, m := range middlewares {
		r.Use(m)
	}
	r.Get("/:id", v.getProfile)
	r.Use(authHandler)
	r.Get("/", v.getMyProfile)
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
	tpl := templates["profile/any.html"]
	viewData := AnyProfileViewData{
		UserName: user.Username,
	}
	if err := tpl.Execute(c.Context().Response.BodyWriter(), viewData); err != nil {
		return &entity.ErrResponse{
			StatusCode: fiber.StatusInternalServerError,
			Err:        errors.New("unable to send page"),
		}
	}
	return c.SendStatus(fiber.StatusOK)
}

func (v *View) getMyProfile(c *fiber.Ctx) error {
	user, ok := c.Locals("user_entity").(*entity.User)
	if !ok {
		return &entity.ErrResponse{
			StatusCode: fiber.StatusUnauthorized,
			Err:        errors.New("авторизация неудачна"),
		}
	}
	user, err := v.controller.GetUser(c.Context(), user.ID)
	if err != nil {
		if err == pg.ErrNoRows {
			return &entity.ErrResponse{
				StatusCode: fiber.StatusNotFound,
				Err:        errors.New("пользователь не найден"),
			}
		}
		return &entity.ErrResponse{
			StatusCode: fiber.StatusInternalServerError,
			Err:        err,
		}
	}
	tpl := templates["profile/my.html"]
	data := MyProfileViewData{
		UserName: user.Username,
	}
	if err := tpl.Execute(c.Response().BodyWriter(), data); err != nil {
		return &entity.ErrResponse{
			StatusCode: fiber.StatusInternalServerError,
			Err:        err,
		}
	}
	return c.SendStatus(fiber.StatusOK)
}
