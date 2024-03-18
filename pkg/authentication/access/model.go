package access

type rule struct {
	PType string
	V0    string
	V1    string
	V2    string
	V3    string
	V4    string
	V5    string
}

func (rule rule) Data() []string {
	s := []string{rule.PType, rule.V0, rule.V1, rule.V2, rule.V3, rule.V4, rule.V5}
	data := make([]string, 0, maxParameterCount)

	for _, val := range s {
		if val == "" {
			break
		}
		data = append(data, val)
	}

	return data
}

// Filter define the filtering rules for a FilteredAdapter's policy.
// Empty values are ignored, but all others must match the Filter.
type Filter struct {
	PType []string
	V0    []string
	V1    []string
	V2    []string
	V3    []string
	V4    []string
	V5    []string
}

type filterData struct {
	fieldName string
	arg       []string
}

func (filter Filter) genData() [maxParameterCount]filterData {
	return [maxParameterCount]filterData{
		{"p_type", filter.PType},
		{"v0", filter.V0},
		{"v1", filter.V1},
		{"v2", filter.V2},
		{"v3", filter.V3},
		{"v4", filter.V4},
		{"v5", filter.V5},
	}
}
