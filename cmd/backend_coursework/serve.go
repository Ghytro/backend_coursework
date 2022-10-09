package main

import (
	"backend_coursework/internal/common"
	profileController "backend_coursework/internal/controller/profile"
	"backend_coursework/internal/entity"
	profileModel "backend_coursework/internal/model/profile"
	"backend_coursework/internal/view"
	profileView "backend_coursework/internal/view/profile"
	"errors"
	"os"

	"github.com/go-pg/pg/v10"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
)

func serve() {
	dbOpts, err := pg.ParseURL(os.Getenv("DB_URL"))
	common.LogFatalErr(err)
	db := pg.Connect(dbOpts)

	profModel := profileModel.NewProfileModel(db)
	profController := profileController.NewController(profModel)
	profView := profileView.NewView(profController)
	NewApp("", db, profView).Listen(":3001")
}

func NewApp(token interface{}, db *pg.DB, views ...view.View) *fiber.App {
	r := fiber.New()
	authHandler := jwtware.New(jwtware.Config{
		SigningKey:  token,
		TokenLookup: "cookie:jwt",
		ContextKey:  "user_jwt",
		SuccessHandler: func(c *fiber.Ctx) error {
			var u entity.User
			jwtToken, ok := c.Locals("user_jwt").(*jwt.Token)
			if !ok {
				return &entity.ErrResponse{
					StatusCode: fiber.StatusUnauthorized,
					Err:        errors.New("unable to get jwt"),
				}
			}
			claims, ok := jwtToken.Claims.(jwt.MapClaims)
			if !ok {
				return &entity.ErrResponse{
					StatusCode: fiber.StatusUnauthorized,
					Err:        errors.New("unable to get claims from jwt"),
				}
			}
			userIdClaims, ok := claims["id"].(float64)
			if !ok {
				return &entity.ErrResponse{
					StatusCode: fiber.StatusUnauthorized,
					Err:        errors.New("unable to get 'id' from claims"),
				}
			}
			userId := entity.PK(userIdClaims)
			userPass, ok := claims["pass"].(string)
			if !ok {
				return &entity.ErrResponse{
					StatusCode: fiber.StatusUnauthorized,
					Err:        errors.New("unable to get 'pass' from claims"),
				}
			}
			if err := db.ModelContext(c.Context(), &u).
				Where("id = ? AND password = crypt(?, password)", userId, userPass).
				Select(); err != nil {
				if err == pg.ErrNoRows {
					return &entity.ErrResponse{
						StatusCode: fiber.StatusUnauthorized,
						Err:        errors.New("incorrect token, auth again"),
					}
				}
				return &entity.ErrResponse{
					StatusCode: fiber.StatusInternalServerError,
					Err:        err,
				}
			}
			c.Locals("user_entity", &u)
			return c.Next()
		},
	})
	middlewares := []fiber.Handler{
		func(c *fiber.Ctx) error {
			c.Set("Content-Type", "text/html;charset=utf-8")
			return c.Next()
		},
	}
	for _, v := range views {
		v.Routers(r, authHandler, middlewares...)
	}
	return r
}
