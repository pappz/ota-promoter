package promoter

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

const (
	ignoreFile = ".promoterignore"
)

var (
	log = logrus.WithField("tag", "promoter")
)

type Promoter struct {
	*sync.Mutex
	promotedFolder string
	promotedFiles  []*File
	ignoredFiles   []string
	version        string
}

func NewPromoter(promotedFolder string) *Promoter {
	p := Promoter{
		Mutex:          &sync.Mutex{},
		promotedFolder: filepath.Clean(promotedFolder),
	}

	return &p
}

func (p *Promoter) ReadFiles() (err error) {
	p.Lock()
	defer p.Unlock()

	p.promotedFiles = p.promotedFiles[:0]

	p.ignoredFiles, err = p.readIgnoreList()
	if err != nil {
		return err
	}

	err = filepath.Walk(p.promotedFolder, p.walkFn)
	if err != nil {
		return err
	}

	p.version = p.calcVersion()
	for _, p := range p.promotedFiles {
		log.Debugf("file: %s - %s", p.Checksum, p.PromotedPath)
	}
	log.Debugf("promoted version is: %s", p.version)
	return nil
}

func (p *Promoter) PromotedFileByChecksum(s string) (*File, bool) {
	p.Lock()
	defer p.Unlock()
	for _, p := range p.promotedFiles {
		if p.Checksum == s {
			return p, true
		}
	}
	return nil, false
}

func (p *Promoter) PromotedFiles() []*File {
	p.Lock()
	defer p.Unlock()

	return p.promotedFiles
}

func (p *Promoter) Version() string {
	p.Lock()
	defer p.Unlock()
	return p.version
}

func (p *Promoter) readIgnoreList() ([]string, error) {
	ignoreFile := path.Join(p.promotedFolder, ignoreFile)
	ignoredFiles := make([]string, 0, 0)
	f, err := os.Stat(ignoreFile)
	if os.IsNotExist(err) {
		return ignoredFiles, nil
	}
	if f.IsDir() {
		return ignoredFiles, nil
	}

	file, err := os.Open(ignoreFile)
	if err != nil {
		return ignoredFiles, err
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
	return ignoredFiles, scanner.Err()
}

func (p *Promoter) calcVersion() string {
	hasher := sha1.New()
	for _, h := range p.promotedFiles {
		hasher.Write([]byte(h.Checksum))
	}
	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash[:])
}

func (p *Promoter) walkFn(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if path == p.promotedFolder {
		return nil
	}

	if info.IsDir() {
		return nil
	}

	if info.Mode()&os.ModeSymlink == os.ModeSymlink {
		return nil
	}

	promotedPath := p.getPromotedPath(path)
	if p.fileIsOnIgnoreList(promotedPath) {
		return nil
	}

	pf, err := newFile(path, promotedPath)
	if err != nil {
		return err
	}

	p.promotedFiles = append(p.promotedFiles, pf)
	return nil
}

func (p *Promoter) fileIsOnIgnoreList(name string) bool {
	for _, i := range p.ignoredFiles {
		regString := fmt.Sprintf("^%s$", i)
		match, _ := regexp.MatchString(regString, name)
		if match {
			return true
		}
	}
	return false
}

func (p *Promoter) getPromotedPath(path string) string {
	cp := strings.TrimPrefix(path, p.promotedFolder)
	return fmt.Sprintf("%s", cp[1:])
}
