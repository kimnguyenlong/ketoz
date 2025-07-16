import { Context, Namespace, SubjectSet } from "@ory/keto-namespace-types";

class Identity implements Namespace {
    related: {
        self: Identity[];
        children: Identity[];
    };
}

class Resource implements Namespace {
    related: {
        self: Resource[];
        parents: Resource[];
        child_creators: (Identity | SubjectSet<Identity, "children">)[]; // identities allowed to create children resource
        viewers: (Identity | SubjectSet<Identity, "children">)[];
        editors: (Identity | SubjectSet<Identity, "children">)[];
        owners: (Identity | SubjectSet<Identity, "children">)[];
        denied_child_creators: (Identity | SubjectSet<Identity, "children">)[];
        denied_editors: (Identity | SubjectSet<Identity, "children">)[];
        denied_viewers: (Identity | SubjectSet<Identity, "children">)[];
    };

    permits = {
        create_child: (ctx: Context): boolean =>
            !this.permits.deny_create_child(ctx) &&
            this.permits.allow_create_child(ctx),

        view: (ctx: Context): boolean =>
            !this.permits.deny_view(ctx) &&
            this.permits.allow_view(ctx),

        edit: (ctx: Context): boolean =>
            !this.permits.deny_edit(ctx) &&
            this.permits.allow_edit(ctx),

        delete: (ctx: Context): boolean =>
            this.related.owners.includes(ctx.subject),

        allow_create_child: (ctx: Context): boolean =>
            this.related.child_creators.includes(ctx.subject) || // subject is child creator
            this.related.owners.includes(ctx.subject) || // subject is owner
            this.related.parents.traverse((p) =>
                p.permits.allow_create_child(ctx)
            ), // subject is allowed to create child for parent resource

        allow_view: (ctx: Context): boolean =>
            this.related.viewers.includes(ctx.subject) || // subject is viewer
            this.related.editors.includes(ctx.subject) || // subject is editor
            this.related.owners.includes(ctx.subject) || // subject is owner
            this.related.parents.traverse((p) => p.permits.allow_view(ctx)), // subject is allowed to view parent resource

        allow_edit: (ctx: Context): boolean =>
            this.related.editors.includes(ctx.subject) || // subject is editor
            this.related.owners.includes(ctx.subject) || // subject is owner
            this.related.parents.traverse((p) => p.permits.allow_edit(ctx)), // subject is allowed to edit parent resource

        deny_create_child: (ctx: Context): boolean =>
            this.related.denied_child_creators.includes(ctx.subject) || // subject is in denied child creators
            this.related.parents.traverse((p) =>
                p.permits.deny_create_child(ctx)
            ), // subject is denied to create child for parent resource

        deny_view: (ctx: Context): boolean =>
            this.related.denied_viewers.includes(ctx.subject) || // subject is in denied viewers
            this.related.parents.traverse((p) => p.permits.deny_view(ctx)), // subject is denied to view parent resource

        deny_edit: (ctx: Context): boolean =>
            this.related.denied_editors.includes(ctx.subject) || // subject is in denied editors
            this.related.denied_viewers.includes(ctx.subject) || // subject is in denied viewers
            this.related.parents.traverse((p) => p.permits.deny_edit(ctx)), // subject is denied to edit parent resource
    };
}
