package transport

import (
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type GrpcServer struct {
	router grpctransport.Handler
}

/*func NewGRPCServer(ep endpoints.Set) *GrpcServer {
	return &GrpcServer{
		router: grpctransport.NewServer(
			ep.RouterEndpoint,
			decodeGRPCRouterRequest,
			decodeGRPCRouterResponse,
		),
	}
}

func (g *GrpcServer) Get(ctx context.Context, r endpoints.RouterRequest) (*api.GetReply, error) {
	_, rep, err := g.router.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*watermark.GetReply), nil
}

func decodeGRPCRouterRequest(_ context.Context) (interface{}, error) {
	return endpoints.RouterRequest{}, nil
}

func decodeGRPCRouterResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*watermark.GetReply)
	var docs []internal.Document
	for _, d := range reply.Documents {
		doc := internal.Document{
			Content:   d.Content,
			Title:     d.Title,
			Author:    d.Author,
			Topic:     d.Topic,
			Watermark: d.Watermark,
		}
		docs = append(docs, doc)
	}
	return endpoints.GetResponse{Documents: docs, Err: reply.Err}, nil
}*/
