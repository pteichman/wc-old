package wc

import "testing"

func TestReservoirs(t *testing.T) {
	m := newOilMap(Size{5, 5})

	for i, loc := range m.Locs {
		if loc.Oil && (loc.Reservoir == nil) {
			t.Errorf("Location needs reservoir: %d", i)
		}
	}
}
