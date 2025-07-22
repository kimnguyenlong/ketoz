package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kimnguyenlong/ketoz/internal/entity"
	"github.com/kimnguyenlong/ketoz/internal/repository"
)

type identity struct {
	idRepo repository.Identity
}

func NewIdentity(idRepo repository.Identity) *identity {
	return &identity{
		idRepo: idRepo,
	}
}

func (i *identity) RegisterRoutes(r fiber.Router) {
	g := r.Group("/identities")
	// Identities
	g.Post("/", i.Create)
	g.Get("/", i.List)
	g.Get("/:id", i.Get)
	// Children
	g.Get("/:id/children", i.ListChildren)
	g.Post("/:id/children", i.AddChild)
	// Permissions
	g.Get("/:id/permissions", i.ListPermissions)
}

type CreateIdentityRequest struct {
	Id string `json:"id"`
}

func (i *identity) Create(c *fiber.Ctx) error {
	req := new(CreateIdentityRequest)
	if err := c.BodyParser(req); err != nil {
		return responseError(c, entity.NewInvalidParamsError(err.Error()))
	}

	if err := i.idRepo.Create(c.Context(), &entity.Identity{
		Id: req.Id,
	}); err != nil {
		return responseError(c, err)
	}

	return responseNilData(c, fiber.StatusOK)
}

func (i *identity) List(c *fiber.Ctx) error {
	list, err := i.idRepo.List(c.Context())
	if err != nil {
		return responseError(c, err)
	}

	return responseRecords(c, fiber.StatusOK, list)
}

func (i *identity) Get(c *fiber.Ctx) error {
	idn, err := i.idRepo.Get(c.Context(), c.Params("id"))
	if err != nil {
		return responseError(c, err)
	}

	return responseData(c, fiber.StatusOK, idn)
}

func (i *identity) ListChildren(c *fiber.Ctx) error {
	list, err := i.idRepo.ListChildren(c.Context(), c.Params("id"))
	if err != nil {
		return responseError(c, err)
	}

	return responseRecords(c, fiber.StatusOK, list)
}

type AddChildIdentityRequest struct {
	ChildId string `json:"child_id"`
}

func (i *identity) AddChild(c *fiber.Ctx) error {
	req := new(AddChildIdentityRequest)
	if err := c.BodyParser(req); err != nil {
		return responseError(c, entity.NewInvalidParamsError(err.Error()))
	}

	if err := i.idRepo.AddChild(c.Context(), c.Params("id"), req.ChildId); err != nil {
		return responseError(c, err)
	}

	return responseNilData(c, fiber.StatusOK)
}

func (i *identity) ListPermissions(c *fiber.Ctx) error {
	list, err := i.idRepo.ListPermissions(c.Context(), c.Params("id"))
	if err != nil {
		return responseError(c, err)
	}

	return responseRecords(c, fiber.StatusOK, list)
}
