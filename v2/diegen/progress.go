/*
Copyright 2026 the original author or authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// progress reports what diegen is doing, on stderr: each phase as it starts,
// each die package as its dies are generated, and each file that changes.
var progress = newReporter(os.Stderr)

type reporter struct {
	w     io.Writer
	start time.Time
	// cwd is where paths are reported relative to
	cwd                         string
	pkgs, dies                  int
	written, unchanged, removed int
}

func newReporter(w io.Writer) *reporter {
	cwd, _ := os.Getwd()
	return &reporter{w: w, start: time.Now(), cwd: cwd}
}

func (r *reporter) logf(format string, args ...any) {
	fmt.Fprintf(r.w, "diegen: "+format+"\n", args...)
}

// generated reports a die package as its dies are generated.
func (r *reporter) generated(pkgPath string, dies int) {
	r.pkgs++
	r.dies += dies
	r.logf("%s: %d %s", pkgPath, dies, plural(dies, "die", "dies"))
}

func (r *reporter) wrote(file string, changed bool) {
	if !changed {
		r.unchanged++
		return
	}
	r.written++
	r.logf("wrote %s", r.rel(file))
}

func (r *reporter) deleted(file string) {
	r.removed++
	r.logf("removed %s", r.rel(file))
}

func (r *reporter) done() {
	r.logf("%d %s in %d %s; %d %s written, %d unchanged, %d removed, in %s",
		r.dies, plural(r.dies, "die", "dies"), r.pkgs, plural(r.pkgs, "package", "packages"),
		r.written, plural(r.written, "file", "files"), r.unchanged, r.removed,
		time.Since(r.start).Round(10*time.Millisecond))
}

func (r *reporter) rel(file string) string {
	if rel, err := filepath.Rel(r.cwd, file); err == nil {
		return rel
	}
	return file
}

func plural(n int, one, many string) string {
	if n == 1 {
		return one
	}
	return many
}
