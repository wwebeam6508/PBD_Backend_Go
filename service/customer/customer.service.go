package service

import (
	"PBD_backend_go/common"
	"PBD_backend_go/configuration"
	"PBD_backend_go/exception"
	model "PBD_backend_go/model/customer"
	"context"
	"os"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetCustomerService(input model.GetCustomerInput, searchPipeline model.SearchPipeline) ([]model.GetCustomerResult, error) {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return nil, err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("customers")
	matchState := bson.D{{Key: "$match", Value: bson.D{{Key: "status", Value: bson.D{{Key: "$eq", Value: 1}}}}}}
	if input.Page > 0 {
		input.Page = input.Page - 1
	}
	skipStage := bson.D{{Key: "$skip", Value: input.Page * input.PageSize}}
	limitStage := bson.D{{Key: "$limit", Value: input.PageSize}}
	projectStage := bson.D{{Key: "$project", Value: bson.D{
		{Key: "customerID", Value: "$_id"},
		{Key: "name", Value: 1},
		{Key: "address", Value: 1},
		{Key: "taxID", Value: 1},
	}}}
	pipeline := bson.A{matchState, projectStage, skipStage, limitStage}
	if searchPipeline.Search != "" {
		pipeline = append(pipeline, searchPipeline.SearchPipeline...)
	}
	if input.SortTitle != "" && input.SortType != "" {
		var sortValue int
		if input.SortType == "desc" {
			sortValue = -1
		} else {
			sortValue = 1
		}
		sortStage := bson.D{{Key: "$sort", Value: bson.D{{Key: input.SortTitle, Value: sortValue}}}}
		pipeline = append(pipeline, sortStage)
	}

	var result []model.GetCustomerResult
	cursor, err := ref.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.Background(), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetCustomerCountService(searchPipeline model.SearchPipeline) (int32, error) {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return 0, err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("customers")
	matchState := bson.D{{Key: "$match", Value: bson.D{{Key: "status", Value: bson.D{{Key: "$eq", Value: 1}}}}}}
	pipeline := bson.A{matchState}
	if searchPipeline.Search != "" {
		pipeline = append(pipeline, searchPipeline.SearchPipeline...)
	}
	groupStage := bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: nil}, {Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}}}}}
	pipeline = append(pipeline, groupStage)
	var result []bson.M
	cursor, err := ref.Aggregate(context.Background(), pipeline)
	if err != nil {
		return 0, err
	}
	err = cursor.All(context.Background(), &result)
	if err != nil {
		return 0, err
	}

	if len(result) == 0 {
		return 0, nil
	}
	return result[0]["count"].(int32), nil
}

func GetCustomerByIDService(input model.GetCustomerByIDInput) (model.GetCustomerByIDResult, error) {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return model.GetCustomerByIDResult{}, err
	}
	customerIDObjectID, err := primitive.ObjectIDFromHex(input.CustomerID)
	if err != nil {
		return model.GetCustomerByIDResult{}, exception.ValidationError{Message: "invalid customerID"}
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("customers")
	// aggregate
	matchState := bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: customerIDObjectID}, {Key: "status", Value: bson.D{{Key: "$eq", Value: 1}}}}}}
	projectStage := bson.D{{Key: "$project", Value: bson.D{
		{Key: "customerID", Value: "$_id"},
		{Key: "name", Value: 1},
		{Key: "address", Value: 1},
		{Key: "taxID", Value: 1},
		{Key: "emails", Value: 1},
		{Key: "phones", Value: 1},
	}}}
	pipeline := bson.A{matchState, projectStage}
	cursor, err := ref.Aggregate(context.Background(), pipeline)
	if err != nil {
		return model.GetCustomerByIDResult{}, err
	}
	var result []model.GetCustomerByIDResult
	err = cursor.All(context.Background(), &result)
	if err != nil {
		return model.GetCustomerByIDResult{}, err
	}

	//check result empty
	if len(result) == 0 {
		return model.GetCustomerByIDResult{}, exception.NotFoundError{Message: "customer not found"}
	}
	return result[0], nil
}

func AddCustomerService(input model.AddCustomerInput) (primitive.ObjectID, error) {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return primitive.NilObjectID, err
	}
	input.Status = 1
	input.CreatedAt = time.Now()
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("customers")
	insertResult, err := ref.InsertOne(context.Background(), input)
	if err != nil {
		return primitive.NilObjectID, err
	}
	if insertResult.InsertedID == nil {
		return primitive.NilObjectID, exception.ValidationError{Message: "add customer failed"}
	}
	return insertResult.InsertedID.(primitive.ObjectID), nil
}

func UpdateCustomerService(input model.UpdateCustomerInput, updateCustomerID model.UpdateCustomerID) error {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("customers")
	customerIDObjectID, err := primitive.ObjectIDFromHex(updateCustomerID.CustomerID)
	if err != nil {
		return exception.ValidationError{Message: "invalid customerID"}
	}
	input.UpdatedAt = time.Now()
	// check each field that not empty of input
	updateField := bson.D{}
	//dynamic check by for loop
	refValue := reflect.ValueOf(input)
	for i := 0; i < refValue.NumField(); i++ {
		if !common.IsEmpty(refValue.Field(i).Interface()) {
			if reflect.ValueOf(input).Type().Field(i).Name == "RemoveEmails" || reflect.ValueOf(input).Type().Field(i).Name == "addEmails" {
				continue
			}
			updateField = append(updateField, bson.E{Key: refValue.Type().Field(i).Tag.Get("json"), Value: refValue.Field(i).Interface()})
		}
	}
	filter := bson.D{{Key: "_id", Value: customerIDObjectID}, {Key: "status", Value: 1}}
	updateResult, err := ref.UpdateOne(context.Background(), filter, bson.D{{Key: "$set", Value: updateField}})
	if err != nil {
		return err
	}
	if updateResult.MatchedCount == 0 {
		return exception.NotFoundError{Message: "customer not found"}
	}
	if !common.IsEmpty(input.RemoveEmails) {
		updateField := bson.M{"$pull": bson.M{
			"emails": bson.M{
				"$in": input.RemoveEmails,
			},
		}}
		_, err := ref.UpdateOne(context.Background(), filter, updateField)
		if err != nil {
			return err
		}
	}
	if !common.IsEmpty(input.AddEmails) {
		updateField := bson.M{"$push": bson.M{"emails": bson.M{
			"$each": input.AddEmails,
		}}}
		_, err := ref.UpdateOne(context.Background(), filter, updateField)
		if err != nil {
			return err
		}
	}
	if !common.IsEmpty(input.RemovePhones) {
		updateField := bson.M{"$pull": bson.M{
			"phones": bson.M{
				"$in": input.RemovePhones,
			},
		}}
		_, err := ref.UpdateOne(context.Background(), filter, updateField)
		if err != nil {
			return err
		}
	}
	if !common.IsEmpty(input.AddPhones) {
		updateField := bson.M{"$push": bson.M{
			"phones": bson.M{
				"$each": input.AddPhones,
			}}}
		_, err := ref.UpdateOne(context.Background(), filter, updateField)
		if err != nil {
			return err
		}
	}

	return nil
}

func DeleteCustomerService(input model.DeleteCustomerInput) error {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("customers")
	customerIDObjectID, err := primitive.ObjectIDFromHex(input.CustomerID)
	if err != nil {
		return exception.ValidationError{Message: "invalid customerID"}
	}
	deleteResult, err := ref.UpdateOne(context.Background(), bson.D{{Key: "_id", Value: customerIDObjectID}, {Key: "status", Value: 1}}, bson.D{{Key: "$set", Value: bson.D{{Key: "status", Value: 0}}}})
	if err != nil {
		return err
	}
	if deleteResult.MatchedCount == 0 {
		return exception.NotFoundError{Message: "customer not found"}
	}
	return nil
}
