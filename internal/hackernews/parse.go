package hackernews

import (
	"fmt"
	"io"

	"golang.org/x/net/html"
)

// parseHiddenInputs walks an HTML document and returns the value attribute of
// each <input> whose name matches one of the requested names. It errors if any
// requested name is missing, mirroring Hacker News's requirement that fnid and
// fnop be present on the submit form.
func parseHiddenInputs(r io.Reader, names ...string) (map[string]string, error) {
	want := make(map[string]bool, len(names))
	for _, n := range names {
		want[n] = true
	}

	doc, err := html.Parse(r)
	if err != nil {
		return nil, fmt.Errorf("parse html: %w", err)
	}

	found := make(map[string]string, len(names))

	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "input" {
			var name, value string
			for _, attr := range n.Attr {
				switch attr.Key {
				case "name":
					name = attr.Val
				case "value":
					value = attr.Val
				}
			}
			if want[name] {
				found[name] = value
			}
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}
	}
	walk(doc)

	for _, n := range names {
		if _, ok := found[n]; !ok {
			return nil, fmt.Errorf("%w: %q", ErrFormFieldMissing, n)
		}
	}
	return found, nil
}
