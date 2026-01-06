package commands

import (
	"log/slog"

	"github.com/quar15/qq-go/internal/cursor"
	"github.com/quar15/qq-go/internal/format"
	"github.com/quar15/qq-go/internal/mode"
	"golang.design/x/clipboard"
)

type CopyToClipboardSpreadsheet struct{}
type CopyToClipboardEditor struct{}

func (CopyToClipboardSpreadsheet) Execute(ctx *mode.Context) error {
	if ctx.Cursor.Type != cursor.TypeSpreadsheet {
		return nil
	}

	var dataString string = ""
	c := ctx.Cursor
	dg := ctx.DataGrid

	switch ctx.Cursor.Common.Mode {
	case cursor.ModeVisual:
		for row := c.Position.SelectStartRow; row <= c.Position.SelectEndRow; row++ {
			for col := int32(0); col < dg.Cols; col++ {
				if c.IsSelected(col, row) {
					dataString += format.GetValueAsString(dg.Data[row][dg.Headers[col]]) + ","
				}
			}
			dataString += "\n"
		}
	case cursor.ModeVLine:
		for row := c.Position.SelectStartRow; row <= c.Position.SelectEndRow; row++ {
			for col := int32(0); col < dg.Cols; col++ {
				dataString += format.GetValueAsString(dg.Data[row][dg.Headers[col]]) + ","
			}
			dataString += format.GetValueAsString(dg.Data[row][dg.Headers[c.Position.SelectEndCol]]) + "\n"
		}
	case cursor.ModeVBlock:
		for row := c.Position.SelectStartRow; row <= c.Position.SelectEndRow; row++ {
			for col := c.Position.SelectStartCol; col < c.Position.SelectEndCol; col++ {
				dataString += format.GetValueAsString(dg.Data[row][dg.Headers[col]]) + ","
			}
			dataString += format.GetValueAsString(dg.Data[row][dg.Headers[c.Position.SelectEndCol]]) + "\n"
		}
	case cursor.ModeNormal:
		dataString = format.GetValueAsString(dg.Data[c.Position.Row][dg.Headers[c.Position.Col]])
	}

	slog.Debug("Copied to clipboard from spreadsheet", slog.String("dataString", dataString))
	clipboard.Write(clipboard.FmtText, []byte(dataString))

	return nil
}

func (CopyToClipboardEditor) Execute(ctx *mode.Context) error {
	if ctx.Cursor.Type != cursor.TypeEditor {
		return nil
	}

	var dataString string = ""
	c := ctx.Cursor
	eg := ctx.EditorGrid

	switch ctx.Cursor.Common.Mode {
	case cursor.ModeVisual:
		for row := c.Position.SelectStartRow; row <= c.Position.SelectEndRow; row++ {
			if eg.Cols[row] > 0 {
				for col := int32(0); col < eg.Cols[row]; col++ {
					if c.IsSelected(col, row) {
						dataString += string(eg.Text[row][col])
					}
				}
			}
			dataString += "\n"
		}
	case cursor.ModeVLine:
		for row := c.Position.SelectStartRow; row <= c.Position.SelectEndRow; row++ {
			dataString += eg.Text[row] + "\n"
		}
	case cursor.ModeVBlock:
		for row := c.Position.SelectStartRow; row <= c.Position.SelectEndRow; row++ {
			if eg.Cols[row] > 0 {
				endCol := min(c.Position.SelectEndCol, eg.Cols[row])
				for col := c.Position.SelectStartCol; col <= endCol; col++ {
					dataString += string(eg.Text[row][col])
				}
			}
			dataString += "\n"
		}
	default:
		dataString = eg.Text[c.Position.Row]
	}

	slog.Debug("Copied to clipboard from editor", slog.String("dataString", dataString))
	clipboard.Write(clipboard.FmtText, []byte(dataString))

	return nil
}
