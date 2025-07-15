package keto

type Action string

const (
	ActionCreate Action = "create"
	ActionView   Action = "view"
	ActionUpdate Action = "update"
	ActionDelete Action = "delete"
)

const (
	NamespaceIdentity string = "Identity"
	NamespaceRole     string = "Role"
	NamespaceResource string = "Resource"
)

const (
	RelationChildren      string = "children"
	RelationMembers       string = "members"
	RelationParents       string = "parents"
	RelationViewers       string = "viewers"
	RelationSelf          string = "self"
	RelationDeniedViewers string = "denied_viewers"
	RelationEmpty         string = ""
)

var ActionToRelation = map[Action]string{
	ActionView: RelationViewers,
}

var RelationToAction = map[string]Action{
	RelationViewers: ActionView,
}

var DeniedActionToRelation = map[Action]string{
	ActionView: RelationDeniedViewers,
}
