package components

import (
	"github.com/danielkraic/kjfttlib/pkg/book"
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

func PageBooks(books []*book.Model) (string, g.Node) {
	return "KJFTT books wishlist", Div(
		H1(g.Text("Welcome to KJFTT books wishlist")),
		P(g.Text("I hope it will make you happy. 😄 It's using TailwindCSS for styling.")),
		Books(books),
	)
}
