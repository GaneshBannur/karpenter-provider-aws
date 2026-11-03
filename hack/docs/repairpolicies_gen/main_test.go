/*
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
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestExitsWithErrorWhenWriteFails(t *testing.T) {
	dir := t.TempDir()
	gen := filepath.Join(dir, "repairpolicies_gen")
	if out, err := exec.Command("go", "build", "-o", gen, ".").CombinedOutput(); err != nil {
		t.Fatalf("building the generator: %v\n%s", err, out)
	}
	doc, err := os.ReadFile("../../../website/content/en/preview/concepts/disruption.md")
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "disruption.md")
	if err := os.WriteFile(path, doc, 0o644); err != nil {
		t.Fatal(err)
	}
	// Cap the file size at 512 bytes so the write fails partway, as on a full disk.
	if err := exec.Command("sh", "-c", `ulimit -f 1 && exec "$0" "$1"`, gen, path).Run(); err == nil {
		written, _ := os.ReadFile(path)
		t.Errorf("generator exited 0 after a failed write cut the file from %d to %d bytes", len(doc), len(written))
	}
}
