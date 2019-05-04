package ods

import (
	"errors"
	"sync"

	"github.com/knieriem/tblutil"

	"github.com/knieriem/odf/ods"
)

func init() {
	tblutil.Register(".ods", &tblutil.Impl{
		ReadAll:    ReadAll,
		ReadSheets: ReadSheets,
	})
}

var lk sync.Mutex
var docMap map[string]*ods.Doc

func readDoc(filename string) (*ods.Doc, error) {
	lk.Lock()
	defer lk.Unlock()
	if docMap != nil {
		if d, ok := docMap[filename]; ok {
			return d, nil
		}
	}
	f, err := ods.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	d := new(ods.Doc)
	if err = f.ParseContent(d); err != nil {
		return nil, err
	}
	if docMap == nil {
		docMap = make(map[string]*ods.Doc, 8)
	}
	docMap[filename] = d
	return d, nil
}

func ReadAll(filename string) ([]tblutil.Table, error) {
	d, err := readDoc(filename)
	if err != nil {
		return nil, err
	}

	var tables = make([]tblutil.Table, len(d.Table))
	for i := range d.Table {
		dt := &d.Table[i]
		var t tblutil.Table
		t.Name = dt.Name
		t.Rows = dt.Strings()
		tables[i] = t
	}
	return tables, nil
}

func ReadSheets(filename string, sheetNames ...string) ([]tblutil.Table, error) {
	if sheetNames == nil {
		return nil, nil
	}

	d, err := readDoc(filename)
	if err != nil {
		return nil, err
	}
	dtables := make(map[string]*ods.Table, len(d.Table))
	for i := range d.Table {
		dtables[d.Table[i].Name] = &d.Table[i]
	}
	var tables = make([]tblutil.Table, len(sheetNames))
	for i, name := range sheetNames {
		dt, ok := dtables[name]
		if !ok {
			return nil, errors.New("sheet not found: " + name)
		}
		var t tblutil.Table
		t.Name = dt.Name
		t.Rows = dt.Strings()
		tables[i] = t
	}
	return tables, nil
}
