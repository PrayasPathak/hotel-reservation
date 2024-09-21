package db

import (
	"context"

	"github.com/PrayasPathak/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const HotelCollection = "hotels"

type HotelStore interface {
	GetHotels(context.Context, bson.M) ([]*types.Hotel, error)
	GetHotelById(context.Context, string) (*types.Hotel, error)
	InsertHotel(context.Context, *types.Hotel) (*types.Hotel, error)
	UpdateHotel(context.Context, bson.M, bson.M) error
}

type MongoHotelStore struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client) *MongoHotelStore {
	return &MongoHotelStore{
		client:     client,
		collection: client.Database(DBNAME).Collection(HotelCollection),
	}
}

func (s *MongoHotelStore) InsertHotel(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
	res, err := s.collection.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}
	hotel.ID = res.InsertedID.(primitive.ObjectID)
	return hotel, nil
}

func (s *MongoHotelStore) UpdateHotel(ctx context.Context, filter, update bson.M) error {
	_, err := s.collection.UpdateOne(ctx, filter, update)
	return err
}

func (s *MongoHotelStore) GetHotels(ctx context.Context, filter bson.M) ([]*types.Hotel, error) {
	cur, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var hotels []*types.Hotel
	if err := cur.All(ctx, &hotels); err != nil {
		return nil, err
	}
	return hotels, nil
}

func (s *MongoHotelStore) GetHotelById(ctx context.Context, id string) (*types.Hotel, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var hotel types.Hotel
	if err := s.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&hotel); err != nil {
		return nil, err
	}
	return &hotel, nil

}
