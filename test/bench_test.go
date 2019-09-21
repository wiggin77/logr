package test

import (
	"io/ioutil"
	"testing"

	"github.com/wiggin77/logr"
	"github.com/wiggin77/logr/format"
	"github.com/wiggin77/logr/level"
	"github.com/wiggin77/logr/target"
)

// Enabled avoids compiler optimization.
var Enabled bool

// Stacktrace avoids compiler optimization.
var Stacktrace bool

func BenchmarkFilterOut(b *testing.B) {
	for i := 0; i < 10; i++ {
		target := &target.Writer{Level: level.Error, Fmtr: &format.Plain{Delim: " | "}, Out: ioutil.Discard, MaxQueued: 1000}
		logr.AddTarget(target)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		status := logr.IsLevelEnabled(level.Debug)
		Enabled = status.Enabled
		Stacktrace = status.Stacktrace
	}
	b.StopTimer()
	logr.Shutdown()
}

// BenchmarkLog measures adding a log record to the queue. It does not measure
// how long the record takes to be output as that happens async.
// Level caching is enabled.
// This is how long you can expect logging to tie up the calling thread.
func BenchmarkLog(b *testing.B) {
	for i := 0; i < 10; i++ {
		target := &target.Writer{Level: level.Warn, Fmtr: &format.Plain{Delim: " | "}, Out: ioutil.Discard, MaxQueued: 1000}
		logr.AddTarget(target)
	}

	logger := logr.NewLogger()
	logger.Errorln("log entry cache primer")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Errorf("log entry %d", b.N)
	}
	b.StopTimer()
	logr.Shutdown()
}

// BenchmarkLogFiltered measures a logging call for a level that has no
// targets matching the level.  Level caching is enabled.
// This is how long you can expect logging to tie up the calling thread.
func BenchmarkLogFiltered(b *testing.B) {
	for i := 0; i < 10; i++ {
		target := &target.Writer{Level: level.Fatal, Fmtr: &format.Plain{Delim: " | "}, Out: ioutil.Discard, MaxQueued: 1000}
		logr.AddTarget(target)
	}

	logger := logr.NewLogger()
	logger.Errorln("log entry cache primer")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Log(level.Error, "blap bleep bloop")
	}
	b.StopTimer()
	logr.Shutdown()
}

// BenchmarkLogStacktrace measures adding a log record to the queue.
// It does not measure how long the record takes to be output as that happens async.
// Level caching is enabled.
// This is how long you can expect logging to tie up the calling thread when a stack
// trace is generated.
func BenchmarkLogStacktrace(b *testing.B) {
	oldStackLevel := level.StdStacktraceLevel
	level.StdStacktraceLevel = level.Error
	defer func() { level.StdStacktraceLevel = oldStackLevel }()

	for i := 0; i < 10; i++ {
		target := &target.Writer{Level: level.Error, Fmtr: &format.Plain{Delim: " | "}, Out: ioutil.Discard, MaxQueued: 1000}
		logr.AddTarget(target)
	}

	logger := logr.NewLogger()
	logger.Errorln("log entry cache primer")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Errorf("log entry with stack trace %d", b.N)
	}
	b.StopTimer()
	logr.Shutdown()
}
