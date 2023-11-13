package everglow

import "math"

func hertzToMel(freq float64) float64 {
	return 2595.0 * math.Log10(1+(freq/700.0))
}

func melToHertz(mel float64) float64 {
	return 700.0*(math.Pow(10, mel/2595.0)) - 700.0
}

func melfrequenciesMelFilterbank(numBands int, freqMin float64, freqMax float64, numFFTBands int) ([]float64, []float64, []float64) {
	melMax := hertzToMel(freqMax)
	melMin := hertzToMel(freqMin)
	deltaMel := math.Abs(melMax-melMin) / float64(numBands+1.0)
	var frequenciesMel []float64
	for i := 0; i < numBands+2; i++ {
		frequenciesMel = append(frequenciesMel, melMin+deltaMel*float64(i))
	}
	lowerEdgesMel := frequenciesMel[:len(frequenciesMel)-2]
	upperEdgesMel := frequenciesMel[2:]
	centerFrequenciesMel := frequenciesMel[1 : len(frequenciesMel)-1]
	return centerFrequenciesMel, lowerEdgesMel, upperEdgesMel
}

func computeMelmat(numMelBands int, freqMin float64, freqMax float64, numFFTBands int, sampleRate float64) ([][]float64, []float64, []float64) {
	centerFrequenciesMel, lowerEdgesMel, upperEdgesMel := melfrequenciesMelFilterbank(numMelBands, freqMin, freqMax, numFFTBands)

	var centerFrequenciesHz []float64
	var lowerEdgesHz []float64
	var upperEdgesHz []float64
	for _, mel := range centerFrequenciesMel {
		centerFrequenciesHz = append(centerFrequenciesHz, melToHertz(mel))
	}
	for _, mel := range lowerEdgesMel {
		lowerEdgesHz = append(lowerEdgesHz, melToHertz(mel))
	}
	for _, mel := range upperEdgesMel {
		upperEdgesHz = append(upperEdgesHz, melToHertz(mel))
	}

	freqs := make([]float64, numFFTBands)
	for i := 0; i < numFFTBands; i++ {
		freqs[i] = float64(i) * (sampleRate / 2.0) / float64(numFFTBands-1)
	}

	var melmat [][]float64
	for i := 0; i < numMelBands; i++ {
		row := make([]float64, numFFTBands)
		for j, freq := range freqs {
			center, lower, upper := centerFrequenciesHz[i], lowerEdgesHz[i], upperEdgesHz[i]
			leftSlope := (freq >= lower) && (freq <= center)
			rightSlope := (freq >= center) && (freq <= upper)

			if leftSlope {
				row[j] = (freq - lower) / (center - lower)
			} else if rightSlope {
				row[j] = (upper - freq) / (upper - center)
			}
		}
		melmat = append(melmat, row)
	}

	return melmat, centerFrequenciesMel, freqs
}
