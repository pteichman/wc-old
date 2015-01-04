package wc

import (
	"math"
	"math/rand"
)

type OilMap struct {
	Size
	Locs []Secret
}

func (om *OilMap) Loc(x, y int) *Secret {
	return &om.Locs[om.Index(x, y)]
}

type Secret struct {
	Survey
	Oil       bool
	OilDepth  int8
	Reservoir *reservoir
}

func randPeaks(n, w, h int) []Point {
	var peaks []Point
	for i := 0; i < n; i++ {
		peaks = append(peaks, Point{X: rand.Intn(w), Y: rand.Intn(h)})
	}

	return peaks
}

func dist(p1, p2 Point) float64 {
	dx := float64(p1.X - p2.X)
	dy := float64(p1.Y - p2.Y)

	return math.Sqrt(dx*dx + dy*dy)
}

func closestPeak(peaks []Point, p Point) (float64, int) {
	var idx int
	var min = math.Inf(0)

	for i, peak := range peaks {
		d := dist(peak, p)
		if d < min {
			idx = i
			min = d
		}
	}

	return min, idx
}

func clamp(v, min, max float64) float64 {
	return math.Min(math.Max(v, min), max)
}

func newOilMap(s Size) *OilMap {
	m := &OilMap{Size: s, Locs: make([]Secret, s.W*s.H)}

	m.fill(
		peakSpec{min: 1, max: 5, decay: 0.15, fuzz: 0.15},
		func(s *Secret, v float64) {
			s.Prob = int8(100 * v)
			s.Oil = rand.Float64() < v

			if s.Oil {
				s.Reservoir = &reservoir{oil: int(math.Max(0.1, gauss(1, 1)) * 666)}
			}
		},
	)

	m.fill(
		peakSpec{min: 5, max: 10, decay: 0.1, fuzz: 0.25},
		func(s *Secret, v float64) { s.DrillCost = int(v) },
	)

	m.fill(
		peakSpec{min: 1, max: 1, decay: 0.1, fuzz: 0.5},
		func(s *Secret, v float64) { s.OilDepth = int8((1.0 - v) * 10.0) },
	)

	m.fill(
		peakSpec{min: 10, max: 20, decay: 0.1, fuzz: 0.5},
		func(s *Secret, v float64) { s.Tax = int(v) },
	)

	m.mergeReservoirs()

	return m
}

func gauss(mu, sigma float64) float64 {
	return rand.NormFloat64()*sigma + mu
}

type peakSpec struct {
	min int
	max int

	decay float64
	fuzz  float64
}

func (om *OilMap) fill(spec peakSpec, fill func(s *Secret, v float64)) {
	peaks := randPeaks(randRange(spec.min, spec.max+1), om.W, om.H)

	for y := 0; y < om.H; y++ {
		for x := 0; x < om.W; x++ {
			// Start with the distance to the nearest peak.
			v, idx := closestPeak(peaks, Point{x, y})

			// Convert to a ratio of the distance to the
			// longest diagonal distance in the field.
			v /= dist(Point{0, 0}, Point{om.W, om.H})

			// Double the value for a better input into
			// log. :/ This should be distilled to some
			// sane 0.0 to 1.0 param which allows for
			// controlling the steepness of the peaks.
			v *= 2.0

			// Logarithmically adjust the value, shifting
			// and dividing to get a nice curve roughly in
			// the range of 0 to 1.
			v = 1.0 - ((math.Log(v) + 4.0) / 4.0)

			// Adjust for subsequent peaks being
			// progressively lower.
			v *= math.Pow(1.0-spec.decay, float64(idx))

			// Apply some random fuzz to keep everyone guessing.
			v += 2.0 * (rand.Float64() - 0.5) * spec.fuzz

			// Contain the final value.
			v = clamp(v, 0.0, 1.0)

			fill(om.Loc(x, y), v)
		}
	}
}

type reservoir struct {
	oil int
}

func merge(a, b *reservoir) *reservoir {
	return &reservoir{oil: a.oil + b.oil}
}

func (om *OilMap) mergeReservoirs() {
	for y := 0; y < om.H-1; y++ {
		for x := 0; x < om.W-1; x++ {
			site := om.Loc(x, y)
			if !site.Oil {
				continue
			}

			// See if we should link to our neighbor to the right.
			rt := om.Loc(x+1, y)
			if rt.Oil && rt.OilDepth == site.OilDepth {
				r := merge(site.Reservoir, rt.Reservoir)
				rt.Reservoir = r
				site.Reservoir = r
			}

			// See if we should link to our neighbor to the bottom.
			dn := om.Loc(x, y+1)
			if dn.Oil && dn.OilDepth == site.OilDepth {
				r := merge(site.Reservoir, dn.Reservoir)
				dn.Reservoir = r
				site.Reservoir = r
			}
		}
	}
}
