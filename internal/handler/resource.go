package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kimnguyenlong/ketoz/internal/entity"
	"github.com/kimnguyenlong/ketoz/internal/repository"
)

type resource struct {
	rscRepo repository.Resource
}

func NewResource(rscRepo repository.Resource) *resource {
	return &resource{
		rscRepo: rscRepo,
	}
}

func (i *resource) RegisterRoutes(r fiber.Router) {
	g := r.Group("/resources")
	// Resources
	g.Post("/", i.Create)
	g.Get("/", i.List)
	g.Get("/:id", i.Get)
	// Children
	g.Get("/:id/children", i.ListChildren)
	g.Post("/:id/children", i.AddChild)
}

type CreateResourceRequest struct {
	Id string `json:"id"`
}

func (r *resource) Create(c *fiber.Ctx) error {
	req := new(CreateResourceRequest)
	if err := c.BodyParser(req); err != nil {
		return responseError(c, entity.NewInvalidParamsError(err.Error()))
	}

	if err := r.rscRepo.Create(c.Context(), &entity.Resource{
		Id: req.Id,
	}); err != nil {
		return responseError(c, err)
	}

	return responseNilData(c, fiber.StatusCreated)
}

func (r *resource) List(c *fiber.Ctx) error {
	list, err := r.rscRepo.List(c.Context())
	if err != nil {
		return responseError(c, err)
	}

	return responseRecords(c, fiber.StatusOK, list)
}

func (r *resource) Get(c *fiber.Ctx) error {
	rsc, err := r.rscRepo.Get(c.Context(), c.Params("id"))
	if err != nil {
		return responseError(c, err)
	}

	return responseData(c, fiber.StatusOK, rsc)
}

func (r *resource) ListChildren(c *fiber.Ctx) error {
	list, err := r.rscRepo.ListChildren(c.Context(), c.Params("id"))
	if err != nil {
		return responseError(c, err)
	}

	return responseRecords(c, fiber.StatusOK, list)
}

type AddChildResourceRequest struct {
	ChildId string `json:"child_id"`
}

func (r *resource) AddChild(c *fiber.Ctx) error {
	req := new(AddChildResourceRequest)
	if err := c.BodyParser(req); err != nil {
		return responseError(c, entity.NewInvalidParamsError(err.Error()))
	}

	if err := r.rscRepo.AddChild(c.Context(), c.Params("id"), req.ChildId); err != nil {
		return responseError(c, err)
	}

	return responseNilData(c, fiber.StatusCreated)
}
