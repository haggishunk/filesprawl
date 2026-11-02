package modals

import "testing"

func TestNewKeepModal(t *testing.T) {
	modal := NewKeepModal()
	if modal == nil {
		t.Fatal("expected non-nil KeepModal")
	}
	if modal.title != "Flag as KEEP?" {
		t.Errorf("expected title 'Flag as KEEP?', got %q", modal.title)
	}
}

func TestNewArchiveModal(t *testing.T) {
	modal := NewArchiveModal()
	if modal == nil {
		t.Fatal("expected non-nil ArchiveModal")
	}
	if modal.title != "Archive File" {
		t.Errorf("expected title 'Archive File', got %q", modal.title)
	}
}

func TestNewDeleteModal_Normal(t *testing.T) {
	modal := NewDeleteModal(false)
	if modal == nil {
		t.Fatal("expected non-nil DeleteModal")
	}
	if modal.title != "Delete File" {
		t.Errorf("expected title 'Delete File', got %q", modal.title)
	}
	if modal.isCritical {
		t.Errorf("expected isCritical=false")
	}
}

func TestNewDeleteModal_Critical(t *testing.T) {
	modal := NewDeleteModal(true)
	if modal == nil {
		t.Fatal("expected non-nil DeleteModal")
	}
	if !modal.isCritical {
		t.Errorf("expected isCritical=true")
	}
	if modal.message == "" {
		t.Errorf("expected non-empty warning message for critical delete")
	}
}

func TestNewSkipModal(t *testing.T) {
	modal := NewSkipModal()
	if modal == nil {
		t.Fatal("expected non-nil SkipModal")
	}
	if modal.title != "Remove Flag?" {
		t.Errorf("expected title 'Remove Flag?', got %q", modal.title)
	}
}

func TestNewQuitModal(t *testing.T) {
	modal := NewQuitModal(5, 3, 2)
	if modal == nil {
		t.Fatal("expected non-nil QuitModal")
	}
	if modal.keepCount != 5 {
		t.Errorf("expected keepCount=5, got %d", modal.keepCount)
	}
	if modal.archiveCount != 3 {
		t.Errorf("expected archiveCount=3, got %d", modal.archiveCount)
	}
	if modal.deleteCount != 2 {
		t.Errorf("expected deleteCount=2, got %d", modal.deleteCount)
	}
}

func TestNewErrorModal(t *testing.T) {
	modal := NewErrorModal("Test Error", "Something went wrong")
	if modal == nil {
		t.Fatal("expected non-nil ErrorModal")
	}
	if modal.title != "Test Error" {
		t.Errorf("expected title 'Test Error', got %q", modal.title)
	}
	if modal.message != "Something went wrong" {
		t.Errorf("expected message 'Something went wrong', got %q", modal.message)
	}
}

func TestModalResult_Confirmed(t *testing.T) {
	result := ModalResult{Confirmed: true, Value: "test"}
	if !result.Confirmed {
		t.Errorf("expected Confirmed=true")
	}
	if result.Value != "test" {
		t.Errorf("expected Value='test', got %q", result.Value)
	}
}

func TestModalResult_Cancelled(t *testing.T) {
	result := ModalResult{Confirmed: false, Error: nil}
	if result.Confirmed {
		t.Errorf("expected Confirmed=false")
	}
}

func TestKeepModalCancel(t *testing.T) {
	modal := NewKeepModal()
	modal.Cancel()
	// Should not panic and should be callable
}

func TestArchiveModalCancel(t *testing.T) {
	modal := NewArchiveModal()
	modal.Cancel()
	// Should not panic and should be callable
}
