package parsers

import (
	"os"
	"path/filepath"
	"strings"
)

func DailyFileCleaner(pathOfFile string) (bool, error) {
	isEmpty, err := isEmptyFile(pathOfFile)
	if isEmpty {
		err = os.Remove(pathOfFile)
		if err != nil {
			return false, err
		}
		return true, nil
	}

	noteTask, err := DailyFileParser(pathOfFile)
	if err != nil {
		return false, err
	}

	isAnyPending := false

	for _, t := range noteTask.Tasks {
		if !t.Done {
			isAnyPending = true
			break
		}
	}

	if !isAnyPending {
		dir := filepath.Dir(pathOfFile)
		base := filepath.Base(pathOfFile)
		ext := filepath.Ext(pathOfFile)

		justName := strings.TrimSuffix(base, ext)
		newName := justName + "-X" + ext

		newPathOfFile := filepath.Join(dir, newName)

		err := os.Rename(pathOfFile, newPathOfFile)
		if err != nil {
			return false, err
		}
	}

	return !isAnyPending, nil
}

func isEmptyFile(filepath string) (bool, error) {
	stat, err := os.Stat(filepath)
	if err != nil {
		return false, err
	}
	return stat.Size() == 0, nil
}
