package repository

import (
	"context"
	"fmt"
	"slices"

	"github.com/kimnguyenlong/ketoz/internal/entity"
	"github.com/kimnguyenlong/ketoz/pkg/keto"
	"github.com/kimnguyenlong/ketoz/pkg/util"
	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
)

type resource struct {
	list []*entity.Resource
	keto *keto.Keto
}

func NewResource(keto *keto.Keto) Resource {
	return &resource{
		keto: keto,
	}
}

func (r *resource) List(ctx context.Context, offset, limit int) ([]*entity.Resource, error) {
	return r.list, nil
}

func (r *resource) Get(ctx context.Context, id string) (*entity.Resource, error) {
	idx := slices.IndexFunc(r.list, func(item *entity.Resource) bool {
		return item.Id == id
	})
	if idx == -1 {
		return nil, fmt.Errorf("resource with id %s not found", id)
	}

	return r.list[idx], nil
}

func (r *resource) Insert(ctx context.Context, resource *entity.Resource) error {
	idx := slices.IndexFunc(r.list, func(item *entity.Resource) bool {
		return item.Id == resource.Id
	})
	if idx != -1 {
		return fmt.Errorf("resource with id %s already exists", resource.Id)
	}
	r.list = append(r.list, resource)
	return nil
}

func (r *resource) AddChild(ctx context.Context, parentId, childId string) error {
	req := &rts.TransactRelationTuplesRequest{
		RelationTupleDeltas: []*rts.RelationTupleDelta{
			{
				Action: rts.RelationTupleDelta_ACTION_INSERT,
				RelationTuple: &rts.RelationTuple{
					Namespace: keto.NamespaceResource,
					Object:    childId,
					Relation:  keto.RelationParents,
					Subject: &rts.Subject{
						Ref: &rts.Subject_Set{
							Set: &rts.SubjectSet{
								Namespace: keto.NamespaceResource,
								Object:    parentId,
								Relation:  keto.RelationEmpty,
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
			Namespace: util.StringPointer(keto.NamespaceResource),
			Relation:  util.StringPointer(keto.RelationParents),
			Subject: &rts.Subject{
				Ref: &rts.Subject_Set{
					Set: &rts.SubjectSet{
						Namespace: keto.NamespaceResource,
						Object:    id,
						Relation:  keto.RelationEmpty,
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
