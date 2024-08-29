package components

import (
	"github.com/danielkraic/kjfttlib/pkg/book"
	g "github.com/maragudk/gomponents"
	html "github.com/maragudk/gomponents/html"
	b "github.com/willoma/bulma-gomponents"
	e "github.com/willoma/gomplements"
)

type PageAddBookNotification struct {
	BookID  string
	Book    book.Model
	UserErr string
}

func (n *PageAddBookNotification) ToNode() g.Node {
	if n.UserErr != "" {
		if n.BookID == "" {
			return b.Notification(
				b.Danger,
				"Unable to add book to wishlist. ", n.UserErr,
			)
		}

		return b.Notification(
			b.Danger,
			"Unable to add book with ID ", e.Strong(n.BookID), " to wishlist. ", n.UserErr,
		)
	}

	return b.Notification(
		b.Success,
		"Book with ID ", e.Strong(n.BookID), " was added to wishlist.",
	)
}

func (n *PageAddBookNotification) Color() b.Color {
	if n.UserErr != "" {
		return b.Danger
	}
	return b.Success
}

func (n *PageAddBookNotification) Msg() string {
	if n.UserErr != "" {
		return n.UserErr
	}
	return "Book added to wishlist."
}

func PageAddBook(notifications ...PageAddBookNotification) (string, g.Node) {
	return "Add book to wishlist", e.Div(
		b.Title("Add book to wishlist"),
		g.Group(g.Map(notifications, func(n PageAddBookNotification) g.Node {
			return n.ToNode()
		})),
		html.Form(
			html.Method("POST"),
			b.Field(
				b.Label("Book"),
				b.Control(
					b.InputText(e.Placeholder("Book ID or URL"), html.Name("bookid")),
				),
				b.Help("Use ID of book from library or URL of book."),
			),
			b.Control(
				b.Button(b.Primary, "Add book", html.Type("submit")),
			),
		),
	)
}
