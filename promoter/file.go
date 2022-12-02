package promoter

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"os"
)

// File describe the information of a promoted file
type File struct {
	PromotedPath string `json:"path"`
	Checksum     string `json:"checksum"`
	LocalPath    string `json:"-"`
	Size         int64  `json:"-"`
}

func newFile(localPath string, promotedPath string) (*File, error) {
	var err error
	var pf = &File{
		PromotedPath: promotedPath,
		LocalPath:    localPath,
	}

	if err = pf.calcChecksum(); err != nil {
		return pf, err
	}

	if err = pf.setSizeOfFile(); err != nil {
		return pf, err
	}

	return pf, nil
}

func (pf *File) calcChecksum() error {
	f, err := os.Open(pf.LocalPath)
	if err != nil {
		return err
	}
	defer f.Close()
	hasher := sha1.New()
	if _, err := io.Copy(hasher, f); err != nil {
		return err
	}

	hasher.Write([]byte(pf.PromotedPath))

	hash := hasher.Sum(nil)
	pf.Checksum = hex.EncodeToString(hash[:])
	return nil
}

func (pf *File) setSizeOfFile() error {
	openFile, err := os.Open(pf.LocalPath)
	defer openFile.Close()
	if err != nil {
		return err
	}
	fileStat, _ := openFile.Stat()
	pf.Size = fileStat.Size()
	return nil
}
