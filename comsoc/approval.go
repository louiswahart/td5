package comsoc

import "td5/vtypes"

func ApprovalSWF(p vtypes.Profile, thresholds []int) (count vtypes.Count, err error) {
	count = make(vtypes.Count)
	err = CheckProfile(p)
	if err == nil {
		for i, prof := range p {
			for j, val := range prof {
				if j < thresholds[i] {
					count[val] += 1
				}
			}
		}
	}
	return
}
func ApprovalSCF(p vtypes.Profile, thresholds []int) (bestAlts []vtypes.Alternative, err error) {
	count, err := BordaSWF(p)
	if err == nil {
		bestAlts = maxCount(count)
	}
	return
}
