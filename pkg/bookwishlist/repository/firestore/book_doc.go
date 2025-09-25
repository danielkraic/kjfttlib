package firestore

import (
	"time"

	"github.com/danielkraic/kjfttlib/pkg/book"
)

type bookDoc struct {
	ID          string             `firestore:"id"`
	Title       string             `firestore:"title"`
	Author      string             `firestore:"author"`
	Instances   []*bookInstanceDoc `firestore:"instances"`
	URL         string             `firestore:"url"`
	UpdatedTime time.Time          `firestore:"updated_time,omitempty"`
}

type bookInstanceDoc struct {
	Location string `firestore:"location"`
	Status   string `firestore:"status"`
}

func newDoc(b *book.Model) *bookDoc {
	instancesDocs := make([]*bookInstanceDoc, 0, len(b.Instances))

	for _, instance := range b.Instances {
		instancesDocs = append(instancesDocs, newBookInstanceDoc(instance))
	}

	return &bookDoc{
		ID:        b.ID,
		Title:     b.Title,
		Author:    b.Author,
		URL:       b.URL,
		Instances: instancesDocs,
	}
}

func (d *bookDoc) toBook() *book.Model {
	instances := make([]*book.Instance, 0, len(d.Instances))

	for _, instance := range d.Instances {
		instances = append(instances, instance.toBookInstance())
	}

	return &book.Model{
		ID:        d.ID,
		Title:     d.Title,
		Author:    d.Author,
		URL:       d.URL,
		Instances: instances,
	}
}

func newBookInstanceDoc(i *book.Instance) *bookInstanceDoc {
	return &bookInstanceDoc{
		Location: i.Location,
		Status:   i.Status,
	}
}

func (d *bookInstanceDoc) toBookInstance() *book.Instance {
	return &book.Instance{
		Location: d.Location,
		Status:   d.Status,
	}
}
