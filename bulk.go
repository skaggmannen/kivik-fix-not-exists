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
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/go-kivik/kivik/v4/driver"
)

// BulkResults is an iterator over the results of a BulkDocs query.
type BulkResults struct {
	*iter
}

type bulkIterator struct{ driver.BulkResults }

var _ iterator = &bulkIterator{}

func (r *bulkIterator) Next(i interface{}) error {
	return r.BulkResults.Next(i.(*driver.BulkResult))
}

func newBulkResults(ctx context.Context, onClose func(), bulki driver.BulkResults) *BulkResults {
	return &BulkResults{
		iter: newIterator(ctx, onClose, &bulkIterator{bulki}, &driver.BulkResult{}),
	}
}

// ID returns the document ID name for the current result.
func (r *BulkResults) ID() string {
	runlock, err := r.rlock()
	if err != nil {
		return ""
	}
	defer runlock()
	return r.curVal.(*driver.BulkResult).ID
}

// Rev returns the revision of the current curResult.
func (r *BulkResults) Rev() string {
	runlock, err := r.rlock()
	if err != nil {
		return ""
	}
	defer runlock()
	return r.curVal.(*driver.BulkResult).Rev
}

// UpdateErr returns the error associated with the current result, or nil
// if none. Do not confuse this with [BulkResults.Err], which returns an error
// for the iterator itself.
func (r *BulkResults) UpdateErr() error {
	runlock, err := r.rlock()
	if err != nil {
		return nil
	}
	defer runlock()
	return r.curVal.(*driver.BulkResult).Error
}

// BulkDocs allows you to create and update multiple documents at the same time
// within a single request. This function returns an iterator over the results
// of the bulk operation.
//
// See http://docs.couchdb.org/en/2.0.0/api/database/bulk-api.html#db-bulk-docs
//
// As with [DB.Put], each individual document may be a JSON-marshable object, or
// a raw JSON string in a [encoding/json.RawMessage], or [io.Reader].
func (db *DB) BulkDocs(ctx context.Context, docs []interface{}, options ...Options) *BulkResults {
	if db.err != nil {
		return &BulkResults{errIterator(db.err)}
	}
	docsi, err := docsInterfaceSlice(docs)
	if err != nil {
		return &BulkResults{errIterator(err)}
	}
	if len(docsi) == 0 {
		return &BulkResults{errIterator(&Error{Status: http.StatusBadRequest, Err: errors.New("kivik: no documents provided")})}
	}
	if err := db.startQuery(); err != nil {
		return &BulkResults{errIterator(err)}
	}
	opts := mergeOptions(options...)
	if bulkDocer, ok := db.driverDB.(driver.BulkDocer); ok {
		bulki, err := bulkDocer.BulkDocs(ctx, docsi, opts)
		if err != nil {
			return &BulkResults{errIterator(err)}
		}
		return newBulkResults(ctx, db.endQuery, bulki)
	}
	var results []driver.BulkResult
	for _, doc := range docsi {
		var err error
		var id, rev string
		if docID, ok := extractDocID(doc); ok {
			id = docID
			rev, err = db.Put(ctx, id, doc, opts)
		} else {
			id, rev, err = db.CreateDoc(ctx, doc, opts)
		}
		results = append(results, driver.BulkResult{
			ID:    id,
			Rev:   rev,
			Error: err,
		})
	}
	return newBulkResults(ctx, db.endQuery, &emulatedBulkResults{results})
}

type emulatedBulkResults struct {
	results []driver.BulkResult
}

var _ driver.BulkResults = &emulatedBulkResults{}

func (r *emulatedBulkResults) Close() error {
	r.results = nil
	return nil
}

func (r *emulatedBulkResults) Next(res *driver.BulkResult) error {
	if len(r.results) == 0 {
		return io.EOF
	}
	*res = r.results[0]
	r.results = r.results[1:]
	return nil
}

func docsInterfaceSlice(docsi []interface{}) ([]interface{}, error) {
	for i, doc := range docsi {
		x, err := normalizeFromJSON(doc)
		if err != nil {
			return nil, &Error{Status: http.StatusBadRequest, Err: err}
		}
		docsi[i] = x
	}
	return docsi, nil
}
