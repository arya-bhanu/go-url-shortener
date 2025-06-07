package routes

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"

	"github.com/arya-bhanu/go-url-shortener/database"
	"github.com/arya-bhanu/go-url-shortener/helpers"
)

type Request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"custom_short"`
	Expiry      time.Duration `json:"expiry"`
}

type Response struct {
	URL            string        `json:"url"`
	CustomShort    string        `json:"custom_short"`
	Expiry         time.Duration `json:"expiry"`
	XRateRemaining int           `json:"rate_remaining"`
	XRateLimitRest time.Duration `json:"rate_limit_rest"`
}

func ShortenUrlHandler(c *fiber.Ctx) error {
	body := new(Request)
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "failed parse JSON"})
	}

	// implement rate limiting
	rDb := database.CreateClientRedis(2)
	defer rDb.Close()
	val, err := rDb.Get(database.CtxDb, c.IP()).Result()

	if err != nil {
		if err == redis.Nil {
			_ = rDb.Set(database.CtxDb, c.IP(), os.Getenv("API_QUOTA"), 30*time.Minute)
		} else {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error": fmt.Sprintf("something wrong with redis: %s", err.Error()),
			})
		}
	} else {
		intVal, err := strconv.Atoi(val)
		if err != nil {
			log.Fatalf("failed convert string into integer: %s\n", err.Error())
		}
		if intVal <= 0 {
			c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error": "to many request for processing the url",
			})
		}
	}

	// check the url, is it valid url
	if ok := govalidator.IsURL(body.URL); !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid url"})
	}

	// check for domain error
	if ok := helpers.RemoveDomainError(body.URL); !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid domain"})
	}

	// enforce http
	validUrl, err := helpers.EnforceHTTP(body.URL)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	body.URL = validUrl

	// storing the url in redis database

	rDb.Decr(database.CtxDb, c.IP())
	return nil
}
