package pow

import "testing"

func TestChallenge(t *testing.T) {
	base := "fgh5er6h5ehdfghe4r5yhg4w"
	difficulty := 5
	challenge := Make(base, difficulty)

	if len(challenge) != ChallengeSize {
		t.Errorf("invalid challenge size %d expected %d", len(challenge), ChallengeSize)
	}

	solution := Solve(challenge)
	if solution < 0 {
		t.Errorf("solution not found")
	}
	if !Verify(challenge, solution) {
		t.Errorf("verification failed")
	}
	invalidChallenge := challenge[:]
	invalidChallenge[0] = 0
	if Verify(invalidChallenge, solution) || Verify(challenge, solution-1) || Verify(challenge, solution+1) {
		t.Errorf("verification is unreliable")
	}
}

func benchmarkSolve(b *testing.B, difficulty int) {
	for i := 0; i < b.N; i++ {
		Solve(Make("fgh5er6h5ehdfghe4r5yhg4w", difficulty))
	}
}

func BenchmarkSolve5(b *testing.B)  { benchmarkSolve(b, 5) }
func BenchmarkSolve10(b *testing.B) { benchmarkSolve(b, 10) }
func BenchmarkSolve15(b *testing.B) { benchmarkSolve(b, 15) }
func BenchmarkSolve20(b *testing.B) { benchmarkSolve(b, 20) }
func BenchmarkSolve25(b *testing.B) { benchmarkSolve(b, 25) }
