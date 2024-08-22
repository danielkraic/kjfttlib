package components

import (
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

func PageAddBook() (string, g.Node) {
	return "KJFTT books wishlist", Div(
		H1(g.Text("Add book to KJFTT books wishlist")),
	)
}
