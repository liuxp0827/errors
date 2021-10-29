package errorctx

import "context"

type errorsKey struct {
}

// NewErrorsContext returns a new Context that carries value.
func NewErrorsContext(ctx context.Context, lang string) context.Context {
	return context.WithValue(ctx, errorsKey{}, lang)
}

// FromErrorsContext returns the Transport value stored in errorctx, if any.
func FromErrorsContext(ctx context.Context) (lang string, ok bool) {
	lang, ok = ctx.Value(errorsKey{}).(string)
	return
}
