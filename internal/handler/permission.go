package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kimnguyenlong/ketoz/internal/entity"
	"github.com/kimnguyenlong/ketoz/internal/repository"
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
	g.Get("/check", p.Check)
}

type CheckRequest struct {
	IdentityId string        `query:"identity_id"`
	ResourceId string        `query:"resource_id"`
	Action     entity.Action `query:"action"`
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
