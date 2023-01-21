package measurementapi

import (
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	measurementdb "github.com/hieutrtr/influxdb-poc/services/measurements/db"
	"go.mongodb.org/mongo-driver/bson/primitive"

	docs "github.com/hieutrtr/influxdb-poc/services/measurements/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	s.setupAuth(creds)
	s.setupSwagger()
	s.setupRoutes()
	return s
}

func (s *Server) setupSwagger() {
	docs.SwaggerInfo.BasePath = "/"
	s.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

type Credentials struct {
	// Username for basic auth
	Username string `json:"username"`
	// Password for basic auth
	Password string `json:"password"`
} // @name Credentials

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

type GetMeasurementResponse struct {
	// ID of the measurement
	ID string `json:"id" example:"5f1f2c5b5b9b9b0b5c1c1c1c" bson:"_id"`
	// Name of the measurement
	Name string `json:"name" example:"temperature" bson:"name"`
	// Description of the measurement
	Description string `json:"description" example:"Temperature measurement in Celsius" bson:"description"`
	// Status of the measurement
	Status string `json:"status" example:"active" bson:"status"`
} // @name GetMeasurementResponse

type listMeasurementsResponse []GetMeasurementResponse // @name ListMeasurementsResponse

// @Summary List all measurements in the database
// @Description List all measurements in the database with pagination support
// @Tags measurements
// @Accept  json
// @Produce  json
// @Param limit query int false "Limit the number of measurements returned" default(10)
// @Param offset query int false "Offset the number of measurements returned" default(0)
// @Success 200 {object} ListMeasurementsResponse
// @Failure 500 {object} ErrorResponse
// @Router /measurements [get]
func (s *Server) handleListMeasurements(c *gin.Context) {
	limit, err := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
	if err != nil {
		c.JSON(400, errorResponse{Error: "invalid limit"})
		return
	}
	offset, err := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64)
	if err != nil {
		c.JSON(400, errorResponse{Error: "invalid offset"})
		return
	}
	measurements, err := s.store.ListMeasurements(c, limit, offset)
	if err != nil {
		c.JSON(500, errorResponse{Error: err.Error()})
		return
	}
	c.JSON(200, measurements)
}

type errorResponse struct {
	// Error message
	Error string `json:"error" example:"error message"`
} // @name ErrorResponse

type createMeasurementRequest struct {
	// Name of the measurement
	Name string `json:"name" binding:"required" example:"temperature"`
	// Description of the measurement
	Description string `json:"description" binding:"required" example:"Temperature measurement in Celsius"`
} // @name CreateMeasurementRequest

type createMeasurementResponse struct {
	// ID of the measurement
	ID string `json:"id" example:"5f1f2c5b5b9b9b0b5c1c1c1c"`
} // @name CreateMeasurementResponse

// @Summary Create a new measurement
// @Description Create a new measurement in the database with the given name and description
// @Tags measurements
// @Accept  json
// @Produce  json
// @Param body body CreateMeasurementRequest true "Create measurement request"
// @Success 200 {object} CreateMeasurementResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /measurements [post]
func (s *Server) handleCreateMeasurement(c *gin.Context) {
	var req createMeasurementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, errorResponse{Error: err.Error()})
		return
	}
	measurement := measurementdb.Measurement{
		Name:        req.Name,
		Description: req.Description,
	}
	id, err := s.store.CreateMeasurement(c, measurement)
	if err != nil {
		c.JSON(500, errorResponse{Error: err.Error()})
		return
	}
	c.JSON(200, id)
}

// @Summary Get a measurement by ID
// @Description Get a measurement by ID from the database
// @Tags measurements
// @Accept  json
// @Produce  json
// @Param id path string true "Measurement ID"
// @Success 200 {object} GetMeasurementResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /measurements/{id} [get]
func (s *Server) handleGetMeasurement(c *gin.Context) {
	id := c.Param("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(400, errorResponse{Error: err.Error()})
		return
	}
	measurement, err := s.store.GetMeasurement(c, measurementdb.MeasurementID{ID: oid})
	if err != nil {
		if reflect.TypeOf(err) == reflect.TypeOf(&measurementdb.ErrMeasurementNotFound{}) {
			c.JSON(404, errorResponse{Error: err.Error()})
			return
		}
		c.JSON(500, errorResponse{Error: err.Error()})
		return
	}
	c.JSON(200, measurement)
}

type archivedMeasurementResponse struct {
	// ID of the measurement
	ID string `json:"id" example:"5f1f2c5b5b9b9b0b5c1c1c1c" bson:"_id"`
	// Message of the archived measurement
	Message string `json:"message" example:"Archived measurement"`
} // @name ArchivedMeasurementResponse

// @Summary Archive a measurement by ID
// @Description Archive a measurement by ID from the database
// @Tags measurements
// @Accept  json
// @Produce  json
// @Param id path string true "Measurement ID"
// @Success 200 {object} ArchivedMeasurementResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /measurements/{id} [delete]
func (s *Server) handleArchiveMeasurement(c *gin.Context) {
	id := c.Param("id")
	err := s.store.ArchiveMeasurement(c, id)
	if err != nil {
		c.JSON(500, errorResponse{Error: err.Error()})
		return
	}
	c.JSON(200, archivedMeasurementResponse{ID: id, Message: "Archived measurement"})
}
