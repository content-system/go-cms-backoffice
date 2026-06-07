package slug

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

func Slugify(title string, uuid string, wordLimit int, maxLength int) string {
	if wordLimit <= 0 {
		wordLimit = 10
	}

	if maxLength <= 0 {
		maxLength = 60
	}

	// lowercase
	s := strings.ToLower(title)

	// normalize NFD
	s = norm.NFD.String(s)

	// remove accents
	var b strings.Builder
	for _, r := range s {
		// remove unicode combining marks
		if unicode.Is(unicode.Mn, r) {
			continue
		}

		// Vietnamese special case
		if r == 'đ' {
			r = 'd'
		}

		b.WriteRune(r)
	}

	s = b.String()

	// treat "-" as space
	s = strings.ReplaceAll(s, "-", " ")

	// keep ASCII letters, numbers, spaces only
	re := regexp.MustCompile(`[^a-z0-9\s]`)
	s = re.ReplaceAllString(s, "")

	// trim + split spaces
	words := strings.Fields(strings.TrimSpace(s))

	// limit words
	if len(words) > wordLimit {
		words = words[:wordLimit]
	}

	// join with "-"
	slug := strings.Join(words, "-")

	// limit max length
	if len(slug) > maxLength {
		slug = slug[:maxLength]
		slug = strings.TrimRight(slug, "-")
	}

	// fallback
	if len(slug) == 0 {
		return uuid
	}

	return slug + "-" + uuid
}
