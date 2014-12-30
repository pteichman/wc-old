package wc

// Single registry point for ecs.World tags used in wildcatting.

type tag int

const (
	taxTag tag = iota
	posTag
	siteTag
)
