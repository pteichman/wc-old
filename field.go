package wc

import (
	"math"
	"math/rand"

	"yasty.org/peter/wc/ecs"
)

// point is a coordinate in oilspace.
type point struct {
	x, y int
}

type site struct {
	prob, cost, tax, depth int

	oil bool
}

type field struct {
	w, h  int
	sites []site
}

func makeField(world *ecs.World, w, h int) ecs.Entity {
	f := newField(w, h)
	for _, site := range f.sites {
		world.AddTag(world.NewEntity(), siteTag, site)
	}

	return world.NewEntity()
}

func randPeaks(n, w, h int) []point {
	var peaks []point
	for i := 0; i < n; i++ {
		peaks = append(peaks, point{x: rand.Intn(w), y: rand.Intn(h)})
	}

	return peaks
}

func dist(p1, p2 point) float64 {
	dx := float64(p1.x - p2.x)
	dy := float64(p1.y - p2.y)

	return math.Sqrt(dx*dx + dy*dy)
}

func closestPeak(peaks []point, p point) (float64, int) {
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

func newField(w, h int) field {
	f := field{w: w, h: h, sites: make([]site, w*h)}

	f.fill(
		peakSpec{min: 1, max: 5, decay: 0.15, fuzz: 0.15},
		func(s site, v float64) {
			s.prob = int(100 * v)
			s.oil = rand.Float64() < v
		},
	)

	f.fill(
		peakSpec{min: 5, max: 10, decay: 0.1, fuzz: 0.25},
		func(s site, v float64) { s.cost = int(v) },
	)

	f.fill(
		peakSpec{min: 1, max: 1, decay: 0.1, fuzz: 0.5},
		func(s site, v float64) { s.depth = int((1.0 - v) * 10.0) },
	)

	f.fill(
		peakSpec{min: 10, max: 20, decay: 0.1, fuzz: 0.5},
		func(s site, v float64) { s.tax = int(v) },
	)

	return f
}

type peakSpec struct {
	min int
	max int

	decay float64
	fuzz  float64
}

func (f field) fill(spec peakSpec, fill func(s site, v float64)) {
	peaks := randPeaks(randRange(spec.min, spec.max+1), f.w, f.h)

	for y := 0; y < f.h; y++ {
		for x := 0; x < f.w; x++ {
			// Start with the distance to the nearest peak.
			v, idx := closestPeak(peaks, point{x, y})

			// Convert to a ratio of the distance to the
			// longest diagonal distance in the field.
			v /= dist(point{0, 0}, point{f.w, f.h})

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

			fill(f.sites[x+f.w*y], v)
		}
	}
}

type reservoir struct {
	oil int
	loc []point
}

type link struct {
	p, q point
}

func linkedSites(s []site, w, h int) []link {
	var links []link

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			site := s[x+w*y]

			if !site.oil {
				continue
			}

			// See if we should link to our neighbor above.
			if y > 0 {
				up := s[x+w*(y-1)]
				if up.oil && up.depth == site.depth {
					l := link{point{x, y}, point{x, y - 1}}
					links = append(links, l)
				}
			}

			// See if we should link to our neighbor to the left.
			if x > 0 {
				left := s[(x-1)+w*y]
				if left.oil && left.depth == site.depth {
					l := link{point{x, y}, point{x - 1, y}}
					links = append(links, l)
				}
			}
		}
	}

	return links
}
