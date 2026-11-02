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

package apis

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// URL is cloned, but its JSON methods are not, so they are written here:
// upstream encodes a URL as its string form.

func (u URL) MarshalJSON() ([]byte, error) {
	uu := url.URL(u)
	return []byte(fmt.Sprintf("%q", uu.String())), nil
}

func (u *URL) UnmarshalJSON(b []byte) error {
	var ref string
	if err := json.Unmarshal(b, &ref); err != nil {
		return err
	}
	if ref == "" {
		*u = URL{}
		return nil
	}
	parsed, err := url.Parse(ref)
	if err != nil {
		return err
	}
	*u = URL(*parsed)
	return nil
}
