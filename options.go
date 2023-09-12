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

package kivik

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/go-kivik/kivik/v4/driver"
)

// Option wraps a Kivik or backend option.
type Option interface {
	// Apply applies the option to target, if target is of the expected type.
	// Unexpected/recognized target types should be ignored.
	Apply(target interface{})
}

var _ Option = (driver.Options)(nil)

type allOptions []Option

var _ Option = (allOptions)(nil)

func (o allOptions) Apply(t interface{}) {
	for _, opt := range o {
		if opt != nil {
			opt.Apply(t)
		}
	}
}

func (o allOptions) String() string {
	parts := make([]string, 0, len(o))
	for _, opt := range o {
		if part := fmt.Sprintf("%s", opt); part != "" {
			parts = append(parts, part)
		}
	}
	return strings.Join(parts, ",")
}

// Params is a collection of options. The keys and values are backend specific.
type Params map[string]interface{}

// Apply applies o to target. The following target types are supported:
//
//   - map[string]interface{}
//   - *url.Values
func (o Params) Apply(target interface{}) {
	switch t := target.(type) {
	case map[string]interface{}:
		for k, v := range o {
			t[k] = v
		}
	case *url.Values:
		for key, i := range o {
			var values []string
			switch v := i.(type) {
			case string:
				values = []string{v}
			case []string:
				values = v
			case bool:
				values = []string{fmt.Sprintf("%t", v)}
			case int, uint, uint8, uint16, uint32, uint64, int8, int16, int32, int64:
				values = []string{fmt.Sprintf("%d", v)}
			}
			for _, value := range values {
				t.Add(key, value)
			}
		}
	}
}

func (o Params) String() string {
	if len(o) == 0 {
		return ""
	}
	return fmt.Sprintf("%v", map[string]interface{}(o))
}

// Param sets a single key/value pair as a query parameter.
func Param(key string, value interface{}) Option {
	return Params{key: value}
}

// Rev is a convenience function to set the revision. A less verbose alternative
// to Param("rev", rev).
func Rev(rev string) Option {
	return Params{"rev": rev}
}

// IncludeDocs instructs the query to include documents. A less verbose
// alternative to Param("include_docs", true).
func IncludeDocs() Option {
	return Params{"include_docs": true}
}
