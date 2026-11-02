package widgets

import (
	"fmt"

	"github.com/metaspartan/gotui/v5/widgets"
)

// BannerInfo holds information to display in the banner
type BannerInfo struct {
	Version       string
	ReviewMode    bool
	ExecuteMode   bool
	Remotes       []string
	TotalDuplicates int
	FlaggedCount   int
}

// BannerWidget displays a top banner with essential information
type BannerWidget struct {
	paragraph *widgets.Paragraph
	info      *BannerInfo
}

// NewBannerWidget creates a new banner widget
func NewBannerWidget(info *BannerInfo) *BannerWidget {
	widget := &BannerWidget{
		paragraph: widgets.NewParagraph(),
		info:      info,
	}
	
	if widget.paragraph != nil {
		widget.updateBannerText()
	}
	
	return widget
}

// GetWidget returns the underlying gotui Paragraph widget
func (w *BannerWidget) GetWidget() interface{} {
	if w == nil || w.paragraph == nil {
		return nil
	}
	return w.paragraph
}

// updateBannerText generates the banner text
func (w *BannerWidget) updateBannerText() {
	if w.paragraph == nil || w.info == nil {
		return
	}

	// Mode indicator
	modeStr := "REVIEW MODE"
	if w.info.ExecuteMode {
		modeStr = "EXECUTE MODE"
	}

	// Remotes summary
	remotesStr := "no remotes"
	if len(w.info.Remotes) > 0 {
		remotesStr = fmt.Sprintf("%d remote(s): %v", len(w.info.Remotes), w.info.Remotes)
	}

	// Status summary
	statusStr := fmt.Sprintf("Duplicates: %d | Flagged: %d", 
		w.info.TotalDuplicates, w.info.FlaggedCount)

	// Build banner text
	text := fmt.Sprintf(
		"filesprawl v%s | %s | %s | %s | [?] Help",
		w.info.Version,
		modeStr,
		remotesStr,
		statusStr,
	)

	w.paragraph.Text = text
}

// SetRect sets the drawing rectangle
func (w *BannerWidget) SetRect(x, y, width, height int) {
	if w == nil || w.paragraph == nil {
		return
	}
	w.paragraph.SetRect(x, y, width, height)
}

// UpdateInfo updates the banner with new information
func (w *BannerWidget) UpdateInfo(info *BannerInfo) {
	if w == nil {
		return
	}
	w.info = info
	w.updateBannerText()
}

// UpdateFlagCount updates just the flag count
func (w *BannerWidget) UpdateFlagCount(count int) {
	if w == nil || w.info == nil {
		return
	}
	w.info.FlaggedCount = count
	w.updateBannerText()
}
