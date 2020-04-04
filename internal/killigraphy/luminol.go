package killigraphy

import "strings"

type luminol struct {
	words []string
}

func newLuminol(words []string) *luminol {
	return &luminol{words}
}

func (l *luminol) Reveal() string {
	contents := strings.Join(l.words, " ")
	contents = l.scrub(contents)
	return l.swapYear(contents)
}

func (l *luminol) scrub(contents string) string {
	contents = strings.ReplaceAll(contents, "  ", " ")
	contents = strings.ReplaceAll(contents, " %HESITATION", "")
	return contents
}

func (l *luminol) swapYear(contents string) string {
	contents = strings.ReplaceAll(contents, "nineteen seventy", "1970")
	contents = strings.ReplaceAll(contents, "nineteen seventy", "1970")
	contents = strings.ReplaceAll(contents, "nineteen seventy one", "1971")
	contents = strings.ReplaceAll(contents, "nineteen seventy two", "1972")
	contents = strings.ReplaceAll(contents, "nineteen seventy three", "1973")
	contents = strings.ReplaceAll(contents, "nineteen seventy four", "1974")
	contents = strings.ReplaceAll(contents, "nineteen seventy five", "1975")
	contents = strings.ReplaceAll(contents, "nineteen seventy six", "1976")
	contents = strings.ReplaceAll(contents, "nineteen seventy seven", "1977")
	contents = strings.ReplaceAll(contents, "nineteen seventy eight", "1978")
	contents = strings.ReplaceAll(contents, "nineteen seventy nine", "1979")
	contents = strings.ReplaceAll(contents, "nineteen eighty", "1980")
	contents = strings.ReplaceAll(contents, "nineteen eighty one", "1981")
	contents = strings.ReplaceAll(contents, "nineteen eighty two", "1982")
	contents = strings.ReplaceAll(contents, "nineteen eighty three", "1983")
	contents = strings.ReplaceAll(contents, "nineteen eighty four", "1984")
	contents = strings.ReplaceAll(contents, "nineteen eighty five", "1985")
	contents = strings.ReplaceAll(contents, "nineteen eighty six", "1986")
	contents = strings.ReplaceAll(contents, "nineteen eighty seven", "1987")
	contents = strings.ReplaceAll(contents, "nineteen eighty eight", "1988")
	contents = strings.ReplaceAll(contents, "nineteen eighty nine", "1989")
	contents = strings.ReplaceAll(contents, "nineteen ninety", "1990")
	contents = strings.ReplaceAll(contents, "nineteen ninety one", "1991")
	contents = strings.ReplaceAll(contents, "nineteen ninety two", "1992")
	contents = strings.ReplaceAll(contents, "nineteen ninety three", "1993")
	contents = strings.ReplaceAll(contents, "nineteen ninety four", "1994")
	contents = strings.ReplaceAll(contents, "nineteen ninety five", "1995")
	contents = strings.ReplaceAll(contents, "nineteen ninety six", "1996")
	contents = strings.ReplaceAll(contents, "nineteen ninety seven", "1997")
	contents = strings.ReplaceAll(contents, "nineteen ninety eight", "1998")
	contents = strings.ReplaceAll(contents, "nineteen ninety nine", "1999")
	contents = strings.ReplaceAll(contents, "two thousand one", "2001")
	contents = strings.ReplaceAll(contents, "two thousand two", "2002")
	contents = strings.ReplaceAll(contents, "two thousand three", "2003")
	contents = strings.ReplaceAll(contents, "two thousand four", "2004")
	contents = strings.ReplaceAll(contents, "two thousand five", "2005")
	contents = strings.ReplaceAll(contents, "two thousand six", "2006")
	contents = strings.ReplaceAll(contents, "two thousand seven", "2007")
	contents = strings.ReplaceAll(contents, "two thousand eight", "2008")
	contents = strings.ReplaceAll(contents, "two thousand nine", "2009")
	return contents
}
