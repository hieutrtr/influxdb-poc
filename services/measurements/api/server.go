package measurementapi

import (
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	measurementdb "github.com/hieutrtr/influxdb-poc/services/measurements/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Server struct {
	engine *gin.Engine
	store  measurementdb.Store
}

func NewServer(store measurementdb.Store, creds *Credentials) *Server {
	s := &Server{
		engine: gin.Default(),
		store:  store,
	}
	s.setupRoutes()
	s.setupAuth(creds)
	return s
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s *Server) setupAuth(creds *Credentials) {
	s.engine.Use(gin.BasicAuth(gin.Accounts{
		creds.Username: creds.Password,
	}))
}

func (s *Server) setupRoutes() {
	s.engine.GET("/measurements", s.handleListMeasurements)
	s.engine.POST("/measurements", s.handleCreateMeasurement)
	s.engine.GET("/measurements/:id", s.handleGetMeasurement)
	s.engine.DELETE("/measurements/:id", s.handleArchiveMeasurement)
}

func (s *Server) Run(addr string) error {
	return s.engine.Run(addr)
}

func (s *Server) handleListMeasurements(c *gin.Context) {
	limit, err := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid limit"})
		return
	}
	offset, err := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid offset"})
		return
	}
	measurements, err := s.store.ListMeasurements(c, limit, offset)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, measurements)
}

type createMeasurementRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

func (s *Server) handleCreateMeasurement(c *gin.Context) {
	var req createMeasurementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	measurement := measurementdb.Measurement{
		Name:        req.Name,
		Description: req.Description,
	}
	id, err := s.store.CreateMeasurement(c, measurement)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, id)
}

func (s *Server) handleGetMeasurement(c *gin.Context) {
	id := c.Param("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	measurement, err := s.store.GetMeasurement(c, measurementdb.MeasurementID{ID: oid})
	if err != nil {
		if reflect.TypeOf(err) == reflect.TypeOf(&measurementdb.ErrMeasurementNotFound{}) {
			c.JSON(404, gin.H{"error": err.Error()})
			return
		}
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, measurement)
}

func (s *Server) handleArchiveMeasurement(c *gin.Context) {
	id := c.Param("id")
	err := s.store.ArchiveMeasurement(c, id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, nil)
}
