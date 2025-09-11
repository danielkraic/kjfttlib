package components

import (
	"fmt"
	"regexp"

	"github.com/danielkraic/kjfttlib/pkg/book"
	g "github.com/maragudk/gomponents"
	html "github.com/maragudk/gomponents/html"
	b "github.com/willoma/bulma-gomponents"
	"github.com/willoma/bulma-gomponents/fa"
	e "github.com/willoma/gomplements"
)

const _instanceStatusAvailable = "Voľný"

var dateRegex = regexp.MustCompile(` do \d{2}\.\d{2}\.\d{4}`)

func PageBooks(books []*book.Model) (string, g.Node) {
	return "KJFTT books wishlist",
		e.Div(
			b.Title(
				"KJFTT books wishlist",
			),
			// Search input field
			b.Field(
				b.Control(
					b.InputText(
						e.Placeholder("Search books by title, author, or ID..."),
						e.ID("book-search"),
						g.Attr("oninput", "filterBooks()"),
					),
				),
			),
			b.Table(
				b.Striped,
				e.Class("sortable"),
				e.ID("books-table"),
				b.Hoverable,
				b.FullWidth,
				b.HeadRow(
					e.Span("ID", e.Styles{"cursor": "pointer"}),
					e.Span("Book title", e.Styles{"cursor": "pointer"}),
					e.Span("Author", e.Styles{"cursor": "pointer"}),
					e.Span("Instances", e.Styles{"cursor": "pointer"}),
					b.Dropdown(
						b.Clickable,
						b.Right,
						e.Class("no-sort"),
						b.PulledRight,
						e.ID("dropdown-menu-books-all"),
						b.OnTrigger(
							b.Button(
								e.AriaHasPopupTrue,
								e.AriaControlsID("dropdown-menu-books-all"),
								fa.Icon(fa.Solid, "ellipsis-v", b.Small),
							),
						),
						b.DropdownAHref("/books/refresh", "Refresh instances of all books in wishlist"),
					),
				),
				g.Group(
					g.Map(books, func(book *book.Model) g.Node {
						return b.Row(
							b.Td(book.ID),
							b.Td(
								e.A(e.Href(book.URL), book.Title),
							),
							b.Td(book.Author),
							getBookInstances(book),

							b.Dropdown(
								b.Clickable,
								b.Right,
								b.PulledRight,

								e.ID("dropdown-menu-book-"+book.ID),
								b.OnTrigger(
									b.Button(
										e.AriaHasPopupTrue,
										e.AriaControlsID("dropdown-menu-book-"+book.ID),
										fa.Icon(fa.Solid, "ellipsis-v", b.Small),
									),
								),
								b.DropdownAHref("/books/refresh/"+book.ID, "Refresh book instances"),
								b.DropdownAHref("/books/delete/"+book.ID, "Delete book from wishlist"),
							),
						)
					}),
				),
			),
			// JavaScript for search functionality
			html.Script(g.Raw(`
				function filterBooks() {
					const searchInput = document.getElementById('book-search');
					const table = document.getElementById('books-table');
					const rows = table.getElementsByTagName('tbody')[0].getElementsByTagName('tr');
					const searchTerm = searchInput.value.toLowerCase();

					for (let i = 0; i < rows.length; i++) {
						const cells = rows[i].getElementsByTagName('td');
						let rowText = '';

						// Concatenate text from ID, Title, and Author columns (first 3 columns)
						for (let j = 0; j < Math.min(3, cells.length); j++) {
							rowText += cells[j].textContent.toLowerCase() + ' ';
						}

						if (rowText.includes(searchTerm)) {
							rows[i].style.display = '';
						} else {
							rows[i].style.display = 'none';
						}
					}
				}
			`)),
		)
}

func getBookInstances(libBook *book.Model) g.Node {
	instanceCountByStatus := getBookInstanceCountByStatus(libBook)

	if len(instanceCountByStatus) == 0 {
		return b.Tag(b.Grey, "No instances")
	}

	instances := []g.Node{}

	availableCount, ok := instanceCountByStatus[_instanceStatusAvailable]
	if ok && availableCount > 0 {
		instances = append(instances, b.Tag(b.Success, fmt.Sprintf("%s: %d", _instanceStatusAvailable, availableCount)))
	}

	for status, count := range instanceCountByStatus {
		if status != _instanceStatusAvailable && count > 0 {
			instances = append(instances, b.Tag(b.Grey, fmt.Sprintf("%s: %d", status, count)))
		}
	}

	return g.Group(instances)
}

func getBookInstanceCountByStatus(book *book.Model) map[string]int {
	instancesByStatus := make(map[string]int)
	for _, instance := range book.Instances {
		status := trimDateFromStatus(instance.Status)
		instancesByStatus[status]++
	}
	return instancesByStatus
}

func trimDateFromStatus(status string) string {
	return dateRegex.ReplaceAllString(status, "")
}
