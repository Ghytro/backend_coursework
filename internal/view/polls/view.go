package polls

import (
	"backend_coursework/internal/entity"
	"backend_coursework/internal/view"
	"errors"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
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
	r.Post("/:id/vote", v.vote)
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
	user, ok := c.Locals("user_entity").(*entity.User)
	if !ok {
		return entity.ErrRespUnauthorized(errors.New("авторизуйтесь заново"))
	}
	poll, currentUserVotes, err := v.service.GetPollWithVotesAmount(c.Context(), entity.PK(pollID), user.ID)
	if err != nil {
		return entity.ErrRespBadRequest(err)
	}
	viewData := GetPollViewData{
		PollID:           fmt.Sprint(poll.ID),
		Topic:            poll.Topic,
		UserID:           fmt.Sprint(poll.Creator.ID),
		Username:         poll.Creator.Username,
		IsAnonymous:      poll.IsAnonymous,
		MultipleChoice:   poll.MultipleChoice,
		RevoteAbility:    poll.RevoteAbility,
		CurrentUserVoted: len(currentUserVotes) != 0,
	}
	for _, o := range poll.Options {
		viewData.Options = append(viewData.Options, Option{
			Option:      o.Option,
			VotesNumber: fmt.Sprint(o.VotesAmount),
		})
	}

	if viewData.CurrentUserVoted {
		viewData.CurrentUserVotes = make([]bool, len(poll.Options))
		for _, v := range currentUserVotes {
			_, optIdx, _ := lo.FindIndexOf(poll.Options, func(o *entity.PollOption) bool {
				return o.ID == v.OptionID
			})
			viewData.CurrentUserVotes[optIdx] = true
		}
	}

	tpl := templates.MustGet("polls/get.html")
	return view.SendTemplate(c, tpl, viewData)
}

func (v *View) vote(c *fiber.Ctx) error {
	pollID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return entity.ErrRespBadRequest(err)
	}
	f, err := c.MultipartForm()
	if err != nil {
		return entity.ErrRespInternalServerError(err)
	}
	idxs := lo.Map(f.Value["votes"], func(strIdx string, _ int) int {
		i, _ := strconv.Atoi(strIdx)
		return i
	})
	if lo.Contains(idxs, 0) {
		return entity.ErrRespBadRequest(errors.New("значение выбранной опции может быть только числовым"))
	}
	user, ok := c.Locals("user_entity").(*entity.User)
	if !ok {
		return entity.ErrRespUnauthorized(errors.New("авторизуйтесь заново"))
	}
	if err := v.service.Vote(c.Context(), user.ID, entity.PK(pollID), idxs...); err != nil {
		return entity.ErrRespBadRequest(err)
	}
	return c.Redirect(fmt.Sprintf("/votes/%d", pollID), fiber.StatusSeeOther)
}
