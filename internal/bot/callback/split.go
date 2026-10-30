package callback

import "strings"

func Split(data string) (head, rest string) {
	head, rest, _ = strings.Cut(data, "-")
	return head, rest
}
