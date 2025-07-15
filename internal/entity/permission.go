package entity

import "github.com/kimnguyenlong/ketoz/pkg/keto"

type Permission struct {
	ResourceId string      `json:"resource_id"`
	Action     keto.Action `json:"action"`
}
