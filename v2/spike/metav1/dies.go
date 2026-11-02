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

package metav1

import (
	"reconciler.io/dies/v2/die"
	"reconciler.io/dies/v2/die/base"
	"reconciler.io/dies/v2/die/fields"
)

// ObjectMeta

type ObjectMetaDie struct {
	base.DieBase[ObjectMeta, *ObjectMetaDie]
}

func (*ObjectMetaDie) DieFrom(b base.DieBase[ObjectMeta, *ObjectMetaDie]) *ObjectMetaDie {
	return &ObjectMetaDie{b}
}

var ObjectMetaBlank = base.Blank[ObjectMeta, *ObjectMetaDie]()

var objectMetaFields = struct {
	Name        fields.Field[ObjectMeta, string]
	Namespace   fields.Field[ObjectMeta, string]
	Generation  fields.Field[ObjectMeta, int64]
	Labels      fields.MapField[ObjectMeta, string, string]
	Annotations fields.MapField[ObjectMeta, string, string]
}{
	Name:        fields.NewField[ObjectMeta, string]("name"),
	Namespace:   fields.NewField[ObjectMeta, string]("namespace"),
	Generation:  fields.NewField[ObjectMeta, int64]("generation"),
	Labels:      fields.NewMap[ObjectMeta, string, string]("labels"),
	Annotations: fields.NewMap[ObjectMeta, string, string]("annotations"),
}

func (d *ObjectMetaDie) Name() die.Value[string, *ObjectMetaDie] {
	return fields.ValueOf(d, objectMetaFields.Name)
}

func (d *ObjectMetaDie) Namespace() die.Value[string, *ObjectMetaDie] {
	return fields.ValueOf(d, objectMetaFields.Namespace)
}

func (d *ObjectMetaDie) Generation() die.Value[int64, *ObjectMetaDie] {
	return fields.ValueOf(d, objectMetaFields.Generation)
}

func (d *ObjectMetaDie) Labels() *die.Map[string, string] {
	return fields.MapOf(d, objectMetaFields.Labels)
}

func (d *ObjectMetaDie) Annotations() *die.Map[string, string] {
	return fields.MapOf(d, objectMetaFields.Annotations)
}

// Condition

type ConditionDie struct {
	base.DieBase[Condition, *ConditionDie]
}

func (*ConditionDie) DieFrom(b base.DieBase[Condition, *ConditionDie]) *ConditionDie {
	return &ConditionDie{b}
}

var ConditionBlank = base.Blank[Condition, *ConditionDie]()

var conditionFields = struct {
	Type               fields.Field[Condition, string]
	Status             fields.Field[Condition, string]
	ObservedGeneration fields.Field[Condition, int64]
	Reason             fields.Field[Condition, string]
	Message            fields.Field[Condition, string]
}{
	Type:               fields.NewField[Condition, string]("type"),
	Status:             fields.NewField[Condition, string]("status"),
	ObservedGeneration: fields.NewField[Condition, int64]("observedGeneration"),
	Reason:             fields.NewField[Condition, string]("reason"),
	Message:            fields.NewField[Condition, string]("message"),
}

func (d *ConditionDie) Type() die.Value[string, *ConditionDie] {
	return fields.ValueOf(d, conditionFields.Type)
}

func (d *ConditionDie) Status() die.Value[string, *ConditionDie] {
	return fields.ValueOf(d, conditionFields.Status)
}

func (d *ConditionDie) ObservedGeneration() die.Value[int64, *ConditionDie] {
	return fields.ValueOf(d, conditionFields.ObservedGeneration)
}

func (d *ConditionDie) Reason() die.Value[string, *ConditionDie] {
	return fields.ValueOf(d, conditionFields.Reason)
}

func (d *ConditionDie) Message() die.Value[string, *ConditionDie] {
	return fields.ValueOf(d, conditionFields.Message)
}

// True sets the condition to "True" with a reason.
func (d *ConditionDie) True(reason string) *ConditionDie {
	return d.Status().Set("True").Reason().Set(reason)
}
