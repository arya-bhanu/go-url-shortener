package routes

import (
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"

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

	return nil
}
