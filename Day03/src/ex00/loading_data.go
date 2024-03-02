package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
)

const scheme = `
{
  "mappings": {
    "properties": {
      "name": {
        "type":  "text"
      },
      "address": {
          "type":  "text"
      },
      "phone": {
         "type":  "text"
      },
      "location": {
        "type": "geo_point"
      }
    }
  }
}
`

type Location struct {
	Latitude   float64 `json:"lat"`
	Longtitude float64 `json:"lon"`
}

type Place struct {
	ID       uint64   `json:"id"`
	Name     string   `json:"name"`
	Address  string   `json:"address"`
	Phone    string   `json:"phone"`
	Location Location `json:"location"`
}

var (
	indexName  string
	fileName   string
	numWorkers int
	flushBytes int
)

func init() {
	flag.StringVar(&indexName, "index", "places", "Index name")
	flag.StringVar(&fileName, "f", "", "Data file path")
	flag.IntVar(&numWorkers, "workers", runtime.NumCPU(), "Number of indexer workers")
	flag.IntVar(&flushBytes, "flush", 5e+6, "Flush threshold in bytes")
	flag.Parse()
}

func main() {
	log.SetFlags(0)

	var (
		places []*Place
		err    error
	)

	log.Printf(
		"\x1b[1mBulkIndexer\x1b[0m: workers [%d] flush [%d]",
		numWorkers,
		flushBytes,
	)
	log.Println(strings.Repeat("▁", 65))

	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	// Create the BulkIndexer
	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:         indexName,        // The default index name
		Client:        es,               // The Elasticsearch client
		NumWorkers:    numWorkers,       // The number of worker goroutines
		FlushBytes:    flushBytes,       // The flush threshold in bytes
		FlushInterval: 30 * time.Second, // The periodic flush interval
	})
	if err != nil {
		log.Fatalf("Error creating the indexer: %s", err)
	}

	places, err = readPlacesData(fileName)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("→ Read %d places\n", len(places))

	if err = recreateIndex(es); err != nil {
		log.Fatal(err)
	}

	start := time.Now().UTC()

	_, err = addData(es, bi, places)
	if err != nil {
		log.Fatal(err)
	}

	biStats := bi.Stats()

	// Report the results: number of indexed docs, number of errors, duration, indexing rate
	//
	log.Println(strings.Repeat("▔", 65))

	dur := time.Since(start)

	if biStats.NumFailed > 0 {
		log.Fatalf(
			"Indexed [%d] documents with [%d] errors in %s (%d docs/sec)",
			int64(biStats.NumFlushed),
			int64(biStats.NumFailed),
			dur.Truncate(time.Millisecond),
			int64(1000.0/float64(dur/time.Millisecond)*float64(biStats.NumFlushed)),
		)
	} else {
		log.Printf(
			"Sucessfuly indexed [%d] documents in %s (%d docs/sec)",
			int64(biStats.NumFlushed),
			dur.Truncate(time.Millisecond),
			int64(1000.0/float64(dur/time.Millisecond)*float64(biStats.NumFlushed)),
		)
	}
}

func readPlacesData(filename string) ([]*Place, error) {
	var places []*Place
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	csvReader.Comma = '\t'
	_, err = csvReader.Read()
	if err != nil && err != io.EOF {
		return nil, err
	}

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		id, err := strconv.ParseUint(record[0], 10, 64)
		if err != nil {
			return nil, errors.New("Cannot convert id:" + err.Error())
		}
		lon, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			return nil, errors.New("Cannot convert id:" + err.Error())
		}
		lat, err := strconv.ParseFloat(record[5], 64)
		if err != nil {
			return nil, errors.New("Cannot convert id:" + err.Error())
		}

		places = append(places, &Place{
			ID:       id + 1,
			Name:     record[1],
			Address:  record[2],
			Phone:    record[3],
			Location: Location{Latitude: lat, Longtitude: lon},
		})
	}

	return places, nil
}

func recreateIndex(es *elasticsearch.Client) error {
	res, err := es.Indices.Delete(
		[]string{indexName},
		es.Indices.Delete.WithIgnoreUnavailable(true),
	)
	if err != nil || res.IsError() {
		return errors.New("Cannot delete index: " + err.Error())
	}
	res.Body.Close()
	res, err = es.Indices.Create(
		indexName,
		es.Indices.Create.WithBody(strings.NewReader(string(scheme))),
	)
	if err != nil {
		return errors.New("Cannot create index: " + err.Error())
	}
	if res.IsError() {
		return errors.New("Cannot create index: " + res.String())
	}
	res.Body.Close()
	return nil
}

func addData(
	es *elasticsearch.Client,
	bi esutil.BulkIndexer,
	places []*Place,
) (uint64, error) {

	var count uint64

	for _, place := range places {
		data, err := json.Marshal(place)
		if err != nil {
			return count, errors.New(fmt.Sprintf("Cannot encode place %d: %s", place.ID, err))
		}

		// Add an item to the BulkIndexer
		err = bi.Add(
			context.Background(),
			esutil.BulkIndexerItem{
				// Action field configures the operation to perform (index, create, delete, update)
				Action: "index",
				// DocumentID is the (optional) document ID
				DocumentID: strconv.FormatUint(place.ID, 10),
				// Body is an `io.Reader` with the payload
				Body: bytes.NewReader(data),
				// OnSuccess is called for each successful operation
				OnSuccess: func(
					ctx context.Context,
					item esutil.BulkIndexerItem,
					res esutil.BulkIndexerResponseItem,
				) {
					atomic.AddUint64(&count, 1)
				},
				// OnFailure is called for each failed operation
				OnFailure: func(
					ctx context.Context,
					item esutil.BulkIndexerItem,
					res esutil.BulkIndexerResponseItem,
					err error,
				) {
					if err != nil {
						log.Printf("ERROR: %s", err)
					} else {
						log.Printf("ERROR: %s: %s", res.Error.Type, res.Error.Reason)
					}
				},
			},
		)
		if err != nil {
			return count, errors.New(fmt.Sprintf("Unexpected error: %s", err))
		}
	}

	// Close the indexer
	if err := bi.Close(context.Background()); err != nil {
		return count, errors.New(fmt.Sprintf("Unexpected error: %s", err))
	}
	return count, nil
}
