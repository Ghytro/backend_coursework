package profile

import (
	"backend_coursework/internal/common"
	"backend_coursework/internal/entity"
	"backend_coursework/internal/view"
	"errors"
	"strconv"

	"github.com/go-pg/pg/v10"
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
	r.Get("/:id", v.getProfile)
	r.Use(authHandler)
	r.Get("/", v.getMyProfile)
	app.Mount("/profile", r)
}

func (v *View) getProfile(c *fiber.Ctx) error {
	idParam, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return entity.ErrRespBadRequest(err)
	}
	userID := entity.PK(idParam)
	user, err := v.service.GetUser(c.Context(), userID)
	if err != nil {
		if err == pg.ErrNoRows {
			return entity.ErrRespNotFound(err)
		}
		return entity.ErrRespBadRequest(err)
	}
	tpl := templates.MustGet("profile/any.html")
	viewData := AnyProfileViewData{
		Username: user.Username,
	}
	if user.FirstName != nil {
		viewData.FullName += *user.FirstName
	}
	if user.LastName != nil {
		viewData.FullName += *user.LastName
	}
	if user.Country != nil {
		viewData.CountryCode = *user.Country
		viewData.CountryFullName = *user.Country // TODO country mapping
	} else {
		viewData.CountryCode = "AQ" // unknown
		viewData.CountryFullName = "Unknown"
	}
	if user.Bio != nil {
		viewData.Bio = *user.Bio
	} else {
		viewData.Bio = "<Статус пуст>"
	}
	return view.SendTemplate(c, tpl, viewData)
}

func (v *View) getMyProfile(c *fiber.Ctx) error {
	user, ok := c.Locals("user_entity").(*entity.User)
	if !ok {
		return entity.ErrRespUnauthorized(errors.New("некорректный токен авторизации"))
	}
	user, err := v.service.GetUser(c.Context(), user.ID)
	if err != nil {
		if err == pg.ErrNoRows {
			return entity.ErrRespNotFound(errors.New("пользователь не найден"))
		}
		return entity.ErrRespInternalServerError(err)
	}
	tpl := templates.MustGet("profile/my.html")
	viewData := MyProfileViewData{
		Username: user.Username,
	}
	if user.FirstName != nil {
		viewData.FullName += *user.FirstName
	}
	if user.LastName != nil {
		viewData.FullName += *user.LastName
	}
	viewData.CountryCode = "AQ" // unknown
	viewData.CountryFullName = "&lt;Unknown&gt;"
	if user.Country != nil {
		if c := common.GetCountryByAlpha2(*user.Country); c != nil {
			viewData.CountryCode = *user.Country
			viewData.CountryFullName = c.Code.StringRus()
		}
	}
	if user.Bio != nil {
		viewData.Bio = *user.Bio
	} else {
		viewData.Bio = "<Статус пуст>"
	}
	return view.SendTemplate(c, tpl, viewData)
}
