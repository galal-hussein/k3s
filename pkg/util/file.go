// +build !windows

package util

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

func SetFileModeForPath(name string, mode os.FileMode) error {
	return os.Chmod(name, mode)
}

func SetFileModeForFile(file *os.File, mode os.FileMode) error {
	return file.Chmod(mode)
}

// BackupDirWithRetention will move the dir to a backup dir
// and will keep only maxBackupRetention of dirs
func BackupDirWithRetention(dir string, maxBackupRetention int) (string, error) {
	backupDir := dir + "-backup-" + strconv.Itoa(int(time.Now().Unix()))
	if _, err := os.Stat(dir); err != nil {
		return "", nil
	}
	files, err := ioutil.ReadDir(filepath.Dir(dir))
	if err != nil {
		return "", err
	}
	sort.Slice(files, func(i,j int) bool{
		return files[i].ModTime().After(files[j].ModTime())
	})
	count := 0
	for _, f := range files {
		count++
		if strings.HasPrefix(f.Name(), filepath.Base(dir)+"-backup") && f.IsDir() {
			if count > maxBackupRetention {
				if err := os.RemoveAll(filepath.Join(filepath.Dir(dir), f.Name())); err != nil {
					return "", err
				}
			}
		}
	}
	// move the directory to a temp path
	if err := os.Rename(dir, backupDir); err != nil {
		return "", err
	}
	return backupDir, nil
}