package data

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Collection struct {
	*mongo.Collection
}

func Ctx() context.Context {
	return ctx()
}

func ctx() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	return ctx
}

// FindByID method finds a doc and decodes it to a model, otherwise returns an error.
// The id field can be any value that if passed to the `PrepareID` method, it returns
// a valid ID (e.g string, bson.ObjectId).
func (coll *Collection) FindByID(id interface{}, model Model, opts ...options.Lister[options.FindOneOptions]) error {
	return coll.FindByIDWithCtx(ctx(), id, model, opts...)
}

// FindByIDWithCtx method finds a doc and decodes it to a model, otherwise returns an error.
// The id field can be any value that if passed to the `PrepareID` method, it returns
// a valid ID (e.g string, bson.ObjectId).
func (coll *Collection) FindByIDWithCtx(ctx context.Context, id interface{}, model Model, opts ...options.Lister[options.FindOneOptions]) error {
	id, err := model.PrepareID(id)

	if err != nil {
		return err
	}

	return first(ctx, coll, bson.M{ID: id}, model, opts...)
}

// First method searches and returns the first document in the search results.
func (coll *Collection) First(filter interface{}, model Model, opts ...options.Lister[options.FindOneOptions]) error {
	return coll.FirstWithCtx(ctx(), filter, model, opts...)
}

// FirstWithCtx method searches and returns the first document in the search results.
func (coll *Collection) FirstWithCtx(ctx context.Context, filter interface{}, model Model, opts ...options.Lister[options.FindOneOptions]) error {
	return coll.FindOne(ctx, filter, opts...).Decode(model)
}

// SimpleFind finds, decodes and returns the results.
func (coll *Collection) SimpleFind(results interface{}, filter interface{}, opts ...options.Lister[options.FindOptions]) error {
	return coll.SimpleFindWithCtx(ctx(), results, filter, opts...)
}

// SimpleFindWithCtx finds, decodes and returns the results using the specified context.
func (coll *Collection) SimpleFindWithCtx(ctx context.Context, results interface{}, filter interface{}, opts ...options.Lister[options.FindOptions]) error {
	cur, err := coll.Find(ctx, filter, opts...)

	if err != nil {
		return err
	}

	return cur.All(ctx, results)
}

// Create method inserts a new model into the database.
func (coll *Collection) Create(model Model, opts ...options.Lister[options.InsertOneOptions]) error {
	return coll.CreateWithCtx(ctx(), model, opts...)
}

// CreateWithCtx method inserts a new model into the database.
func (coll *Collection) CreateWithCtx(ctx context.Context, model Model, opts ...options.Lister[options.InsertOneOptions]) error {
	return create(ctx, coll, model, opts...)
}

// Delete method deletes a model (doc) from a collection.
// To perform additional operations when deleting a model
// you should use hooks rather than overriding this method.
func (coll *Collection) Delete(model Model) error {
	return del(ctx(), coll, model)
}

// DeleteWithCtx method deletes a model (doc) from a collection using the specified context.
// To perform additional operations when deleting a model
// you should use hooks rather than overriding this method.
func (coll *Collection) DeleteWithCtx(ctx context.Context, model Model) error {
	return del(ctx, coll, model)
}

// Update function persists the changes made to a model to the database.
// Calling this method also invokes the model's updating, updated,
// saving, and saved hooks.
func (coll *Collection) Update(model Model, opts ...options.Lister[options.UpdateOneOptions]) error {
	return coll.UpdateWithCtx(ctx(), model, opts...)
}

// UpdateWithCtx function persists the changes made to a model to the database using the specified context.
// Calling this method also invokes the model's updating, updated,
// saving, and saved hooks.
func (coll *Collection) UpdateWithCtx(ctx context.Context, model Model, opts ...options.Lister[options.UpdateOneOptions]) error {
	return update(ctx, coll, model, opts...)
}
