package components

import (
	"time"

	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

func PageFooter() g.Node {
	return Footer(Class("prose prose-sm prose-indigo"),
		P(g.Textf("Rendered %v. ", time.Now().Format(time.RFC3339))),
		P(A(Href("https://www.gomponents.com"), g.Text("gomponents"))),
	)
}
