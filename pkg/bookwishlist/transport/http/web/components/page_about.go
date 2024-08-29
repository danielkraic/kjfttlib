package components

import (
	g "github.com/maragudk/gomponents"
	html "github.com/maragudk/gomponents/html"
	b "github.com/willoma/bulma-gomponents"
	e "github.com/willoma/gomplements"
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
