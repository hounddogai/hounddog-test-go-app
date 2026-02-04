package main

import (
	"context"
	"net"
	"sync"

	"google.golang.org/grpc"
	"example.com/hounddog-test-go-app/proto/beerpb"
	"example.com/hounddog-test-go-app/utils/logging"
)

type beerServer struct {
	beerpb.UnimplementedBeerServiceServer
	mu    sync.RWMutex
	beers map[string]*beerpb.Beer
}

func newBeerServer() *beerServer {
	return &beerServer{beers: make(map[string]*beerpb.Beer)}
}

func (s *beerServer) CreateBeer(ctx context.Context, req *beerpb.CreateBeerRequest) (*beerpb.BeerResponse, error) {
	span := logging.Span{TraceID: "srv-trace", SpanID: "srv-create", Baggage: map[string]string{"route": "create"}}
	ctx = logging.ContextWithSpan(ctx, span)
	beer := req.GetBeer()
	beerName := ""
	if beer != nil {
		beerName = beer.GetName()
	}
	logging.MyCustomlog.Infof(ctx, "create beer request: %s", beerName)

	if beer == nil || beer.Id == "" {
		return &beerpb.BeerResponse{}, nil
	}
	clone := *beer
	
	s.mu.Lock()
	s.beers[beer.Id] = &clone
	s.mu.Unlock()

	return &beerpb.BeerResponse{Beer: &clone}, nil
}

func (s *beerServer) GetBeer(ctx context.Context, req *beerpb.GetBeerRequest) (*beerpb.BeerResponse, error) {
	span := logging.Span{TraceID: "srv-trace", SpanID: "srv-get", Baggage: map[string]string{"route": "get"}}
	ctx = logging.ContextWithSpan(ctx, span)
	logging.MyCustomlog.Infof(ctx, "get beer request: %s", req.GetId())

	s.mu.RLock()
	beer := s.beers[req.GetId()]
	s.mu.RUnlock()
	if beer == nil {
		return &beerpb.BeerResponse{}, nil
	}
	clone := *beer
	return &beerpb.BeerResponse{Beer: &clone}, nil
}

func (s *beerServer) ListBeers(ctx context.Context, _ *beerpb.ListBeersRequest) (*beerpb.ListBeersResponse, error) {
	span := logging.Span{TraceID: "srv-trace", SpanID: "srv-list", Baggage: map[string]string{"route": "list"}}
	ctx = logging.ContextWithSpan(ctx, span)
	logging.MyCustomlog.Infof(ctx, "list beers request")

	s.mu.RLock()
	beers := make([]*beerpb.Beer, 0, len(s.beers))
	for _, beer := range s.beers {
		clone := *beer
		beers = append(beers, &clone)
	}
	s.mu.RUnlock()

	return &beerpb.ListBeersResponse{Beers: beers}, nil
}

func (s *beerServer) DeleteBeer(ctx context.Context, req *beerpb.DeleteBeerRequest) (*beerpb.BeerResponse, error) {
	span := logging.Span{TraceID: "srv-trace", SpanID: "srv-delete", Baggage: map[string]string{"route": "delete"}}
	ctx = logging.ContextWithSpan(ctx, span)
	logging.MyCustomlog.Infof(ctx, "delete beer request: %s", req.GetId())

	s.mu.Lock()
	beer := s.beers[req.GetId()]
	delete(s.beers, req.GetId())
	s.mu.Unlock()
	if beer == nil {
		return &beerpb.BeerResponse{}, nil
	}
	clone := *beer
	return &beerpb.BeerResponse{Beer: &clone}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	beerpb.RegisterBeerServiceServer(grpcServer, newBeerServer())

	if err := grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}
