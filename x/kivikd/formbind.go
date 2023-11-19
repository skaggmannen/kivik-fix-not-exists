// Licensed under the Apache License, Version 2.0 (the "License"); you may not
// use this file except in compliance with the License. You may obtain a copy of
// the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations under
// the License.

package kivikd

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/http"

	"github.com/ajg/form"
)

// BindParams binds the request form or JSON body to the provided struct.
func BindParams(r *http.Request, i interface{}) error {
	mtype, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
	switch mtype {
	case typeJSON:
		defer r.Body.Close() // nolint: errcheck
		return json.NewDecoder(r.Body).Decode(i)
	case typeForm:
		defer r.Body.Close() // nolint: errcheck
		return form.NewDecoder(r.Body).Decode(i)
	}
	return fmt.Errorf("unable to bind media type %s", mtype)
}
