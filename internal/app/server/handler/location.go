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

// LocationsAddOne godoc
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
func LocationsAddOne(c *fiber.Ctx) error {
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
			Message: "there is an error validating the location data provided",
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

// LocationsGetOne godoc
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
func LocationsGetOne(c *fiber.Ctx) error {
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

// LocationsGetAll godoc
//
//	@Summary		Get all locations data
//	@Description	Get all locations data from database based on query parameters
//	@Tags			Locations
//	@Produce		json
//	@Param			q	query		dto.QueryLocationRequest		false	"Query parameters for filtering locations"
//	@Success		200	{object}	dto.QueryLocationResponse		"located documents"
//	@Failure		400	{object}	GlobalErrorHandlerResp			"validation error"
//	@Failure		404	{object}	dto.DefaultResponseMessageOut	"no locations found"
//	@Failure		500	{object}	GlobalErrorHandlerResp			"internal server error"
//	@Router			/api/v1/locations [get]
func LocationsGetAll(c *fiber.Ctx) error {
	queryParams := new(dto.QueryLocationRequest)

	if err := c.QueryParser(queryParams); err != nil {
		slog.Error("error parsing query parameters", "error", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(GlobalErrorHandlerResp{
			Success: false,
			Message: "there is an error processing the location data provided",
			Error:   err.Error(),
		})
	}

	if err := makeValidation(queryParams); err != nil {
		slog.Error("error validating query parameters", "error", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(GlobalErrorHandlerResp{
			Success: false,
			Message: "query parameters are not valid",
			Error:   err.Error(),
		})
	}

	locationsDataOut, err := usecase.GetAllLocations(queryParams)
	if err != nil {
		slog.Error("error getting all locations", "error", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(GlobalErrorHandlerResp{
			Success: false,
			Message: "there is an error getting all location data",
			Error:   err.Error(),
		})
	}

	if len(locationsDataOut.Data) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(dto.DefaultResponseMessageOut{
			Message: "no locations found",
		})
	}

	return c.JSON(locationsDataOut)
}

// LocationsUpdateOne godoc
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
func LocationsUpdateOne(c *fiber.Ctx) error {
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
			Message: "there is an error validating the location data provided",
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

// LocationsDeleteOne godoc
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
func LocationsDeleteOne(c *fiber.Ctx) error {
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
