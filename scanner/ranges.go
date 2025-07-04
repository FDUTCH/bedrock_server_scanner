package scanner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"net/http"
	"net/netip"
)

func GetRangesByToken(token, host string) []netip.Prefix {
	var ranges struct {
		R []string `json:"ranges"`
	}

	url := fmt.Sprintf("https://ipinfo.io/ranges/%s?token=%s", host, token)

	req, err := http.NewRequest("GET", url, bytes.NewBuffer(nil))
	if err != nil {
		panic(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	err = json.NewDecoder(resp.Body).Decode(&ranges)
	if err != nil {
		panic(err)
	}

	var prefixes = make([]netip.Prefix, 0, len(ranges.R))

	for _, prefix := range ranges.R {
		p, err := netip.ParsePrefix(prefix)
		if err != nil {
			continue
		}
		prefixes = append(prefixes, p)
	}
	return prefixes
}

func GetRangesScraping(as string) []netip.Prefix {
	var prefixes []netip.Prefix
	c := colly.NewCollector()
	c.OnHTML("a[class]", func(e *colly.HTMLElement) {
		if e.Attr("class") == "charcoal-link " {
			prefix, err := netip.ParsePrefix(e.Text)
			if err == nil {
				addr := prefix.Addr()
				if addr.Is4() {
					prefixes = append(prefixes, prefix)
				}
			}
		}
	})
	c.Visit(fmt.Sprintf("https://ipinfo.io/%s", as))
	return prefixes
}
