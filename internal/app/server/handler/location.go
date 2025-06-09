package handler

import (
	"errors"
	"fmt"
	"github.com/allansbo/goapi/internal/app/server/dto"
	"github.com/allansbo/goapi/internal/domain/usecase"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"log/slog"
	"strings"
)

func makeValidation(data any) *fiber.Error {
	dataValidator := &XValidator{validator: validate}

	if errs := dataValidator.Validate(data); len(errs) > 0 && errs[0].Error {
		errMsgs := make([]string, 0)

		for _, err := range errs {
			errMsgs = append(errMsgs, fmt.Sprintf(
				"[%s]: '%v' | Needs to implement '%s'",
				err.FailedField,
				err.Value,
				err.Tag,
			))
		}

		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: strings.Join(errMsgs, " and "),
		}
	}

	return nil
}

func LocationsAdd(c *fiber.Ctx) error {
	locationDataIn := new(dto.LocationInApp)
	if err := c.BodyParser(locationDataIn); err != nil {
		slog.Error("error parsing locationDataIn", "error", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(GlobalErrorHandlerResp{
			Success: false,
			Message: "there is an error processing the location data provided",
			Error:   err.Error(),
		})
	}

	if err := makeValidation(locationDataIn); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(GlobalErrorHandlerResp{
			Success: false,
			Message: "There is an error validating the location data provided",
			Error:   err.Error(),
		})
	}

	locationDataOut, err := usecase.SaveLocation(locationDataIn)
	if err != nil {
		slog.Error("error saving location", "error", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(GlobalErrorHandlerResp{
			Success: false,
			Message: "there is an error saving the location data provided",
			Error:   err.Error(),
		})
	}

	return c.JSON(fiber.Map{"document_id": locationDataOut.ID})
}

func LocationsGet(c *fiber.Ctx) error {
	locationID := c.Params("id")

	locationDataOut, err := usecase.GetLocationById(locationID)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": fmt.Sprintf("The location ID %s does not exist", locationID),
		})
	} else if err != nil {
		slog.Error("error getting location", "error", err.Error(), "locationID", locationID)
		return c.Status(fiber.StatusInternalServerError).JSON(GlobalErrorHandlerResp{})
	}

	return c.JSON(locationDataOut)
}

func LocationsUpdate(c *fiber.Ctx) error {
	locationID := c.Params("id")
	locationDataIn := new(dto.LocationInApp)
	if err := c.BodyParser(locationDataIn); err != nil {
		slog.Error("error parsing locationDataInUpdate", "error", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(GlobalErrorHandlerResp{
			Success: false,
			Message: "there is an error processing the location data provided",
			Error:   err.Error(),
		})
	}

	if err := makeValidation(locationDataIn); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(GlobalErrorHandlerResp{
			Success: false,
			Message: "There is an error validating the location data provided",
			Error:   err.Error(),
		})
	}

	locationUpdated, err := usecase.UpdateLocation(locationID, locationDataIn)
	if err != nil {
		slog.Error("error updating location", "error", err.Error(), "locationID", locationID, "locationData", locationDataIn)
		return c.Status(fiber.StatusInternalServerError).JSON(GlobalErrorHandlerResp{
			Success: false,
			Message: "there is an error updating the location data provided",
			Error:   err.Error(),
		})
	}

	if locationUpdated {
		return c.JSON(fiber.Map{"message": "location updated successfully"})
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "no location was updated"})
}

func LocationsDelete(c *fiber.Ctx) error {
	locationID := c.Params("id")

	locationDeleted, err := usecase.DeleteLocation(locationID)
	if err != nil {
		slog.Error("error deleting location", "error", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(GlobalErrorHandlerResp{})
	}

	if locationDeleted {
		return c.JSON(fiber.Map{"message": "location deleted successfully"})
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "no location was deleted"})
}
