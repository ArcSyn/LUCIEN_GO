package ui

import (
	"github.com/luciendev/lucien-core/internal/completion"
	"github.com/luciendev/lucien-core/internal/history"
)

// HistoryAdapter adapts the history manager to work with the completion engine
type HistoryAdapter struct {
	historyMgr *history.Manager
}

// NewHistoryAdapter creates a new history adapter
func NewHistoryAdapter(historyMgr *history.Manager) *HistoryAdapter {
	return &HistoryAdapter{
		historyMgr: historyMgr,
	}
}

// GetRecent implements completion.HistoryProvider
func (h *HistoryAdapter) GetRecent(n int) []completion.HistoryEntry {
	if h.historyMgr == nil {
		return []completion.HistoryEntry{}
	}
	
	entries := h.historyMgr.GetRecent(n)
	result := make([]completion.HistoryEntry, len(entries))
	
	for i, entry := range entries {
		result[i] = completion.HistoryEntry{
			Command: entry.Command,
		}
	}
	
	return result
}

// Search implements completion.HistoryProvider
func (h *HistoryAdapter) Search(query string, limit int) []completion.HistoryEntry {
	if h.historyMgr == nil {
		return []completion.HistoryEntry{}
	}
	
	entries := h.historyMgr.Search(query, limit)
	result := make([]completion.HistoryEntry, len(entries))
	
	for i, entry := range entries {
		result[i] = completion.HistoryEntry{
			Command: entry.Command,
		}
	}
	
	return result
}