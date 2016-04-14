package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/blevesearch/bleve"
	"github.com/yanyiwu/gojieba"
	_ "github.com/yanyiwu/gojieba/bleve"
)

func Example() {
	INDEX_DIR := "gojieba.bleve"
	messages := []struct {
		Id   string
		Body string
	}{
		{
			Id:   "1",
			Body: "你好",
		},
		{
			Id:   "2",
			Body: "世界",
		},
		{
			Id:   "3",
			Body: "亲口",
		},
		{
			Id:   "4",
			Body: "交代",
		},
	}

	indexMapping := bleve.NewIndexMapping()
	os.RemoveAll(INDEX_DIR)
	// clean index when example finished
	defer os.RemoveAll(INDEX_DIR)

	err := indexMapping.AddCustomTokenizer("gojieba",
		map[string]interface{}{
			"dictpath":     gojieba.DICT_PATH,
			"hmmpath":      gojieba.HMM_PATH,
			"userdictpath": gojieba.USER_DICT_PATH,
			"type":         "gojieba",
		},
	)
	if err != nil {
		panic(err)
	}
	err = indexMapping.AddCustomAnalyzer("gojieba",
		map[string]interface{}{
			"type":      "gojieba",
			"tokenizer": "gojieba",
		},
	)
	if err != nil {
		panic(err)
	}
	indexMapping.DefaultAnalyzer = "gojieba"

	index, err := bleve.New(INDEX_DIR, indexMapping)
	if err != nil {
		panic(err)
	}
	for _, msg := range messages {
		if err := index.Index(msg.Id, msg); err != nil {
			panic(err)
		}
	}

	querys := []string{
		"你好世界",
		"亲口交代",
	}

	for _, q := range querys {
		req := bleve.NewSearchRequest(bleve.NewQueryStringQuery(q))
		req.Highlight = bleve.NewHighlight()
		res, err := index.Search(req)
		if err != nil {
			panic(err)
		}
		x, err := json.Marshal(res.Hits)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(x))
	}

	// Output:
	// [{"id":"2","score":0.4232867878957415,"locations":{"Body":{"世界":[{"pos":1,"start":0,"end":6,"array_positions":null}]}},"fragments":{"Body":["\u003cmark\u003e世界\u003c/mark\u003e"]}},{"id":"1","score":0.4232867878957415,"locations":{"Body":{"你好":[{"pos":1,"start":0,"end":6,"array_positions":null}]}},"fragments":{"Body":["\u003cmark\u003e你好\u003c/mark\u003e"]}}]
	// [{"id":"4","score":0.4232867878957415,"locations":{"Body":{"交代":[{"pos":1,"start":0,"end":6,"array_positions":null}]}},"fragments":{"Body":["\u003cmark\u003e交代\u003c/mark\u003e"]}},{"id":"3","score":0.4232867878957415,"locations":{"Body":{"亲口":[{"pos":1,"start":0,"end":6,"array_positions":null}]}},"fragments":{"Body":["\u003cmark\u003e亲口\u003c/mark\u003e"]}}]
}
