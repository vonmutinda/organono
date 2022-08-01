package forms

type Filter struct {
	Page   int
	Per    int
	From   string
	To     string
	Term   string
	Status string
}

func (f *Filter) NoPagination() *Filter {
	return &Filter{
		From: f.From,
		To:   f.To,
	}
}
