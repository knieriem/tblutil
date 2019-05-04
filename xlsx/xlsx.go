package xlsx

import (
	"errors"
	"strings"

	"github.com/knieriem/tblutil"

	"github.com/tealeg/xlsx"
)

func init() {
	tblutil.Register(".xlsx", &tblutil.Impl{
		ReadAll:    ReadAll,
		ReadSheets: ReadSheets,
	})
}

func ReadAll(filename string) ([]tblutil.Table, error) {
	xlFile, err := xlsx.OpenFile(filename)
	if err != nil {
		return nil, err
	}
	sheets := xlFile.Sheets
	var tables = make([]tblutil.Table, len(sheets))
	for i, sheet := range sheets {
		var t tblutil.Table
		readSheet(&t, sheet)
		tables[i] = t
	}
	return tables, nil
}

func ReadSheets(filename string, sheetNames ...string) ([]tblutil.Table, error) {
	xlFile, err := xlsx.OpenFile(filename)
	if err != nil {
		return nil, err
	}
	if sheetNames == nil {
		return nil, nil
	}
	var tables = make([]tblutil.Table, len(sheetNames))
	for i, name := range sheetNames {
		var t tblutil.Table
		sheet, ok := xlFile.Sheet[name]
		if !ok {
			return nil, errors.New("sheet not found: " + name)
		}
		readSheet(&t, sheet)
		tables[i] = t
	}
	return tables, nil
}

func readSheet(t *tblutil.Table, sheet *xlsx.Sheet) {
	t.Name = sheet.Name
	rows := sheet.Rows
	t.Rows = make([][]string, 0, len(rows))
	for _, row := range rows {
		var r = make([]string, 0, len(row.Cells))
		for _, cell := range row.Cells {
			r = append(r, strings.TrimSpace(cell.String()))
		}
		t.Rows = append(t.Rows, r)
	}
}
