package txqr

// Note that these CDFs (cumulative distribution function) will be used for
// selecting source blocks for code block generation. A typical algorithm
// is to choose a number from a distribution, then pick uniformly that many
// source blocks to XOR into a code block. To use CDF mapping values, pick a
// random number r (0 <= r < 1) and then find the smallest i such that
// CDF[i] >= r.

// solitonDistribution returns a CDF mapping for the soliton distribution.
// N (the number of elements in the CDF) cannot be less than 1
// The CDF is one-based: the probability of picking 1 from the distribution
// is CDF[1].
func solitonDistribution(n int) []float64 {
	cdf := make([]float64, n+1)
	cdf[1] = 1 / float64(n)
	for i := 2; i < len(cdf); i++ {
		cdf[i] = cdf[i-1] + (1 / (float64(i) * float64(i-1)))
	}
	return cdf
}

// ids create slice with IDs for 0..n values.
func ids(n int) []int64 {
	ids := make([]int64, n)
	for i := int64(0); i < int64(n); i++ {
		ids[i] = i
	}
	return ids
}
