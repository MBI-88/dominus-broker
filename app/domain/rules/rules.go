package rules

import (
	"regexp"
)

type rules struct {
	re *regexp.Regexp
}


func (r *rules) CheckURI(uri string) bool {
	return r.re.MatchString(uri)
}


type RulesInt interface {
	CheckURI(uri string) bool
}

func NewRules() RulesInt {
	return &rules{
		re: regexp.MustCompile(`^(?:(?:(?:25[0-5]|2[0-4]\d|[01]?\d\d?)(?:\.(?:25[0-5]|2[0-4]\d|[01]?\d\d?)){3}:\d{1,5})|(?:(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}))(?:\/[^\s]*)?$`),
	}
}
