package scanners

import (
	"github.com/sauravmarvani/nextcode/internal/codereview"
)

// Registry holds all available scanners
type Registry struct {
	scanners map[codereview.ScannerType]codereview.Scanner
}

// NewRegistry creates a new scanner registry
func NewRegistry() *Registry {
	return &Registry{
		scanners: make(map[codereview.ScannerType]codereview.Scanner),
	}
}

// RegisterAllScanners registers all built-in scanners
func (r *Registry) RegisterAllScanners(analyzer *codereview.Analyzer) {
	// Register security scanner
	securityScanner := NewSecurityScanner()
	r.scanners[codereview.ScannerSecurity] = securityScanner
	analyzer.RegisterScanner(securityScanner)

	// Register performance scanner
	performanceScanner := NewPerformanceScanner()
	r.scanners[codereview.ScannerPerformance] = performanceScanner
	analyzer.RegisterScanner(performanceScanner)

	// Register correctness scanner
	correctnessScanner := NewCorrectnessScanner()
	r.scanners[codereview.ScannerCorrectness] = correctnessScanner
	analyzer.RegisterScanner(correctnessScanner)

	// Register style scanner
	styleScanner := NewStyleScanner()
	r.scanners[codereview.ScannerStyle] = styleScanner
	analyzer.RegisterScanner(styleScanner)
}

// GetScanner retrieves a scanner by type
func (r *Registry) GetScanner(scannerType codereview.ScannerType) (codereview.Scanner, bool) {
	scanner, exists := r.scanners[scannerType]
	return scanner, exists
}

// ListScanners returns all registered scanners
func (r *Registry) ListScanners() map[codereview.ScannerType]codereview.Scanner {
	return r.scanners
}

// EnableScanner enables a specific scanner
func (r *Registry) EnableScanner(scannerType codereview.ScannerType) bool {
	if scanner, exists := r.scanners[scannerType]; exists {
		scanner.SetEnabled(true)
		return true
	}
	return false
}

// DisableScanner disables a specific scanner
func (r *Registry) DisableScanner(scannerType codereview.ScannerType) bool {
	if scanner, exists := r.scanners[scannerType]; exists {
		scanner.SetEnabled(false)
		return true
	}
	return false
}
