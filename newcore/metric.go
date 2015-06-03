package newcore

type Metric string

func NewMetric(m string) Metric {
	mt := Metric(NormalizeMetricKey(m[:]))
	return mt
}

//TODO: add more tests
func (m Metric) Clean() string {
	return NormalizeMetricKey(string(m))
}

func (m Metric) CleanWithTags(tpl string, tags TagSet) (string, error) {
	return FlatMetricKeyAndTags(tpl, string(m), tags)
}
