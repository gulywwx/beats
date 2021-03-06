// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package filestream

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	loginp "github.com/elastic/beats/v7/filebeat/input/filestream/internal/input-logfile"
	"github.com/elastic/beats/v7/libbeat/common/match"
	"github.com/elastic/beats/v7/libbeat/logp"
)

var (
	excludedFilePath = filepath.Join("testdata", "excluded_file")
	includedFilePath = filepath.Join("testdata", "included_file")
	directoryPath    = filepath.Join("testdata", "unharvestable_dir")
)

func TestFileScanner(t *testing.T) {
	testCases := map[string]struct {
		paths         []string
		excludedFiles []match.Matcher
		symlinks      bool
		expectedFiles []string
	}{
		"select all files": {
			paths: []string{excludedFilePath, includedFilePath},
			expectedFiles: []string{
				mustAbsPath(excludedFilePath),
				mustAbsPath(includedFilePath),
			},
		},
		"skip excluded files": {
			paths: []string{excludedFilePath, includedFilePath},
			excludedFiles: []match.Matcher{
				match.MustCompile("excluded_file"),
			},
			expectedFiles: []string{
				mustAbsPath(includedFilePath),
			},
		},
		"skip directories": {
			paths:         []string{directoryPath},
			expectedFiles: []string{},
		},
	}

	setupFilesForScannerTest(t)
	defer removeFilesOfScannerTest(t)

	for name, test := range testCases {
		test := test

		t.Run(name, func(t *testing.T) {
			cfg := fileScannerConfig{
				ExcludedFiles: test.excludedFiles,
				Symlinks:      test.symlinks,
				RecursiveGlob: false,
			}
			fs, err := newFileScanner(test.paths, cfg)
			if err != nil {
				t.Fatal(err)
			}
			files := fs.GetFiles()
			paths := make([]string, 0)
			for p, _ := range files {
				paths = append(paths, p)
			}
			assert.ElementsMatch(t, paths, test.expectedFiles)
		})
	}
}

func setupFilesForScannerTest(t *testing.T) {
	err := os.MkdirAll(directoryPath, 0750)
	if err != nil {
		t.Fatal(t)
	}

	for _, path := range []string{excludedFilePath, includedFilePath} {
		f, err := os.Create(path)
		if err != nil {
			t.Fatalf("file %s, error %v", path, err)
		}
		f.Close()
	}
}

func removeFilesOfScannerTest(t *testing.T) {
	err := os.RemoveAll("testdata")
	if err != nil {
		t.Fatal(err)
	}
}

func TestFileWatchNewDeleteModified(t *testing.T) {
	oldTs := time.Now()
	newTs := oldTs.Add(5 * time.Second)
	testCases := map[string]struct {
		prevFiles      map[string]os.FileInfo
		nextFiles      map[string]os.FileInfo
		expectedEvents []loginp.FSEvent
	}{
		"one new file": {
			prevFiles: map[string]os.FileInfo{},
			nextFiles: map[string]os.FileInfo{
				"new_path": testFileInfo{"new_path", 5, oldTs},
			},
			expectedEvents: []loginp.FSEvent{
				loginp.FSEvent{Op: loginp.OpCreate, OldPath: "", NewPath: "new_path", Info: testFileInfo{"new_path", 5, oldTs}},
			},
		},
		"one deleted file": {
			prevFiles: map[string]os.FileInfo{
				"old_path": testFileInfo{"old_path", 5, oldTs},
			},
			nextFiles: map[string]os.FileInfo{},
			expectedEvents: []loginp.FSEvent{
				loginp.FSEvent{Op: loginp.OpDelete, OldPath: "old_path", NewPath: "", Info: testFileInfo{"old_path", 5, oldTs}},
			},
		},
		"one modified file": {
			prevFiles: map[string]os.FileInfo{
				"path": testFileInfo{"path", 5, oldTs},
			},
			nextFiles: map[string]os.FileInfo{
				"path": testFileInfo{"path", 10, newTs},
			},
			expectedEvents: []loginp.FSEvent{
				loginp.FSEvent{Op: loginp.OpWrite, OldPath: "path", NewPath: "path", Info: testFileInfo{"path", 10, newTs}},
			},
		},
		"two modified files": {
			prevFiles: map[string]os.FileInfo{
				"path1": testFileInfo{"path1", 5, oldTs},
				"path2": testFileInfo{"path2", 5, oldTs},
			},
			nextFiles: map[string]os.FileInfo{
				"path1": testFileInfo{"path1", 10, newTs},
				"path2": testFileInfo{"path2", 10, newTs},
			},
			expectedEvents: []loginp.FSEvent{
				loginp.FSEvent{Op: loginp.OpWrite, OldPath: "path1", NewPath: "path1", Info: testFileInfo{"path1", 10, newTs}},
				loginp.FSEvent{Op: loginp.OpWrite, OldPath: "path2", NewPath: "path2", Info: testFileInfo{"path2", 10, newTs}},
			},
		},
		"one modified file, one new file": {
			prevFiles: map[string]os.FileInfo{
				"path1": testFileInfo{"path1", 5, oldTs},
			},
			nextFiles: map[string]os.FileInfo{
				"path1": testFileInfo{"path1", 10, newTs},
				"path2": testFileInfo{"path2", 10, newTs},
			},
			expectedEvents: []loginp.FSEvent{
				loginp.FSEvent{Op: loginp.OpWrite, OldPath: "path1", NewPath: "path1", Info: testFileInfo{"path1", 10, newTs}},
				loginp.FSEvent{Op: loginp.OpCreate, OldPath: "", NewPath: "path2", Info: testFileInfo{"path2", 10, newTs}},
			},
		},
		"one new file, one deleted file": {
			prevFiles: map[string]os.FileInfo{
				"path_deleted": testFileInfo{"path_deleted", 5, oldTs},
			},
			nextFiles: map[string]os.FileInfo{
				"path_new": testFileInfo{"path_new", 10, newTs},
			},
			expectedEvents: []loginp.FSEvent{
				loginp.FSEvent{Op: loginp.OpDelete, OldPath: "path_deleted", NewPath: "", Info: testFileInfo{"path_deleted", 5, oldTs}},
				loginp.FSEvent{Op: loginp.OpCreate, OldPath: "", NewPath: "path_new", Info: testFileInfo{"path_new", 10, newTs}},
			},
		},
	}

	for name, test := range testCases {
		test := test

		t.Run(name, func(t *testing.T) {
			w := fileWatcher{
				log:     logp.L(),
				prev:    test.prevFiles,
				scanner: &mockScanner{test.nextFiles},
				events:  make(chan loginp.FSEvent),
			}

			go w.watch(context.Background())

			count := len(test.expectedEvents)
			actual := make([]loginp.FSEvent, count)
			for i := 0; i < count; i++ {
				actual[i] = w.Event()
			}

			assert.ElementsMatch(t, actual, test.expectedEvents)
		})
	}
}

type mockScanner struct {
	files map[string]os.FileInfo
}

func (m *mockScanner) GetFiles() map[string]os.FileInfo {
	return m.files
}

type testFileInfo struct {
	path string
	size int64
	time time.Time
}

func (t testFileInfo) Name() string       { return t.path }
func (t testFileInfo) Size() int64        { return t.size }
func (t testFileInfo) Mode() os.FileMode  { return 0 }
func (t testFileInfo) ModTime() time.Time { return t.time }
func (t testFileInfo) IsDir() bool        { return false }
func (t testFileInfo) Sys() interface{}   { return nil }

func mustAbsPath(path string) string {
	p, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	return p
}

func mustDuration(durStr string) time.Duration {
	dur, err := time.ParseDuration(durStr)
	if err != nil {
		panic(err)
	}
	return dur
}
