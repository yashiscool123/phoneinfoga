package filter

type Filter interface {
	IsScannerIgnored(string) bool
}

type Engine struct {
	rules []string
}

func NewFilter() *Engine {
	return &Engine{}
}

func (e *Engine) AddRule(r ...string) {
	for _, rule := range r {
		e.rules = append(e.rules, rule)
	}
}

func (e *Engine) IsScannerIgnored(r string) bool {
	for _, rule := range e.rules {
		if rule == r {
			return true
		}
	}
	return false
}
