package handler

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/allansbo/goapi/internal/app/server/dto"
	"github.com/allansbo/goapi/internal/domain/usecase"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/mongo"
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

// LocationsAdd godoc
//
//	@Summary		Insert location data
//	@Description	Insert location data into database
//	@Tags			Locations
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.LocationInApp				true	"Request of creating location object"
//	@Success		201		{object}	dto.LocationCreatedResponseOut	"document created"
//	@Failure		400		{object}	GlobalErrorHandlerResp			"validation error"
//	@Failure		500		{object}	GlobalErrorHandlerResp			"internal server error"
//	@Router			/api/v1/locations [post]
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

	return c.Status(fiber.StatusCreated).JSON(dto.LocationCreatedResponseOut{DocumentID: locationDataOut.ID})
}

// LocationsGet godoc
//
//	@Summary		Get location data
//	@Description	Get location data from database based on a document_id
//	@Tags			Locations
//	@Param			id	path	string	true	"id from document"
//	@Produce		json
//	@Success		200	{object}	dto.LocationOutApp				"located document"
//	@Success		404	{object}	dto.DefaultResponseMessageOut	"document not found"
//	@Failure		500	{object}	GlobalErrorHandlerResp			"internal server error"
//	@Router			/api/v1/locations/{id} [get]
func LocationsGet(c *fiber.Ctx) error {
	locationID := c.Params("id")

	locationDataOut, err := usecase.GetLocationById(locationID)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return c.Status(fiber.StatusNotFound).JSON(dto.DefaultResponseMessageOut{
				Message: fmt.Sprintf("the location ID %s does not exist", locationID),
		})
	} else if err != nil {
		slog.Error("error getting location", "error", err.Error(), "locationID", locationID)
		return c.Status(fiber.StatusInternalServerError).JSON(GlobalErrorHandlerResp{
			Success: false,
			Message: fmt.Sprintf("there is an error getting the location data by provided id %s", locationID),
			Error:   err.Error(),
		})
	}

	return c.JSON(locationDataOut)
}

// LocationsUpdate godoc
//
//	@Summary		Update location data
//	@Description	Update location data into database based on a document_id
//	@Tags			Locations
//	@Param			id	path	string	true	"id from document"
//	@Produce		json
//	@Param			request	body		dto.LocationInApp				true	"Request of updating location object"
//	@Success		200		{object}	dto.DefaultResponseMessageOut	"updated document"
//	@Failure		400		{object}	GlobalErrorHandlerResp			"validation error"
//	@Success		404		{object}	dto.DefaultResponseMessageOut	"document not found"
//	@Failure		500		{object}	GlobalErrorHandlerResp			"internal server error"
//	@Router			/api/v1/locations/{id} [put]
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
		return c.JSON(dto.DefaultResponseMessageOut{
			Message: fmt.Sprintf("the location ID %s has been updated", locationID),
		})
	}

	return c.Status(fiber.StatusNotFound).JSON(dto.DefaultResponseMessageOut{
		Message: fmt.Sprintf("the location ID %s does not exist", locationID),
	})
}

// LocationsDelete godoc
//
//	@Summary		Delete location data
//	@Description	Delete location data from database based on a document_id
//	@Tags			Locations
//	@Param			id	path	string	true	"id from document"
//	@Produce		json
//	@Success		200	{object}	dto.DefaultResponseMessageOut	"deleted document"
//	@Success		404	{object}	dto.DefaultResponseMessageOut	"document not found"
//	@Success		500	{object}	GlobalErrorHandlerResp			"internal server error"
//	@Router			/api/v1/locations/{id} [delete]
func LocationsDelete(c *fiber.Ctx) error {
	locationID := c.Params("id")

	locationDeleted, err := usecase.DeleteLocation(locationID)
	if err != nil {
		slog.Error("error deleting location", "error", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(GlobalErrorHandlerResp{
			Success: false,
			Message: "there is an error deleting the location data provided",
			Error:   err.Error(),
		})
	}

	if locationDeleted {
		return c.JSON(dto.DefaultResponseMessageOut{
			Message: fmt.Sprintf("the location ID %s has been deleted", locationID),
		})
	}

	return c.Status(fiber.StatusNotFound).JSON(dto.DefaultResponseMessageOut{
		Message: fmt.Sprintf("the location ID %s does not exist", locationID),
	})
}
