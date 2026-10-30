package linkedin

import "net/url"

type Part string

const (
	Posts   Part = "posts"
	Profile Part = "profile"
	Contact Part = "contact"
)

func (p Part) Path() string {
	switch p {
	case Posts:
		return "profile/posts"
	case Profile:
		return "profile/detail"
	case Contact:
		return "profile/contact"
	default:
		return ""
	}
}

func (p Part) Query(username string) url.Values {
	query := url.Values{}
	query.Set("username", username)
	if p == Posts {
		query.Set("page_number", "1")
	}
	return query
}
