package middleware

import (
	"PBD_backend_go/exception"
	model "PBD_backend_go/model"
	service "PBD_backend_go/service"
	authservice "PBD_backend_go/service/auth"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func RankCheck(c *fiber.Ctx) error {
	//get body userID and userTypeID from body
	var body struct {
		UserID     string `json:"userID" bson:"userID"`
		UserTypeID string `json:"userTypeID" bson:"userTypeID"`
		Rank       int32  `json:"rank" bson:"rank"`
	}
	err := c.BodyParser(&body)
	if err != nil {
		err := c.QueryParser(&body)
		if err != nil {
			return exception.ErrorHandler(c, err)
		}
	}
	token := c.Get("Authorization")
	splitToken := strings.Split(token, " ")
	if len(splitToken) < 2 {
		return exception.ValidationError{Message: "invalid token"}
	}
	claims, err := authservice.VerifyJWT(splitToken[1])
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	userData := claims.Claims.(jwt.MapClaims)["data"].(map[string]interface{})
	//check if userID is empty even it spacing
	if body.UserID != "" {
		err := againistOther(userData, body.UserID)
		if err != nil {
			return exception.ErrorHandler(c, err)
		}
	} else if body.UserTypeID != "" {
		//get userTypeID from claim
		err := againistOtherType(userData, body.UserTypeID)
		if err != nil {
			return exception.ErrorHandler(c, err)
		}
	} else if body.Rank != 0 {
		//get rank from claim
		err := againistOtherRank(userData, body.Rank)
		if err != nil {
			return exception.ErrorHandler(c, err)
		}
	}
	return c.Next()
}

func againistOther(userData map[string]interface{}, userID string) error {
	selfUserID := userData["userID"].(string)
	//check if userID from body is equal to userID from claim
	if userID == selfUserID {
		return exception.ValidationError{Message: "cannot change your own data"}
	}
	//get userTypeID from input.SelfID
	rank, err := service.GetUserRankByUserIDService(model.GetUserTypeByUserIDInput{UserID: userID})
	if err != nil {
		return err
	}
	selfRank, err := service.GetUserRankByUserIDService(model.GetUserTypeByUserIDInput{UserID: selfUserID})
	if err != nil {
		return err
	}
	if selfRank.Rank >= rank.Rank {
		return exception.ValidationError{Message: "cannot change rank higher than or equal to your rank"}
	}

	return nil
}

func againistOtherType(userData map[string]interface{}, userTypeID string) error {
	selfUserTypeID := userData["userType"].(map[string]interface{})["userTypeID"].(string)
	selfUserID := userData["userID"].(string)
	//get rank from userTypeID
	if userTypeID == selfUserTypeID {
		return exception.ValidationError{Message: "cannot change your own data"}
	}
	//get userTypeID from input.SelfID
	rank, err := service.GetUserRankByUserTypeIDService(model.GetUserRankByUserTypeIDInput{UserTypeID: userTypeID})
	if err != nil {
		return err
	}
	selfRank, err := service.GetUserRankByUserIDService(model.GetUserTypeByUserIDInput{UserID: selfUserID})
	if err != nil {
		return err
	}
	if selfRank.Rank >= rank.Rank {
		return exception.ValidationError{Message: "cannot access or change data rank higher than or equal to your rank"}
	}

	return nil
}

func againistOtherRank(userData map[string]interface{}, rank int32) error {
	selfUserID := userData["userID"].(string)

	selfRank, err := service.GetUserRankByUserIDService(model.GetUserTypeByUserIDInput{UserID: selfUserID})
	if err != nil {
		return err
	}
	if selfRank.Rank >= rank {
		return exception.ValidationError{Message: "cannot access or change data rank higher than or equal to your rank"}
	}
	return nil
}
