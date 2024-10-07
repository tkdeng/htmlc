package htmlcfiber

import (
	"testing"

	"github.com/gofiber/fiber/v3"
)

func Test(t *testing.T) {
	engine, err := New("../test/templates")
	if err != nil {
		panic(err)
	}

	app := fiber.New(fiber.Config{
		Views:             engine,
		PassLocalsToViews: true,
	})

	app.Get("/", func(c fiber.Ctx) error {
		return c.Render("index", fiber.Map{"title": "Test"}, "layout")
	})

	app.Listen(":3000")
}
