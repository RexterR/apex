package apex

import "net/http"

// Apex is this package's main type. It contains the root route group.
type Apex struct {
	*Group
}

// New takes a Handler and initializes a new Apex instance with a root route
// group.
func New(handler Handler) *Apex {
	return &Apex{
		Group: &Group{Path: "/", handler: handler},
	}
}

func (m *Apex) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.handler.ServeHTTP(w, r)
}
