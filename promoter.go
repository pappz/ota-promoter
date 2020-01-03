package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	ignoreFile = ".promoterignore"
)

var (
	srcPath       = filepath.Clean("/tmp/test")
	promotedFiles []*PromotedFile
	ignoredFiles  []string
	version       string
)

type PromotedFile struct {
	CanonicalName string `json:"path"`
	Checksum      string `json:"checksum"`
	localPath     string
	size          int64
}

func newPromotedFile(path string) (*PromotedFile, error) {
	var err error
	var pf = &PromotedFile{
		CanonicalName: getCanonicalName(path),
		localPath:     path,
	}

	if err = pf.calcChecksum(); err != nil {
		return pf, err
	}

	if err = pf.setSizeOfFile(); err != nil {
		return pf, err
	}

	return pf, nil
}

func (pf *PromotedFile) calcChecksum() error {
	f, err := os.Open(pf.localPath)
	if err != nil {
		return err
	}
	defer f.Close()

	hasher := sha1.New()
	if _, err := io.Copy(hasher, f); err != nil {
		return err
	}

	hasher.Write([]byte(pf.CanonicalName))

	hash := hasher.Sum(nil)
	pf.Checksum = hex.EncodeToString(hash[:])
	return nil
}

func (pf *PromotedFile) setSizeOfFile() error {
	openFile, err := os.Open(pf.localPath)
	defer openFile.Close() //Close after function return
	if err != nil {
		return err
	}
	fileStat, _ := openFile.Stat()
	pf.size = fileStat.Size()
	return nil
}

func readFiles() error {
	promotedFiles = promotedFiles[:0]
	srcPath = filepath.Clean(promotedFolder)

	if err := readIgnoreList(); err != nil {
		return err
	}

	err := filepath.Walk(srcPath, processFile)
	if err != nil {
		return err
	}

	calcVersion()
	for _, p := range promotedFiles {
		log.Printf("file: %s - %s", p.Checksum, p.CanonicalName)
	}
	log.Printf("promoted version is: %s", version)
	return nil
}

func readIgnoreList() error {
	ignoreFile := path.Join(srcPath, ignoreFile)
	f, err := os.Stat(ignoreFile)
	if os.IsNotExist(err) {
		ignoredFiles = make([]string, 0, 0)
		return nil
	}
	if f.IsDir() {
		ignoredFiles = make([]string, 0, 0)
		return nil
	}

	file, err := os.Open(ignoreFile)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		txt := strings.TrimSpace(scanner.Text())
		if txt == "" {
			continue
		}
		ignoredFiles = append(ignoredFiles, txt)
	}
	return scanner.Err()
}

func getPromotedFileByChecksum(s string) (*PromotedFile, bool) {
	for _, p := range promotedFiles {
		if p.Checksum == s {
			return p, true
		}
	}
	return nil, false
}

func processFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if path == srcPath {
		return nil
	}

	if info.IsDir() {
		return nil
	}

	if info.Mode()&os.ModeSymlink == os.ModeSymlink {
		return nil
	}

	if fileIsOnIgnoreList(getCanonicalName(path)) {
		return nil
	}

	pf, err := newPromotedFile(path)
	if err != nil {
		return err
	}

	promotedFiles = append(promotedFiles, pf)
	return nil
}

func fileIsOnIgnoreList(name string) bool {
	for _, i := range ignoredFiles {
		regString := fmt.Sprintf("^%s$", i)
		match, _ := regexp.MatchString(regString, name)
		if match {
			return true
		}
	}
	return false
}

func getCanonicalName(p string) string {
	cp := strings.TrimPrefix(p, srcPath)
	return fmt.Sprintf("%s", cp[1:])
}

func calcVersion() {
	hasher := sha1.New()
	for _, h := range promotedFiles {
		hasher.Write([]byte(h.Checksum))
	}
	hash := hasher.Sum(nil)
	version = hex.EncodeToString(hash[:])
}
