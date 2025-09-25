package components

import (
	b "github.com/willoma/bulma-gomponents"
	e "github.com/willoma/gomplements"
	g "maragu.dev/gomponents"
	html "maragu.dev/gomponents/html"
)

func PageAbout() (string, g.Node) {
	return "About",
		e.Div(
			b.Title("About"),
			html.P(
				g.Text("KJFTTLIB is a website for managing a book wishlist to "),
				html.A(
					html.Href("https://www.kniznicatrnava.sk/"),
					g.Text("KJFTT")),
				g.Text(" library."),
			),
		)
}
