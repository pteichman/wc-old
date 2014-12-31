package wc

import (
	"math"
	"math/rand"

	"yasty.org/peter/wc/ecs"
)

// point is a coordinate in oilspace.
type Point struct {
	X, Y int
}

type site struct {
	prob, cost, tax, depth int

	Oil bool `json:"oil"`
}

type field struct {
	W     int `json:"w"`
	H     int `json:"h"`
	Sites []site
}

func makeField(world *ecs.World, w, h int) ecs.Entity {
	f := newField(w, h)
	for _, site := range f.Sites {
		world.AddTag(world.NewEntity(), siteTag, site)
	}

	return world.NewEntity()
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

func newField(w, h int) field {
	f := field{W: w, H: h, Sites: make([]site, w*h)}

	f.fill(
		peakSpec{min: 1, max: 5, decay: 0.15, fuzz: 0.15},
		func(s *site, v float64) {
			s.prob = int(100 * v)
			s.Oil = rand.Float64() < v
		},
	)

	f.fill(
		peakSpec{min: 5, max: 10, decay: 0.1, fuzz: 0.25},
		func(s *site, v float64) { s.cost = int(v) },
	)

	f.fill(
		peakSpec{min: 1, max: 1, decay: 0.1, fuzz: 0.5},
		func(s *site, v float64) { s.depth = int((1.0 - v) * 10.0) },
	)

	f.fill(
		peakSpec{min: 10, max: 20, decay: 0.1, fuzz: 0.5},
		func(s *site, v float64) { s.tax = int(v) },
	)

	return f
}

type peakSpec struct {
	min int
	max int

	decay float64
	fuzz  float64
}

func (f field) fill(spec peakSpec, fill func(s *site, v float64)) {
	peaks := randPeaks(randRange(spec.min, spec.max+1), f.W, f.H)

	for y := 0; y < f.H; y++ {
		for x := 0; x < f.W; x++ {
			// Start with the distance to the nearest peak.
			v, idx := closestPeak(peaks, Point{x, y})

			// Convert to a ratio of the distance to the
			// longest diagonal distance in the field.
			v /= dist(Point{0, 0}, Point{f.W, f.H})

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

			fill(&f.Sites[x+f.W*y], v)
		}
	}
}

type reservoir struct {
	oil int
	loc []Point
}

type link struct {
	p, q Point
}

func linkedSites(s []site, w, h int) []link {
	var links []link

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			site := s[x+w*y]

			if !site.Oil {
				continue
			}

			// See if we should link to our neighbor above.
			if y > 0 {
				up := s[x+w*(y-1)]
				if up.Oil && up.depth == site.depth {
					l := link{Point{x, y}, Point{x, y - 1}}
					links = append(links, l)
				}
			}

			// See if we should link to our neighbor to the left.
			if x > 0 {
				left := s[(x-1)+w*y]
				if left.Oil && left.depth == site.depth {
					l := link{Point{x, y}, Point{x - 1, y}}
					links = append(links, l)
				}
			}
		}
	}

	return links
}
