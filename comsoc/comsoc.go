package comsoc

import (
	"errors"
	"td5/vtypes"
)

func rank(alt vtypes.Alternative, prefs []vtypes.Alternative) int {
	for i := 0; i < len(prefs); i++ {
		if prefs[i] == alt {
			return i
		}
	}
	return -1
}

func maxCount(count vtypes.Count) (bestAlts []vtypes.Alternative) {
	max := -1
	for k := range count {
		if count[k] > max {
			max = count[k]
			bestAlts = make([]vtypes.Alternative, 1)
			bestAlts[0] = k
		} else if count[k] == max {
			bestAlts = append(bestAlts, k)
		}
	}
	return bestAlts
}

func CheckProfile(prefs vtypes.Profile) error {
	nbrAlts := len(prefs[0])
	for i := 1; i < len(prefs); i++ {
		if len(prefs[i]) != nbrAlts {
			return errors.New("Profil incomplet")
		}

		for j := 0; j < len(prefs[i]); j++ {
			for k := j + 1; k < len(prefs[i]); k++ {
				if prefs[i][j] == prefs[i][k] {
					return errors.New("Alternative dupliquée")
				}
			}
		}
	}
	return nil
}

func CheckPrefs(prefs []vtypes.Alternative, nbrAlts int) error {
	if len(prefs) != nbrAlts {
		return errors.New("Profil incomplet ou surchargé")
	}
	for i := 0; i < len(prefs)-1; i++ {
		for j := i + 1; j < len(prefs); j++ {
			if prefs[i] == prefs[j] {
				return errors.New("Alternative dupliquée")
			}
		}
	}
	return nil
}

func SortCount(count vtypes.Count) (alts []vtypes.Alternative) {
	alts = make([]vtypes.Alternative, len(count))
	var i int = 0
	for k := range count {
		alts[i] = k
		i++
	}
	var permut bool = true
	var passage int = 0
	for permut {
		permut = false
		passage += 1
		for i := 0; i < len(count)-passage; i++ {
			if count[alts[i]] < count[alts[i+1]] {
				permut = true
				temp := alts[i]
				alts[i] = alts[i+1]
				alts[i+1] = temp
			}
		}
	}
	return alts
}

func contains(alts []vtypes.Alternative, alt vtypes.Alternative) bool {
	for _, a := range alts {
		if a == alt {
			return true
		}
	}
	return false
}
