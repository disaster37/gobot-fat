package helper

import (
	"encoding/json"
	"io/ioutil"
	"reflect"

	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type ElasticsearchResponseSearch struct {
	Took     int64                       `json:"took"`
	TimedOut bool                        `json:"timed_out"`
	Shard    *ElasticsearchResponseShard `json:"_shards"`
	Hits     *ElasticsearchHits          `json:"hits"`
}

type ElasticsearchResponseShard struct {
	Total      int64 `json:"total"`
	Successful int64 `json:"successful"`
	Skipped    int64 `json:"skipped"`
	Failed     int64 `json:"failed"`
}

type ElasticsearchHits struct {
	Total    *ElasticsearchHitsTotal `json:"total"`
	MaxScore float64                 `json:"max_score"`
	Hits     []ElasticsearchHit      `json:"hits"`
}

type ElasticsearchHitsTotal struct {
	Value    int64  `json:"value"`
	Relation string `json:"eq"`
}

type ElasticsearchHit struct {
	Index        string                 `json:"_index"`
	DocumentType string                 `json:"_type"`
	ID           string                 `json:"_id"`
	Score        float64                `json:"_score"`
	Source       map[string]interface{} `json:"_source"`
}

type ElasticsearchResponseGet struct {
	Index        string                 `json:"_index"`
	DocumentType string                 `json:"_type"`
	ID           string                 `json:"_id"`
	Found        bool                   `json:"found"`
	Version      int64                  `json:"_version"`
	Source       map[string]interface{} `json:"_source"`
}

type ElasticsearchResponseIndex struct {
	Index        string                      `json:"_index"`
	DocumentType string                      `json:"_type"`
	ID           string                      `json:"_id"`
	Version      int64                       `json:"_version"`
	Shard        *ElasticsearchResponseShard `json:"_shards"`
}

func ProcessElasticsearchGet(res *esapi.Response, data interface{}) error {

	if res == nil {
		return errors.New("Res can't be null")
	}

	defer res.Body.Close()

	// Check if query found
	if res.IsError() {
		return errors.Errorf("Error when read response: %s", res.String())
	}

	// Need to convert on object
	esResponse := &ElasticsearchResponseGet{}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	log.Debugf("Response:\n%s", string(body))

	err = json.Unmarshal(body, esResponse)
	if err != nil {
		return err
	}

	log.Debugf("esResponse: %+v", esResponse)

	// When no result found
	if !esResponse.Found {
		return nil
	}

	// Compute data
	t := reflect.TypeOf(data)
	switch t.Kind() {
	case reflect.Ptr:
		log.Debugf("Detect pointer")
		switch t.Elem().Kind() {
		case reflect.Struct:
			log.Debugf("Detect pointer")

			esResponse.Source["id"] = esResponse.ID
			err = Decode(esResponse.Source, data)
			if err != nil {
				return err
			}
		default:
			return errors.New("You must provide pointer of struct")
		}

	default:
		return errors.New("You must provide pointer of struct")

	}

	log.Debugf("Data: %+v", data)

	return nil

}

func ProcessElasticsearchIndex(res *esapi.Response, data interface{}) error {

	if res == nil {
		return errors.New("Res can't be null")
	}

	defer res.Body.Close()

	// Check if query found
	if res.IsError() {
		return errors.Errorf("Error when read response: %s", res.String())
	}

	// Need to convert on object
	esResponse := &ElasticsearchResponseIndex{}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	log.Debugf("Response:\n%s", string(body))

	err = json.Unmarshal(body, esResponse)
	if err != nil {
		return err
	}
	log.Debugf("esResponse: %+v", esResponse)

	// Add id on object
	t := reflect.TypeOf(data)
	switch t.Kind() {
	case reflect.Ptr:
		log.Debugf("Detect pointer")
		switch t.Elem().Kind() {
		case reflect.Struct:
			log.Debugf("Detect pointer")
			value := reflect.ValueOf(data).Elem()
			idField := value.FieldByName("ID")
			if idField.CanSet() {
				idField.Set(reflect.ValueOf(esResponse.ID))
			}
		default:
			return errors.New("You must provide pointer of struct")
		}

	default:
		return errors.New("You must provide pointer of struct")

	}

	log.Debugf("Data: %+v", data)

	return nil

}

func ProcessElasticsearchSearch(res *esapi.Response, data interface{}, minimalScoring float64) error {

	if res == nil {
		return errors.New("Res can't be null")
	}

	defer res.Body.Close()

	if res.IsError() {
		return errors.Errorf("Error when read response: %s", res.String())
	}

	// Need to convert on object
	esResponse := &ElasticsearchResponseSearch{}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	log.Debugf("Response:\n%s", string(body))

	err = json.Unmarshal(body, esResponse)
	if err != nil {
		return err
	}

	log.Debugf("esResponse: %+v", esResponse.Shard)

	// Check scoring and number of result
	if esResponse.Hits.Total.Value <= 0 || esResponse.Hits.MaxScore < minimalScoring {
		return nil
	}

	// Search if data is array or object
	t := reflect.TypeOf(data)
	switch t.Kind() {
	case reflect.Ptr:
		log.Debugf("Detect pointer")
		switch t.Elem().Kind() {
		case reflect.Slice:
			log.Debugf("Detect slice")
			value := reflect.ValueOf(data).Elem()
			for _, hit := range esResponse.Hits.Hits {

				if hit.Score >= minimalScoring {
					hit.Source["id"] = hit.ID
					obj := reflect.New(t.Elem().Elem())

					err = Decode(hit.Source, obj.Interface())
					if err != nil {
						return err
					}

					value.Set(reflect.Append(value, obj.Elem()))
				} else {
					break
					log.Debugf("Score is %f <= %f", hit.Score, minimalScoring)
				}

			}
		default:
			return errors.New("You must provide pointer of slice or pointer of struct")
		}

	default:
		return errors.New("You must provide pointer of slice or pointer of struct")

	}

	log.Debugf("Data: %+v", data)

	return nil

}

func ProcessElasticsearchFetch(res *esapi.Response, data interface{}) error {

	if res == nil {
		return errors.New("Res can't be null")
	}

	defer res.Body.Close()

	if res.IsError() {
		return errors.Errorf("Error when read response: %s", res.String())
	}

	// Need to convert on object
	esResponse := &ElasticsearchResponseSearch{}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	log.Debugf("Response:\n%s", string(body))

	err = json.Unmarshal(body, esResponse)
	if err != nil {
		return err
	}

	log.Debugf("esResponse: %+v", esResponse)

	// Check if document
	if esResponse.Hits.Total.Value <= 0 {
		return nil
	}

	// Search if data is array or object
	t := reflect.TypeOf(data)
	switch t.Kind() {
	case reflect.Ptr:
		log.Debugf("Detect pointer")
		switch t.Elem().Kind() {
		case reflect.Slice:
			log.Debugf("Detect slice")
			value := reflect.ValueOf(data).Elem()
			for _, hit := range esResponse.Hits.Hits {
				hit.Source["id"] = hit.ID
				obj := reflect.New(t.Elem().Elem())

				err = Decode(hit.Source, obj.Interface())
				if err != nil {
					return err
				}

				value.Set(reflect.Append(value, obj.Elem()))
			}
		default:
			return errors.New("You must provide pointer of slice or pointer of struct")
		}

	default:
		return errors.New("You must provide pointer of slice or pointer of struct")

	}

	log.Debugf("Data: %+v", data)

	return nil

}
