package comsoc

import "td5/vtypes"

func CopelandSWF(p vtypes.Profile) (count vtypes.Count, err error) {
	count = make(vtypes.Count)
	err = CheckProfile(p)

	if err == nil {
		for i, candidate := range p[0][:len(p[0])-1] {
			for _, opponent := range p[0][i+1:] {
				var candidateCount, opponentCount int
				for _, prefs := range p {
					for _, alt := range prefs {
						if alt == candidate {
							candidateCount += 1
							break
						} else if alt == opponent {
							opponentCount += 1
							break
						}
					}
				}
				if candidateCount > opponentCount {
					count[candidate] += 1
				} else if opponentCount > candidateCount {
					count[opponent] += 1
				}
			}
		}
	}
	return count, err
}

func CopelandSCF(p vtypes.Profile) (bestAlts []vtypes.Alternative, err error) {
	count, err := CopelandSWF(p)

	if err == nil {
		bestAlts = maxCount(count)
	}
	return bestAlts, err
}
