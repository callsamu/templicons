package templicons

func iconURL(api string, set, name string, p *Parameters) string {
	url := api + "/" + set + "/" + name
	if p != nil {
		url += "?" + p.asQueryString()
	}
	return url
}
