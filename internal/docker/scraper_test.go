package docker

import (
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
)

func TestScraperFetch(t *testing.T) {
	s := &Scraper{
		settings: ScraperSettings{BufferSize: 5}.withDefaults(),
		apps:     make(map[string]*appData),
	}

	// Record some samples
	for i := range 3 {
		s.recordSample("myapp", Sample{
			Timestamp:   time.Now(),
			CPUPercent:  float64(i + 1),
			MemoryBytes: uint64((i + 1) * 1000),
		})
	}

	// Fetch all samples
	samples := s.Fetch("myapp", 10)
	if len(samples) != 3 {
		t.Errorf("expected 3 samples, got %d", len(samples))
	}

	// Verify newest-to-oldest order
	if samples[0].CPUPercent != 3.0 {
		t.Errorf("expected newest sample first (CPU 3.0), got %f", samples[0].CPUPercent)
	}
	if samples[2].CPUPercent != 1.0 {
		t.Errorf("expected oldest sample last (CPU 1.0), got %f", samples[2].CPUPercent)
	}
}

func TestScraperFetchWithLimit(t *testing.T) {
	s := &Scraper{
		settings: ScraperSettings{BufferSize: 5}.withDefaults(),
		apps:     make(map[string]*appData),
	}

	for i := range 5 {
		s.recordSample("myapp", Sample{CPUPercent: float64(i + 1)})
	}

	samples := s.Fetch("myapp", 2)
	if len(samples) != 2 {
		t.Errorf("expected 2 samples, got %d", len(samples))
	}
	if samples[0].CPUPercent != 5.0 || samples[1].CPUPercent != 4.0 {
		t.Errorf("expected [5, 4], got [%f, %f]", samples[0].CPUPercent, samples[1].CPUPercent)
	}
}

func TestScraperFetchUnknownApp(t *testing.T) {
	s := &Scraper{
		settings: ScraperSettings{}.withDefaults(),
		apps:     make(map[string]*appData),
	}

	samples := s.Fetch("unknown", 10)
	if samples != nil {
		t.Errorf("expected nil for unknown app, got %v", samples)
	}
}

func TestScraperRingBufferWrap(t *testing.T) {
	s := &Scraper{
		settings: ScraperSettings{BufferSize: 3}.withDefaults(),
		apps:     make(map[string]*appData),
	}

	// Record more samples than buffer size
	for i := range 5 {
		s.recordSample("myapp", Sample{CPUPercent: float64(i + 1)})
	}

	samples := s.Fetch("myapp", 10)
	if len(samples) != 3 {
		t.Errorf("expected 3 samples (buffer size), got %d", len(samples))
	}

	// Should have the 3 most recent: 5, 4, 3
	if samples[0].CPUPercent != 5.0 {
		t.Errorf("expected 5.0, got %f", samples[0].CPUPercent)
	}
	if samples[1].CPUPercent != 4.0 {
		t.Errorf("expected 4.0, got %f", samples[1].CPUPercent)
	}
	if samples[2].CPUPercent != 3.0 {
		t.Errorf("expected 3.0, got %f", samples[2].CPUPercent)
	}
}

func TestCalculateCPUPercent(t *testing.T) {
	stats := &container.StatsResponse{
		CPUStats: container.CPUStats{
			CPUUsage: container.CPUUsage{
				TotalUsage: 200000000,
			},
			SystemUsage: 1000000000,
			OnlineCPUs:  4,
		},
		PreCPUStats: container.CPUStats{
			CPUUsage: container.CPUUsage{
				TotalUsage: 100000000,
			},
			SystemUsage: 500000000,
		},
	}

	// CPU delta = 100000000, System delta = 500000000
	// Percent = (100000000 / 500000000) * 4 * 100 = 80%
	percent := calculateCPUPercent(stats)
	if percent != 80.0 {
		t.Errorf("expected 80.0%%, got %f%%", percent)
	}
}

func TestCalculateCPUPercentZeroDelta(t *testing.T) {
	stats := &container.StatsResponse{
		CPUStats: container.CPUStats{
			CPUUsage: container.CPUUsage{
				TotalUsage: 100000000,
			},
			SystemUsage: 500000000,
			OnlineCPUs:  4,
		},
		PreCPUStats: container.CPUStats{
			CPUUsage: container.CPUUsage{
				TotalUsage: 100000000,
			},
			SystemUsage: 500000000,
		},
	}

	percent := calculateCPUPercent(stats)
	if percent != 0.0 {
		t.Errorf("expected 0.0%% for zero delta, got %f%%", percent)
	}
}
