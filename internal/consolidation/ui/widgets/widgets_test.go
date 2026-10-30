package widgets

import (
	"testing"
)

func TestNewBaseListWidget(t *testing.T) {
	widget := NewBaseListWidget("Test Widget")

	if widget.GetTitle() != "Test Widget" {
		t.Errorf("expected title 'Test Widget', got %s", widget.GetTitle())
	}

	if widget.GetWidget() == nil {
		t.Errorf("expected non-nil widget")
	}

	if !widget.IsEmpty() {
		t.Errorf("expected empty widget initially")
	}
}

func TestSetAndGetRows(t *testing.T) {
	widget := NewBaseListWidget("Test")
	rows := []string{"Row 1", "Row 2", "Row 3"}

	err := widget.SetRows(rows)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if widget.Count() != 3 {
		t.Errorf("expected 3 rows, got %d", widget.Count())
	}

	if widget.GetRows()[0] != "Row 1" {
		t.Errorf("expected 'Row 1', got %s", widget.GetRows()[0])
	}
}

func TestSelectedIndex(t *testing.T) {
	widget := NewBaseListWidget("Test")
	rows := []string{"Row 1", "Row 2", "Row 3"}
	widget.SetRows(rows)

	if widget.GetSelectedIndex() != 0 {
		t.Errorf("expected initial index 0, got %d", widget.GetSelectedIndex())
	}

	widget.SetSelectedIndex(2)
	if widget.GetSelectedIndex() != 2 {
		t.Errorf("expected index 2, got %d", widget.GetSelectedIndex())
	}

	if widget.GetSelectedRow() != "Row 3" {
		t.Errorf("expected 'Row 3', got %s", widget.GetSelectedRow())
	}
}

func TestScrolling(t *testing.T) {
	widget := NewBaseListWidget("Test")
	rows := []string{"Row 1", "Row 2", "Row 3"}
	widget.SetRows(rows)

	// Test scroll down
	widget.ScrollDown()
	if widget.GetSelectedIndex() != 1 {
		t.Errorf("expected index 1 after scroll down, got %d", widget.GetSelectedIndex())
	}

	widget.ScrollDown()
	if widget.GetSelectedIndex() != 2 {
		t.Errorf("expected index 2 after scroll down, got %d", widget.GetSelectedIndex())
	}

	// Should not scroll past end
	widget.ScrollDown()
	if widget.GetSelectedIndex() != 2 {
		t.Errorf("expected index to stay at 2, got %d", widget.GetSelectedIndex())
	}

	// Test scroll up
	widget.ScrollUp()
	if widget.GetSelectedIndex() != 1 {
		t.Errorf("expected index 1 after scroll up, got %d", widget.GetSelectedIndex())
	}
}

func TestClear(t *testing.T) {
	widget := NewBaseListWidget("Test")
	rows := []string{"Row 1", "Row 2"}
	widget.SetRows(rows)

	if widget.Count() != 2 {
		t.Errorf("expected 2 rows, got %d", widget.Count())
	}

	widget.Clear()
	if widget.Count() != 0 {
		t.Errorf("expected 0 rows after clear, got %d", widget.Count())
	}

	if !widget.IsEmpty() {
		t.Errorf("expected empty widget after clear")
	}
}

func TestSetTitle(t *testing.T) {
	widget := NewBaseListWidget("Initial")

	if widget.GetTitle() != "Initial" {
		t.Errorf("expected 'Initial', got %s", widget.GetTitle())
	}

	widget.SetTitle("Updated")
	if widget.GetTitle() != "Updated" {
		t.Errorf("expected 'Updated', got %s", widget.GetTitle())
	}
}

func TestSetRect(t *testing.T) {
	widget := NewBaseListWidget("Test")
	widget.SetRect(10, 20, 80, 40)

	x, y, w, h := widget.GetRect()
	if x != 10 || y != 20 || w != 80 || h != 40 {
		t.Errorf("expected rect (10,20,80,40), got (%d,%d,%d,%d)", x, y, w, h)
	}
}

func TestGetSelectedRowOnEmpty(t *testing.T) {
	widget := NewBaseListWidget("Test")

	row := widget.GetSelectedRow()
	if row != "" {
		t.Errorf("expected empty string for empty widget, got %s", row)
	}
}
