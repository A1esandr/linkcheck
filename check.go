package linkcheck

type (
	checker struct {
	}

	Checker interface {
		Check(url string) (map[string]string, error)
	}
)
