import { Namespace, SubjectSet, Context } from "@ory/keto-namespace-types"

class Principal implements Namespace {
    related: {
        members: Principal[]
    }
}

class Role implements Namespace {
    related: {
        members: (Principal | SubjectSet<Principal, "members">)[]
    }
}

class Resource implements Namespace {
    related: {
        parents: Resource[]
        viewers: SubjectSet<Role, "members">[]
        denied_viewers: (Principal | SubjectSet<Principal, "members">)[]
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

