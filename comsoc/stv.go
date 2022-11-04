package comsoc

import "td5/vtypes"

func StvSWF(p vtypes.Profile) (count vtypes.Count, err error) {
	count = make(vtypes.Count)
	err = CheckProfile(p)

	if err == nil {
		var outCandidates []vtypes.Alternative
		for i := 0; i < len(p[0])-1; i++ {
			var low, high vtypes.Alternative
			var c vtypes.Count = make(vtypes.Count)
			for _, prefs := range p {
				for _, alt := range prefs {
					if !contains(outCandidates, alt) {
						c[alt]++
						low = alt
						high = alt
						break
					}
				}
			}
			for k, v := range c {
				if v > c[high] {
					high = k
				} else if v < c[low] {
					low = k
				}
			}
			if c[high] >= len(c)/2 {
				count[high] = i + 2
				for _, alt := range p[0] {
					if count[alt] == 0 {
						count[alt] = i + 1
					}
				}
				break
			} else {
				count[low] = i + 1
				outCandidates = append(outCandidates, low)
			}
		}
	}
	return count, err
}

func StvSCF(p vtypes.Profile) (bestAlts []vtypes.Alternative, err error) {
	count, err := StvSWF(p)

	if err == nil {
		bestAlts = maxCount(count)
	}
	return bestAlts, err
}
