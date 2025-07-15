package repository

import (
	"context"
	"fmt"

	"github.com/kimnguyenlong/ketoz/internal/entity"
	"github.com/kimnguyenlong/ketoz/pkg/keto"
	"github.com/kimnguyenlong/ketoz/pkg/util"
	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
)

type role struct {
	keto *keto.Keto
}

func NewRole(keto *keto.Keto) Role {
	return &role{
		keto: keto,
	}
}

func permissionToRelation(roleId string, p *entity.Permission) *rts.RelationTuple {
	if p == nil {
		return nil
	}
	return &rts.RelationTuple{
		Namespace: keto.NamespaceResource,
		Object:    p.ResourceId,
		Relation:  keto.ActionToRelation[p.Action],
		Subject: &rts.Subject{
			Ref: &rts.Subject_Set{
				Set: &rts.SubjectSet{
					Namespace: keto.NamespaceRole,
					Object:    roleId,
					Relation:  keto.RelationMembers,
				},
			},
		},
	}
}

func (r *role) List(ctx context.Context) ([]*entity.Role, error) {
	req := &rts.ListRelationTuplesRequest{
		RelationQuery: &rts.RelationQuery{
			Namespace: util.StringPointer(keto.NamespaceRole),
			Relation:  util.StringPointer(keto.RelationSelf),
		},
	}
	res, err := r.keto.Read.ListRelationTuples(ctx, req)
	if err != nil {
		return nil, entity.NewInternalError(err.Error())
	}

	list := make([]*entity.Role, 0, len(res.GetRelationTuples()))
	for _, r := range res.GetRelationTuples() {
		list = append(list, &entity.Role{
			Id: r.GetObject(),
		})
	}

	return list, nil
}

func (r *role) Get(ctx context.Context, id string) (*entity.Role, error) {
	req := &rts.CheckRequest{
		Namespace: keto.NamespaceRole,
		Object:    id,
		Relation:  keto.RelationSelf,
		Subject: &rts.Subject{
			Ref: &rts.Subject_Id{
				Id: id,
			},
		},
	}
	res, err := r.keto.Check.Check(ctx, req)
	if err != nil {
		return nil, entity.NewInternalError(err.Error())
	}

	if !res.GetAllowed() {
		return nil, entity.NewNotFoundError(fmt.Sprintf("no role with id %s", id))
	}

	return &entity.Role{Id: id}, nil
}

func (r *role) Create(ctx context.Context, role *entity.Role) error {
	existing, err := r.Get(ctx, role.Id)
	if err != nil && !entity.IsNotFoundError(err) {
		return err
	}
	if existing != nil {
		return entity.NewInvalidParamsError(fmt.Sprintf("role with id %s already exists", role.Id))
	}

	req := &rts.TransactRelationTuplesRequest{
		RelationTupleDeltas: []*rts.RelationTupleDelta{
			{
				Action: rts.RelationTupleDelta_ACTION_INSERT,
				RelationTuple: &rts.RelationTuple{
					Namespace: keto.NamespaceRole,
					Object:    role.Id,
					Relation:  keto.RelationSelf,
					Subject: &rts.Subject{
						Ref: &rts.Subject_Id{
							Id: role.Id,
						},
					},
				},
			},
		},
	}
	if _, err := r.keto.Write.TransactRelationTuples(ctx, req); err != nil {
		return entity.NewInternalError(err.Error())
	}

	return nil
}

func (r *role) AddPermissions(ctx context.Context, id string, permissions []*entity.Permission) error {
	relations := make([]*rts.RelationTupleDelta, 0, len(permissions))
	for _, p := range permissions {
		relations = append(relations, &rts.RelationTupleDelta{
			Action:        rts.RelationTupleDelta_ACTION_INSERT,
			RelationTuple: permissionToRelation(id, p),
		})
	}
	req := &rts.TransactRelationTuplesRequest{
		RelationTupleDeltas: relations,
	}
	if _, err := r.keto.Write.TransactRelationTuples(ctx, req); err != nil {
		return entity.NewInternalError(err.Error())
	}

	return nil
}

func (r *role) ListPermissions(ctx context.Context, id string) ([]*entity.Permission, error) {
	req := &rts.ListRelationTuplesRequest{
		RelationQuery: &rts.RelationQuery{
			Namespace: util.StringPointer(keto.NamespaceResource),
			Subject: &rts.Subject{
				Ref: &rts.Subject_Set{
					Set: &rts.SubjectSet{
						Namespace: keto.NamespaceRole,
						Object:    id,
						Relation:  keto.RelationMembers,
					},
				},
			},
		},
	}
	res, err := r.keto.Read.ListRelationTuples(ctx, req)
	if err != nil {
		return nil, entity.NewInternalError(err.Error())
	}

	list := make([]*entity.Permission, 0, len(res.GetRelationTuples()))
	for _, r := range res.GetRelationTuples() {
		list = append(list, &entity.Permission{
			ResourceId: r.GetObject(),
			Action:     keto.RelationToAction[r.GetRelation()],
		})
	}

	return list, nil
}

func (r *role) AddMembers(ctx context.Context, id string, identities []*entity.Identity) error {
	relations := make([]*rts.RelationTupleDelta, 0, len(identities))
	for _, i := range identities {
		base := &rts.RelationTupleDelta{ // for the identity
			Action: rts.RelationTupleDelta_ACTION_INSERT,
			RelationTuple: &rts.RelationTuple{
				Namespace: keto.NamespaceRole,
				Object:    id,
				Relation:  keto.RelationMembers,
				Subject: &rts.Subject{
					Ref: &rts.Subject_Set{
						Set: &rts.SubjectSet{
							Namespace: keto.NamespaceIdentity,
							Object:    i.Id,
							Relation:  keto.RelationEmpty,
						},
					},
				},
			},
		}
		sp := &rts.RelationTupleDelta{ // for the children of the identity
			Action: rts.RelationTupleDelta_ACTION_INSERT,
			RelationTuple: &rts.RelationTuple{
				Namespace: keto.NamespaceRole,
				Object:    id,
				Relation:  keto.RelationMembers,
				Subject: &rts.Subject{
					Ref: &rts.Subject_Set{
						Set: &rts.SubjectSet{
							Namespace: keto.NamespaceIdentity,
							Object:    i.Id,
							Relation:  keto.RelationChildren,
						},
					},
				},
			},
		}
		relations = append(relations, base, sp)
	}
	req := &rts.TransactRelationTuplesRequest{
		RelationTupleDeltas: relations,
	}
	if _, err := r.keto.Write.TransactRelationTuples(ctx, req); err != nil {
		return entity.NewInternalError(err.Error())
	}

	return nil
}

func (r *role) ListMembers(ctx context.Context, id string) ([]*entity.Identity, error) {
	req := &rts.ListRelationTuplesRequest{
		RelationQuery: &rts.RelationQuery{
			Namespace: util.StringPointer(keto.NamespaceRole),
			Object:    util.StringPointer(id),
			Relation:  util.StringPointer(keto.RelationMembers),
		},
	}
	res, err := r.keto.Read.ListRelationTuples(ctx, req)
	if err != nil {
		return nil, entity.NewInternalError(err.Error())
	}

	list := make([]*entity.Identity, 0, len(res.GetRelationTuples()))
	for _, r := range res.GetRelationTuples() {
		if r.GetSubject().GetSet().GetRelation() == keto.RelationEmpty {
			list = append(list, &entity.Identity{
				Id: r.GetSubject().GetSet().GetObject(),
			})
		}
	}

	return list, nil
}
