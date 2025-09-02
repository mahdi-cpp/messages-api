package application

import (
	"context"
	_ "image/jpeg"
	_ "image/png"
	"sync"
	"time"
)

type ChatManager struct {
	mu                sync.RWMutex // Protects all indexes and maps
	lastID            int
	lastRebuild       time.Time
	maintenanceCtx    context.Context
	cancelMaintenance context.CancelFunc
	statsMu           sync.Mutex
}
