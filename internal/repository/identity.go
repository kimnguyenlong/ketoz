package repository

import (
	"context"
	"fmt"

	"github.com/kimnguyenlong/ketoz/internal/entity"
	"github.com/kimnguyenlong/ketoz/pkg/keto"
	"github.com/kimnguyenlong/ketoz/pkg/util"
	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
)

type identity struct {
	keto *keto.Keto
}

func NewIdentity(keto *keto.Keto) Identity {
	return &identity{
		keto: keto,
	}
}

func (i *identity) List(ctx context.Context) ([]*entity.Identity, error) {
	req := &rts.ListRelationTuplesRequest{
		RelationQuery: &rts.RelationQuery{
			Namespace: util.StringPointer(keto.NamespaceIdentity),
			Relation:  util.StringPointer(keto.RelationSelf),
		},
	}
	res, err := i.keto.Read.ListRelationTuples(ctx, req)
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

func (i *identity) Get(ctx context.Context, id string) (*entity.Identity, error) {
	req := &rts.CheckRequest{
		Namespace: keto.NamespaceIdentity,
		Object:    id,
		Relation:  keto.RelationSelf,
		Subject: &rts.Subject{
			Ref: &rts.Subject_Id{
				Id: id,
			},
		},
	}
	res, err := i.keto.Check.Check(ctx, req)
	if err != nil {
		return nil, entity.NewInternalError(err.Error())
	}

	if !res.GetAllowed() {
		return nil, entity.NewNotFoundError(fmt.Sprintf("no identity with id %s", id))
	}

	return &entity.Identity{Id: id}, nil
}

func (i *identity) Create(ctx context.Context, identity *entity.Identity) error {
	existing, err := i.Get(ctx, identity.Id)
	if err != nil && !entity.IsNotFoundError(err) {
		return err
	}
	if existing != nil {
		return entity.NewInvalidParamsError(fmt.Sprintf("identity with id %s already exists", identity.Id))
	}

	req := &rts.TransactRelationTuplesRequest{
		RelationTupleDeltas: []*rts.RelationTupleDelta{
			{
				Action: rts.RelationTupleDelta_ACTION_INSERT,
				RelationTuple: &rts.RelationTuple{
					Namespace: keto.NamespaceIdentity,
					Object:    identity.Id,
					Relation:  keto.RelationSelf,
					Subject: &rts.Subject{
						Ref: &rts.Subject_Id{
							Id: identity.Id,
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
