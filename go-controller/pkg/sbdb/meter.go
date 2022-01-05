// Code generated by "libovsdb.modelgen"
// DO NOT EDIT.

package sbdb

import "github.com/ovn-org/libovsdb/model"

type (
	MeterUnit = string
)

var (
	MeterUnitKbps  MeterUnit = "kbps"
	MeterUnitPktps MeterUnit = "pktps"
)

// Meter defines an object in Meter table
type Meter struct {
	UUID  string    `ovsdb:"_uuid"`
	Bands []string  `ovsdb:"bands"`
	Name  string    `ovsdb:"name"`
	Unit  MeterUnit `ovsdb:"unit"`
}

func copyMeterBands(a []string) []string {
	if a == nil {
		return nil
	}
	b := make([]string, len(a))
	copy(b, a)
	return b
}

func equalMeterBands(a, b []string) bool {
	if (a == nil) != (b == nil) {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if b[i] != v {
			return false
		}
	}
	return true
}

func (a *Meter) DeepCopyInto(b *Meter) {
	*b = *a
	b.Bands = copyMeterBands(a.Bands)
}

func (a *Meter) DeepCopy() *Meter {
	b := new(Meter)
	a.DeepCopyInto(b)
	return b
}

func (a *Meter) CloneModelInto(b model.Model) {
	c := b.(*Meter)
	a.DeepCopyInto(c)
}

func (a *Meter) CloneModel() model.Model {
	return a.DeepCopy()
}

func (a *Meter) Equals(b *Meter) bool {
	return a.UUID == b.UUID &&
		equalMeterBands(a.Bands, b.Bands) &&
		a.Name == b.Name &&
		a.Unit == b.Unit
}

func (a *Meter) EqualsModel(b model.Model) bool {
	c := b.(*Meter)
	return a.Equals(c)
}

var _ model.CloneableModel = &Meter{}
var _ model.ComparableModel = &Meter{}
