package apex_test

import (
	"testing"

	"github.com/RexterR/apex"
	"github.com/stretchr/testify/require"
)

func TestNewSubGroup(t *testing.T) {
	m := apex.New(nil)
	g := m.NewGroup("/test")

	require.NotNil(t, g)
	require.Equal(t, m.Group, g.Parent())
}
