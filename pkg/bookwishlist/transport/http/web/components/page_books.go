package components

import (
	"fmt"
	"regexp"

	"github.com/danielkraic/kjfttlib/pkg/book"
	b "github.com/willoma/bulma-gomponents"
	"github.com/willoma/bulma-gomponents/fa"
	e "github.com/willoma/gomplements"
	g "maragu.dev/gomponents"
	html "maragu.dev/gomponents/html"
)

const _instanceStatusAvailable = "Voľný"

var dateRegex = regexp.MustCompile(` do \d{2}\.\d{2}\.\d{4}`)

func PageBooks(books []*book.Model) (string, g.Node) {
	return "KJFTT books wishlist",
		e.Div(
			// Add CSS for custom tooltips
			html.StyleEl(g.Raw(`
				.tooltip {
					position: relative;
					cursor: help;
				}

				.tooltip .tooltip-text {
					visibility: hidden;
					width: 350px;
					max-width: 90vw;
					background-color: #2c3e50;
					color: #ecf0f1;
					text-align: left;
					border-radius: 8px;
					padding: 10px 14px;
					position: absolute;
					z-index: 1000;
					bottom: 125%;
					left: 50%;
					margin-left: -175px;
					opacity: 0;
					transition: opacity 0.3s;
					font-size: 13px;
					line-height: 1.4;
					white-space: pre-wrap;
					box-shadow: 0 4px 12px rgba(0,0,0,0.3);
					border: 1px solid #34495e;
					font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
				}

				.tooltip .tooltip-text::after {
					content: "";
					position: absolute;
					top: 100%;
					left: 50%;
					margin-left: -8px;
					border-width: 8px;
					border-style: solid;
					border-color: #2c3e50 transparent transparent transparent;
				}

				.tooltip:hover .tooltip-text {
					visibility: visible;
					opacity: 1;
				}

				@media (max-width: 768px) {
					.tooltip .tooltip-text {
						width: 280px;
						margin-left: -140px;
						font-size: 12px;
					}
				}
			`)),
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
		tooltipText := getInstanceTooltipText(libBook, _instanceStatusAvailable)
		instances = append(instances,
			e.Span(
				e.Class("tooltip"),
				b.Tag(b.Success, fmt.Sprintf("%s: %d", _instanceStatusAvailable, availableCount)),
				e.Span(
					e.Class("tooltip-text"),
					g.Text(tooltipText),
				),
			),
		)
	}

	for status, count := range instanceCountByStatus {
		if status != _instanceStatusAvailable && count > 0 {
			tooltipText := getInstanceTooltipText(libBook, status)
			instances = append(instances,
				e.Span(
					e.Class("tooltip"),
					b.Tag(b.Grey, fmt.Sprintf("%s: %d", status, count)),
					e.Span(
						e.Class("tooltip-text"),
						g.Text(tooltipText),
					),
				),
			)
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

func getInstanceTooltipText(book *book.Model, targetStatus string) string {
	entryCounts := make(map[string]int)

	for _, instance := range book.Instances {
		trimmedStatus := trimDateFromStatus(instance.Status)
		if trimmedStatus == targetStatus {
			// Create entry in format "location:status"
			entry := fmt.Sprintf("%s: %s", instance.Location, instance.Status)
			entryCounts[entry]++
		}
	}

	if len(entryCounts) == 0 {
		return "No instances found"
	}

	var details []string
	totalCount := 0
	for entry, count := range entryCounts {
		totalCount += count
		if count > 1 {
			details = append(details, fmt.Sprintf("%s (%d)", entry, count))
		} else {
			details = append(details, entry)
		}
	}

	tooltipText := fmt.Sprintf("📚 %s (%d):\n", targetStatus, totalCount)
	for i, detail := range details {
		tooltipText += fmt.Sprintf("• %s", detail)
		if i < len(details)-1 {
			tooltipText += "\n"
		}
	}

	return tooltipText
}
