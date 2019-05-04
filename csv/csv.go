package csv

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/knieriem/tblutil"

	"encoding/csv"
)

func init() {
	tblutil.Register(".csv", &tblutil.Impl{
		ReadAll:    ReadAll,
		ReadSheets: ReadSheets,
	})
}

func ReadAll(filename string) ([]tblutil.Table, error) {
	return ReadSheets(filename, "")
}

func ReadSheets(filename string, sheetNames ...string) ([]tblutil.Table, error) {
	switch len(sheetNames) {
	case 0:
		return nil, nil
	default:
		return nil, errors.New("only one sheet per CSV file supported")
	case 1:
	}
	opts := ""
	if i := strings.Index(filename, ".csv,"); i != -1 {
		opts = filename[i+5:]
		filename = filename[:i+4]
	}
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := csv.NewReader(f)
	if strings.HasPrefix(opts, "comma=") && len(opts) > 6 {
		r.Comma = rune(opts[6])
	}
	var t tblutil.Table
	t.Name = filepath.Base(filename)
	t.Name = strings.TrimSuffix(t.Name, filepath.Ext(t.Name))
	t.Rows, err = r.ReadAll()
	if err != nil {
		return nil, err
	}
	return []tblutil.Table{t}, nil
}
