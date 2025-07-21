package repository

import (
	"context"

	"github.com/kimnguyenlong/ketoz/internal/entity"
	"github.com/kimnguyenlong/ketoz/pkg/keto"
)

type Identity interface {
	Create(ctx context.Context, id *entity.Identity) error
	List(ctx context.Context) ([]*entity.Identity, error)
	Get(ctx context.Context, id string) (*entity.Identity, error)

	AddChild(ctx context.Context, parentId, childId string) error
	ListChildren(ctx context.Context, id string) ([]*entity.Identity, error)
}

type Resource interface {
	Create(ctx context.Context, id *entity.Resource) error
	List(ctx context.Context) ([]*entity.Resource, error)
	Get(ctx context.Context, id string) (*entity.Resource, error)

	AddChild(ctx context.Context, parentId, childId string) error
	ListChildren(ctx context.Context, id string) ([]*entity.Resource, error)
}

type Permission interface {
	GrantPermission(ctx context.Context, identityId, resourceId string, permission keto.Permission) error
	RevokePermission(ctx context.Context, identityId, resourceId string, permission keto.Permission) error
	DenyPermission(ctx context.Context, identityId, resourceId string, permission keto.Permission) error
	DeleteDeniedPermission(ctx context.Context, identityId, resourceId string, permission keto.Permission) error
	IsPermitted(ctx context.Context, identityId, resourceId string, action keto.Action) (bool, error)
}
