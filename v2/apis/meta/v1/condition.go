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

package v1

// IsStale reports whether the condition observed a generation other than
// that of the object it is of. A generation of 0 is none: a condition that
// observed none, as the wire can't tell it from one that observed 0, is
// fresh, as is one of an object the API server hasn't seen, or not of an
// object at all.
//
// +die:shared
func (d *ConditionDie) IsStale() bool {
	observed := d.ObservedGeneration().Get()
	if observed == 0 {
		return false
	}
	root, ok := d.DieRoot().(Object)
	if !ok {
		// not rooted in an object, presume fresh
		return false
	}
	generation := root.GetGeneration()
	return generation != 0 && generation != observed
}

// IsTrue reports whether the condition's status is True.
//
// +die:shared
func (d *ConditionDie) IsTrue() bool {
	s := d.Status()
	return s.Has() && s.Get() == ConditionTrue
}

// IsFalse reports whether the condition's status is False.
//
// +die:shared
func (d *ConditionDie) IsFalse() bool {
	s := d.Status()
	return s.Has() && s.Get() == ConditionFalse
}

// IsUnknown reports whether the condition's status is neither True nor False.
//
// +die:shared
func (d *ConditionDie) IsUnknown() bool {
	return !d.IsTrue() && !d.IsFalse()
}
