package browsing

import (
	"testing"
	"github.com/haggishunk/filesprawl/internal/consolidation/flags"
)

func TestSummaryRendererCreation(t *testing.T) {
	sr := NewSummaryRenderer()

	if sr.presenceAnalyzer == nil {
		t.Errorf("expected presence analyzer to be initialized")
	}
}

func TestCountDecisions(t *testing.T) {
	sr := NewSummaryRenderer()

	flagMap := map[int]*flags.Flag{
		1: {FileID: 1, FlagType: flags.FlagTypeKeep},
		2: {FileID: 2, FlagType: flags.FlagTypeKeep},
		3: {FileID: 3, FlagType: flags.FlagTypeArchive},
		4: {FileID: 4, FlagType: flags.FlagTypeDelete},
		5: {FileID: 5, FlagType: flags.FlagTypeDelete},
		6: {FileID: 6, FlagType: flags.FlagTypeDelete},
	}

	keep, archive, delete := sr.countDecisions(flagMap)

	if keep != 2 {
		t.Errorf("expected 2 KEEP, got %d", keep)
	}
	if archive != 1 {
		t.Errorf("expected 1 ARCHIVE, got %d", archive)
	}
	if delete != 3 {
		t.Errorf("expected 3 DELETE, got %d", delete)
	}
}

func TestFilterByType(t *testing.T) {
	sr := NewSummaryRenderer()

	flagMap := map[int]*flags.Flag{
		1: {FileID: 1, FlagType: flags.FlagTypeKeep},
		2: {FileID: 2, FlagType: flags.FlagTypeArchive},
		3: {FileID: 3, FlagType: flags.FlagTypeDelete},
	}

	keepFlags := sr.filterByType(flagMap, flags.FlagTypeKeep)
	if len(keepFlags) != 1 {
		t.Errorf("expected 1 KEEP flag, got %d", len(keepFlags))
	}

	archiveFlags := sr.filterByType(flagMap, flags.FlagTypeArchive)
	if len(archiveFlags) != 1 {
		t.Errorf("expected 1 ARCHIVE flag, got %d", len(archiveFlags))
	}

	deleteFlags := sr.filterByType(flagMap, flags.FlagTypeDelete)
	if len(deleteFlags) != 1 {
		t.Errorf("expected 1 DELETE flag, got %d", len(deleteFlags))
	}
}

func TestRenderCompactSummary(t *testing.T) {
	sr := NewSummaryRenderer()

	session := &BrowsingSession{
		Flags: map[int]*flags.Flag{
			1: {FileID: 1, FlagType: flags.FlagTypeKeep},
			2: {FileID: 2, FlagType: flags.FlagTypeArchive},
			3: {FileID: 3, FlagType: flags.FlagTypeDelete},
		},
	}

	compact := sr.RenderCompactSummary(session)

	if compact == "" {
		t.Errorf("expected non-empty compact summary")
	}
	if !contains(compact, "K:") {
		t.Errorf("expected K: in compact summary")
	}
	if !contains(compact, "A:") {
		t.Errorf("expected A: in compact summary")
	}
	if !contains(compact, "D:") {
		t.Errorf("expected D: in compact summary")
	}
}

func TestRenderEmptySummary(t *testing.T) {
	sr := NewSummaryRenderer()

	session := &BrowsingSession{
		Flags: make(map[int]*flags.Flag),
	}

	summary := sr.RenderDecisionSummary(session)

	if summary == "" {
		t.Errorf("expected non-empty summary")
	}
	if !contains(summary, "DECISION SUMMARY") {
		t.Errorf("expected title in summary")
	}
}

func TestRenderFullSummary(t *testing.T) {
	sr := NewSummaryRenderer()

	reason := "duplicate"
	archiveDest := "s3-backup"

	session := &BrowsingSession{
		Flags: map[int]*flags.Flag{
			1: {FileID: 1, FlagType: flags.FlagTypeKeep},
			2: {FileID: 2, FlagType: flags.FlagTypeArchive, ArchiveDest: &archiveDest},
			3: {FileID: 3, FlagType: flags.FlagTypeDelete, Reason: &reason},
		},
	}

	summary := sr.RenderDecisionSummary(session)

	if !contains(summary, "KEEP") {
		t.Errorf("expected KEEP section in summary")
	}
	if !contains(summary, "ARCHIVE") {
		t.Errorf("expected ARCHIVE section in summary")
	}
	if !contains(summary, "DELETE") {
		t.Errorf("expected DELETE section in summary")
	}
	if !contains(summary, "s3-backup") {
		t.Errorf("expected archive destination in summary")
	}
}

func TestRenderWarnings(t *testing.T) {
	sr := NewSummaryRenderer()

	session := &BrowsingSession{
		Flags: map[int]*flags.Flag{
			1: {FileID: 1, FlagType: flags.FlagTypeDelete},
			2: {FileID: 2, FlagType: flags.FlagTypeDelete},
		},
	}

	summary := sr.RenderDecisionSummary(session)

	if !contains(summary, "WARNING") {
		t.Errorf("expected WARNING in summary when there are deletes")
	}
	if !contains(summary, "permanent") {
		t.Errorf("expected warning about permanent deletion")
	}
}
