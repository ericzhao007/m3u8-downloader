package spiders

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/duke-git/lancet/v2/fileutil"
)

type DiskStoreEngine struct {
	workDir string
}

func NewDiskStoreEngine(workDir string) *DiskStoreEngine {
	if exists := fileutil.IsExist(workDir); !exists {
		if err := fileutil.CreateDir(workDir); err != nil {
			return nil
		}
	}
	return &DiskStoreEngine{
		workDir: workDir,
	}
}

func (s *DiskStoreEngine) Load(key string) (value []byte, exists bool) {
	datafile := filepath.Join(s.workDir, key)
	if exists := fileutil.IsExist(datafile); !exists {
		return nil, false
	}
	data, err := ioutil.ReadFile(datafile)
	if err != nil {
		return nil, false
	}
	return data, true

}

func (s *DiskStoreEngine) Store(key string, value []byte) error {
	datafile := filepath.Join(s.workDir, key)
	return ioutil.WriteFile(datafile, value, os.ModePerm)
}

func (s *DiskStoreEngine) Clear() error {
	os.RemoveAll(s.workDir)
	return nil
}
