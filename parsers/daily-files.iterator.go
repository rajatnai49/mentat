package parsers

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
)

type FileFunc func(string) error

func DailyFilesIterator(dir string, fn FileFunc) error {
	entries, err := filepath.Glob(dir + "/*.md")
	if err != nil {
		return err
	}

	re := regexp.MustCompile(`^\d{8}\.md$`)

	var errs []error

	for _, v := range entries {
		name := filepath.Base(v)

		if !re.MatchString(name) {
			continue
		}

		if err := fn(v); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", name, err))
			continue
		}
	}

	return errors.Join(errs...)
}
