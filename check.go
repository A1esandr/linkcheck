package linkcheck

type (
	checker struct {
	}

	Checker interface {
		Check(url string) (map[string]string, error)
	}
)

func New() Checker {
	return &checker{}
}

func (c *checker) Check(url string) (map[string]string, error) {
	return nil, nil
}
