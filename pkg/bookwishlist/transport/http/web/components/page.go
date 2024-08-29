package components

import (
	g "github.com/maragudk/gomponents"
	html "github.com/maragudk/gomponents/html"
	b "github.com/willoma/bulma-gomponents"
	e "github.com/willoma/gomplements"
)

func Page(title, path string, pageBody g.Node) g.Node {
	return b.HTML(
		b.HTitle(title),
		b.Language("en"),
		b.Head(
			html.Meta(html.Charset("utf-8")),
			html.Link(html.Rel("stylesheet"), html.Href("https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.6.0/css/all.min.css")),
		),
		e.Div(
			e.Class("container hero is-fullheight"),
			e.Div(
				e.Class("hero-head"),
				navbar(),
			),
			e.Div(
				e.Class("hero-body"),
				b.Container(
					pageBody,
				),
			),
			e.Div(
				e.Class("hero-foot"),
				footer(),
			),
		),
		b.Script("https://cdn.jsdelivr.net/gh/tofsjonas/sortable@latest/sortable.min.js"),
	)
}

func navbar() g.Node {
	return b.Navbar(
		b.NavbarBrand(
			b.NavbarAHref(
				"/",
				e.Strong("KJFTTLIB"),
			),
		),
		b.NavbarStart(
			b.NavbarAHref("/", "Home"),
			b.NavbarAHref("/add-book", "Add book"),
			b.NavbarAHref("/about", "About"),
		),
		b.NavbarEnd(
			b.NavbarAHref("https://github.com/danielkraic/kjfttlib", "Github"),
		),
	)
}

func footer() g.Node {
	return b.Footer(
		e.Div(
			e.Strong("KJFTTLIB"), " created by ", e.AHref("https://github.com/danielkraic", "Daniel Kraic"), ".",
		),
	)
}
