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

package server

import (
	"net/http"

	"gitlab.com/flimzy/httpe"
)

func (s *Server) startSession() httpe.HandlerWithError {
	return httpe.HandlerWithErrorFunc(func(w http.ResponseWriter, r *http.Request) error {
		var req struct {
			Name     *string `json:"name" form:"name"`
			Password string  `json:"password" form:"password"`
		}
		if err := s.bind(r, &req); err != nil {
			return err
		}
		if req.Name == nil {
			return &couchError{status: http.StatusBadRequest, Err: "bad_request", Reason: "request body must contain a username"}
		}
		return nil
	})
}
