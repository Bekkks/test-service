package handlers

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"subscription-service/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriptionHandler struct {
	db *gorm.DB
}

func NewSubscriptionHandler(db *gorm.DB) *SubscriptionHandler {
	return &SubscriptionHandler{db: db}
}


// parseMonthYear парсит строку формата "MM-YYYY" в time.Time
func parseMonthYear(dateStr string) (time.Time, error) {
	parts := strings.Split(dateStr, "-")
	if len(parts) != 2 {
		return time.Time{}, gin.Error{Err: nil, Type: gin.ErrorTypeBind, Meta: "invalid date format, expected MM-YYYY"}
	}

	month, err := strconv.Atoi(parts[0])
	if err != nil || month < 1 || month > 12 {
		return time.Time{}, gin.Error{Err: nil, Type: gin.ErrorTypeBind, Meta: "invalid month"}
	}

	year, err := strconv.Atoi(parts[1])
	if err != nil {
		return time.Time{}, gin.Error{Err: nil, Type: gin.ErrorTypeBind, Meta: "invalid year"}
	}

	return time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC), nil
}

// CreateSubscription создает новую подписку
// @Summary Создать подписку
// @Description Создает новую запись о подписке
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body CreateSubscriptionRequest true "Данные подписки"
// @Success 201 {object} models.Subscription
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions [post]
func (h *SubscriptionHandler) CreateSubscription(c *gin.Context) {
	var req CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		log.Printf("Error parsing user_id: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id format"})
		return
	}

	startDate, err := parseMonthYear(req.StartDate)
	if err != nil {
		log.Printf("Error parsing start_date: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format, expected MM-YYYY"})
		return
	}

	subscription := models.Subscription{
		ServiceName: req.ServiceName,
		Price:        req.Price,
		UserID:       userID,
		StartDate:    startDate,
	}

	if req.EndDate != "" {
		endDate, err := parseMonthYear(req.EndDate)
		if err != nil {
			log.Printf("Error parsing end_date: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format, expected MM-YYYY"})
			return
		}
		subscription.EndDate = &endDate
	}

	if err := h.db.Create(&subscription).Error; err != nil {
		log.Printf("Error creating subscription: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create subscription"})
		return
	}

	log.Printf("Created subscription with ID: %s", subscription.ID)
	c.JSON(http.StatusCreated, subscription)
}

// GetSubscription получает подписку по ID
// @Summary Получить подписку
// @Description Возвращает подписку по её ID
// @Tags subscriptions
// @Produce json
// @Param id path string true "ID подписки"
// @Success 200 {object} models.Subscription
// @Failure 404 {object} map[string]string
// @Router /subscriptions/{id} [get]
func (h *SubscriptionHandler) GetSubscription(c *gin.Context) {
	id := c.Param("id")
	subscriptionID, err := uuid.Parse(id)
	if err != nil {
		log.Printf("Error parsing subscription ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription ID format"})
		return
	}

	var subscription models.Subscription
	if err := h.db.Where("id = ?", subscriptionID).First(&subscription).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("Subscription not found: %s", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
			return
		}
		log.Printf("Error getting subscription: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get subscription"})
		return
	}

	log.Printf("Retrieved subscription: %s", id)
	c.JSON(http.StatusOK, subscription)
}

// UpdateSubscription обновляет подписку
// @Summary Обновить подписку
// @Description Обновляет существующую подписку
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "ID подписки"
// @Param subscription body UpdateSubscriptionRequest true "Данные для обновления"
// @Success 200 {object} models.Subscription
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /subscriptions/{id} [put]
func (h *SubscriptionHandler) UpdateSubscription(c *gin.Context) {
	id := c.Param("id")
	subscriptionID, err := uuid.Parse(id)
	if err != nil {
		log.Printf("Error parsing subscription ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription ID format"})
		return
	}

	var req UpdateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var subscription models.Subscription
	if err := h.db.Where("id = ?", subscriptionID).First(&subscription).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("Subscription not found: %s", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
			return
		}
		log.Printf("Error getting subscription: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get subscription"})
		return
	}

	// Обновление полей
	if req.ServiceName != "" {
		subscription.ServiceName = req.ServiceName
	}
	if req.Price != nil {
		subscription.Price = *req.Price
	}
	if req.UserID != "" {
		userID, err := uuid.Parse(req.UserID)
		if err != nil {
			log.Printf("Error parsing user_id: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id format"})
			return
		}
		subscription.UserID = userID
	}
	if req.StartDate != "" {
		startDate, err := parseMonthYear(req.StartDate)
		if err != nil {
			log.Printf("Error parsing start_date: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format, expected MM-YYYY"})
			return
		}
		subscription.StartDate = startDate
	}
	if req.EndDate != "" {
		endDate, err := parseMonthYear(req.EndDate)
		if err != nil {
			log.Printf("Error parsing end_date: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format, expected MM-YYYY"})
			return
		}
		subscription.EndDate = &endDate
	} else if req.EndDate == "" && subscription.EndDate != nil {
		// Если передана пустая строка, удаляем end_date
		subscription.EndDate = nil
	}

	if err := h.db.Save(&subscription).Error; err != nil {
		log.Printf("Error updating subscription: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update subscription"})
		return
	}

	log.Printf("Updated subscription: %s", id)
	c.JSON(http.StatusOK, subscription)
}

// DeleteSubscription удаляет подписку
// @Summary Удалить подписку
// @Description Удаляет подписку по её ID
// @Tags subscriptions
// @Param id path string true "ID подписки"
// @Success 204 "No Content"
// @Failure 404 {object} map[string]string
// @Router /subscriptions/{id} [delete]
func (h *SubscriptionHandler) DeleteSubscription(c *gin.Context) {
	id := c.Param("id")
	subscriptionID, err := uuid.Parse(id)
	if err != nil {
		log.Printf("Error parsing subscription ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription ID format"})
		return
	}

	var subscription models.Subscription
	if err := h.db.Where("id = ?", subscriptionID).First(&subscription).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("Subscription not found: %s", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
			return
		}
		log.Printf("Error getting subscription: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get subscription"})
		return
	}

	if err := h.db.Delete(&subscription).Error; err != nil {
		log.Printf("Error deleting subscription: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete subscription"})
		return
	}

	log.Printf("Deleted subscription: %s", id)
	c.Status(http.StatusNoContent)
}

// ListSubscriptions возвращает список подписок
// @Summary Список подписок
// @Description Возвращает список всех подписок с пагинацией
// @Tags subscriptions
// @Produce json
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Количество записей на странице" default(10)
// @Success 200 {object} map[string]interface{}
// @Router /subscriptions [get]
func (h *SubscriptionHandler) ListSubscriptions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var subscriptions []models.Subscription
	var total int64

	if err := h.db.Model(&models.Subscription{}).Count(&total).Error; err != nil {
		log.Printf("Error counting subscriptions: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count subscriptions"})
		return
	}

	if err := h.db.Offset(offset).Limit(limit).Find(&subscriptions).Error; err != nil {
		log.Printf("Error listing subscriptions: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list subscriptions"})
		return
	}

	log.Printf("Listed subscriptions: page=%d, limit=%d, total=%d", page, limit, total)
	c.JSON(http.StatusOK, gin.H{
		"data": subscriptions,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// CalculateTotalCost рассчитывает суммарную стоимость подписок
// @Summary Рассчитать стоимость подписок
// @Description Рассчитывает суммарную стоимость всех подписок за выбранный период с фильтрацией
// @Tags subscriptions
// @Produce json
// @Param start_date query string false "Начало периода (MM-YYYY)"
// @Param end_date query string false "Конец периода (MM-YYYY)"
// @Param user_id query string false "ID пользователя (UUID)"
// @Param service_name query string false "Название сервиса"
// @Success 200 {object} map[string]interface{}
// @Router /subscriptions/total-cost [get]
func (h *SubscriptionHandler) CalculateTotalCost(c *gin.Context) {
	var req TotalCostRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		log.Printf("Error binding query: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := h.db.Model(&models.Subscription{})

	// Фильтр по user_id
	if req.UserID != "" {
		userID, err := uuid.Parse(req.UserID)
		if err != nil {
			log.Printf("Error parsing user_id: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id format"})
			return
		}
		query = query.Where("user_id = ?", userID)
	}

	// Фильтр по service_name
	if req.ServiceName != "" {
		query = query.Where("service_name = ?", req.ServiceName)
	}

	// Фильтр по периоду
	if req.StartDate != "" && req.EndDate != "" {
		// Если указаны оба периода, используем диапазон
		startDate, err := parseMonthYear(req.StartDate)
		if err != nil {
			log.Printf("Error parsing start_date: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format, expected MM-YYYY"})
			return
		}
		endDate, err := parseMonthYear(req.EndDate)
		if err != nil {
			log.Printf("Error parsing end_date: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format, expected MM-YYYY"})
			return
		}

		startOfPeriod := time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, time.UTC)
		endOfPeriod := time.Date(endDate.Year(), endDate.Month(), 1, 0, 0, 0, 0, time.UTC)
		endOfPeriod = endOfPeriod.AddDate(0, 1, -1)

		// Подписка активна в периоде, если она началась до конца периода и не закончилась до начала периода
		query = query.Where("start_date <= ? AND (end_date IS NULL OR end_date >= ?)", endOfPeriod, startOfPeriod)
	} else if req.StartDate != "" {
		startDate, err := parseMonthYear(req.StartDate)
		if err != nil {
			log.Printf("Error parsing start_date: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format, expected MM-YYYY"})
			return
		}
		startOfMonth := time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, time.UTC)
		endOfMonth := startOfMonth.AddDate(0, 1, -1)

		query = query.Where("start_date <= ? AND (end_date IS NULL OR end_date >= ?)", endOfMonth, startOfMonth)
	} else if req.EndDate != "" {
		endDate, err := parseMonthYear(req.EndDate)
		if err != nil {
			log.Printf("Error parsing end_date: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format, expected MM-YYYY"})
			return
		}
		startOfMonth := time.Date(endDate.Year(), endDate.Month(), 1, 0, 0, 0, 0, time.UTC)
		endOfMonth := startOfMonth.AddDate(0, 1, -1)

		query = query.Where("start_date <= ? AND (end_date IS NULL OR end_date >= ?)", endOfMonth, startOfMonth)
	}

	var totalCost int64
	if err := query.Select("COALESCE(SUM(price), 0)").Scan(&totalCost).Error; err != nil {
		log.Printf("Error calculating total cost: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to calculate total cost"})
		return
	}

	log.Printf("Calculated total cost: %d", totalCost)
	c.JSON(http.StatusOK, gin.H{
		"total_cost": totalCost,
		"filters": gin.H{
			"start_date":   req.StartDate,
			"end_date":     req.EndDate,
			"user_id":      req.UserID,
			"service_name": req.ServiceName,
		},
	})
}
