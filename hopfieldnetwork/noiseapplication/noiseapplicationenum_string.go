// Code generated by "stringer -type NoiseApplicationEnum"; DO NOT EDIT.

package noiseapplication

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[None-0]
	_ = x[ExactRatioInversion-1]
	_ = x[UniformRatioInversion-2]
	_ = x[GaussianApplication-3]
}

const _NoiseApplicationEnum_name = "NoneExactRatioInversionUniformRatioInversionGaussianApplication"

var _NoiseApplicationEnum_index = [...]uint8{0, 4, 23, 44, 63}

func (i NoiseApplicationEnum) String() string {
	if i < 0 || i >= NoiseApplicationEnum(len(_NoiseApplicationEnum_index)-1) {
		return "NoiseApplicationEnum(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _NoiseApplicationEnum_name[_NoiseApplicationEnum_index[i]:_NoiseApplicationEnum_index[i+1]]
}
