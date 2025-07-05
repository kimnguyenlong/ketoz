package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kimnguyenlong/ketoz/internal/entity"
	"github.com/kimnguyenlong/ketoz/internal/repository"
)

type role struct {
	roleRepo repository.Role
}

func NewRole(roleRepo repository.Role) *role {
	return &role{
		roleRepo: roleRepo,
	}
}

func (i *role) RegisterRoutes(r fiber.Router) {
	g := r.Group("/roles")
	// Roles
	g.Get("/", i.List)
	g.Get("/:id", i.Get)
	// Permissions
	g.Get("/:id/permissions", i.ListPermissions)
	g.Post("/:id/permissions", i.AddPermissions)
	// Members
	g.Get("/:id/members", i.ListMembers)
	g.Post("/:id/members", i.AddMembers)
}

func (i *role) List(c *fiber.Ctx) error {
	return nil
}

func (i *role) Get(c *fiber.Ctx) error {
	return nil
}

type AddPermissionsRequest struct {
	Permissions []*entity.Permission `json:"permissions"`
}

func (r *role) AddPermissions(c *fiber.Ctx) error {
	req := new(AddPermissionsRequest)
	if err := c.BodyParser(req); err != nil {
		return responseError(c, entity.NewInvalidParamsError(err.Error()))
	}

	if err := r.roleRepo.AddPermissions(c.Context(), c.Params("id"), req.Permissions); err != nil {
		return responseError(c, err)
	}

	return responseNilData(c, fiber.StatusCreated)
}

func (r *role) ListPermissions(c *fiber.Ctx) error {
	list, err := r.roleRepo.ListPermissions(c.Context(), c.Params("id"))
	if err != nil {
		return responseError(c, err)
	}

	return responseRecords(c, fiber.StatusOK, list)
}

type AddRoleMembersRequest struct {
	Identities []*entity.Identity `json:"identities"`
}

func (r *role) AddMembers(c *fiber.Ctx) error {
	req := new(AddRoleMembersRequest)
	if err := c.BodyParser(req); err != nil {
		return responseError(c, entity.NewInvalidParamsError(err.Error()))
	}

	if err := r.roleRepo.AddMembers(c.Context(), c.Params("id"), req.Identities); err != nil {
		return responseError(c, err)
	}

	return responseNilData(c, fiber.StatusCreated)
}

func (r *role) ListMembers(c *fiber.Ctx) error {
	list, err := r.roleRepo.ListMembers(c.Context(), c.Params("id"))
	if err != nil {
		return responseError(c, err)
	}

	return responseRecords(c, fiber.StatusOK, list)
}
