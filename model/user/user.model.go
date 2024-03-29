package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetUserControllerInput struct {
	Page         int    `json:"page" bson:"page"`
	PageSize     int    `json:"pageSize" bson:"pageSize"`
	SortTitle    string `json:"sortTitle" bson:"sortTitle"`
	SortType     string `json:"sortType" bson:"sortType"`
	Search       string `json:"search" bson:"search"`
	SearchFilter string `json:"searchFilter" bson:"searchFilter"`
}

type SearchPipeline struct {
	Search         string `json:"search" bson:"search"`
	SearchPipeline bson.A `json:"searchPipeline" bson:"searchPipeline"`
}

type GetUserServiceResult struct {
	UserID   primitive.ObjectID `json:"userID" bson:"userID"`
	UserType string             `json:"userType" bson:"userType"`
	Username string             `json:"username" bson:"username"`
	Date     time.Time          `json:"date" bson:"date"`
}

type GetUserServiceInput struct {
	Page      int    `json:"page" bson:"page"`
	PageSize  int    `json:"pageSize" bson:"pageSize"`
	SortTitle string `json:"sortTitle" bson:"sortTitle"`
	SortType  string `json:"sortType" bson:"sortType"`
}

type GetUserByIDServiceResult struct {
	UserID   primitive.ObjectID `json:"userID" bson:"userID"`
	Username string             `json:"username" bson:"username"`
	UserType primitive.ObjectID `json:"userType" bson:"userType"`
}

type GetUserByIDInput struct {
	UserID string `json:"userID" bson:"userID"`
}

type AddUserInput struct {
	Username   string    `json:"username" bson:"username" validate:"required"`
	Password   string    `json:"password" bson:"password" validate:"required"`
	UserTypeID string    `json:"userType" bson:"userType" validate:"required"`
	Status     int       `json:"status" bson:"status"`
	CreatedAt  time.Time `json:"createdAt" bson:"createdAt"`
}

type UpdateUserID struct {
	UserID string `json:"userID" bson:"userID" validate:"required"`
}

type UpdateUserInput struct {
	Username   string    `json:"username" bson:"username"`
	UserTypeID string    `json:"userType" bson:"userType"`
	Password   string    `json:"password" bson:"password"`
	SelfID     string    `json:"itSelftID" bson:"itSelftID"`
	UpdatedAt  time.Time `json:"updatedAt" bson:"updatedAt"`
}

type DeleteUserInput struct {
	UserID string `json:"userID" bson:"userID"`
	SelfID string `json:"itSelftID" bson:"itSelftID"`
}

type GetUserTypeServiceResult struct {
	ID   primitive.ObjectID `json:"id" bson:"id"`
	Name string             `json:"name" bson:"name"`
	Rank int32              `json:"rank" bson:"rank"`
}

type GetUserTypeNameResult struct {
	ID   primitive.ObjectID `json:"id" bson:"id"`
	Name string             `json:"name" bson:"name"`
}
