package utils

func SplitByNumberOfCharacters(s string, max int) []string {
	var split []string = []string{}
	for len(s) > max {
		split = append(split, s[:max])
		s = s[max:]
	}
	split = append(split, s)
	return split
}
