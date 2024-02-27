package models

type CreatePostRequest struct {
	Title       string `json:"title" binding:"required"`
	UserOwnerID string `json:"user_owner_id" binding:"required"` // TODO
	Amount      string `json:"amount" binding:"required"`
	PostType    string `json:"post_type" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type CreatePostResult struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Amount      string `json:"amount"`
	PostType    string `json:"post_type"`
	Description string `json:"description"`
	Status      string `json:"status"`
}
