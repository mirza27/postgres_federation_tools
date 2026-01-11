package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"db_migrate_server/internal/mapping"
)

type GetMappingNameParamsRequest struct {
	Name string `form:"name" binding:"required"`
}
type GetMappingListResponse struct {
	Files    []string         `json:"files"`
	Entities []mapping.Entity `json:"entities"`
	Message  string           `json:"message"`
	Status   string           `json:"status"`
}

// list all mapping files
func (server *Server) GetMappingList(c *gin.Context) {

	path := server.Config.MappingPath
	listJSONFiles, err := mapping.ListJSONFiles(path)
	if err != nil {
		c.JSON(400, ErrorResponse{
			Error:   err.Error(),
			Status:  "error",
			Message: "Failed to list mapping files",
		})
		return
	}

	Entities, err := mapping.GetEntitiesContentFromMappingList(listJSONFiles)
	if err != nil {
		c.JSON(400, ErrorResponse{
			Error:   err.Error(),
			Status:  "error",
			Message: "Failed to get entities content from mapping list",
		})
		return
	}

	c.JSON(http.StatusOK, GetMappingListResponse{
		Files:    listJSONFiles,
		Entities: Entities,
		Status:   "Success",
		Message:  "Successsfully retrieved mapping files",
	})
}

type GetMappingResponse struct {
	File    string         `json:"file"`
	Entity  mapping.Entity `json:"entity"`
	Message string         `json:"message"`
	Status  string         `json:"status"`
}

// get entity mapping by name (filename)
func (server *Server) GetMappingByName(c *gin.Context) {
	var req GetMappingNameParamsRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   err.Error(),
			Status:  "error",
			Message: "Invalid request parameters",
		})
		return
	}

	path := server.Config.MappingPath
	Entities, err := mapping.GetEntitiesContentFromMappingList([]string{path + "/" + req.Name + ".json"})
	if err != nil {
		c.JSON(400, ErrorResponse{
			Error:   err.Error(),
			Status:  "error",
			Message: "Failed to get entity content from " + req.Name + ".json",
		})
		return
	}

	c.JSON(http.StatusOK, GetMappingResponse{
		File:    req.Name + ".json",
		Entity:  Entities[0],
		Status:  "Success",
		Message: "Successfully retrieved mapping file",
	})
}

type CreateMappingRequest struct {
	Mapping mapping.Entity `json:"mapping" binding:"required"`
}

func (server *Server) CreateMappingByName(c *gin.Context) {
	var entity mapping.Entity

	if err := c.ShouldBindJSON(&entity); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   err.Error(),
			Status:  "error",
			Message: "Invalid request body",
		})
		return
	}

	path := server.Config.MappingPath

	// create mapping file
	err := mapping.CreateEntityMappingFile(path+"/"+entity.Entity+".json", entity)
	if err != nil {
		c.JSON(500, ErrorResponse{
			Error:   err.Error(),
			Status:  "error",
			Message: "Failed to create mapping file",
		})
		return
	}

	c.JSON(http.StatusOK, DefaultResponse{
		Status:  "Success",
		Message: "Mapping file created successfully",
	})
}

// update mapping by name by delete then create new mapping file
func (server *Server) UpdateMappingByName(c *gin.Context) {
	var req GetMappingNameParamsRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   err.Error(),
			Status:  "error",
			Message: "Invalid request parameters",
		})
		return
	}

	path := server.Config.MappingPath

	// delete existing mapping file
	err = mapping.DeleteEntityMappingFile(path + "/" + req.Name + ".json")
	if err != nil {
		c.JSON(500, ErrorResponse{
			Error:   err.Error(),
			Status:  "error",
			Message: "Failed to delete existing mapping file",
		})
		return
	}

	// get entity from request body
	var entity mapping.Entity
	if err := c.ShouldBindJSON(&entity); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   err.Error(),
			Status:  "error",
			Message: "Invalid request body",
		})
		return
	}

	// create new mapping file from request body
	err = mapping.CreateEntityMappingFile(path+"/"+entity.Entity+".json", entity)
	if err != nil {
		c.JSON(500, ErrorResponse{
			Error:   err.Error(),
			Status:  "error",
			Message: "Failed to create mapping file",
		})
		return
	}

	c.JSON(http.StatusOK, DefaultResponse{
		Status:  "Success",
		Message: "Mapping file updated (recreate) successfully",
	})
}

// delete mapping by name (filename)
func (server *Server) DeleteMappingByName(c *gin.Context) {
	var req GetMappingNameParamsRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   err.Error(),
			Status:  "error",
			Message: "Invalid request parameters",
		})
		return
	}

	path := server.Config.MappingPath
	err = mapping.DeleteEntityMappingFile(path + "/" + req.Name + ".json")
	if err != nil {
		c.JSON(500, ErrorResponse{
			Error:   err.Error(),
			Status:  "error",
			Message: "Failed to delete mapping file",
		})
		return
	}

	c.JSON(http.StatusOK, DefaultResponse{
		Status:  "Success",
		Message: "Mapping file deleted successfully",
	})

}
