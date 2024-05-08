package controller

type ValuesView interface {
	Get(key string) string
	Has(key string) bool
}

type ParametersView struct {
	Data map[string]string
}

func (p *ParametersView) Get(key string) string {
	return p.Data[key]
}
func (p *ParametersView) Has(key string) bool {
	_, ok := p.Data[key]
	return ok
}
