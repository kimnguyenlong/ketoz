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

type identity struct {
	list []*entity.Identity
	keto *keto.Keto
}

func NewIdentity(keto *keto.Keto) Identity {
	return &identity{
		keto: keto,
	}
}

func (i *identity) List(ctx context.Context, offset, limit int) ([]*entity.Identity, error) {
	return i.list, nil
}

func (i *identity) Get(ctx context.Context, id string) (*entity.Identity, error) {
	idx := slices.IndexFunc(i.list, func(item *entity.Identity) bool {
		return item.Id == id
	})
	if idx == -1 {
		return nil, fmt.Errorf("identity with id %s not found", id)
	}

	return i.list[idx], nil
}

func (i *identity) Insert(ctx context.Context, identity *entity.Identity) error {
	idx := slices.IndexFunc(i.list, func(item *entity.Identity) bool {
		return item.Id == identity.Id
	})
	if idx != -1 {
		return fmt.Errorf("identity with id %s already exists", identity.Id)
	}
	i.list = append(i.list, identity)
	return nil
}

func (i *identity) AddChild(ctx context.Context, parentId, childId string) error {
	req := &rts.TransactRelationTuplesRequest{
		RelationTupleDeltas: []*rts.RelationTupleDelta{
			{
				Action: rts.RelationTupleDelta_ACTION_INSERT,
				RelationTuple: &rts.RelationTuple{
					Namespace: keto.NamespaceIdentity,
					Object:    parentId,
					Relation:  keto.RelationChildren,
					Subject: &rts.Subject{
						Ref: &rts.Subject_Set{
							Set: &rts.SubjectSet{
								Namespace: keto.NamespaceIdentity,
								Object:    childId,
								Relation:  keto.RelationEmpty,
							},
						},
					},
				},
			},
		},
	}
	if _, err := i.keto.Write.TransactRelationTuples(ctx, req); err != nil {
		return entity.NewInternalError(err.Error())
	}

	return nil
}

func (i *identity) ListChildren(ctx context.Context, id string) ([]*entity.Identity, error) {
	req := &rts.ListRelationTuplesRequest{
		RelationQuery: &rts.RelationQuery{
			Namespace: util.StringPointer(keto.NamespaceIdentity),
			Object:    util.StringPointer(id),
			Relation:  util.StringPointer(keto.RelationChildren),
		},
	}
	res, err := i.keto.Read.ListRelationTuples(ctx, req)
	if err != nil {
		return nil, entity.NewInternalError(err.Error())
	}

	list := make([]*entity.Identity, 0, len(res.GetRelationTuples()))
	for _, r := range res.GetRelationTuples() {
		list = append(list, &entity.Identity{
			Id: r.GetSubject().GetSet().GetObject(),
		})
	}

	return list, nil
}
