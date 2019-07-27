package main

import (
	"html/template"
	"io"
)

const (
	templateData = `
package {{ .PackageName }}

import (
	"fmt"

	"github.com/samonzeweb/godb"
	"github.com/samonzeweb/godb/adapters/{{ .Driver }}"
)

type {{ .Type }}Repository struct {
	db *godb.DB
}

func (entity *{{ .Type }}) TableName() string {
	return "{{ .TableName }}"
}

func New{{ .Type }}Repository() *{{ .Type }}Repository {
	db, err := godb.Open({{ .Driver }}.Adapter, "{{ .URL }}")
	if err != nil {
		panic(err)
	}
	return &{{ .Type }}Repository{
		db: db,
	}
}

func (this *{{ .Type }}Repository) FindAll() (result []{{ .Type }}, err error) {
	result = make([]{{ .Type }}, 0, 0)
	err = this.db.Select(&result).Do()
	return
}

func (this *{{ .Type }}Repository) FindAllWithIterator() (result godb.Iterator, err error) {
	entity := &{{ .Type }}{}
	result ,err = this.db.SelectFrom(entity.TableName()).DoWithIterator()
	return
}

func (this *{{ .Type }}Repository) FindOneById(idField string, idValue interface{}) (result *{{ .Type }}, err error) {
	result = &{{ .Type }}{}
	err = this.db.Select(result).
		Where(fmt.Sprintf("%s = ?", idField), idValue).
		Do()
	return
}

func (this *{{ .Type }}Repository) RawSQL(queryBuffer *godb.SQLBuffer) (result []{{ .Type }}, err error) {
	result = make([]{{ .Type }}, 0, 0)
	err = this.db.RawSQL(queryBuffer.SQL(), queryBuffer.Arguments()...).Do(&result)
	return
}

func (this *{{ .Type }}Repository) BeginTx() (err error) {
	err = this.db.Begin()
	return
}

func (this *{{ .Type }}Repository) CommitTx() (err error) {
	err = this.db.Commit()
	return
}

func (this *{{ .Type }}Repository) RollbackTx() (err error) {
	err = this.db.Rollback()
	return
}

func (this *{{ .Type }}Repository) CreateOne(new{{ .Type }} *{{ .Type }}) (err error) {
	err = this.db.Insert(new{{ .Type }}).Do()
	return
}

func (this *{{ .Type }}Repository) UpdateOne(new{{ .Type }} *{{ .Type }}) (err error) {
	err = this.db.Update(new{{ .Type }}).Do()
	return
}

func (this *{{ .Type }}Repository) DeleteOne(new{{ .Type }} *{{ .Type }}) (err error) {
	_ ,err = this.db.Delete(new{{ .Type }}).Do()
	return
}
`
)

type Metadata struct {
	PackageName string
	Type        string
	Driver      string
	URL         string
	TableName   string
}

type Generator struct {
}

func (g *Generator) Generate(writer io.Writer, metadata Metadata) error {
	tmpl, err := g.template()
	if err != nil {
		return err
	}
	return tmpl.Execute(writer, metadata)
}

func (g *Generator) template() (*template.Template, error) {
	tmpl := template.New("template")
	return tmpl.Parse(string(templateData))
}
