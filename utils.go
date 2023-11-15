package templicons

import "strings"

func iconURL(api string, set, name string, p *Parameters) string {
	url := api + "/" + set + "/" + name
	if p != nil {
		url += "?" + p.asQueryString()
	}
	return url
}

func parseName(name string) (string, string) {
	parsed := strings.Split(name, ":")
	if len(parsed) != 2 {
		panic("invalid icon name")
	}
	return parsed[0], parsed[1]
}
