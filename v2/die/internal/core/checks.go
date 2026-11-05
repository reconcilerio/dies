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

package core

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Checks are the checks of a value from its validation markers, as die tag
// options: Kubernetes' declarative validation (+k8s:minimum=0) and
// kubebuilder's (+kubebuilder:validation:Minimum=0) alike.
type Checks struct {
	// Enum is the values a string may have.
	Enum []string
	// Format is a string's format: one of Kubernetes', such as
	// k8s-short-name, or of OpenAPI's, such as date-time.
	Format string
	// Pattern is a regular expression a string must match.
	Pattern *regexp.Regexp

	Minimum, Maximum                   *float64
	ExclusiveMinimum, ExclusiveMaximum bool
	MultipleOf                         *float64

	MinLength, MaxLength, MaxBytes *int
	MinItems, MaxItems             *int
	MinProperties, MaxProperties   *int
	// Unique is set for a list whose items must be unique, such as a set.
	Unique bool

	// Immutable is set for a value that may not change in an update.
	Immutable bool
	// Update are the constraints of an update: NoSet, NoUnset, NoModify,
	// NoAddItem and NoRemoveItem.
	Update []string

	// Each are the checks of each item of a list, or value of a map.
	Each *Checks
}

// parse sets a check from a die tag option: enum=["a","b"], format=…,
// pattern="…", minimum=0, exclusiveMinimum, maxLength=63, unique, or any of
// them prefixed with each. for the items of a list or values of a map.
func (c *Checks) parse(k, v string) error {
	if rest, ok := strings.CutPrefix(k, "each."); ok {
		if c.Each == nil {
			c.Each = &Checks{}
		}
		return c.Each.parse(rest, v)
	}
	switch k {
	case "enum":
		return json.Unmarshal([]byte(v), &c.Enum)
	case "format":
		c.Format = v
	case "pattern":
		var p string
		if err := json.Unmarshal([]byte(v), &p); err != nil {
			return err
		}
		re, err := regexp.Compile(p)
		if err != nil {
			return err
		}
		c.Pattern = re
	case "minimum":
		return parseNumber(v, &c.Minimum)
	case "maximum":
		return parseNumber(v, &c.Maximum)
	case "multipleOf":
		return parseNumber(v, &c.MultipleOf)
	case "exclusiveMinimum":
		c.ExclusiveMinimum = true
	case "exclusiveMaximum":
		c.ExclusiveMaximum = true
	case "minLength":
		return parseCount(v, &c.MinLength)
	case "maxLength":
		return parseCount(v, &c.MaxLength)
	case "maxBytes":
		return parseCount(v, &c.MaxBytes)
	case "minItems":
		return parseCount(v, &c.MinItems)
	case "maxItems":
		return parseCount(v, &c.MaxItems)
	case "minProperties":
		return parseCount(v, &c.MinProperties)
	case "maxProperties":
		return parseCount(v, &c.MaxProperties)
	case "unique":
		c.Unique = true
	case "immutable":
		c.Immutable = true
	case "update":
		for u := range strings.SplitSeq(v, ",") {
			switch u {
			case "NoSet", "NoUnset", "NoModify", "NoAddItem", "NoRemoveItem":
				c.Update = append(c.Update, u)
			default:
				return fmt.Errorf("unknown update constraint %q", u)
			}
		}
	default:
		return fmt.Errorf("unknown option")
	}
	return nil
}

func parseNumber(v string, to **float64) error {
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return err
	}
	*to = &f
	return nil
}

func parseCount(v string, to **int) error {
	n, err := strconv.Atoi(v)
	if err != nil {
		return err
	}
	*to = &n
	return nil
}
