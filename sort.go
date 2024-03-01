package main

import (
	"sort"
)

//function to sort videos

type VideoWrapper struct {
	videos []VideoReturn
	by     func(p, q *VideoReturn) bool
}

type SortBy func(p, q *VideoReturn) bool

func (vw VideoWrapper) Len() int {
	return len(vw.videos)
}

func (vw VideoWrapper) Swap(i, j int) {
	vw.videos[i], vw.videos[j] = vw.videos[j], vw.videos[i]
}

func (vw VideoWrapper) Less(i, j int) bool {
	return vw.by(&vw.videos[i], &vw.videos[j])
}

func SortVideo(videos []VideoReturn, by SortBy) {
	sort.Sort(VideoWrapper{videos, by})
}

//
