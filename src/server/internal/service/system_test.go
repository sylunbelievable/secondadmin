package service

import (
	"testing"

	"github.com/sylunbelievable/secondadmin/server/internal/dto"
	"github.com/sylunbelievable/secondadmin/server/internal/entity"
)

func TestBuildMenuTree(t *testing.T) {
	items := []entity.Menu{
		{ID: 1, Name: "root"},
		{ID: 2, ParentID: 1, Name: "child"},
		{ID: 3, ParentID: 2, Name: "leaf"},
	}
	tree := buildMenuTree(items)
	if len(tree) != 1 || len(tree[0].Children) != 1 || tree[0].Children[0].Children[0].ID != 3 {
		t.Fatal("menu tree was not built correctly")
	}
}

func TestBuildMenuTreeExcludesButtonsFromNothing(t *testing.T) {
	items := []entity.Menu{{ID: 1, Type: "menu", Name: "page"}, {ID: 2, ParentID: 1, Type: "button", Name: "save"}}
	tree := buildMenuTree(items[:1])
	if len(tree) != 1 || len(tree[0].Children) != 0 {
		t.Fatal("route tree should contain only supplied route nodes")
	}
}

func TestMenuReachableRequiresActiveAncestors(t *testing.T) {
	root := entity.Menu{ID: 1, Status: 0, Visible: true}
	button := entity.Menu{ID: 2, ParentID: 1, Type: "button", Status: 1}
	if menuReachable(button, map[uint64]entity.Menu{1: root, 2: button}, false) {
		t.Fatal("button below disabled parent must not be reachable")
	}
}

func TestPruneEmptyDirectories(t *testing.T) {
	items := []dto.Menu{{ID: 1, Type: "directory"}, {ID: 2, Type: "menu"}}
	items = pruneEmptyDirectories(items)
	if len(items) != 1 || items[0].ID != 2 {
		t.Fatal("empty directory was not pruned")
	}
}
