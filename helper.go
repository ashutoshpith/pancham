package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

func GetParseBodyJsonPayload(c *fiber.Ctx, dto interface{}) error {
	body := c.Body()
	if !json.Valid(body) {
		log.Println("Invalid JSON data")
		return fiber.ErrBadRequest
	}

	err := json.Unmarshal(body, dto)
	if err != nil {
		log.Println("Error unmarshaling JSON:", err)
		return fiber.ErrBadRequest
	}
	return nil
}

func GetErrorResponse(c *fiber.Ctx, value interface{}) error {
	log.Println(value)
	if err, ok := value.(error); ok {
		log.Println("erorr coming")
		return c.SendString(err.Error())

	} else {
		str, _ := convertToString(value)
		log.Println("string error coming", str)
		return c.SendString(str)
	}

}

func GetAuthorizationHeader(c *fiber.Ctx) (string, error) {
	if len(c.Get("Authorization")) > 0 {
		return c.Get("Authorization"), nil
	}
	return "", &fiber.Error{
		Code:    404,
		Message: "Auth Token Not Found In Header",
	}

}

type AppDetailsDto struct {
	AppVersion string
	AppName    string
}

func GetAppDetailsFromHeader(c *fiber.Ctx) (AppDetailsDto, error) {
	if len(c.Get("AppVersion")) > 0 {
		return AppDetailsDto{
			AppVersion: c.Get("AppVersion"),
			AppName:    c.Get("AppName"),
		}, nil
	}
	return AppDetailsDto{}, &fiber.Error{
		Code:    404,
		Message: "Tenant Detail Not Found In Header",
	}
}

func convertToString(value interface{}) (string, error) {
	switch v := value.(type) {
	case string:
		return v, nil
	default:
		return "", fmt.Errorf("value is not of type string")
	}
}
