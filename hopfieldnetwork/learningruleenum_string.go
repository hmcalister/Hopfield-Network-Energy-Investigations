// Code generated by "stringer -type LearningRuleEnum"; DO NOT EDIT.

package hopfieldnetwork

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[HebbianLearningRule-0]
	_ = x[BipolarMappedHebbianLearningRule-1]
	_ = x[DeltaLearningRule-2]
	_ = x[BipolarMappedDeltaLearningRule-3]
	_ = x[ThermalDeltaLearningRule-4]
	_ = x[BipolarMappedThermalDeltaLearningRule-5]
}

const _LearningRuleEnum_name = "HebbianLearningRuleBipolarMappedHebbianLearningRuleDeltaLearningRuleBipolarMappedDeltaLearningRuleThermalDeltaLearningRuleBipolarMappedThermalDeltaLearningRule"

var _LearningRuleEnum_index = [...]uint8{0, 19, 51, 68, 98, 122, 159}

func (i LearningRuleEnum) String() string {
	if i < 0 || i >= LearningRuleEnum(len(_LearningRuleEnum_index)-1) {
		return "LearningRuleEnum(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _LearningRuleEnum_name[_LearningRuleEnum_index[i]:_LearningRuleEnum_index[i+1]]
}
