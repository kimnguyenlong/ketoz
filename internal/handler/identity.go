package handler

import "github.com/gofiber/fiber/v2"

type identity struct {
}

func NewIdentity() *identity {
	return &identity{}
}

func (i *identity) RegisterRoutes(r fiber.Router) {
	g := r.Group("/identities")
	g.Get("/", i.List)
}

func (i *identity) List(c *fiber.Ctx) error {
	return c.SendString("List of identities")
}
