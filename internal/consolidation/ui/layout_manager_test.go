package ui

import "testing"

func TestNewLayoutManager(t *testing.T) {
	lm := NewLayoutManager(100, 30)
	if lm == nil {
		t.Fatal("expected non-nil layout manager")
	}
	if lm.maxX != 100 || lm.maxY != 30 {
		t.Errorf("expected (100, 30), got (%d, %d)", lm.maxX, lm.maxY)
	}
}

func TestLayoutManagerMinSize(t *testing.T) {
	// Test that minimum size is enforced
	lm := NewLayoutManager(50, 20)
	if lm.maxX < 80 || lm.maxY < 24 {
		t.Errorf("expected minimum (80, 24), got (%d, %d)", lm.maxX, lm.maxY)
	}
}

func TestGetListPaneRect(t *testing.T) {
	lm := NewLayoutManager(100, 30)
	x, y, w, h := lm.GetListPaneRect()
	
	if x != 0 || y != 0 {
		t.Errorf("expected origin (0, 0), got (%d, %d)", x, y)
	}
	if w != lm.leftPaneW {
		t.Errorf("expected width %d, got %d", lm.leftPaneW, w)
	}
	if h <= 0 {
		t.Errorf("expected positive height, got %d", h)
	}
}

func TestGetDetailsPaneRect(t *testing.T) {
	lm := NewLayoutManager(100, 30)
	x, y, w, h := lm.GetDetailsPaneRect()
	
	if x <= 0 {
		t.Errorf("expected positive x, got %d", x)
	}
	if y != 0 {
		t.Errorf("expected y=0, got %d", y)
	}
	if w <= 0 {
		t.Errorf("expected positive width, got %d", w)
	}
	if h <= 0 {
		t.Errorf("expected positive height, got %d", h)
	}
}

func TestGetStatusBarRect(t *testing.T) {
	lm := NewLayoutManager(100, 30)
	x, y, w, h := lm.GetStatusBarRect()
	
	if x != 0 {
		t.Errorf("expected x=0, got %d", x)
	}
	if y <= 0 {
		t.Errorf("expected positive y, got %d", y)
	}
	if w != lm.maxX {
		t.Errorf("expected width %d, got %d", lm.maxX, w)
	}
	if h != lm.statusH {
		t.Errorf("expected height %d, got %d", lm.statusH, h)
	}
}

func TestLayoutManagerUpdateSize(t *testing.T) {
	lm := NewLayoutManager(100, 30)
	oldW := lm.leftPaneW
	
	lm.UpdateSize(200, 50)
	newW := lm.leftPaneW
	
	if lm.maxX != 200 || lm.maxY != 50 {
		t.Errorf("expected (200, 50), got (%d, %d)", lm.maxX, lm.maxY)
	}
	if oldW == newW {
		t.Errorf("expected pane width to change, got same %d", newW)
	}
}

func TestPaneAlignmentAndNoOverlap(t *testing.T) {
	lm := NewLayoutManager(100, 30)

	lx, _, lw, lh := lm.GetListPaneRect()
	dx, _, _, dh := lm.GetDetailsPaneRect()

	// List pane should come before details pane horizontally
	if lx+lw >= dx {
		t.Errorf("panes overlap: list ends at %d, details starts at %d", lx+lw, dx)
	}

	// Both should have same height (content height)
	if lh != dh {
		t.Errorf("expected same content height, got list=%d, details=%d", lh, dh)
	}
}
