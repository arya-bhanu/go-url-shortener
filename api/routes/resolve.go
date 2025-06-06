package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/arya-bhanu/go-url-shortener/database"

)

func ResolveURL(c *fiber.Ctx) error {
	url := c.Params("url")

	r := database.CreateClientRedis(0)
	defer r.Close()

	r.Get(database.CtxDb, url)

	return nil
}
