package repository

import (
	"context"

	"github.com/kimnguyenlong/ketoz/internal/entity"
	"github.com/kimnguyenlong/ketoz/pkg/keto"
	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
)

type permission struct {
	keto *keto.Keto
}

func NewPermission(keto *keto.Keto) Permission {
	return &permission{
		keto: keto,
	}
}

func (p *permission) GrantPermission(ctx context.Context, identityId, resourceId string, permission keto.Permission) error {
	req := &rts.TransactRelationTuplesRequest{
		RelationTupleDeltas: []*rts.RelationTupleDelta{
			{ // resource -> identity
				Action: rts.RelationTupleDelta_ACTION_INSERT,
				RelationTuple: &rts.RelationTuple{
					Namespace: keto.NamespaceResource.String(),
					Object:    resourceId,
					Relation:  keto.PermissionToRelation[permission].String(),
					Subject: &rts.Subject{
						Ref: &rts.Subject_Set{
							Set: &rts.SubjectSet{
								Namespace: keto.NamespaceIdentity.String(),
								Object:    identityId,
								Relation:  keto.RelationEmpty.String(),
							},
						},
					},
				},
			},
			{ // resource -> children of identity
				Action: rts.RelationTupleDelta_ACTION_INSERT,
				RelationTuple: &rts.RelationTuple{
					Namespace: keto.NamespaceResource.String(),
					Object:    resourceId,
					Relation:  keto.PermissionToRelation[permission].String(),
					Subject: &rts.Subject{
						Ref: &rts.Subject_Set{
							Set: &rts.SubjectSet{
								Namespace: keto.NamespaceIdentity.String(),
								Object:    identityId,
								Relation:  keto.RelationChildren.String(),
							},
						},
					},
				},
			},
		},
	}
	if _, err := p.keto.Write.TransactRelationTuples(ctx, req); err != nil {
		return entity.NewInternalError(err.Error())
	}

	return nil
}

func (p *permission) DenyPermission(ctx context.Context, identityId, resourceId string, permission keto.Permission) error {
	relations := make([]*rts.RelationTupleDelta, 0, 2)
	id := &rts.RelationTupleDelta{ // for the identity
		Action: rts.RelationTupleDelta_ACTION_INSERT,
		RelationTuple: &rts.RelationTuple{
			Namespace: keto.NamespaceResource.String(),
			Object:    resourceId,
			Relation:  keto.DeniedPermissionToRelation[permission].String(),
			Subject: &rts.Subject{
				Ref: &rts.Subject_Set{
					Set: &rts.SubjectSet{
						Namespace: keto.NamespaceIdentity.String(),
						Object:    identityId,
						Relation:  keto.RelationEmpty.String(),
					},
				},
			},
		},
	}
	children := &rts.RelationTupleDelta{ // for the children of the identity
		Action: rts.RelationTupleDelta_ACTION_INSERT,
		RelationTuple: &rts.RelationTuple{
			Namespace: keto.NamespaceResource.String(),
			Object:    resourceId,
			Relation:  keto.DeniedPermissionToRelation[permission].String(),
			Subject: &rts.Subject{
				Ref: &rts.Subject_Set{
					Set: &rts.SubjectSet{
						Namespace: keto.NamespaceIdentity.String(),
						Object:    identityId,
						Relation:  keto.RelationChildren.String(),
					},
				},
			},
		},
	}
	relations = append(relations, id, children)
	req := &rts.TransactRelationTuplesRequest{
		RelationTupleDeltas: relations,
	}
	if _, err := p.keto.Write.TransactRelationTuples(ctx, req); err != nil {
		return entity.NewInternalError(err.Error())
	}

	return nil
}

func (p *permission) IsPermitted(ctx context.Context, identityId, resourceId string, action keto.Action) (bool, error) {
	req := &rts.CheckRequest{
		Namespace: keto.NamespaceResource.String(),
		Object:    resourceId,
		Relation:  action.String(),
		Subject: &rts.Subject{
			Ref: &rts.Subject_Set{
				Set: &rts.SubjectSet{
					Namespace: keto.NamespaceIdentity.String(),
					Object:    identityId,
					Relation:  keto.RelationEmpty.String(),
				},
			},
		},
	}
	res, err := p.keto.Check.Check(ctx, req)
	if err != nil {
		return false, entity.NewInternalError(err.Error())
	}

	return res.GetAllowed(), nil
}
