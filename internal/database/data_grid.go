package database

import (
	"encoding/csv"
	"errors"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/quar15/qq-go/internal/assets"
	"github.com/quar15/qq-go/internal/format"
)

type DataGrid struct {
	Data         []map[string]any
	Headers      []string
	ColumnsWidth []int32
	Rows         int32
	Cols         int8
}

func (dg *DataGrid) FakeInit(appAssets *assets.Assets) {
	dg.Data = []map[string]any{
		{"ID": 1, "name": "Alice", "created_at": "2025-11-25"},
		{"ID": 2, "name": "Bob", "created_at": "2025-11-25"},
		{"ID": 3, "name": "Charlie", "created_at": "2025-11-25T00:00:00.000Z"},
	}
	dg.Headers = []string{"ID", "name", "created_at"}
	dg.Rows = int32(len(dg.Data))
	dg.Cols = int8(len(dg.Data[0]))
	dg.UpdateColumnsWidth(appAssets)
}

func (dg *DataGrid) UpdateColumnsWidth(appAssets *assets.Assets) {
	const maximumColWidth int32 = 600
	const minimumColWidth int32 = 50
	const textPadding int32 = 8
	dg.ColumnsWidth = nil
	dg.ColumnsWidth = make([]int32, len(dg.Headers))
	for i, h := range dg.Headers {
		var headerWidth int32 = int32(rl.MeasureTextEx(appAssets.MainFont, h, appAssets.MainFontSize, appAssets.MainFontSpacing).X) + (textPadding * 3)
		headerWidth = min(max(headerWidth, minimumColWidth), maximumColWidth)
		for row := range dg.Rows {
			val := dg.Data[row][h]
			var textWidth int32 = int32(rl.MeasureTextEx(appAssets.MainFont, format.GetValueAsString(val), appAssets.MainFontSize, appAssets.MainFontSpacing).X) + (textPadding * 3)
			if textWidth > headerWidth {
				headerWidth = textWidth
			}
		}
		dg.ColumnsWidth[i] = min(headerWidth, maximumColWidth)
	}
}

func LoadDataGridFromCSV(path string, appAssets *assets.Assets) (*DataGrid, error) {
	var dg *DataGrid = &DataGrid{}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	dg.Headers, err = csvReader.Read()
	if err != nil {
		return nil, err
	}
	dg.Cols = int8(len(dg.Headers))

	dg.Data = []map[string]any{}
	dg.Rows = 0
	for {
		record, err := csvReader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, errors.New("Failed to parse provided CSV")
		}
		newData := make(map[string]any)
		for i, key := range dg.Headers {
			if i < len(record) {
				newData[key] = record[i]
			}
		}

		dg.Data = append(dg.Data, newData)
		dg.Rows++
	}

	dg.UpdateColumnsWidth(appAssets)
	return dg, nil
}
