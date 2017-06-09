package rotator

import (
	"fmt"
	"io"
	"sync/atomic"
	"time"
)

const (
	defaultSize   = 100 << 20 // 100MB
	defaultFormat = "2006-01-02_150405"
)

type RotatorWriter interface {
	io.Writer
}

type Spliter interface {
	reachLimit(wBytes int) bool
	getNextName(prefix, ext string) string
	reset()
}

// SizeSpliter rotates log file by size
type SizeSpliter struct {
	currSize   uint64
	limitSize  uint64
	timeFormat string
}

// NewSizeSpliter return a new SizeSpliter
func NewSizeSpliter(limit uint64) *SizeSpliter {
	if limit == 0 {
		limit = defaultSize
	}

	return &SizeSpliter{
		limitSize:  limit,
		timeFormat: defaultFormat,
	}
}

// ReachLimit checks if current size is bigger than limit size
func (ss *SizeSpliter) reachLimit(n int) bool {
	atomic.AddUint64(&ss.currSize, uint64(n))
	if ss.currSize > ss.limitSize {
		return true
	}
	return false
}

func (ss *SizeSpliter) getNextName(prefix, ext string) string {
	timeStr := time.Now().Format(ss.timeFormat)
	return fmt.Sprintf("%s_%s_%d.%s", prefix, timeStr, ss.currSize, ext)
}

func (ss *SizeSpliter) reset() {
	ss.currSize = 0
}

// DaySpliter rotates log file by day
type DaySpliter struct {
	timeFormat string
	currDay    string
}

// NewDaySpliter gets a new DaySpliter
func NewDaySpliter() *DaySpliter {
	return &DaySpliter{
		timeFormat: "2006-01-02",
	}
}

func (ds *DaySpliter) reachLimit(n int) bool {
	timeStr := time.Now().Format(ds.timeFormat)
	if timeStr != ds.currDay {
		return true
	}
	return false
}

func (ds *DaySpliter) getNextName(prefix, ext string) string {
	timeStr := time.Now().Format(ds.timeFormat)
	return fmt.Sprintf("%s_%s.%s", prefix, timeStr, ext)
}

func (ds *DaySpliter) reset() {
}
