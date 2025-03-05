package rag_sql

import (
	"context"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/kitex/pkg/klog"
)

func init() {
	Text2SQL, err = buildText2sql(context.Background())
	if err != nil {
		klog.Fatalf("build graph failed: %v", err)
		return
	}
}

var Text2SQL compose.Runnable[*string, *string]
var err error

func buildText2sql(ctx context.Context) (r compose.Runnable[*string, *string], err error) {
	const (
		DDLRetriever           = "DDLRetriever"
		ChatTemplate1          = "ChatTemplate1"
		DBDescriptionRetriever = "DBDescriptionRetriever"
		Q2SQLRetriever         = "Q2SQLRetriever"
		ThesaurusRetriever     = "ThesaurusRetriever"
		ChatModel2             = "ChatModel2"
	)
	g := compose.NewGraph[*string, *string]()
	dDLRetrieverKeyOfRetriever, err := newRetriever(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddRetrieverNode(DDLRetriever, dDLRetrieverKeyOfRetriever, compose.WithNodeName("DDL"), compose.WithInputKey(""), compose.WithOutputKey("DDL"))
	chatTemplate1KeyOfChatTemplate, err := newChatTemplate(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddChatTemplateNode(ChatTemplate1, chatTemplate1KeyOfChatTemplate)
	dBDescriptionRetrieverKeyOfRetriever, err := newRetriever1(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddRetrieverNode(DBDescriptionRetriever, dBDescriptionRetrieverKeyOfRetriever, compose.WithNodeName("DBDescription"), compose.WithOutputKey("DBDescription"))
	q2SQLRetrieverKeyOfRetriever, err := newRetriever2(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddRetrieverNode(Q2SQLRetriever, q2SQLRetrieverKeyOfRetriever, compose.WithNodeName("Q2SQL"), compose.WithOutputKey("Q2SQL"))
	thesaurusRetrieverKeyOfRetriever, err := newRetriever3(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddRetrieverNode(ThesaurusRetriever, thesaurusRetrieverKeyOfRetriever, compose.WithNodeName("Thesaurus"), compose.WithInputKey(""), compose.WithOutputKey("Thesaurus"))
	chatModel2KeyOfChatModel, err := newChatModel(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddChatModelNode(ChatModel2, chatModel2KeyOfChatModel)
	_ = g.AddEdge(compose.START, DDLRetriever)
	_ = g.AddEdge(compose.START, DBDescriptionRetriever)
	_ = g.AddEdge(compose.START, Q2SQLRetriever)
	_ = g.AddEdge(compose.START, ThesaurusRetriever)
	_ = g.AddEdge(ChatModel2, compose.END)
	_ = g.AddEdge(DDLRetriever, ChatTemplate1)
	_ = g.AddEdge(DBDescriptionRetriever, ChatTemplate1)
	_ = g.AddEdge(ThesaurusRetriever, ChatTemplate1)
	_ = g.AddEdge(Q2SQLRetriever, ChatTemplate1)
	_ = g.AddEdge(ChatTemplate1, ChatModel2)
	r, err = g.Compile(ctx, compose.WithGraphName("text2sql"), compose.WithNodeTriggerMode(compose.AllPredecessor))
	if err != nil {
		return nil, err
	}
	return r, err
}
