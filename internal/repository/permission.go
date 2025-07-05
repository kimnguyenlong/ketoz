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
