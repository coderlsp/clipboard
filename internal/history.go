package internal

import (
	"container/list"
)

const RecentDefaultNum = 10

type History interface {
	Len() int
	Add(string)
	Del(string)
	Has(string) bool
}

// HistoryLRU of clipboard
type HistoryLRU struct {
	num  int // max history items' num
	list *list.List
}

func newHistoryLRU() *HistoryLRU {
	return &HistoryLRU{
		list: list.New(),
		num:  RecentDefaultNum,
	}
}

func (h *HistoryLRU) Len() int {
	return h.list.Len()
}

func (h *HistoryLRU) Add(record string) {
	if h.list.Len() < h.num {
		h.list.PushFront(record)
		return
	}
	h.list.Remove(h.list.Back())
	h.list.PushFront(record)
}

func (h *HistoryLRU) Has(record string) bool {
	return h.get(record) != nil
}

func (h *HistoryLRU) get(record string) *list.Element {
	for e := h.list.Front(); e != nil; e = e.Next() {
		if e.Value.(string) == record {
			return e
		}
	}
	return nil
}

func (h *HistoryLRU) Del(record string) {
	e := h.get(record)
	if e == nil {
		return
	}
	h.list.Remove(e)
}

// HistoryMgr clipboard history manager
type HistoryMgr struct {
	Recent   History
	Resident History
}

func newMgr() *HistoryMgr {
	return &HistoryMgr{
		Recent:   newHistoryLRU(),
		Resident: newHistoryLRU(),
	}
}

// Add a new record
func (mgr *HistoryMgr) Add(record string) {
	if mgr.has(record) {
		// is already in history
		return
	}
	mgr.Recent.Add(record)
}

// Collect a record by index of history
func (mgr *HistoryMgr) Collect(record string) {
	if mgr.Recent.Has(record) {
		mgr.Recent.Del(record)
	}
	mgr.Resident.Add(record)
}

func (mgr *HistoryMgr) has(record string) bool {
	return mgr.Recent.Has(record) || mgr.Resident.Has(record)
}
