package keto

type Namespace string
type Relation string
type Permission string
type Action string

func (n Namespace) String() string {
	return string(n)
}

func (r Relation) String() string {
	return string(r)
}

func (p Permission) String() string {
	return string(p)
}

func (a Action) String() string {
	return string(a)
}

const (
	NamespaceIdentity Namespace = "Identity"
	NamespaceResource Namespace = "Resource"
)

const (
	RelationEmpty               Relation = ""
	RelationSelf                Relation = "self"
	RelationChildren            Relation = "children"
	RelationParents             Relation = "parents"
	RelationChildCreators       Relation = "child_creators"
	RelationViewers             Relation = "viewers"
	RelationEditors             Relation = "editors"
	RelationOwners              Relation = "owners"
	RelationDeniedChildCreators Relation = "denied_child_creators"
	RelationDeniedViewers       Relation = "denied_viewers"
	RelationDeniedEditors       Relation = "denied_editors"
)

const (
	PermissionChildCreators Permission = "child_creators"
	PermissionViewers       Permission = "viewers"
	PermissionEditors       Permission = "editors"
	PermissionOwners        Permission = "owners"
)

const (
	ActionCreateChild Action = "create_child"
	ActionView        Action = "view"
	ActionEdit        Action = "edit"
	ActionDelete      Action = "delete"
)

var PermissionToRelation = map[Permission]Relation{
	PermissionChildCreators: RelationChildCreators,
	PermissionViewers:       RelationViewers,
	PermissionEditors:       RelationEditors,
	PermissionOwners:        RelationOwners,
}

var RelationToPermission = map[Relation]Permission{
	RelationChildCreators: PermissionChildCreators,
	RelationViewers:       PermissionViewers,
	RelationEditors:       PermissionEditors,
	RelationOwners:        PermissionOwners,
}

var DeniedPermissionToRelation = map[Permission]Relation{
	PermissionChildCreators: RelationDeniedChildCreators,
	PermissionViewers:       RelationDeniedViewers,
	PermissionEditors:       RelationDeniedEditors,
}
