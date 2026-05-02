package parsers

import "bytes"

func TaskScanner(data []byte, atEOF bool) (advance int, token []byte, err error){
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if i := bytes.Index(data[1:], []byte("\n- [")); i >= 0 {
		return i+1, data[:i+1], nil
	}

	if atEOF {
		return len(data), data, nil
	}

	return 0, nil, nil
}
