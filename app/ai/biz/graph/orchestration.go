package graph

import (
	"context"

	"github.com/cloudwego/eino/compose"
)

func Buildtext2sql(ctx context.Context) (r compose.Runnable[*string, *string], err error) {
	const (
		DDLRetriever           = "DDLRetriever"
		ChatTemplate1          = "ChatTemplate1"
		DBDescriptionRetriever = "DBDescriptionRetriever"
		Q2SQLRetriever         = "Q2SQLRetriever"
		ThesaurusRetriever     = "ThesaurusRetriever"
		ChatModel2             = "ChatModel2"
		ExeSQL                 = "ExeSQL"
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
	_ = g.AddLambdaNode(ExeSQL, compose.InvokableLambda(newLambda))
	_ = g.AddEdge(compose.START, DDLRetriever)
	_ = g.AddEdge(compose.START, DBDescriptionRetriever)
	_ = g.AddEdge(compose.START, Q2SQLRetriever)
	_ = g.AddEdge(compose.START, ThesaurusRetriever)
	_ = g.AddEdge(ExeSQL, compose.END)
	_ = g.AddEdge(DDLRetriever, ChatTemplate1)
	_ = g.AddEdge(DBDescriptionRetriever, ChatTemplate1)
	_ = g.AddEdge(ThesaurusRetriever, ChatTemplate1)
	_ = g.AddEdge(Q2SQLRetriever, ChatTemplate1)
	_ = g.AddEdge(ChatTemplate1, ChatModel2)
	_ = g.AddEdge(ChatModel2, ExeSQL)
	r, err = g.Compile(ctx, compose.WithGraphName("text2sql"), compose.WithNodeTriggerMode(compose.AnyPredecessor))
	if err != nil {
		return nil, err
	}
	return r, err
}
