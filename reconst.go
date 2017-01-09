package reedsolomon

import (
	"sort"
)

// dp : data+parity shards, all shards size must be equal
// lost : row number in dp
func (r *rs) Reconst(dp matrix, lost []int, repairParity bool) error {
	if len(dp) != r.shards {
		return ErrTooFewShards
	}
	size, err := checkShardSize(dp)
	if err != nil {
		return err
	}
	if len(lost) == 0 {
		return nil
	}
	if len(lost) > r.parity {
		return ErrTooFewShards
	}
	dataLost, parityLost := splitLost(lost, r.data, r.parity)
	sort.Ints(dataLost)
	sort.Ints(parityLost)
	dataRepaired := newMatrix(len(dataLost), size)
	if len(dataLost) > 0 {
		err = reconstData(r.m, dp, dataRepaired, dataLost, parityLost, r.data, size)
		if err != nil {
			return err
		}
	}
	for i, l := range dataLost {
		dp[l] = dataRepaired[i]
	}
	parityRepaired := newMatrix(len(parityLost), size)
	if len(parityLost) > 0 && repairParity {
		reconstParity(r.m, dp, parityRepaired, parityLost, r.data, r.parity, size)
	}
	for i, l := range parityLost {
		dp[l] = parityRepaired[i]
	}
	return nil
}

func reconstData(encodeMatrix, dp, output matrix, dataLost, parityLost []int, numData, size int) error {
	decodeMatrix := newMatrix(numData, numData)
	survivedDP := newMatrix(numData, size)
	numShards := len(encodeMatrix)
	// fill with survived data
	for i := 0; i < numData; i++ {
		if survived(i, dataLost) {
			decodeMatrix[i] = encodeMatrix[i]
			survivedDP[i] = dp[i]
		}
	}
	// "borrow" from survived parity
	k := numData
	for _, dl := range dataLost {
		for j := k; j < numShards; j++ {
			if survived(j, parityLost) {
				decodeMatrix[dl] = encodeMatrix[j]
				survivedDP[dl] = dp[j]
				k++
				break
			}
		}
	}
	var err error
	decodeMatrix, err = decodeMatrix.invert()
	if err != nil {
		return err
	}
	numdl := len(dataLost)
	gen := newMatrix(numdl, numData)
	for i, l := range dataLost {
		gen[i] = decodeMatrix[l]
	}
	encodeRunner(gen, survivedDP, output, numData, numdl, size)
	return nil
}

func reconstParity(encodeMatrix, dp, output matrix, parityLost []int, numData, numParity, size int) {
	subGen := newMatrix(len(parityLost), numData)
	for i := range subGen {
		l := parityLost[i]
		subGen[i] = encodeMatrix[l]
	}
	encodeRunner(subGen, dp[:numData], output, numData, len(parityLost), size)
}

func splitLost(lost []int, d, p int) ([]int, []int) {
	var dataLost []int
	var parityLost []int
	for _, l := range lost {
		if l < d {
			dataLost = append(dataLost, l)
		} else {
			parityLost = append(parityLost, l)
		}
	}
	return dataLost, parityLost
}

func survived(i int, lost []int) bool {
	for _, l := range lost {
		if i == l {
			return false
		}
	}
	return true
}
