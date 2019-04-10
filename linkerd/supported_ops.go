// Copyright 2019 Layer5.io
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

package linkerd

type supportedOperation struct {
	// a friendly name
	name string
	// the template file name
	templateName string
}

const (
	customOpCommand         = "custom"
	installLinkerdCommand   = "linkerd_install"
	installEmojiVotoCommand = "install_emojivoto"
	installBooksAppCommand  = "install_booksapp"
	cbCommand               = "linkerd_cb1"
)

var supportedOps = map[string]supportedOperation{
	installLinkerdCommand: {
		name: "Install the latest version of Linkerd",
	},
	installEmojiVotoCommand: {
		name: "Install the canonical Emojivoto demo Application",
	},
	installBooksAppCommand: {
		name: "Install the Books demo Application",
	},
	customOpCommand: {
		name: "Custom YAML",
	},
}
