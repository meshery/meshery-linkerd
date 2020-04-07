// Copyright 2020 Layer5.io
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"os"
	"testing"
)

func TestSafeClose(t *testing.T) {
	f, err := os.Create("test_safe_close.txt")

	SafeClose(f, &err)

	if err != nil {
		t.Errorf("Close function is failed")
	}

	err = os.Remove("test_safe_close.txt")

	if err != nil {
		t.Errorf("Remove file test_safe_close.txt failed")
	}

}
