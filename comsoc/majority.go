package comsoc

import "td5/vtypes"

func MajoritySWF(p vtypes.Profile) (count vtypes.Count, err error) {
	count = make(vtypes.Count)
	err = CheckProfile(p)

	if err == nil {
		for _, alt := range p[0] {
			count[alt] = 0
		}
		for _, alts := range p {
			count[alts[0]]++
		}
	}
	return count, err
}

func MajoritySCF(p vtypes.Profile) (bestAlts []vtypes.Alternative, err error) {
	count, err := MajoritySWF(p)

	if err == nil {
		bestAlts = maxCount(count)
	}
	return bestAlts, err
}
