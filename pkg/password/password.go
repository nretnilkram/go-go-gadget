package password

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const (
	lowerChars  = "abcdefghijklmnopqrstuvwxyz"
	upperChars  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digitChars  = "0123456789"
	symbolChars = "!@#$%^&*()-_=+[]{}|;:,.<>?"
)

// Weight controls the relative frequency of each character class.
// A class with weight 0 is omitted from the generated password.
type Weight struct {
	Lower  int
	Upper  int
	Digit  int
	Symbol int
}

// Generate returns a cryptographically random password of the given length.
func Generate(length int, weight Weight) (string, error) {
	if length < 1 {
		return "", fmt.Errorf("length must be at least 1")
	}

	pools, weights, err := buildPools(weight)
	if err != nil {
		return "", err
	}

	totalWeight := 0
	for _, w := range weights {
		totalWeight += w
	}

	out := make([]byte, length)
	for i := range out {
		poolIdx, err := weightedIndex(weights, totalWeight)
		if err != nil {
			return "", err
		}
		ch, err := randomByte(pools[poolIdx])
		if err != nil {
			return "", err
		}
		out[i] = ch
	}
	return string(out), nil
}

func buildPools(weight Weight) ([]string, []int, error) {
	pools := make([]string, 0, 4)
	weights := make([]int, 0, 4)

	add := func(w int, chars string) {
		if w > 0 {
			pools = append(pools, chars)
			weights = append(weights, w)
		}
	}
	add(weight.Lower, lowerChars)
	add(weight.Upper, upperChars)
	add(weight.Digit, digitChars)
	add(weight.Symbol, symbolChars)

	if len(pools) == 0 {
		return nil, nil, fmt.Errorf("at least one character class weight must be positive")
	}
	return pools, weights, nil
}

func weightedIndex(weights []int, total int) (int, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(total)))
	if err != nil {
		return 0, fmt.Errorf("select character class: %w", err)
	}
	v := int(n.Int64())
	for i, w := range weights {
		v -= w
		if v < 0 {
			return i, nil
		}
	}
	return len(weights) - 1, nil
}

func randomByte(pool string) (byte, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(pool))))
	if err != nil {
		return 0, fmt.Errorf("select character: %w", err)
	}
	return pool[n.Int64()], nil
}
