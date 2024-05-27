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

func (vw VideoWrapper) Len() int { return len(vw.videos) }

func (vw VideoWrapper) Swap(i, j int) { vw.videos[i], vw.videos[j] = vw.videos[j], vw.videos[i] }

func (vw VideoWrapper) Less(i, j int) bool { return vw.by(&vw.videos[i], &vw.videos[j]) }

func SortVideo(videos []VideoReturn, by SortBy) { sort.Sort(VideoWrapper{videos, by}) }

//function to sort circle affairs

type CircleAffairsSliceDecrement []CircleAffairMessage

func (s CircleAffairsSliceDecrement) Len() int { return len(s) }

func (s CircleAffairsSliceDecrement) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s CircleAffairsSliceDecrement) Less(i, j int) bool { return s[i].Time > s[j].Time }

//function to sort check message

type CheckMessageSliceDecrement []CheckMessage

func (s CheckMessageSliceDecrement) Len() int { return len(s) }

func (s CheckMessageSliceDecrement) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s CheckMessageSliceDecrement) Less(i, j int) bool { return s[i].Time > s[j].Time }
