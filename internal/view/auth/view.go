package auth

import (
	"backend_coursework/internal/entity"
	"errors"

	"github.com/gofiber/fiber/v2"
)

type View struct {
	service UseCase
}

func NewView(s UseCase) *View {
	return &View{
		service: s,
	}
}

func (v *View) Routers(app fiber.Router, authHandler fiber.Handler, middlewares ...fiber.Handler) {
	r := fiber.New()
	for _, m := range middlewares {
		r.Use(m)
	}
	r.Post("/", v.makeAuth)
	r.Patch("/", v.patchAuth)
	app.Mount("/auth", r)

	r = fiber.New()
	r.Post("/", v.register)
	app.Mount("/register", r)
}

func (v *View) makeAuth(c *fiber.Ctx) error {
	var model MakeAuthRequest
	if err := c.BodyParser(&model); err != nil {
		return entity.ErrRespIncorrectForm()
	}

	token, err := v.service.MakeAuth(c.Context(), model.Username, model.Password)
	if err != nil {
		return entity.ErrRespBadRequest(err)
	}
	c.Cookie(&fiber.Cookie{Name: "jwt", Value: token})
	return c.Send(nil)
}

func (v *View) patchAuth(c *fiber.Ctx) error {
	return errors.New("missing impl") // TODO
}

func (v *View) register(c *fiber.Ctx) error {
	var model entity.User
	if err := c.BodyParser(&model); err != nil {
		return entity.ErrRespIncorrectForm()
	}

	token, err := v.service.Register(c.Context(), &model)
	if err != nil {
		return entity.ErrRespBadRequest(err)
	}
	c.Cookie(&fiber.Cookie{Name: "jwt", Value: token})
	return c.Redirect("/profile", fiber.StatusSeeOther)
}
