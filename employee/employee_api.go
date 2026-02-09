package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"employee/config"
	"employee/elastic"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.elastic.co/apm/module/apmgin/v2"
)

var configFile string

type EmployeeInfo struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	JobRole       string  `json:"job_role"`
	JoiningDate   string  `json:"joining_date"`
	Addresss      string  `json:"address"`
	Location      string  `json:"location"`
	Status        string  `json:"status"`
	EmailID       string  `json:"email"`
	AnnualPackage float64 `json:"annual_package"`
	PhoneNumber   string  `json:"phone_number"`
}

func init() {
	configFile = os.Getenv("CONFIG_FILE")
	if configFile == "" {
		configFile = "./config.yaml"
		logrus.Warn("CONFIG_FILE not set, using default ./config.yaml")
	}
}

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	conf, err := config.ParseFile(configFile)
	if err != nil {
		logrus.Fatalf("Unable to parse config (%s): %v", configFile, err)
	}

	router := gin.Default()
	router.Use(apmgin.Middleware(router))

	corsCfg := cors.DefaultConfig()
	corsCfg.AllowOrigins = []string{"*"}
	corsCfg.AllowMethods = []string{"*"}
	corsCfg.AllowHeaders = []string{"*"}
	router.Use(cors.New(corsCfg))

	router.POST("/employee/create", pushEmployeeData)
	router.GET("/employee/search", fetchEmployeeData)
	router.GET("/employee/search/all", fetchALLEmployeeData)
	router.GET("/employee/search/roles", fetchEmployeeRoles)
	router.GET("/employee/search/location", fetchEmployeeLocation)
	router.GET("/employee/search/status", fetchEmployeeStatus)
	router.GET("/employee/healthz", healthCheck)

	router.Run(":" + conf.Employee.APIPort)
}

func pushEmployeeData(c *gin.Context) {
	var request EmployeeInfo
	if err := c.BindJSON(&request); err != nil {
		errorResponse(c, http.StatusBadRequest, "Malformed request body")
		return
	}

	conf, err := config.ParseFile(configFile)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "Config error")
		return
	}

	// ✅ FIX: Do NOT expect error return
	elastic.PostDataInSearch(conf, request.ID, request, c.Request.Context())

	c.JSON(http.StatusOK, gin.H{"message": "Employee created"})
}

func fetchEmployeeData(c *gin.Context) {
	conf, err := config.ParseFile(configFile)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "Config error")
		return
	}

	var searchValue string
	for _, v := range c.Request.URL.Query() {
		searchValue = strings.Join(v, "")
	}

	data := elastic.SearchDataInElastic(conf, searchValue, c.Request.Context())
	if data == nil || data["hits"] == nil {
		errorResponse(c, http.StatusInternalServerError, "Elasticsearch error")
		return
	}

	hits := data["hits"].(map[string]interface{})["hits"].([]interface{})
	if len(hits) == 0 {
		c.JSON(http.StatusOK, gin.H{})
		return
	}

	empData, _ := json.Marshal(hits[0].(map[string]interface{})["_source"])
	response := &EmployeeInfo{}
	_ = json.Unmarshal(empData, response)

	c.JSON(http.StatusOK, response)
}

func fetchALLEmployeeData(c *gin.Context) {
	conf, err := config.ParseFile(configFile)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "Config error")
		return
	}

	data := elastic.SearchALLDataInElastic(conf, c.Request.Context())
	if data == nil || data["hits"] == nil {
		errorResponse(c, http.StatusInternalServerError, "Elasticsearch error")
		return
	}

	hits := data["hits"].(map[string]interface{})["hits"].([]interface{})
	var employees []EmployeeInfo

	for _, h := range hits {
		emp := &EmployeeInfo{}
		b, _ := json.Marshal(h.(map[string]interface{})["_source"])
		_ = json.Unmarshal(b, emp)
		employees = append(employees, *emp)
	}

	c.JSON(http.StatusOK, employees)
}

func fetchEmployeeRoles(c *gin.Context) {
	employees := fetchAllEmployeesInternal(c)
	result := map[string]int{}
	for _, e := range employees {
		result[e.JobRole]++
	}
	c.JSON(http.StatusOK, result)
}

func fetchEmployeeLocation(c *gin.Context) {
	employees := fetchAllEmployeesInternal(c)
	result := map[string]int{}
	for _, e := range employees {
		result[e.Location]++
	}
	c.JSON(http.StatusOK, result)
}

func fetchEmployeeStatus(c *gin.Context) {
	employees := fetchAllEmployeesInternal(c)
	result := map[string]int{}
	for _, e := range employees {
		result[e.Status]++
	}
	c.JSON(http.StatusOK, result)
}

func healthCheck(c *gin.Context) {
	conf, err := config.ParseFile(configFile)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "Config error")
		return
	}

	status, err := elastic.CheckElasticHealth(conf, c.Request.Context())
	if err != nil || !status {
		errorResponse(c, http.StatusBadRequest, "Elasticsearch is not running")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "up",
		"database": "elasticsearch",
	})
}

func fetchAllEmployeesInternal(c *gin.Context) []EmployeeInfo {
	conf, err := config.ParseFile(configFile)
	if err != nil {
		return nil
	}

	data := elastic.SearchALLDataInElastic(conf, c.Request.Context())
	if data == nil || data["hits"] == nil {
		return nil
	}

	hits := data["hits"].(map[string]interface{})["hits"].([]interface{})
	var employees []EmployeeInfo

	for _, h := range hits {
		emp := &EmployeeInfo{}
		b, _ := json.Marshal(h.(map[string]interface{})["_source"])
		_ = json.Unmarshal(b, emp)
		employees = append(employees, *emp)
	}

	return employees
}

func errorResponse(c *gin.Context, code int, msg string) {
	c.JSON(code, gin.H{"error": msg})
}
