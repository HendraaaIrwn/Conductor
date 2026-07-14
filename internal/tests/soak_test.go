package tests

import (
	"runtime"
	"testing"
	"time"

	"github.com/example/conductor/internal/config"
	"github.com/example/conductor/internal/engine"
	"github.com/gdamore/tcell/v2"
)

// TestSoakLongRunning simulates 60 minutes of runtime at 20 FPS without
// rendering to a real terminal. It detects memory growth, entity leaks, and
// scheduler overflow.
//
// This test is skipped by default because it takes time. Run with:
//
//	go test -run TestSoakLongRunning -timeout 30m ./internal/tests/
//
// or with a shorter duration:
//
//	go test -run TestSoakShort -timeout 5m ./internal/tests/
func TestSoakLongRunning(t *testing.T) {
	soakTest(t, 3600, 20) // 60 minutes at 20 FPS
}

// TestSoakShort is a shorter soak test suitable for quick verification.
func TestSoakShort(t *testing.T) {
	soakTest(t, 300, 20) // 5 minutes at 20 FPS
}

// TestSoakMedium is a medium-length soak test.
func TestSoakMedium(t *testing.T) {
	soakTest(t, 900, 20) // 15 minutes at 20 FPS
}

// soakTest runs the simulation for the given number of seconds at the given
// FPS. It checks for entity leaks, memory growth, and scheduler overflow.
func soakTest(t *testing.T, durationSeconds, fps int) {
	if testing.Short() {
		t.Skip("skipping soak test in short mode")
	}

	// Detect memory growth by sampling before and after.
	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)

	cfg := config.DefaultConfig()
	cfg.Seed = 42
	cfg.FPS = fps
	style := tcell.StyleDefault

	w := engine.NewWorldWithConfig(120, 40, cfg, style)
	if w == nil {
		t.Fatal("world should not be nil")
	}

	frameTime := 1.0 / float64(fps)
	totalFrames := durationSeconds * fps
	entityCounts := make([]int, 0, 100)
	sampleInterval := totalFrames / 100
	if sampleInterval < 1 {
		sampleInterval = 1
	}

	start := time.Now()
	for i := 0; i < totalFrames; i++ {
		w.Update(frameTime)

		// Sample entity count every N frames.
		if i%sampleInterval == 0 {
			entityCounts = append(entityCounts, w.EntityCount())
		}

		// Check for entity leaks: if the count grows unboundedly, it's a leak.
		if len(entityCounts) > 10 {
			last := entityCounts[len(entityCounts)-1]
			first := entityCounts[0]
			// Allow some growth (e.g., more particles or trains), but not
			// unbounded. If the count more than doubles, it's likely a leak.
			if last > first*2 && last > 100 {
				t.Errorf("entity count grew from %d to %d — possible leak", first, last)
				break
			}
		}

		// Periodically check for scheduler overflow.
		if i%1000 == 0 && w.PendingEvents() > 10 {
			t.Errorf("scheduler has %d pending events — possible overflow", w.PendingEvents())
			break
		}
	}
	elapsed := time.Since(start)

	runtime.GC()
	runtime.ReadMemStats(&m2)

	allocMB := float64(m2.TotalAlloc-m1.TotalAlloc) / 1024 / 1024
	t.Logf("soak test: %d frames in %v (%.2f FPS), memory allocated: %.2f MB, entity count: %d",
		totalFrames, elapsed, float64(totalFrames)/elapsed.Seconds(), allocMB, w.EntityCount())

	// Check for memory leaks (more than 100 MB allocated is suspicious).
	if allocMB > 100 {
		t.Errorf("high memory allocation: %.2f MB", allocMB)
	}
}
