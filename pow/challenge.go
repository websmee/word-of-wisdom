package pow

import (
	"crypto/sha256"
	"math"
	"math/big"
)

const (
	ChallengeSize = 33
	SolutionSize  = 8
)

func Make(base string, difficulty int) []byte {
	hash := sha256.Sum256([]byte(base))
	return append(hash[:], byte(difficulty))
}

func Solve(challenge []byte) int64 {
	hash, difficulty := unpack(challenge)
	target := getTarget(difficulty)

	var nonce int64
	for nonce = 0; nonce < math.MaxInt64; nonce++ {
		if check(target, getNext(hash, nonce)) {
			return nonce
		}
	}

	return -1
}

func Verify(challenge []byte, solution int64) bool {
	hash, difficulty := unpack(challenge)
	target := getTarget(difficulty)

	return check(target, getNext(hash, solution))
}

func unpack(challenge []byte) ([]byte, int) {
	return challenge[:32], int(challenge[32])
}

func getTarget(difficulty int) *big.Int {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-difficulty))

	return target
}

func getNext(hash []byte, nonce int64) *big.Int {
	testNum := big.NewInt(0)
	testNum.Add(testNum.SetBytes(hash), big.NewInt(nonce))
	testHash := sha256.Sum256(testNum.Bytes())
	testNum.SetBytes(testHash[:])

	return testNum
}

func check(target, test *big.Int) bool {
	return target.Cmp(test) > 0
}
