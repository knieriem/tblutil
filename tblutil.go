package tblutil

import (
	"io"
	"path/filepath"
	"strings"
)

type Table struct {
	Name string
	Rows [][]string
}

var implMap map[string]*Impl

type Impl struct {
	ReadAll    func(filename string) ([]Table, error)
	ReadSheets func(filename string, sheetNames ...string) ([]Table, error)
	ReadAtAll  func(_ io.ReaderAt, size int64) ([]Table, error)
}

func Register(name string, i *Impl) {
	if implMap == nil {
		implMap = make(map[string]*Impl, 8)
	}
	implMap[name] = i
}

func ReadFile(filename string, sheetNames ...string) ([]Table, error) {
	ext := filepath.Ext(filename)
	if i := strings.IndexByte(ext, ','); i != -1 {
		ext = ext[:i]
	}
	impl, ok := implMap[ext]
	if !ok {
		impl = tsvImpl
	}
	f := impl.ReadSheets
	if sheetNames == nil {
		f = func(filename string, _ ...string) ([]Table, error) { return impl.ReadAll(filename) }
	}
	tables, err := f(filename, sheetNames...)
	if err != nil {
		return nil, err
	}
	return cleanTables(tables), nil
}

func cleanTables(tables []Table) []Table {
	for i := range tables {
		tables[i].clean()
	}
	return tables
}

func (t *Table) clean() {
	regular := true
	lastNonEmpty := -1
	for _, r := range t.Rows {
		for i, cell := range r {
			cell = strings.TrimSpace(cell)
			r[i] = cell
			if cell != "" {
				if i > lastNonEmpty {
					lastNonEmpty = i
				}
			}
		}
		if lastNonEmpty+1 != len(r) {
			regular = false
		}
	}
	if regular {
		return
	}

	n := lastNonEmpty + 1
	tail := make([]string, n)
	for i, r := range t.Rows {
		delta := n - len(r)
		if delta > 0 {
			t.Rows[i] = append(r, tail[:delta]...)
		} else {
			t.Rows[i] = r[:n]
		}
	}
}
