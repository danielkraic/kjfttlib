package kjftt

import (
	"io"
	"strings"

	"github.com/danielkraic/kjfttlib/pkg/book"

	"github.com/PuerkitoBio/goquery"
	jErrors "github.com/juju/errors"
)

func ParseBookFromHTML(reader io.Reader) (*book.Model, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, jErrors.Annotate(err, "parse html document from response body")
	}

	decodedBook := &book.Model{}

	doc.Find(_htmlSelectorTitle).First().Each(func(i int, s *goquery.Selection) {
		decodedBook.Title = parseTitle(s.Text())
	})

	doc.Find(_htmlSelectorAuthor).First().Each(func(i int, s *goquery.Selection) {
		s.Find(`a`).First().Each(func(i int, s *goquery.Selection) {
			decodedBook.Author = fixExtraSpaces(s.Text())
		})
	})

	doc.Find(`#holdlist > tbody > tr`).Each(func(i int, s *goquery.Selection) {
		instance := &book.Instance{}

		s.Find(`td:nth-child(4) > a`).First().Each(func(i int, s *goquery.Selection) {
			instance.Location = fixExtraSpaces(s.Text())
		})
		s.Find(`td:nth-child(5)`).First().Each(func(i int, s *goquery.Selection) {
			instance.Status = s.Text()
		})

		decodedBook.Instances = append(decodedBook.Instances, instance)
	})

	err = validateParsedBook(decodedBook)
	if err != nil {
		return nil, jErrors.Annotate(err, "validating parsed book")
	}

	return decodedBook, nil
}

func parseTitle(str string) string {
	str = strings.Join(strings.Split(str, "\n"), " ")

	matches := _reBookTitle.FindStringSubmatch(str)
	if len(matches) < 2 {
		return ""
	}

	return fixExtraSpaces(matches[1])
}

func fixExtraSpaces(str string) string {
	var words []string

	for _, word := range strings.Split(str, " ") {
		word = strings.TrimSpace(word)
		if word == "" {
			continue
		}

		words = append(words, word)
	}

	return strings.Join(words, " ")
}

func validateParsedBook(b *book.Model) error {
	if b.Author == "" {
		return jErrors.New("author is empty")
	}

	if b.Title == "" {
		return jErrors.New("title is empty")
	}

	if len(b.Instances) == 0 {
		return jErrors.New("no instances found")
	}

	for i, instance := range b.Instances {
		if instance.Location == "" {
			return jErrors.Errorf("location of instance %d is empty", i)
		}

		if instance.Status == "" {
			return jErrors.Errorf("status of instance %d is empty", i)
		}
	}

	return nil
}
