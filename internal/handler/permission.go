package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kimnguyenlong/ketoz/internal/entity"
	"github.com/kimnguyenlong/ketoz/internal/repository"
	"github.com/kimnguyenlong/ketoz/pkg/keto"
)

type permission struct {
	pmRepo repository.Permission
}

func NewPermission(pmRepo repository.Permission) *permission {
	return &permission{
		pmRepo: pmRepo,
	}
}

func (p *permission) RegisterRoutes(r fiber.Router) {
	g := r.Group("/permissions")
	g.Post("/granted", p.GrantPermission)
	g.Delete("/granted", p.RevokePermission)
	g.Post("/denied", p.DenyPermission)
	g.Delete("/denied", p.DeleteDeniedPermission)
	g.Get("/check", p.Check)
}

type GrantPermissionRequest struct {
	IdentityId string          `json:"identity_id"`
	ResourceId string          `json:"resource_id"`
	Permission keto.Permission `json:"permission"`
}

func (p *permission) GrantPermission(c *fiber.Ctx) error {
	req := new(GrantPermissionRequest)
	if err := c.BodyParser(req); err != nil {
		return responseError(c, entity.NewInvalidParamsError(err.Error()))
	}

	if err := p.pmRepo.GrantPermission(c.Context(), req.IdentityId, req.ResourceId, req.Permission); err != nil {
		return responseError(c, err)
	}

	return responseNilData(c, fiber.StatusOK)
}

type RevokePermissionRequest struct {
	IdentityId string          `json:"identity_id"`
	ResourceId string          `json:"resource_id"`
	Permission keto.Permission `json:"permission"`
}

func (p *permission) RevokePermission(c *fiber.Ctx) error {
	req := new(RevokePermissionRequest)
	if err := c.BodyParser(req); err != nil {
		return responseError(c, entity.NewInvalidParamsError(err.Error()))
	}

	if err := p.pmRepo.RevokePermission(c.Context(), req.IdentityId, req.ResourceId, req.Permission); err != nil {
		return responseError(c, err)
	}

	return responseNilData(c, fiber.StatusOK)
}

type DenyPermissionRequest struct {
	IdentityId string          `json:"identity_id"`
	ResourceId string          `json:"resource_id"`
	Permission keto.Permission `json:"permission"`
}

func (p *permission) DenyPermission(c *fiber.Ctx) error {
	req := new(DenyPermissionRequest)
	if err := c.BodyParser(req); err != nil {
		return responseError(c, entity.NewInvalidParamsError(err.Error()))
	}

	if err := p.pmRepo.DenyPermission(c.Context(), req.IdentityId, req.ResourceId, req.Permission); err != nil {
		return responseError(c, err)
	}

	return responseNilData(c, fiber.StatusOK)
}

type DeleteDeniedPermissionRequest struct {
	IdentityId string          `json:"identity_id"`
	ResourceId string          `json:"resource_id"`
	Permission keto.Permission `json:"permission"`
}

func (p *permission) DeleteDeniedPermission(c *fiber.Ctx) error {
	req := new(DeleteDeniedPermissionRequest)
	if err := c.BodyParser(req); err != nil {
		return responseError(c, entity.NewInvalidParamsError(err.Error()))
	}

	if err := p.pmRepo.DeleteDeniedPermission(c.Context(), req.IdentityId, req.ResourceId, req.Permission); err != nil {
		return responseError(c, err)
	}

	return responseNilData(c, fiber.StatusOK)
}

type CheckRequest struct {
	IdentityId string      `query:"identity_id"`
	ResourceId string      `query:"resource_id"`
	Action     keto.Action `query:"action"`
}

type CheckResponse struct {
	IsPermitted bool `json:"is_permitted"`
}

func (p *permission) Check(c *fiber.Ctx) error {
	req := new(CheckRequest)
	if err := c.QueryParser(req); err != nil {
		return responseError(c, entity.NewInvalidParamsError(err.Error()))
	}

	isPermitted, err := p.pmRepo.IsPermitted(c.Context(), req.IdentityId, req.ResourceId, req.Action)
	if err != nil {
		return responseError(c, err)
	}

	return responseData(c, fiber.StatusOK, &CheckResponse{
		IsPermitted: isPermitted,
	})
}
