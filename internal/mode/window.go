package mode

import (
	"log/slog"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/quar15/qq-go/internal/cursor"
	"github.com/quar15/qq-go/internal/motion"
)

type WindowManager struct {
	currCtx        *Context
	editorCtx      *Context
	spreadsheetCtx *Context
	connectionsCtx *Context
}

func (wm *WindowManager) CurrCtx() *Context {
	return wm.currCtx
}

func InitWindowManager(editorCtx *Context, spreadsheetCtx *Context, connectionsCtx *Context) *WindowManager {
	mgr := &WindowManager{
		editorCtx:      editorCtx,
		spreadsheetCtx: spreadsheetCtx,
		connectionsCtx: connectionsCtx,
	}
	mgr.ChangeWindow(cursor.TypeEditor)

	return mgr
}

func (mgr *WindowManager) ChangeWindow(target cursor.Type) {
	if mgr.currCtx != nil {
		mgr.currCtx.Cursor.Deactivate()
	}
	switch target {
	case cursor.TypeEditor:
		mgr.currCtx = mgr.editorCtx
	case cursor.TypeSpreadsheet:
		mgr.currCtx = mgr.spreadsheetCtx
	case cursor.TypeConnections:
		mgr.currCtx = mgr.connectionsCtx
	}
	mgr.currCtx.Cursor.Activate()
	slog.Info("Activated window", slog.String("type", mgr.currCtx.Cursor.Type.String()))
}

func (mgr *WindowManager) SwapWindow() {
	switch mgr.currCtx.Cursor.Type {
	case cursor.TypeEditor:
		mgr.ChangeWindow(cursor.TypeSpreadsheet)
	case cursor.TypeSpreadsheet:
		mgr.ChangeWindow(cursor.TypeEditor)
	case cursor.TypeConnections:
		mgr.ChangeWindow(cursor.TypeEditor)
	}
}

type WindowManagementMode struct{}

func (WindowManagementMode) Handle(ctx *Context, k motion.Key) {
	slog.Debug("Window Management Mode | Handling key", slog.Any("key", k))
	switch k.Rune {
	case rl.KeyEscape, rl.KeyCapsLock:
		ctx.Parser.Reset()
		ctx.Cursor.TransitionMode(cursor.ModeNormal)
	case 'w', 'n', 'p':
		ctx.WindowManager.SwapWindow()
		ctx.Cursor.TransitionMode(cursor.ModeNormal)
	case 'k', rl.KeyUp:
		ctx.WindowManager.ChangeWindow(cursor.TypeEditor)
		ctx.Cursor.TransitionMode(cursor.ModeNormal)
	case 'j', rl.KeyDown:
		ctx.WindowManager.ChangeWindow(cursor.TypeSpreadsheet)
		ctx.Cursor.TransitionMode(cursor.ModeNormal)
	// @TODO: Window split management
	case '=':
	case '+':
	case '-':
	}
}
