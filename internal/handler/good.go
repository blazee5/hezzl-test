package handler

import (
	"errors"
	"github.com/blazee5/hezzl-test/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"net/http"
	"strconv"
)

func (s *Server) CreateGood(c *gin.Context) {
	var input domain.CreateGoodRequest

	projectID, err := strconv.Atoi(c.Query("projectId"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid project id",
		})
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "bad request",
		})
		return
	}

	good, err := s.services.Good.CreateGood(c.Request.Context(), projectID, input)

	if err != nil {
		s.log.Infof("error while create good: %v", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "server error",
		})
		return
	}

	c.JSON(http.StatusOK, good)
}

func (s *Server) UpdateGood(c *gin.Context) {
	var input domain.UpdateGoodRequest

	projectID, err := strconv.Atoi(c.Query("projectId"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid project id",
		})
		return
	}

	id, err := strconv.Atoi(c.Query("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid id",
		})
		return
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "bad request",
		})
		return
	}

	good, err := s.services.Good.UpdateGood(c.Request.Context(), projectID, id, input)

	if errors.Is(err, pgx.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  3,
			"message": "errors.good.notFound",
			"details": gin.H{},
		})
		return
	}

	if err != nil {
		s.log.Infof("error while update good: %v", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "server error",
		})
		return
	}

	c.JSON(http.StatusOK, good)
}

func (s *Server) DeleteGood(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Query("projectId"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid project id",
		})
		return
	}

	id, err := strconv.Atoi(c.Query("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid id",
		})
		return
	}

	good, err := s.services.Good.DeleteGood(c.Request.Context(), projectID, id)

	if errors.Is(err, pgx.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  3,
			"message": "errors.good.notFound",
			"details": gin.H{},
		})
		return
	}

	if err != nil {
		s.log.Infof("error while update good: %v", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "server error",
		})
		return
	}

	c.JSON(http.StatusOK, good)
}

func (s *Server) GetGoods(c *gin.Context) {
	limit, err := strconv.Atoi(c.Query("limit"))

	if err != nil {
		limit = 10
	}

	offset, err := strconv.Atoi(c.Query("offset"))

	if err != nil {
		offset = 0
	}

	goodsList, err := s.services.Good.GetGoods(c.Request.Context(), limit, offset)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "server error",
		})
		return
	}

	c.JSON(http.StatusOK, goodsList)
}

func (s *Server) GetGoodByID(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Query("projectId"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid project id",
		})
		return
	}

	id, err := strconv.Atoi(c.Query("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid id",
		})
		return
	}

	good, err := s.services.Good.GetGoodByID(c.Request.Context(), projectID, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "server error",
		})
		return
	}

	c.JSON(http.StatusOK, good)
}

func (s *Server) ReprioritizeGood(c *gin.Context) {
	var input domain.ReprioritizeRequest

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "bad request",
		})
		return
	}

	projectID, err := strconv.Atoi(c.Query("projectId"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid project id",
		})
		return
	}

	id, err := strconv.Atoi(c.Query("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid id",
		})
		return
	}

	good, err := s.services.Good.ReprioritizeGood(c.Request.Context(), projectID, id, input)

	if errors.Is(err, pgx.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  3,
			"message": "errors.good.notFound",
			"details": gin.H{},
		})
		return
	}

	if err != nil {
		s.log.Infof("error while update good: %v", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "server error",
		})
		return
	}

	c.JSON(http.StatusOK, good)
}
