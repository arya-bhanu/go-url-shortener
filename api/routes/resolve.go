package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"

	"github.com/arya-bhanu/go-url-shortener/database"
)

func ResolveURL(c *fiber.Ctx) error {
	url := c.Params("url")

	r := database.CreateClientRedis(0)
	defer r.Close()

	r.Get(database.CtxDb, url)

	value, err := r.Get(database.CtxDb, url).Result()

	if err != nil {
		if err == redis.Nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": fmt.Errorf("short url not found: %s", err.Error()),
			})
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": fmt.Errorf("internal server error: %s", err.Error()),
			})
		}
	}

	rIncrement := database.CreateClientRedis(1)
	defer rIncrement.Close()

	_ = rIncrement.Incr(database.CtxDb, "counter")

	return c.Redirect(value, 301)
}
