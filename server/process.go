package main

import (
	"github.com/robxu9/kahinah/kahinah"
)

func AdvisoryProcessFunc(a *kahinah.Advisory) kahinah.AdvisoryStatus {
	total := 0

	list := make(map[int64]int)

	for _, v := range a.Comments {
		switch v.Verdict {
		case kahinah.NEUTRAL:
			list[v.UserId] = 0
		case kahinah.NO:
			list[v.UserId] = config.Karma.AddFailKarma
		case kahinah.YES:
			list[v.UserId] = config.Karma.AddPassKarma
		case kahinah.BLOCK:
			list[v.UserId] = config.Karma.AddBlockKarma
		case kahinah.OVERRIDE:
			list[v.UserId] = config.Karma.AddOverrideKarma
		}
	}

	for _, v := range list {
		total += v
	}

	if total >= config.Karma.PassLimit {
		return kahinah.PASS
	} else if total <= config.Karma.FailLimit {
		return kahinah.FAIL
	}

	return kahinah.OPEN
}
