package polls

import (
	"backend_coursework/internal/entity"
	"backend_coursework/internal/view"
	"errors"
	"fmt"
	"strconv"

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
	r.Get("/new/", v.newPollPage)
	r.Get("/:id", v.getPoll)
	r.Post("/", v.postNewPoll)
	router.Mount("/polls", r)
}

func (v *View) newPollPage(c *fiber.Ctx) error {
	tpl := templates.MustGet("polls/new.html")
	return view.SendTemplate(c, tpl, nil)
}

func (v *View) postNewPoll(c *fiber.Ctx) error {
	var model NewPollRequest
	if err := c.BodyParser(&model); err != nil {
		return entity.ErrRespBadRequest(err)
	}
	model.IsAnonymous = c.FormValue("is_anonymous")
	model.MultipleChoice = c.FormValue("multiple_choice")
	model.CantRevote = c.FormValue("cant_revote")
	user, ok := c.Locals("user_entity").(*entity.User)
	if !ok {
		return entity.ErrRespUnauthorized(errors.New("авторизуйтесь заново"))
	}
	poll, err := v.service.CreatePoll(c.Context(), user.ID, &model)
	if err != nil {
		return entity.ErrRespBadRequest(err)
	}
	return c.Redirect(fmt.Sprintf("/polls/%d", poll.ID), fiber.StatusSeeOther)
}

func (v *View) getPoll(c *fiber.Ctx) error {
	pollID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return entity.ErrRespBadRequest(err)
	}
	poll, err := v.service.GetPollWithVotesAmount(c.Context(), entity.PK(pollID))
	if err != nil {
		return entity.ErrRespBadRequest(err)
	}
	viewData := GetPollViewData{
		Topic:    poll.Topic,
		UserID:   fmt.Sprint(poll.Creator.ID),
		Username: poll.Creator.Username,
	}
	for _, o := range poll.Options {
		viewData.Options = append(viewData.Options, Option{
			Option:      o.Option,
			VotesNumber: fmt.Sprint(o.VotesAmount),
		})
	}
	tpl := templates.MustGet("polls/get.html")
	return view.SendTemplate(c, tpl, viewData)
}
