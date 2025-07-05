package repository

import (
	"context"

	"github.com/kimnguyenlong/ketoz/internal/entity"
)

type Identity interface {
	List(ctx context.Context, offset, limit int) ([]*entity.Identity, error)
	Get(ctx context.Context, id string) (*entity.Identity, error)
	Insert(ctx context.Context, id *entity.Identity) error

	AddChild(ctx context.Context, parentId, childId string) error
	ListChildren(ctx context.Context, id string) ([]*entity.Identity, error)
}

type Resource interface {
	AddChild(ctx context.Context, parentId, childId string) error
	ListChildren(ctx context.Context, id string) ([]*entity.Resource, error)
}

type Role interface {
	AddPermissions(ctx context.Context, id string, permissions []*entity.Permission) error
	ListPermissions(ctx context.Context, id string) ([]*entity.Permission, error)

	AddMembers(ctx context.Context, id string, identities []*entity.Identity) error
	ListMembers(ctx context.Context, id string) ([]*entity.Identity, error)
}

type Permission interface {
	IsPermitted(ctx context.Context, identityId, resourceId string, action entity.Action) (bool, error)
}
