package comsoc

import "td5/vtypes"

func BordaSWF(p vtypes.Profile) (count vtypes.Count, err error) {
	count = make(vtypes.Count)
	err = CheckProfile(p)

	if err == nil {
		for _, alts := range p {
			for i, alt := range alts {
				count[alt] += len(alts) - 1 - i
			}
		}
	}
	return count, err
}

func BordaSCF(p vtypes.Profile) (bestAlts []vtypes.Alternative, err error) {
	count, err := BordaSWF(p)

	if err == nil {
		bestAlts = maxCount(count)
	}
	return bestAlts, err
}
