package tblutil

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

var tsvImpl = &Impl{
	ReadAll: func(filename string) ([]Table, error) {
		return readTSV(filename)
	},
	ReadSheets: readTSV,
}

func readTSV(filename string, _ ...string) ([]Table, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	var t Table
	t.Name = filepath.Base(filename)
	t.Name = strings.TrimSuffix(t.Name, filepath.Ext(t.Name))
	for s.Scan() {
		row := strings.Split(s.Text(), "\t")
		t.Rows = append(t.Rows, row)
	}
	return []Table{t}, nil
}
