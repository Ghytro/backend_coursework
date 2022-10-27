package polls

import (
	"backend_coursework/internal/view"
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

func (v *View) Routers(router fiber.Router, authHandler fiber.Handler, middlewares ...fiber.Handler) {
	r := fiber.New()
	for _, m := range middlewares {
		r.Use(m)
	}
	r.Use(authHandler)
	r.Get("/new", v.newPollPage)
	r.Post("/new", v.postNewPoll)
	router.Mount("/polls", r)
}

func (v *View) newPollPage(c *fiber.Ctx) error {
	tpl := templates.MustGet("polls/new.html")
	return view.SendTemplate(c, tpl, nil)
}

func (v *View) postNewPoll(c *fiber.Ctx) error {
	return errors.New("missing impl")
}
