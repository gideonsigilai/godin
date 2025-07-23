package widgets

import (
	"fmt"
	"sync"
	"time"

	"github.com/gideonsigilai/godin/pkg/core"
)

// DialogManager manages modal dialogs and bottom sheets
type DialogManager struct {
	activeDialogs      map[string]*DialogInfo
	activeBottomSheets map[string]*BottomSheetInfo
	mutex              sync.RWMutex
	context            *core.Context
	zIndexCounter      int
	maxZIndex          int
}

// DialogInfo contains information about an active dialog
type DialogInfo struct {
	ID                 string
	Widget             core.Widget
	BarrierDismissible bool
	OnDismiss          func()
	ZIndex             int
	CreatedAt          time.Time
	Result             interface{}
	ResultCallback     func(interface{})
}

// BottomSheetInfo contains information about an active bottom sheet
type BottomSheetInfo struct {
	ID             string
	Widget         core.Widget
	IsModal        bool
	IsDraggable    bool
	OnDismiss      func()
	ZIndex         int
	CreatedAt      time.Time
	Result         interface{}
	ResultCallback func(interface{})
}

// DialogOptions contains options for showing dialogs
type DialogOptions struct {
	BarrierDismissible bool
	OnDismiss          func()
	ResultCallback     func(interface{})
	UseRootNavigator   bool
}

// BottomSheetOptions contains options for showing bottom sheets
type BottomSheetOptions struct {
	IsModal        bool
	IsDraggable    bool
	OnDismiss      func()
	ResultCallback func(interface{})
	EnableDrag     bool
	ShowDragHandle bool
}

// NewDialogManager creates a new DialogManager instance
func NewDialogManager(ctx *core.Context) *DialogManager {
	return &DialogManager{
		activeDialogs:      make(map[string]*DialogInfo),
		activeBottomSheets: make(map[string]*BottomSheetInfo),
		context:            ctx,
		zIndexCounter:      1000, // Start with high z-index for dialogs
		maxZIndex:          9999,
	}
}

// ShowDialog displays a modal dialog
func (dm *DialogManager) ShowDialog(widget core.Widget, options DialogOptions) string {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	// Generate unique dialog ID
	dialogID := fmt.Sprintf("dialog_%d_%d", time.Now().UnixNano(), len(dm.activeDialogs))

	// Increment z-index for stacking
	dm.zIndexCounter++
	if dm.zIndexCounter > dm.maxZIndex {
		dm.zIndexCounter = 1000 // Reset if we hit max
	}

	// Create dialog info
	dialogInfo := &DialogInfo{
		ID:                 dialogID,
		Widget:             widget,
		BarrierDismissible: options.BarrierDismissible,
		OnDismiss:          options.OnDismiss,
		ZIndex:             dm.zIndexCounter,
		CreatedAt:          time.Now(),
		ResultCallback:     options.ResultCallback,
	}

	// Store dialog info
	dm.activeDialogs[dialogID] = dialogInfo

	// Register dismiss callback endpoint if context is available
	if dm.context != nil && dm.context.App != nil {
		callbackRegistry := dm.context.App.CallbackRegistry()
		if callbackRegistry != nil {
			dismissCallback := func() {
				dm.DismissDialog(dialogID)
			}
			callbackRegistry.RegisterCallback(
				dialogID,
				"Dialog",
				"OnDismiss",
				dismissCallback,
				dm.context,
			)
		}
	}

	return dialogID
}

// ShowBottomSheet displays a bottom sheet
func (dm *DialogManager) ShowBottomSheet(widget core.Widget, options BottomSheetOptions) string {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	// Generate unique bottom sheet ID
	sheetID := fmt.Sprintf("bottomsheet_%d_%d", time.Now().UnixNano(), len(dm.activeBottomSheets))

	// Increment z-index for stacking
	dm.zIndexCounter++
	if dm.zIndexCounter > dm.maxZIndex {
		dm.zIndexCounter = 1000 // Reset if we hit max
	}

	// Create bottom sheet info
	sheetInfo := &BottomSheetInfo{
		ID:             sheetID,
		Widget:         widget,
		IsModal:        options.IsModal,
		IsDraggable:    options.IsDraggable,
		OnDismiss:      options.OnDismiss,
		ZIndex:         dm.zIndexCounter,
		CreatedAt:      time.Now(),
		ResultCallback: options.ResultCallback,
	}

	// Store bottom sheet info
	dm.activeBottomSheets[sheetID] = sheetInfo

	// Register dismiss callback endpoint if context is available
	if dm.context != nil && dm.context.App != nil {
		callbackRegistry := dm.context.App.CallbackRegistry()
		if callbackRegistry != nil {
			dismissCallback := func() {
				dm.DismissBottomSheet(sheetID)
			}
			callbackRegistry.RegisterCallback(
				sheetID,
				"BottomSheet",
				"OnDismiss",
				dismissCallback,
				dm.context,
			)
		}
	}

	return sheetID
}

// DismissDialog dismisses a specific dialog
func (dm *DialogManager) DismissDialog(dialogID string) bool {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	dialogInfo, exists := dm.activeDialogs[dialogID]
	if !exists {
		return false
	}

	// Call dismiss callback if provided
	if dialogInfo.OnDismiss != nil {
		go dialogInfo.OnDismiss() // Run in goroutine to avoid blocking
	}

	// Call result callback with nil result if no result was set
	if dialogInfo.ResultCallback != nil && dialogInfo.Result == nil {
		go dialogInfo.ResultCallback(nil)
	}

	// Remove from active dialogs
	delete(dm.activeDialogs, dialogID)

	// Clean up callback registry if context is available
	if dm.context != nil && dm.context.App != nil {
		callbackRegistry := dm.context.App.CallbackRegistry()
		if callbackRegistry != nil {
			callbackRegistry.CleanupCallback(dialogID)
		}
	}

	return true
}

// DismissBottomSheet dismisses a specific bottom sheet
func (dm *DialogManager) DismissBottomSheet(sheetID string) bool {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	sheetInfo, exists := dm.activeBottomSheets[sheetID]
	if !exists {
		return false
	}

	// Call dismiss callback if provided
	if sheetInfo.OnDismiss != nil {
		go sheetInfo.OnDismiss() // Run in goroutine to avoid blocking
	}

	// Call result callback with nil result if no result was set
	if sheetInfo.ResultCallback != nil && sheetInfo.Result == nil {
		go sheetInfo.ResultCallback(nil)
	}

	// Remove from active bottom sheets
	delete(dm.activeBottomSheets, sheetID)

	// Clean up callback registry if context is available
	if dm.context != nil && dm.context.App != nil {
		callbackRegistry := dm.context.App.CallbackRegistry()
		if callbackRegistry != nil {
			callbackRegistry.CleanupCallback(sheetID)
		}
	}

	return true
}

// DismissAll dismisses all active dialogs and bottom sheets
func (dm *DialogManager) DismissAll() {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	// Dismiss all dialogs
	for dialogID := range dm.activeDialogs {
		dm.dismissDialogUnsafe(dialogID)
	}

	// Dismiss all bottom sheets
	for sheetID := range dm.activeBottomSheets {
		dm.dismissBottomSheetUnsafe(sheetID)
	}
}

// dismissDialogUnsafe dismisses a dialog without acquiring the mutex (internal use)
func (dm *DialogManager) dismissDialogUnsafe(dialogID string) {
	dialogInfo, exists := dm.activeDialogs[dialogID]
	if !exists {
		return
	}

	if dialogInfo.OnDismiss != nil {
		go dialogInfo.OnDismiss()
	}

	if dialogInfo.ResultCallback != nil && dialogInfo.Result == nil {
		go dialogInfo.ResultCallback(nil)
	}

	delete(dm.activeDialogs, dialogID)

	if dm.context != nil && dm.context.App != nil {
		callbackRegistry := dm.context.App.CallbackRegistry()
		if callbackRegistry != nil {
			callbackRegistry.CleanupCallback(dialogID)
		}
	}
}

// dismissBottomSheetUnsafe dismisses a bottom sheet without acquiring the mutex (internal use)
func (dm *DialogManager) dismissBottomSheetUnsafe(sheetID string) {
	sheetInfo, exists := dm.activeBottomSheets[sheetID]
	if !exists {
		return
	}

	if sheetInfo.OnDismiss != nil {
		go sheetInfo.OnDismiss()
	}

	if sheetInfo.ResultCallback != nil && sheetInfo.Result == nil {
		go sheetInfo.ResultCallback(nil)
	}

	delete(dm.activeBottomSheets, sheetID)

	if dm.context != nil && dm.context.App != nil {
		callbackRegistry := dm.context.App.CallbackRegistry()
		if callbackRegistry != nil {
			callbackRegistry.CleanupCallback(sheetID)
		}
	}
}

// GetActiveDialogs returns a copy of active dialog IDs
func (dm *DialogManager) GetActiveDialogs() []string {
	dm.mutex.RLock()
	defer dm.mutex.RUnlock()

	dialogIDs := make([]string, 0, len(dm.activeDialogs))
	for id := range dm.activeDialogs {
		dialogIDs = append(dialogIDs, id)
	}
	return dialogIDs
}

// GetActiveBottomSheets returns a copy of active bottom sheet IDs
func (dm *DialogManager) GetActiveBottomSheets() []string {
	dm.mutex.RLock()
	defer dm.mutex.RUnlock()

	sheetIDs := make([]string, 0, len(dm.activeBottomSheets))
	for id := range dm.activeBottomSheets {
		sheetIDs = append(sheetIDs, id)
	}
	return sheetIDs
}

// HasActiveDialogs returns true if there are any active dialogs
func (dm *DialogManager) HasActiveDialogs() bool {
	dm.mutex.RLock()
	defer dm.mutex.RUnlock()
	return len(dm.activeDialogs) > 0
}

// HasActiveBottomSheets returns true if there are any active bottom sheets
func (dm *DialogManager) HasActiveBottomSheets() bool {
	dm.mutex.RLock()
	defer dm.mutex.RUnlock()
	return len(dm.activeBottomSheets) > 0
}

// GetDialogInfo returns information about a specific dialog
func (dm *DialogManager) GetDialogInfo(dialogID string) (*DialogInfo, bool) {
	dm.mutex.RLock()
	defer dm.mutex.RUnlock()

	info, exists := dm.activeDialogs[dialogID]
	if !exists {
		return nil, false
	}

	// Return a copy to prevent external modification
	infoCopy := *info
	return &infoCopy, true
}

// GetBottomSheetInfo returns information about a specific bottom sheet
func (dm *DialogManager) GetBottomSheetInfo(sheetID string) (*BottomSheetInfo, bool) {
	dm.mutex.RLock()
	defer dm.mutex.RUnlock()

	info, exists := dm.activeBottomSheets[sheetID]
	if !exists {
		return nil, false
	}

	// Return a copy to prevent external modification
	infoCopy := *info
	return &infoCopy, true
}

// SetDialogResult sets the result for a dialog (to be returned when dismissed)
func (dm *DialogManager) SetDialogResult(dialogID string, result interface{}) bool {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	dialogInfo, exists := dm.activeDialogs[dialogID]
	if !exists {
		return false
	}

	dialogInfo.Result = result
	return true
}

// SetBottomSheetResult sets the result for a bottom sheet (to be returned when dismissed)
func (dm *DialogManager) SetBottomSheetResult(sheetID string, result interface{}) bool {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	sheetInfo, exists := dm.activeBottomSheets[sheetID]
	if !exists {
		return false
	}

	sheetInfo.Result = result
	return true
}

// GetTopDialog returns the dialog with the highest z-index (most recently shown)
func (dm *DialogManager) GetTopDialog() (*DialogInfo, bool) {
	dm.mutex.RLock()
	defer dm.mutex.RUnlock()

	var topDialog *DialogInfo
	maxZIndex := -1

	for _, dialog := range dm.activeDialogs {
		if dialog.ZIndex > maxZIndex {
			maxZIndex = dialog.ZIndex
			topDialog = dialog
		}
	}

	if topDialog == nil {
		return nil, false
	}

	// Return a copy
	dialogCopy := *topDialog
	return &dialogCopy, true
}

// GetTopBottomSheet returns the bottom sheet with the highest z-index (most recently shown)
func (dm *DialogManager) GetTopBottomSheet() (*BottomSheetInfo, bool) {
	dm.mutex.RLock()
	defer dm.mutex.RUnlock()

	var topSheet *BottomSheetInfo
	maxZIndex := -1

	for _, sheet := range dm.activeBottomSheets {
		if sheet.ZIndex > maxZIndex {
			maxZIndex = sheet.ZIndex
			topSheet = sheet
		}
	}

	if topSheet == nil {
		return nil, false
	}

	// Return a copy
	sheetCopy := *topSheet
	return &sheetCopy, true
}

// Cleanup cleans up all resources and dismisses all dialogs/bottom sheets
func (dm *DialogManager) Cleanup() {
	dm.DismissAll()
}
