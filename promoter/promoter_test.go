package promoter

import (
	"testing"
)

const (
	testFolder = "test"
)

func TestPromoter_ReadFiles(t *testing.T) {
	p := NewPromoter(testFolder)
	err := p.ReadFiles()
	if err != nil {
		t.Errorf("failed to read promoted files")
	}
}

func TestPromoter_Version(t *testing.T) {
	p := NewPromoter(testFolder)
	err := p.ReadFiles()
	if err != nil {
		t.Errorf("failed to read promoted files")
	}

	if p.Version() != "44452108317c0ea6cb3054a69351c32fc36a8663" {
		t.Errorf("invalid version: %s", p.Version())
	}
}

func TestPromoter_PromotedFiles(t *testing.T) {
	p := NewPromoter(testFolder)
	err := p.ReadFiles()
	if err != nil {
		t.Errorf("failed to read promoted files")
	}

	files := p.PromotedFiles()
	if len(files) != 2 {
		t.Errorf("invalid promoted list: %d", len(files))
	}
}

func TestPromoter_PromotedFileByChecksum(t *testing.T) {
	p := NewPromoter(testFolder)
	err := p.ReadFiles()
	if err != nil {
		t.Errorf("failed to read promoted files")
	}

	contentFileChksum := "0dca091da529abc1c507269acd06e619ed395c88"
	f, ok := p.PromotedFileByChecksum(contentFileChksum)
	if !ok {
		t.Fatalf("file not found by checksum: %s", contentFileChksum)
	}

	if f.PromotedPath != "content.txt" {
		t.Errorf("invalid promoted path: %s", f.PromotedPath)
	}

	_, ok = p.PromotedFileByChecksum("invalidchecksum")
	if ok {
		t.Fatalf("invalid checksum match")
	}
}
