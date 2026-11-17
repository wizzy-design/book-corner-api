package models

type CreateTagRequest struct {
	Name string `json:"name" binding:"required"`
}

type ReplaceBookTagsRequest struct {
	TagIDs []int `json:"tag_ids" binding:"required"`
}