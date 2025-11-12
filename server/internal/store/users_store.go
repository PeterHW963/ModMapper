package store

import (
	"context"
	"modmapper/server/internal/models"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UsersStore interface {
	List(ctx context.Context, query string, limit int64, skip int64) ([]models.User, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Create(ctx context.Context, user models.User) (models.User, error)
	Update(ctx context.Context, id primitive.ObjectID, fields bson.M) (*models.User, error)
	Delete(ctx context.Context, id primitive.ObjectID) error
}

type usersMongoStore struct {
	collection *mongo.Collection
}

// Repo constructor / wiring point to the collection. Wraps the handle in the usersMongoStore
// Returns a UsersStore interface
func NewUsersStore(db *mongo.Database) UsersStore {
	return &usersMongoStore{collection: db.Collection("users")}
}

func (s *usersMongoStore) List(ctx context.Context, query string, limit int64, skip int64) ([]models.User, error) {
	filter := bson.M{} // empty filter to match every doc
	if query != "" {
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": query, "$options": "i"}},
			{"email": bson.M{"$regex": query, "$options": "i"}},
		}
	}

	// sort by creation date (descending)
	opts := options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: -1}}).
		SetLimit(limit).
		SetSkip(skip)

	// get cursor: pointer or iterator to the results of a query
	cursor, err := s.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	var out []models.User
	if err := cursor.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *usersMongoStore) GetByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	var user models.User
	if err := s.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (s *usersMongoStore) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	var user models.User
	if err := s.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (s *usersMongoStore) Create(ctx context.Context, user models.User) (models.User, error) {
	if user.ID.IsZero() {
		user.ID = primitive.NewObjectID()
	}
	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now().UTC()
	}
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))
	_, err := s.collection.InsertOne(ctx, user)
	return user, err
}

func (s *usersMongoStore) Update(ctx context.Context, id primitive.ObjectID, fields bson.M) (*models.User, error) {
	allowedFields := bson.M{}
	// val takes the value stored at the key, ok = true if the key exists, false if it doesn't
	val, ok := fields["name"]
	if ok {
		allowedFields["name"] = val
	}
	val, ok = fields["telegramHandle"]
	if ok {
		allowedFields["telegramHandle"] = val
	}
	after := options.After
	opts := options.FindOneAndUpdate().SetReturnDocument(after)
	var out models.User
	if err := s.collection.FindOneAndUpdate(ctx, bson.M{"_id": id}, bson.M{"$set": allowedFields}, opts).Decode(&out); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &out, nil
}

func (s *usersMongoStore) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := s.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
