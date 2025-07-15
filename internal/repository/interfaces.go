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

type Role interface {
	Create(ctx context.Context, id *entity.Role) error
	List(ctx context.Context) ([]*entity.Role, error)
	Get(ctx context.Context, id string) (*entity.Role, error)

	AddPermissions(ctx context.Context, id string, permissions []*entity.Permission) error
	ListPermissions(ctx context.Context, id string) ([]*entity.Permission, error)

	AddMembers(ctx context.Context, id string, identities []*entity.Identity) error
	ListMembers(ctx context.Context, id string) ([]*entity.Identity, error)
}

type Permission interface {
	IsPermitted(ctx context.Context, identityId, resourceId string, action keto.Action) (bool, error)
	AddDeniedPermission(ctx context.Context, identityId, resourceId string, action keto.Action) error
}
