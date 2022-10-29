package promoter

import (
	"testing"
)

func Test_newFile(t *testing.T) {
	f, err := newFile("test/content.txt", "content.txt")
	if err != nil {
		t.Errorf("%s", err.Error())
	}

	if f.PromotedPath != "content.txt" {
		t.Errorf("invalid promoted path: %s", f.PromotedPath)
	}

	if f.Checksum != "0dca091da529abc1c507269acd06e619ed395c88" {
		t.Errorf("invalid checksum: %s", f.Checksum)
	}

	if f.Size != 7 {
		t.Errorf("invalid size: %d", f.Size)
	}
}

func Test_notExist(t *testing.T) {
	_, err := newFile("test/notexist.txt", "notexist.txt")
	if err == nil {
		t.Errorf("err in file exits validation: %s", err)
	}
}

func Test_checksum(t *testing.T) {
	f1, err := newFile("test/sub/content2.txt", "sub/content2.txt")
	if err != nil {
		t.Errorf("%s", err.Error())
	}

	f2, err := newFile("test/sub/content2.txt", "sub/same_content.txt")
	if err != nil {
		t.Errorf("%s", err.Error())
	}

	if f1.Checksum == f2.Checksum {
		t.Errorf("invalid checksum")
	}
}
