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
	"github.com/go-kivik/kivik/v4/internal"
)

const (
	// ErrClientClosed is returned by any client operations after [Client.Close]
	// has been called.
	ErrClientClosed = internal.CompositeError("\x67client closed")
	// ErrDatabaseClosed is returned by any database operations after [DB.Close]
	// has been called.
	ErrDatabaseClosed = internal.CompositeError("\x67database closed")

	// Various not-implemented errors, that are returned, but don't need to be exposed directly.
	findNotImplemented        = internal.CompositeError("\x65driver does not support Find interface")
	clusterNotImplemented     = internal.CompositeError("\x65driver does not support cluster operations")
	openRevsNotImplemented    = internal.CompositeError("\x65driver does not support OpenRevs interface")
	securityNotImplemented    = internal.CompositeError("\x65driver does not support Security interface")
	configNotImplemented      = internal.CompositeError("\x65driver does not support Config interface")
	replicationNotImplemented = internal.CompositeError("\x65driver does not support replication")
)

// HTTPStatus returns the HTTP status code embedded in the error, or 500
// (internal server error), if there was no specified status code.  If err is
// nil, HTTPStatus returns 0. This provides a convenient way to determine the
// precise nature of a Kivik-returned error.
//
// For example, to panic for all but NotFound errors:
//
//	err := db.Get(context.TODO(), "docID").ScanDoc(&doc)
//	if kivik.HTTPStatus(err) == http.StatusNotFound {
//	    return
//	}
//	if err != nil {
//	    panic(err)
//	}
//
// This method uses the statusCoder interface, which is not exported by this
// package, but is considered part of the stable public API.  Driver
// implementations are expected to return errors which conform to this
// interface.
//
//	type statusCoder interface {
//	    HTTPStatus() int
//	}
func HTTPStatus(err error) int {
	return internal.HTTPStatus(err)
}
