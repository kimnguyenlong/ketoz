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

func (p *permission) IsPermitted(ctx context.Context, identityId, resourceId string, action entity.Action) (bool, error) {
	req := &rts.CheckRequest{
		Namespace: keto.NamespaceResource,
		Object:    resourceId,
		Relation:  string(action),
		Subject: &rts.Subject{
			Ref: &rts.Subject_Set{
				Set: &rts.SubjectSet{
					Namespace: keto.NamespaceIdentity,
					Object:    identityId,
					Relation:  keto.RelationEmpty,
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

func (p *permission) AddDeniedPermission(ctx context.Context, identityId, resourceId string, action entity.Action) error {
	relations := make([]*rts.RelationTupleDelta, 0, 2)
	id := &rts.RelationTupleDelta{ // for the identity
		Action: rts.RelationTupleDelta_ACTION_INSERT,
		RelationTuple: &rts.RelationTuple{
			Namespace: keto.NamespaceResource,
			Object:    resourceId,
			Relation:  entity.DeniedActionToRelation[action],
			Subject: &rts.Subject{
				Ref: &rts.Subject_Set{
					Set: &rts.SubjectSet{
						Namespace: keto.NamespaceIdentity,
						Object:    identityId,
						Relation:  keto.RelationEmpty,
					},
				},
			},
		},
	}
	children := &rts.RelationTupleDelta{ // for the children of the identity
		Action: rts.RelationTupleDelta_ACTION_INSERT,
		RelationTuple: &rts.RelationTuple{
			Namespace: keto.NamespaceResource,
			Object:    resourceId,
			Relation:  entity.DeniedActionToRelation[action],
			Subject: &rts.Subject{
				Ref: &rts.Subject_Set{
					Set: &rts.SubjectSet{
						Namespace: keto.NamespaceIdentity,
						Object:    identityId,
						Relation:  keto.RelationChildren,
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
