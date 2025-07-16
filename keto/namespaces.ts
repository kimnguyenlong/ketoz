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
        viewers: (Identity | SubjectSet<Identity, "children">)[];
        denied_viewers: (Identity | SubjectSet<Identity, "children">)[];
    };

    permits = {
        view: (ctx: Context): boolean =>
            !this.permits.deny_view(ctx) &&
            this.permits.allow_view(ctx),
        allow_view: (ctx: Context): boolean =>
            this.related.viewers.includes(ctx.subject) ||
            this.related.parents.traverse((p) => p.permits.allow_view(ctx)),
        deny_view: (ctx: Context): boolean =>
            this.related.denied_viewers.includes(ctx.subject) ||
            this.related.parents.traverse((p) => p.permits.deny_view(ctx)),
    };
}
