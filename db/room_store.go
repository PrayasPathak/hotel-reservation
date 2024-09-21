package db

import (
	"context"

	"github.com/PrayasPathak/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const RoomCollection = "rooms"

type RoomStore interface {
	InsertRoom(context.Context, *types.Room) (*types.Room, error)
}

type MongoRoomStore struct {
	client     *mongo.Client
	collection *mongo.Collection
	HotelStore
}

func NewMongoRoomStore(client *mongo.Client, hotelStore HotelStore) *MongoRoomStore {
	return &MongoRoomStore{
		client:     client,
		collection: client.Database(DBNAME).Collection(RoomCollection),
		HotelStore: hotelStore,
	}
}

func (s *MongoRoomStore) InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error) {
	res, err := s.collection.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}
	room.ID = res.InsertedID.(primitive.ObjectID)
	// update the hotel room with this id
	filter := bson.M{
		"_id": room.HotelID,
	}
	update := bson.M{
		"$push": bson.M{
			"rooms": room.ID,
		},
	}
	if err := s.HotelStore.UpdateHotel(ctx, filter, update); err != nil {
		return nil, err
	}
	return room, nil
}
