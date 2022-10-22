package main

import (
	"backend_coursework/internal/database"
	"backend_coursework/internal/entity"
	"backend_coursework/internal/repository"
	authService "backend_coursework/internal/usecase/auth"
	profileService "backend_coursework/internal/usecase/profile"
	"backend_coursework/internal/view"
	authView "backend_coursework/internal/view/auth"
	profileView "backend_coursework/internal/view/profile"
	"errors"
	"os"

	"github.com/go-pg/pg/v10"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
)

func serve() {
	db := database.NewPGDB(os.Getenv("DB_URL"), &database.PGLogger{})

	profileRepo := repository.NewUserRepo(db)
	profileService := profileService.NewService(profileRepo)
	profileView := profileView.NewView(profileService)

	authService := authService.NewService(profileRepo, "")
	authView := authView.NewView(authService)
	NewApp("", db, profileView, authView).Listen(":3001")
}

func NewApp(token interface{}, db repository.DBI, views ...view.View) *fiber.App {
	r := fiber.New()
	authHandler := jwtware.New(jwtware.Config{
		SigningKey:  token,
		TokenLookup: "cookie:jwt",
		ContextKey:  "user_jwt",
		SuccessHandler: func(c *fiber.Ctx) error {
			var u entity.User
			jwtToken, ok := c.Locals("user_jwt").(*jwt.Token)
			if !ok {
				return entity.ErrRespUnauthorized(errors.New("unable to get jwt"))
			}
			claims, ok := jwtToken.Claims.(jwt.MapClaims)
			if !ok {
				return entity.ErrRespUnauthorized(errors.New("unable to get claims from jwt"))
			}
			userIdClaims, ok := claims["id"].(float64)
			if !ok {
				return entity.ErrRespUnauthorized(errors.New("unable to get 'id' from claims"))
			}
			userId := entity.PK(userIdClaims)
			userPass, ok := claims["pass"].(string)
			if !ok {
				return entity.ErrRespUnauthorized(errors.New("unable to get 'pass' from claims"))
			}
			if err := db.ModelContext(c.Context(), &u).
				Where("id = ? AND password = crypt(?, password)", userId, userPass).
				Select(); err != nil {
				if err == pg.ErrNoRows {
					return entity.ErrRespUnauthorized(errors.New("incorrect token, auth again"))
				}
				return entity.ErrRespInternalServerError(err)
			}
			c.Locals("user_entity", &u)
			return c.Next()
		},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Redirect("/auth", fiber.StatusSeeOther)
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
