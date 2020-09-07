package core

func arrangePattern(patterns string) (map[string][]string, error) {
	m := make(map[string][]string)
	m["token1"] = []string{
		"location/20200901/tokens/token1.tar.gz",
		"location/20200902/tokens/token2.tar.gz",
		"location/20200903/tokens/token3.tar.gz",
		"location/20200903/tokens/token4.tar.gz",
	}
	m["token2"] = []string{
		"location/20200901/tokens/token1.tar.gz",
		"location/20200902/tokens/token2.tar.gz",
		"location/20200903/tokens/token3.tar.gz",
		"location/20200903/tokens/token4.tar.gz",
	}
	m["token3"] = []string{
		"location/20200901/tokens/token1.tar.gz",
		"location/20200902/tokens/token2.tar.gz",
		"location/20200903/tokens/token3.tar.gz",
		"location/20200903/tokens/token4.tar.gz",
	}
	return m, nil
}
