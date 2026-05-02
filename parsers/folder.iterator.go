package parsers

import (
	"log"
	"path/filepath"

	"github.com/rajatnai49/mentat/vault"
)

type ParserFunc func(string) (*vault.NoteTask, error)

func FolderIterator(dir string, parserFunc ParserFunc) []*vault.NoteTask {
	entries, err := filepath.Glob(dir + "/*.md")
	if err != nil {
		log.Fatal(err)
	}

	var nts []*vault.NoteTask

	for _, v := range entries {
		file_path := filepath.Join(dir, v)
		nt, err := parserFunc(file_path)
		if err != nil {
			log.Fatal(err)
			continue
		}
		if nt != nil && len(nt.Tasks) > 0 {
			nts = append(nts, nt)
		}
	}

	return nts
}
