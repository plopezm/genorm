package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

const (
	testFile            = "book-repository.go"
	structName          = "Book"
	driver              = "sqlite3"
	path                = "local.db"
	tableName           = "books"
	expectedPackage     = "package genorm"
	expectedFileContent = `
import (
	"fmt"

	"github.com/samonzeweb/godb"
	"github.com/samonzeweb/godb/adapters/sqlite"
)

type BookRepository struct {
	db *godb.DB
}

func (entity *Book) TableName() string {
	return "books"
}

func NewBookRepository() *BookRepository {
	db, err := godb.Open(sqlite.Adapter, "local.db")
	if err != nil {
		panic(err)
	}
	return &BookRepository{
		db: db,
	}
}

func (this *BookRepository) FindAll() (result []Book, err error) {
	result = make([]Book, 0, 0)
	err = this.db.Select(&result).Do()
	return
}

func (this *BookRepository) FindAllWithIterator() (result godb.Iterator, err error) {
	entity := &Book{}
	result ,err = this.db.SelectFrom(entity.TableName()).DoWithIterator()
	return
}

func (this *BookRepository) FindOneById(idField string, idValue interface{}) (result *Book, err error) {
	result = &Book{}
	err = this.db.Select(result).
		Where(fmt.Sprintf("%s = ?", idField), idValue).
		Do()
	return
}

func (this *BookRepository) RawSQL(queryBuffer *godb.SQLBuffer) (result []Book, err error) {
	result = make([]Book, 0, 0)
	err = this.db.RawSQL(queryBuffer.SQL(), queryBuffer.Arguments()...).Do(&result)
	return
}

func (this *BookRepository) BeginTx() (err error) {
	err = this.db.Begin()
	return
}

func (this *BookRepository) CommitTx() (err error) {
	err = this.db.Commit()
	return
}

func (this *BookRepository) RollbackTx() (err error) {
	err = this.db.Rollback()
	return
}

func (this *BookRepository) CreateOne(newBook *Book) (err error) {
	err = this.db.Insert(newBook).Do()
	return
}

func (this *BookRepository) UpdateOne(newBook *Book) (err error) {
	err = this.db.Update(newBook).Do()
	return
}

func (this *BookRepository) DeleteOne(newBook *Book) (err error) {
	_ ,err = this.db.Delete(newBook).Do()
	return
}
`
)

func executeGeneration() {
	GenerateFiles(structName, driver, path, tableName, "")
}

func clean() {
	os.Remove(testFile)
}

func TestAll(t *testing.T) {
	executeGeneration()

	data, err := ioutil.ReadFile(testFile)
	if err != nil {
		t.Errorf("File was not created")
		return
	}

	dataStr := string(data)
	if !strings.Contains(dataStr, expectedPackage) {
		t.Errorf("wrong package")
		return
	}
	if !strings.Contains(dataStr, expectedFileContent) {
		t.Errorf("wrong package")
		return
	}
	clean()
}
