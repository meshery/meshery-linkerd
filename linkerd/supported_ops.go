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

import "github.com/layer5io/meshery-linkerd/meshes"

type supportedOperation struct {
	// a friendly name
	name string
	// the template file name
	templateName string

	opType meshes.OpCategory
}

const (
	customOpCommand         = "custom"
	installLinkerdCommand   = "linkerd_install"
	installEmojiVotoCommand = "install_emojivoto"
	installBooksAppCommand  = "install_booksapp"
	installHTTPBinApp       = "install_http_bin"
	installIstioBookInfoApp = "install_istio_book_info"
	injectLinkerd="inject_linkerd"
)

var supportedOps = map[string]supportedOperation{
	installLinkerdCommand: {
		name:   "Latest version of Linkerd",
		opType: meshes.OpCategory_INSTALL,
	},
	installEmojiVotoCommand: {
		name:   "Emojivoto Application",
		opType: meshes.OpCategory_SAMPLE_APPLICATION,
	},
	installBooksAppCommand: {
		name:   "Linkerd Books Application",
		opType: meshes.OpCategory_SAMPLE_APPLICATION,
	},
	customOpCommand: {
		name:   "Custom YAML",
		opType: meshes.OpCategory_CUSTOM,
	},
	installHTTPBinApp: {
		name:         "HTTPbin Application",
		templateName: "httpbin.yaml",
		opType:       meshes.OpCategory_SAMPLE_APPLICATION,
	},
	installIstioBookInfoApp: {
		name:         "Istio Book Info Application",
		templateName: "istiobookinfo.yaml",
		opType:       meshes.OpCategory_SAMPLE_APPLICATION,
	},
	injectLinkerd: {
		name:         "Annotate namespace for sidecar proxy injection",
		opType:       meshes.OpCategory_CONFIGURE,
	},
}
