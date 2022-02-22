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

func TestRootPath(t *testing.T) {
	m := apex.New(nil)
	require.Equal(t, "/", m.Path)
}

func TestSubGroupPath(t *testing.T) {
	g := apex.New(nil).NewGroup("/sub")

	require.NotNil(t, g)
	require.Equal(t, "/sub", g.Path)
}

func TestFullPath(t *testing.T) {
	g := apex.New(nil).NewGroup("/sub").NewGroup("/group")

	require.NotNil(t, g)
	require.Equal(t, "/sub/group", g.FullPath())
}
