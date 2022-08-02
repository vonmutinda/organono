package forms

type Filter struct {
	Page   int
	Per    int
	Term   string
	Status string
}

func (f *Filter) NoPagination() *Filter {
	return &Filter{
		Term:   f.Term,
		Status: f.Status,
	}
}
