package repository

import (
	"context"

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
		Relation:  entity.ActionToRelation[p.Action],
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
			Action:     entity.RelationToAction[r.GetRelation()],
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
		list = append(list, &entity.Identity{
			Id: r.GetObject(),
		})
	}

	return list, nil
}
