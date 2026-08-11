//go:build expr_noreflectmethod

package expr

// Eval is a panic-only stub used when building with the expr_noreflectmethod
// tag. The real Eval relies on runtime dispatch on the env, which requires
// reflect-based method resolution. Use Compile + Run instead.
func Eval(input string, env any) (any, error) {
	panic("expr.Eval is not available with the expr_noreflectmethod build tag")
}
