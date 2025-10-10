package pie

import (
	"context"
	"fmt"
	"os"
)

// HookFunc hook function type
type HookFunc func(ctx context.Context, doc interface{}) error

// Hook interface definitions

// BeforeCreator before create hook
type BeforeCreator interface {
	BeforeCreate(ctx context.Context) error
}

// AfterCreator after create hook
type AfterCreator interface {
	AfterCreate(ctx context.Context) error
}

// BeforeUpdater before update hook
type BeforeUpdater interface {
	BeforeUpdate(ctx context.Context) error
}

// AfterUpdater after update hook
type AfterUpdater interface {
	AfterUpdate(ctx context.Context) error
}

// BeforeDeleter before delete hook
type BeforeDeleter interface {
	BeforeDelete(ctx context.Context) error
}

// AfterDeleter after delete hook
type AfterDeleter interface {
	AfterDelete(ctx context.Context) error
}

// AfterFinder after find hook
type AfterFinder interface {
	AfterFind(ctx context.Context) error
}

// BeforeSaver before save hook (insert or update)
type BeforeSaver interface {
	BeforeSave(ctx context.Context) error
}

// AfterSaver after save hook (insert or update)
type AfterSaver interface {
	AfterSave(ctx context.Context) error
}

// HookManager hook manager
type HookManager struct {
	globalBeforeCreate []HookFunc
	globalAfterCreate  []HookFunc
	globalBeforeUpdate []HookFunc
	globalAfterUpdate  []HookFunc
	globalBeforeDelete []HookFunc
	globalAfterDelete  []HookFunc
	globalAfterFind    []HookFunc
	globalBeforeSave   []HookFunc
	globalAfterSave    []HookFunc
}

// NewHookManager create hook manager
func NewHookManager() *HookManager {
	return &HookManager{
		globalBeforeCreate: make([]HookFunc, 0),
		globalAfterCreate:  make([]HookFunc, 0),
		globalBeforeUpdate: make([]HookFunc, 0),
		globalAfterUpdate:  make([]HookFunc, 0),
		globalBeforeDelete: make([]HookFunc, 0),
		globalAfterDelete:  make([]HookFunc, 0),
		globalAfterFind:    make([]HookFunc, 0),
		globalBeforeSave:   make([]HookFunc, 0),
		globalAfterSave:    make([]HookFunc, 0),
	}
}

// RegisterBeforeCreate register global before create hook
func (hm *HookManager) RegisterBeforeCreate(fn HookFunc) {
	hm.globalBeforeCreate = append(hm.globalBeforeCreate, fn)
}

// RegisterAfterCreate register global after create hook
func (hm *HookManager) RegisterAfterCreate(fn HookFunc) {
	hm.globalAfterCreate = append(hm.globalAfterCreate, fn)
}

// RegisterBeforeUpdate register global before update hook
func (hm *HookManager) RegisterBeforeUpdate(fn HookFunc) {
	hm.globalBeforeUpdate = append(hm.globalBeforeUpdate, fn)
}

// RegisterAfterUpdate register global after update hook
func (hm *HookManager) RegisterAfterUpdate(fn HookFunc) {
	hm.globalAfterUpdate = append(hm.globalAfterUpdate, fn)
}

// RegisterBeforeDelete register global before delete hook
func (hm *HookManager) RegisterBeforeDelete(fn HookFunc) {
	hm.globalBeforeDelete = append(hm.globalBeforeDelete, fn)
}

// RegisterAfterDelete register global after delete hook
func (hm *HookManager) RegisterAfterDelete(fn HookFunc) {
	hm.globalAfterDelete = append(hm.globalAfterDelete, fn)
}

// RegisterAfterFind register global after find hook
func (hm *HookManager) RegisterAfterFind(fn HookFunc) {
	hm.globalAfterFind = append(hm.globalAfterFind, fn)
}

// RegisterBeforeSave register global before save hook
func (hm *HookManager) RegisterBeforeSave(fn HookFunc) {
	hm.globalBeforeSave = append(hm.globalBeforeSave, fn)
}

// RegisterAfterSave register global after save hook
func (hm *HookManager) RegisterAfterSave(fn HookFunc) {
	hm.globalAfterSave = append(hm.globalAfterSave, fn)
}

// executeBeforeCreate execute global before create hook
func (hm *HookManager) executeBeforeCreate(ctx context.Context, doc interface{}) error {
	for _, hook := range hm.globalBeforeCreate {
		if err := hook(ctx, doc); err != nil {
			return fmt.Errorf("global before create hook failed: %w", err)
		}
	}
	return nil
}

// executeAfterCreate execute global after create hook
func (hm *HookManager) executeAfterCreate(ctx context.Context, doc interface{}) error {
	for _, hook := range hm.globalAfterCreate {
		if err := hook(ctx, doc); err != nil {
			// After hook errors don't interrupt the flow, just log
			fmt.Fprintf(os.Stderr, "global after create hook failed: %v\n", err)
		}
	}
	return nil
}

// executeBeforeUpdate execute global before update hook
func (hm *HookManager) executeBeforeUpdate(ctx context.Context, doc interface{}) error {
	for _, hook := range hm.globalBeforeUpdate {
		if err := hook(ctx, doc); err != nil {
			return fmt.Errorf("global before update hook failed: %w", err)
		}
	}
	return nil
}

// executeAfterUpdate execute global after update hook
func (hm *HookManager) executeAfterUpdate(ctx context.Context, doc interface{}) error {
	for _, hook := range hm.globalAfterUpdate {
		if err := hook(ctx, doc); err != nil {
			fmt.Fprintf(os.Stderr, "global after update hook failed: %v\n", err)
		}
	}
	return nil
}

// executeBeforeDelete execute global before delete hook
func (hm *HookManager) executeBeforeDelete(ctx context.Context, doc interface{}) error {
	for _, hook := range hm.globalBeforeDelete {
		if err := hook(ctx, doc); err != nil {
			return fmt.Errorf("global before delete hook failed: %w", err)
		}
	}
	return nil
}

// executeAfterDelete execute global after delete hook
func (hm *HookManager) executeAfterDelete(ctx context.Context, doc interface{}) error {
	for _, hook := range hm.globalAfterDelete {
		if err := hook(ctx, doc); err != nil {
			fmt.Fprintf(os.Stderr, "global after delete hook failed: %v\n", err)
		}
	}
	return nil
}

// executeAfterFind execute global after find hook
func (hm *HookManager) executeAfterFind(ctx context.Context, doc interface{}) error {
	for _, hook := range hm.globalAfterFind {
		if err := hook(ctx, doc); err != nil {
			fmt.Fprintf(os.Stderr, "global after find hook failed: %v\n", err)
		}
	}
	return nil
}

// executeBeforeSave execute global before save hook
func (hm *HookManager) executeBeforeSave(ctx context.Context, doc interface{}) error {
	for _, hook := range hm.globalBeforeSave {
		if err := hook(ctx, doc); err != nil {
			return fmt.Errorf("global before save hook failed: %w", err)
		}
	}
	return nil
}

// executeAfterSave execute global after save hook
func (hm *HookManager) executeAfterSave(ctx context.Context, doc interface{}) error {
	for _, hook := range hm.globalAfterSave {
		if err := hook(ctx, doc); err != nil {
			fmt.Fprintf(os.Stderr, "global after save hook failed: %v\n", err)
		}
	}
	return nil
}

// executeModelBeforeCreate execute model before create hook
func (hm *HookManager) executeModelBeforeCreate(ctx context.Context, doc interface{}) error {
	if hook, ok := doc.(BeforeSaver); ok {
		if err := hook.BeforeSave(ctx); err != nil {
			return fmt.Errorf("model before save failed: %w", err)
		}
	}
	if hook, ok := doc.(BeforeCreator); ok {
		if err := hook.BeforeCreate(ctx); err != nil {
			return fmt.Errorf("model before create failed: %w", err)
		}
	}
	return nil
}

// executeModelAfterCreate execute model after create hook
func (hm *HookManager) executeModelAfterCreate(ctx context.Context, doc interface{}) error {
	if hook, ok := doc.(AfterCreator); ok {
		if err := hook.AfterCreate(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "model after create failed: %v\n", err)
		}
	}
	if hook, ok := doc.(AfterSaver); ok {
		if err := hook.AfterSave(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "model after save failed: %v\n", err)
		}
	}
	return nil
}

// executeModelBeforeUpdate execute model before update hook
func (hm *HookManager) executeModelBeforeUpdate(ctx context.Context, doc interface{}) error {
	if hook, ok := doc.(BeforeSaver); ok {
		if err := hook.BeforeSave(ctx); err != nil {
			return fmt.Errorf("model before save failed: %w", err)
		}
	}
	if hook, ok := doc.(BeforeUpdater); ok {
		if err := hook.BeforeUpdate(ctx); err != nil {
			return fmt.Errorf("model before update failed: %w", err)
		}
	}
	return nil
}

// executeModelAfterUpdate execute model after update hook
func (hm *HookManager) executeModelAfterUpdate(ctx context.Context, doc interface{}) error {
	if hook, ok := doc.(AfterUpdater); ok {
		if err := hook.AfterUpdate(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "model after update failed: %v\n", err)
		}
	}
	if hook, ok := doc.(AfterSaver); ok {
		if err := hook.AfterSave(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "model after save failed: %v\n", err)
		}
	}
	return nil
}

// executeModelBeforeDelete execute model before delete hook
func (hm *HookManager) executeModelBeforeDelete(ctx context.Context, doc interface{}) error {
	if hook, ok := doc.(BeforeDeleter); ok {
		if err := hook.BeforeDelete(ctx); err != nil {
			return fmt.Errorf("model before delete failed: %w", err)
		}
	}
	return nil
}

// executeModelAfterDelete execute model after delete hook
func (hm *HookManager) executeModelAfterDelete(ctx context.Context, doc interface{}) error {
	if hook, ok := doc.(AfterDeleter); ok {
		if err := hook.AfterDelete(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "model after delete failed: %v\n", err)
		}
	}
	return nil
}

// executeModelAfterFind execute model after find hook
func (hm *HookManager) executeModelAfterFind(ctx context.Context, doc interface{}) error {
	if hook, ok := doc.(AfterFinder); ok {
		if err := hook.AfterFind(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "model after find failed: %v\n", err)
		}
	}
	return nil
}
