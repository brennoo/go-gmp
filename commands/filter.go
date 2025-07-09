package commands

type Filter struct {
	Term string `xml:"term"`
}

func NewFilter(term string) *Filter {
	return &Filter{Term: term}
}
