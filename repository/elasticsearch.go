package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"time"

	"github.com/disaster37/gobot-fat/models"
	elastic "github.com/elastic/go-elasticsearch/v7"
	olivere "github.com/olivere/elastic/v7"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// ElasticsearchRepositoryGen represent generic repository to request Elasticsearch
type ElasticsearchRepositoryGen struct {
	Conn  *elastic.Client
	Index string
}

// NewElasticsearchRepository create new Elasticsearch repository
func NewElasticsearchRepository(conn *elastic.Client, index string) Repository {
	return &ElasticsearchRepositoryGen{
		Conn:  conn,
		Index: index,
	}
}

// Get return one document from Elasticsearch with ID
func (h *ElasticsearchRepositoryGen) Get(ctx context.Context, id uint, data interface{}) error {

	if data == nil {
		return errors.New("Data can't be null")
	}
	if reflect.TypeOf(data).Kind() != reflect.Ptr {
		return errors.New("Data must be a pointer")
	}

	dataModel := data.(models.Model)

	res, err := h.Conn.Get(
		h.Index,
		fmt.Sprintf("%d", id),
		h.Conn.Get.WithContext(ctx),
		h.Conn.Get.WithPretty(),
	)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	// Check if query found
	if res.IsError() {
		return errors.Errorf("Error when read response: %s", res.String())
	}

	ret := new(olivere.GetResult)
	if err := h.decode(res.Body, ret); err != nil {
		return err
	}

	// When no result found
	if !ret.Found {
		return ErrRecordNotFoundError
	}

	if err = json.Unmarshal(ret.Source, data); err != nil {
		return err
	}
	idDoc, err := strconv.ParseUint(ret.Id, 10, 32)
	if err != nil {
		return err
	}
	dataModel.GetModel().ID = uint(idDoc)

	log.Debugf("Data: %+v", data)

	return nil
}

// List return all document on index
func (h *ElasticsearchRepositoryGen) List(ctx context.Context, listData interface{}) error {

	if listData == nil {
		return errors.New("Data can't be null")
	}
	if reflect.TypeOf(listData).Kind() != reflect.Ptr {
		return errors.New("ListData must be a pointer")
	}
	if reflect.TypeOf(listData).Elem().Kind() != reflect.Slice {
		return errors.New("ListData must contain slice")
	}

	ld := reflect.ValueOf(listData).Elem()

	res, err := h.Conn.Search(
		h.Conn.Search.WithIndex(h.Index),
		h.Conn.Search.WithQuery(`{"query": {"match_all" : {}}}`),
		h.Conn.Search.WithContext(ctx),
		h.Conn.Search.WithPretty(),
	)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	// Check if query found
	if res.IsError() {
		return errors.Errorf("Error when read response: %s", res.String())
	}

	ret := new(olivere.SearchResult)
	if err := h.decode(res.Body, ret); err != nil {
		return err
	}
	for _, hit := range ret.Hits.Hits {
		tmp := reflect.New(reflect.TypeOf(listData).Elem().Elem())
		if err = json.Unmarshal(hit.Source, tmp.Interface()); err != nil {
			return err
		}
		idDoc, err := strconv.ParseUint(hit.Id, 10, 32)
		if err != nil {
			return err
		}
		var data models.Model
		if tmp.Elem().Kind() == reflect.Ptr {
			data = tmp.Elem().Interface().(models.Model)
		} else {
			data = tmp.Interface().(models.Model)
		}

		data.GetModel().ID = uint(idDoc)

		ld.Set(reflect.Append(ld, tmp.Elem()))

	}

	log.Debugf("Data: %+v", listData)

	return nil
}

// Update document on Elasticsearch
func (h *ElasticsearchRepositoryGen) Update(ctx context.Context, data interface{}) error {

	if data == nil {
		return errors.New("Data can't be null")
	}
	if reflect.TypeOf(data).Kind() != reflect.Ptr {
		return errors.New("Data must a pointer")
	}
	log.Debugf("Data: %s", data)

	dataModel := data.(models.Model)

	dataModel.SetUpdatedAt(time.Now())

	sdata, err := json.Marshal(data)
	if err != nil {
		return err
	}
	b := bytes.NewBuffer(sdata)

	res, err := h.Conn.Index(
		h.Index,
		b,
		h.Conn.Index.WithDocumentID(fmt.Sprintf("%d", dataModel.GetModel().ID)),
		h.Conn.Index.WithContext(ctx),
		h.Conn.Index.WithPretty(),
	)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	// Check if query found
	if res.IsError() {
		return errors.Errorf("Error when read response: %s", res.String())
	}

	log.Debugf("Response: %s", res.String())

	return nil
}

// Create add new document on Elasticsearch
func (h *ElasticsearchRepositoryGen) Create(ctx context.Context, data interface{}) error {
	return h.Update(ctx, data)
}

func (h *ElasticsearchRepositoryGen) decode(body io.Reader, ret interface{}) error {

	decoder := json.NewDecoder(body)

	return decoder.Decode(ret)
}
