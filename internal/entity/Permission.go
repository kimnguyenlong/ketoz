package entity

import "github.com/kimnguyenlong/ketoz/pkg/keto"

type Permission struct {
	IdentityID string          `json:"identity_id"`
	ResourceID string          `json:"resource_id"`
	Permission keto.Permission `json:"permission"`
}
