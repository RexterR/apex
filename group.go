package apex

import "path"

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

// FullPath returns the group's full path in the group tree (as opposed to this
// group's sub-path)
func (g *Group) FullPath() string {
	if g.parent == nil {
		return g.Path
	}
	return path.Join(g.parent.FullPath(), g.Path)
}
