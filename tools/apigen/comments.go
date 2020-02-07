package apigen

import "strings"

const maxCommentLength = 80

func splitLongComment(line string) []string {
	line = strings.TrimSpace(line)

	if len(line) < maxCommentLength {
		return []string{line}
	}

	maxPos := -1
	for i, c := range line {
		if i > maxCommentLength && maxPos > 0 {
			break
		}
		if c == ' ' {
			maxPos = i
		}
	}

	if maxPos <= 0 {
		return []string{line}
	}

	res := []string{line[:maxPos]}
	next := splitLongComment(line[maxPos:])

	res = append(res, next...)
	return res
}

func processComments(comments []string) []string {
	var lines []string
	for i, note := range comments {
		if i > 0 {
			lines = append(lines, "")
		}

		arr := strings.Split(note, "\n")
		lines = append(lines, arr...)
	}

	var result []string
	for _, ln := range lines {
		splitted := splitLongComment(ln)
		result = append(result, splitted...)
	}

	return result
}
