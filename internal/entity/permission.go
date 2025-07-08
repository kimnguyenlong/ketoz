package entity

import "github.com/kimnguyenlong/ketoz/pkg/keto"

type Action string

const (
	ActionCreate Action = "create"
	ActionView   Action = "view"
	ActionUpdate Action = "update"
	ActionDelete Action = "delete"
)

type Permission struct {
	ResourceId string `json:"resource_id"`
	Action     Action `json:"action"`
}

var ActionToRelation = map[Action]string{
	ActionView: keto.RelationViewers,
}

var RelationToAction = map[string]Action{
	keto.RelationViewers: ActionView,
}

var DeniedActionToRelation = map[Action]string{
	ActionView: keto.RelationDeniedViewers,
}
