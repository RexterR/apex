package apex

// Group represents a route group that shares a middleware chain and a common
// path.
type Group struct {
	Path    string
	handler Handler
	parent  *Group
}

// NewGroup creates a new subgroup of group g.
func (g *Group) NewGroup(path string) *Group {
	return &Group{
		Path:    path,
		handler: g.handler,
		parent:  g,
	}
}

// Parent gets the group's parent in the group tree.
func (g *Group) Parent() *Group {
	return g.parent
}
