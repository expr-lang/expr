//go:build !expr_noreflectmethod

package nature

// MethodByName looks up a method on a Nature by name. It transitively reaches
// reflect.Type.Method via the methodset cache.
func (n *Nature) MethodByName(c *Cache, name string) (Nature, bool) {
	if s := n.getMethodset(c); s != nil {
		if m := s.method(c, name); m != nil {
			return m.nature, true
		}
	}
	return Nature{}, false
}
