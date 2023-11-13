package everglow

import (
	"fmt"

	"github.com/matrixorigin/simdc/simd"
)

var (
	samples int
	melY    []float64
	melX    []float64
)

type ExpFilter struct {
	alphaDecay float64
	alphaRise  float64
	value      float64
}

func NewExpFilter(val, alphaDecay, alphaRise float64) *ExpFilter {
	assert(0.0 < alphaDecay && alphaDecay < 1.0, "Invalid decay smoothing factor")
	assert(0.0 < alphaRise && alphaRise < 1.0, "Invalid rise smoothing factor")
	return &ExpFilter{
		alphaDecay: alphaDecay,
		alphaRise:  alphaRise,
		value:      val,
	}
}

func (ef *ExpFilter) Update(newValue float64) float64 {
	if ef.valueIsArray() {
		alpha := simd.If(ef.value < newValue, ef.alphaRise, ef.alphaDecay)
		ef.value = alpha*newValue + (1.0-alpha)*ef.value
	} else {
		alpha := simd.If(newValue > ef.value, ef.alphaRise, ef.alphaDecay)
		ef.value = alpha*newValue + (1.0-alpha)*ef.value
	}
	return ef.value
}

func (ef *ExpFilter) valueIsArray() bool {
	_, isArray := ef.value.([]float64)
	return isArray
}

func window(length int) float64 {
	return 1.0
}

func rfft(data []float64, window func(int) float64) ([]float64, []float64) {
	w := window(len(data))
	ys := simd.Abs(simd.FFTReal(simd.MakeGoComplex(data, make([]float64, len(data))), w))
	xs := simd.FFTFreq(len(data), 1.0/MIC_RATE)
	return xs, ys
}

func fft(data []float64, window func(int) float64) ([]float64, []float64) {
	w := window(len(data))
	ys := simd.FFT(simd.MakeGoComplex(data, make([]float64, len(data))), w)
	xs := simd.FFTFreq(len(data), 1.0/MIC_RATE)
	return xs, ys
}

func createMelBank() {
	samples = int(MIC_RATE * N_ROLLING_HISTORY / (2.0 * FPS))
	_, melY, melX = computeMelmat(N_FFT_BINS, MIN_FREQUENCY, MAX_FREQUENCY, samples, MIC_RATE)
}

func assert(condition bool, message string) {
	if !condition {
		panic(fmt.Sprintf("Assertion failed: %s", message))
	}
}
