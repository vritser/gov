package gov

func cleanHeaderFlags(h string) string {
	for i, c := range h {
		if c == ' ' || c == ';' {
			return h[:i]
		}
	}

	return h
}
