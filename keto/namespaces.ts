import { Namespace, SubjectSet, Context } from "@ory/keto-namespace-types"

class Identity implements Namespace {
    related: {
        children: Identity[]
    }
}

class Role implements Namespace {
    related: {
        members: (Identity | SubjectSet<Identity, "children">)[]
    }
}

class Resource implements Namespace {
    related: {
        parents: Resource[]
        viewers: SubjectSet<Role, "members">[]
        denied_viewers: (Identity | SubjectSet<Identity, "children">)[]
    }

    permits = {
        view: (ctx: Context): boolean =>
            !this.permits.deny_view(ctx) &&
            this.permits.allow_view(ctx),
        allow_view: (ctx: Context): boolean =>
            this.related.viewers.includes(ctx.subject) ||
            this.related.parents.traverse((p) => p.permits.view(ctx)),
        deny_view: (ctx: Context): boolean =>
            this.related.denied_viewers.includes(ctx.subject) ||
            this.related.parents.traverse((p) => p.permits.deny_view(ctx))
    }
}

