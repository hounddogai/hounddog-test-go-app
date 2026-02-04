package main

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"example.com/hounddog-test-go-app/proto/beerpb"
	"example.com/hounddog-test-go-app/utils/logging"
)

type BeerRecord struct {
	ID        string `gorm:"primaryKey"`
	Name      string
	Style     string
	ABV       float64
	CreatedAt time.Time
}

type BeerClient struct {
	db     *gorm.DB
	client beerpb.BeerServiceClient
}

func NewBeerClient(conn *grpc.ClientConn) (*BeerClient, error) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&BeerRecord{}); err != nil {
		return nil, err
	}
	return &BeerClient{db: db, client: beerpb.NewBeerServiceClient(conn)}, nil
}

func (c *BeerClient) CreateBeer(ctx context.Context, beer *beerpb.Beer) (*beerpb.Beer, error) {
	span := logging.Span{TraceID: "cli-trace", SpanID: "cli-create", Baggage: map[string]string{"op": "create"}}
	ctx = logging.ContextWithSpan(ctx, span)
	beerName := ""
	if beer != nil {
		beerName = beer.GetName()
	}
	logging.MyCustomlog.Infof(ctx, "client create beer: %s", beerName)

	resp, err := c.client.CreateBeer(ctx, &beerpb.CreateBeerRequest{Beer: beer})
	if err != nil {
		return nil, err
	}
	if resp.GetBeer() != nil {
		rec := BeerRecord{ID: resp.GetBeer().GetId(), Name: resp.GetBeer().GetName(), Style: resp.GetBeer().GetStyle(), ABV: resp.GetBeer().GetAbv()}
		if err := c.db.Create(&rec).Error; err != nil {
			return nil, err
		}
	}
	return resp.GetBeer(), nil
}

func (c *BeerClient) GetBeer(ctx context.Context, id string) (*beerpb.Beer, error) {
	span := logging.Span{TraceID: "cli-trace", SpanID: "cli-get", Baggage: map[string]string{"op": "get"}}
	ctx = logging.ContextWithSpan(ctx, span)
	logging.MyCustomlog.Infof(ctx, "client get beer: %s", id)

	resp, err := c.client.GetBeer(ctx, &beerpb.GetBeerRequest{Id: id})
	if err != nil {
		return nil, err
	}
	return resp.GetBeer(), nil
}

func (c *BeerClient) ListBeers(ctx context.Context) ([]*beerpb.Beer, error) {
	span := logging.Span{TraceID: "cli-trace", SpanID: "cli-list", Baggage: map[string]string{"op": "list"}}
	ctx = logging.ContextWithSpan(ctx, span)
	logging.MyCustomlog.Infof(ctx, "client list beers")

	resp, err := c.client.ListBeers(ctx, &beerpb.ListBeersRequest{})
	if err != nil {
		return nil, err
	}
	return resp.GetBeers(), nil
}

func (c *BeerClient) DeleteBeer(ctx context.Context, id string) (*beerpb.Beer, error) {
	span := logging.Span{TraceID: "cli-trace", SpanID: "cli-delete", Baggage: map[string]string{"op": "delete"}}
	ctx = logging.ContextWithSpan(ctx, span)
	logging.MyCustomlog.Infof(ctx, "client delete beer: %s", id)

	resp, err := c.client.DeleteBeer(ctx, &beerpb.DeleteBeerRequest{Id: id})
	if err != nil {
		return nil, err
	}
	if err := c.db.Delete(&BeerRecord{}, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return resp.GetBeer(), nil
}

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client, err := NewBeerClient(conn)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	_, _ = client.CreateBeer(ctx, &beerpb.Beer{Id: "1", Name: "Hoppy", Style: "IPA", Abv: 6.5})
	_, _ = client.ListBeers(ctx)
	_, _ = client.DeleteBeer(ctx, "1")
}
