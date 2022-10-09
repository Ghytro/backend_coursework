package main

import (
	"backend_coursework/internal/common"
	profileController "backend_coursework/internal/controller/profile"
	profileModel "backend_coursework/internal/model/profile"
	"backend_coursework/internal/view"
	profileView "backend_coursework/internal/view/profile"
	"os"

	"github.com/go-pg/pg/v10"
	"github.com/gofiber/fiber/v2"
)

func serve() {
	dbOpts, err := pg.ParseURL(os.Getenv("DATABASE_URL"))
	common.LogFatalErr(err)
	db := pg.Connect(dbOpts)

	profModel := profileModel.NewProfileModel(db)
	profController := profileController.NewController(profModel)
	profView := profileView.NewView(profController)
	NewApp(profView).Listen(":3001")
}

func NewApp(views ...view.View) *fiber.App {
	r := fiber.New()
	authHandler := func(c *fiber.Ctx) error {
		return nil
	}
	for _, v := range views {
		v.Routers(r, authHandler)
	}
	return r
}
