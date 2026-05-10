package parsers

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	"github.com/rajatnai49/mentat/vault"
)

var reTitle = regexp.MustCompile(`]\s+(.*?)\n`)
var reFile = regexp.MustCompile(`\[\[(.*?)\]\]`)
var reTag = regexp.MustCompile(`[#@]([^\s#@]+)`)

func DailyFileParser(filename string) (*vault.NoteTask, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var nt vault.NoteTask

	info, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}
	nt.Date = info.ModTime().UTC()
	nt.FilePath = filename

	scanner := bufio.NewScanner(f)
	scanner.Split(TaskScanner)

	scanner.Buffer(make([]byte, 0, 1024), 1024*1024)

	for scanner.Scan() {
		block := strings.TrimSpace(scanner.Text())
		lines := strings.Split(block, "\n")

		var t vault.Task

		t.Type = vault.Daily

		l1 := lines[0]

		if strings.Contains(l1, "[x]") {
			t.Done = true
		}

		splits := strings.SplitN(l1, "]", 2)
		if len(splits) < 2 {
			continue
		}

		t.Title = cleanLine(strings.TrimSpace(splits[1]))

		var desc []string

		for _, l := range lines[1:] {
			for _, m := range reTag.FindAllStringSubmatch(l, -1) {
				t.Tags = append(t.Tags, m[1])
			}
			for _, m := range reFile.FindAllStringSubmatch(l, -1) {
				t.LinkedNotes = append(t.LinkedNotes, m[1])
			}

			if clean := cleanLine(l); clean != "" {
				desc = append(desc, l)
			}
		}

		t.Description = strings.Join(desc, "\n")

		nt.Tasks = append(nt.Tasks, t)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &nt, err
}

func cleanLine(line string) string {
	clean := line

	clean = reFile.ReplaceAllString(clean, "")
	clean = reTag.ReplaceAllString(clean, "")
	clean = strings.TrimSpace(clean)

	return clean
}
