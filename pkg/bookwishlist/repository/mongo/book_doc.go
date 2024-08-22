package mongo

import (
	"time"

	"github.com/danielkraic/kjfttlib/pkg/book"
)

type bookDoc struct {
	ID          string             `bson:"_id"`
	Title       string             `bson:"title"`
	Author      string             `bson:"author"`
	Instances   []*bookInstanceDoc `bson:"instances"`
	URL         string             `bson:"url"`
	UpdatedTime time.Time          `bson:"updated_time,omitempty"`
}

type bookInstanceDoc struct {
	Location string `bson:"location"`
	Status   string `bson:"status"`
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
