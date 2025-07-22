package repository

import (
	"context"
	"fmt"

	"github.com/kimnguyenlong/ketoz/internal/entity"
	"github.com/kimnguyenlong/ketoz/pkg/keto"
	"github.com/kimnguyenlong/ketoz/pkg/util"
	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
)

type resource struct {
	keto *keto.Keto
}

func NewResource(keto *keto.Keto) Resource {
	return &resource{
		keto: keto,
	}
}

func (r *resource) List(ctx context.Context) ([]*entity.Resource, error) {
	req := &rts.ListRelationTuplesRequest{
		RelationQuery: &rts.RelationQuery{
			Namespace: util.StringPointer(keto.NamespaceResource.String()),
			Relation:  util.StringPointer(keto.RelationSelf.String()),
		},
	}
	res, err := r.keto.Read.ListRelationTuples(ctx, req)
	if err != nil {
		return nil, entity.NewInternalError(err.Error())
	}

	list := make([]*entity.Resource, 0, len(res.GetRelationTuples()))
	for _, r := range res.GetRelationTuples() {
		list = append(list, &entity.Resource{
			Id: r.GetObject(),
		})
	}

	return list, nil
}

func (r *resource) Get(ctx context.Context, id string) (*entity.Resource, error) {
	req := &rts.CheckRequest{
		Namespace: keto.NamespaceResource.String(),
		Object:    id,
		Relation:  keto.RelationSelf.String(),
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
		return nil, entity.NewNotFoundError(fmt.Sprintf("no resource with id %s", id))
	}

	return &entity.Resource{Id: id}, nil
}

func (r *resource) Create(ctx context.Context, rsc *entity.Resource) error {
	existing, err := r.Get(ctx, rsc.Id)
	if err != nil && !entity.IsNotFoundError(err) {
		return err
	}
	if existing != nil {
		return entity.NewInvalidParamsError(fmt.Sprintf("resource with id %s already exists", rsc.Id))
	}

	req := &rts.TransactRelationTuplesRequest{
		RelationTupleDeltas: []*rts.RelationTupleDelta{
			{
				Action: rts.RelationTupleDelta_ACTION_INSERT,
				RelationTuple: &rts.RelationTuple{
					Namespace: keto.NamespaceResource.String(),
					Object:    rsc.Id,
					Relation:  keto.RelationSelf.String(),
					Subject: &rts.Subject{
						Ref: &rts.Subject_Id{
							Id: rsc.Id,
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

func (r *resource) AddChild(ctx context.Context, parentId, childId string) error {
	req := &rts.TransactRelationTuplesRequest{
		RelationTupleDeltas: []*rts.RelationTupleDelta{
			{
				Action: rts.RelationTupleDelta_ACTION_INSERT,
				RelationTuple: &rts.RelationTuple{
					Namespace: keto.NamespaceResource.String(),
					Object:    childId,
					Relation:  keto.RelationParents.String(),
					Subject: &rts.Subject{
						Ref: &rts.Subject_Set{
							Set: &rts.SubjectSet{
								Namespace: keto.NamespaceResource.String(),
								Object:    parentId,
								Relation:  keto.RelationEmpty.String(),
							},
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

func (r *resource) ListChildren(ctx context.Context, id string) ([]*entity.Resource, error) {
	req := &rts.ListRelationTuplesRequest{
		RelationQuery: &rts.RelationQuery{
			Namespace: util.StringPointer(keto.NamespaceResource.String()),
			Relation:  util.StringPointer(keto.RelationParents.String()),
			Subject: &rts.Subject{
				Ref: &rts.Subject_Set{
					Set: &rts.SubjectSet{
						Namespace: keto.NamespaceResource.String(),
						Object:    id,
						Relation:  keto.RelationEmpty.String(),
					},
				},
			},
		},
	}
	res, err := r.keto.Read.ListRelationTuples(ctx, req)
	if err != nil {
		return nil, entity.NewInternalError(err.Error())
	}

	list := make([]*entity.Resource, 0, len(res.GetRelationTuples()))
	for _, r := range res.GetRelationTuples() {
		list = append(list, &entity.Resource{
			Id: r.GetObject(),
		})
	}

	return list, nil
}
