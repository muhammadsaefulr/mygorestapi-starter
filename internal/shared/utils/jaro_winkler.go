package utils

import (
	"sort"
	"strings"

	"github.com/xrash/smetrics"
)

func JaroWinklerPairIndices(aniTitles, odTitles []string, topN int) (a []int, b []int) {
	type pair struct {
		aniIdx int
		odIdx  int
		score  float64
	}
	var pairs []pair

	norm := func(s string) string {
		return strings.ToLower(strings.TrimSpace(s))
	}

	for i, ani := range aniTitles {
		nAni := norm(ani)
		if nAni == "" {
			continue
		}
		for j, od := range odTitles {
			nOd := norm(od)
			if nOd == "" {
				continue
			}
			score := smetrics.JaroWinkler(nAni, nOd, 0.7, 4)
			// bonus kecil kalau substring
			if strings.Contains(nAni, nOd) || strings.Contains(nOd, nAni) {
				score += 0.05
			}
			if score >= 0.60 {
				pairs = append(pairs, pair{i, j, score})
			}
		}
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].score > pairs[j].score
	})

	if len(pairs) > topN {
		pairs = pairs[:topN]
	}

	for _, p := range pairs {
		a = append(a, p.aniIdx)
		b = append(b, p.odIdx)
	}

	return
}
