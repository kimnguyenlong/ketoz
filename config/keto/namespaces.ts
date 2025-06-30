import { Namespace, SubjectSet, Context } from "@ory/keto-namespace-types"

class Principal implements Namespace {
    related: {
        members: Principal[]
    }
}

class Role implements Namespace {
    related: {
        members: (Principal | SubjectSet<Principal,"members">)[]
    }
}

class Resource implements Namespace {
    related: {
        parents: Resource[]
        viewers: SubjectSet<Role, "members">[]
        denied_viewers: (Principal | SubjectSet<Principal,"members">)[]
        restricted_viewers: (Principal | SubjectSet<Principal,"members">)[]
    }

    permits = {
        view: (ctx: Context): boolean =>
            ( // deny policy
                !this.permits.denied_view(ctx)
            ) &&
            ( // allow policy
                this.related.viewers.includes(ctx.subject) ||
                this.related.parents.traverse((p) => p.permits.view(ctx))
            ),
        denied_view: (ctx: Context): boolean =>
            this.related.denied_viewers.includes(ctx.subject) ||
            this.related.parents.traverse((p) => p.permits.denied_view(ctx))
    }
}

