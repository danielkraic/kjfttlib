package components

import (
	"github.com/danielkraic/kjfttlib/pkg/book"

	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

func Books(books []*book.Model) g.Node {
	return Table(
		Class("border-collapse border border-slate-400"),
		g.Group(
			[]g.Node{
				Th(
					Class("border border-slate-300"),
					g.Text("Title"),
				),
				Th(
					Class("border border-slate-300"),
					g.Text("Author"),
				),
			},
		),
		g.Group(
			g.Map(books, func(b *book.Model) g.Node {
				return Tr(
					Td(
						Class("border border-slate-300"),
						A(
							Href(b.URL),
							g.Text(b.Title),
						),
					),
					Td(
						Class("border border-slate-300"),
						g.Text(b.Author),
					),
				)
			}),
		),
	)
}
