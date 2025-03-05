package rag_sql

import (
	"context"

	"github.com/cloudwego/eino-ext/components/retriever/volc_vikingdb"
	"github.com/cloudwego/eino/components/retriever"
)

var cfg = &volc_vikingdb.RetrieverConfig{
	Host:       "api-vikingdb.volces.com",
	Region:     "cn-beijing",
	AK:         "your-ak",
	SK:         "your-sk",
	Collection: "your-collection",
	Index:      "your-index",
	EmbeddingConfig: volc_vikingdb.EmbeddingConfig{
		UseBuiltin:  true,
		ModelName:   "model-name",
		UseSparse:   true,
		DenseWeight: 0.5,
	},
	TopK: ptrOf(5),
}

// newRetriever component initialization function of node 'DDLRetriever' in graph 'text2sql'
func newRetriever(ctx context.Context) (rtr retriever.Retriever, err error) {
	config := cfg

	rtr, err = volc_vikingdb.NewRetriever(ctx, config)
	if err != nil {
		return nil, err
	}
	return rtr, nil
}

func ptrOf(i int) *int {
	p := &i
	return p
}

// newRetriever1 component initialization function of node 'DBDescriptionRetriever' in graph 'text2sql'
func newRetriever1(ctx context.Context) (rtr retriever.Retriever, err error) {
	config := cfg
	rtr, err = volc_vikingdb.NewRetriever(ctx, config)
	if err != nil {
		return nil, err
	}
	return rtr, nil
}

// newRetriever2 component initialization function of node 'Q2SQLRetriever' in graph 'text2sql'
func newRetriever2(ctx context.Context) (rtr retriever.Retriever, err error) {
	config := cfg
	rtr, err = volc_vikingdb.NewRetriever(ctx, config)
	if err != nil {
		return nil, err
	}
	return rtr, nil
}

// newRetriever3 component initialization function of node 'ThesaurusRetriever' in graph 'text2sql'
func newRetriever3(ctx context.Context) (rtr retriever.Retriever, err error) {
	config := cfg
	rtr, err = volc_vikingdb.NewRetriever(ctx, config)
	if err != nil {
		return nil, err
	}
	return rtr, nil
}
