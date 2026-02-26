package handlers

type CreateSubscriptionRequest struct {
	ServiceName string `json:"service_name" binding:"required" example:"Yandex Plus"`
	Price       int    `json:"price" binding:"required,min=0" example:"400"`
	UserID      string `json:"user_id" binding:"required,uuid" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	StartDate   string `json:"start_date" binding:"required" example:"07-2025"`
	EndDate     string `json:"end_date,omitempty" example:"12-2025"`
}

type UpdateSubscriptionRequest struct {
	ServiceName string `json:"service_name,omitempty" example:"Yandex Plus"`
	Price       *int   `json:"price,omitempty" example:"400"`
	UserID      string `json:"user_id,omitempty" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	StartDate   string `json:"start_date,omitempty" example:"07-2025"`
	EndDate     string `json:"end_date,omitempty" example:"12-2025"`
}

type TotalCostRequest struct {
	StartDate   string `form:"start_date" example:"01-2025"`
	EndDate     string `form:"end_date" example:"12-2025"`
	UserID      string `form:"user_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	ServiceName string `form:"service_name" example:"Yandex Plus"`
}