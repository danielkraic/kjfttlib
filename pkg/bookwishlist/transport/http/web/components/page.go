package components

import (
	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	. "github.com/maragudk/gomponents/html"
)

func Page(title, path string, body g.Node) g.Node {
	return c.HTML5(c.HTML5Props{
		Title:    title,
		Language: "en",
		Head: []g.Node{
			Link(Rel("stylesheet"), Href("https://unpkg.com/tailwindcss@2.1.2/dist/base.min.css")),
			Link(Rel("stylesheet"), Href("https://unpkg.com/tailwindcss@2.1.2/dist/components.min.css")),
			Link(Rel("stylesheet"), Href("https://unpkg.com/@tailwindcss/typography@0.4.0/dist/typography.min.css")),
			Link(Rel("stylesheet"), Href("https://unpkg.com/tailwindcss@2.1.2/dist/utilities.min.css")),
		},
		Body: []g.Node{
			Navbar(path, []PageLink{
				{Path: "/add-book", Name: "Add book"},
				{Path: "/about", Name: "About"},
			}),
			Container(
				Prose(body),
				PageFooter(),
			),
		},
	})
}

func Container(children ...g.Node) g.Node {
	return Div(Class("max-w-7xl mx-auto px-2 sm:px-6 lg:px-8"), g.Group(children))
}

func Prose(children ...g.Node) g.Node {
	return Div(Class("prose"), g.Group(children))
}
