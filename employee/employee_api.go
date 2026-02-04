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

var (
	configFile = os.Getenv("CONFIG_FILE")
)

// EmployeeInfo struct will be the data structure for employee's information
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

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	conf, err := config.ParseFile(configFile)
	if err != nil {
		logrus.Fatalf("Unable to parse configuration file: %v", err)
	}

	logrus.Infof("Running employee-management in webserver mode")
	logrus.Infof("employee-management is listening on port: %v", conf.Employee.APIPort)

	router := gin.Default()
	router.Use(apmgin.Middleware(router))

	corsCfg := cors.DefaultConfig()
	corsCfg.AllowOrigins = []string{"*"}
	corsCfg.AllowMethods = []string{"*"}
	corsCfg.AllowHeaders = []string{"*"}
	corsCfg.AllowCredentials = true
	router.Use(cors.New(corsCfg))

	router.POST("/employee/create", pushEmployeeData)
	router.GET("/employee/search", fetchEmployeeData)
	router.GET("/employee/search/all", fetchALLEmployeeData)
	router.GET("/employee/search/roles", fetchEmployeeRoles)
	router.GET("/employee/search/location", fetchEmployeeLocation)
	router.GET("/employee/search/status", fetchEmployeeStatus)
	router.GET("/employee/healthz", healthCheck)

	if err := router.Run(":" + conf.Employee.APIPort); err != nil {
		logrus.Fatalf("failed to start server: %v", err)
	}
}

func pushEmployeeData(c *gin.Context) {
	var request EmployeeInfo
	if err := c.BindJSON(&request); err != nil {
		errorResponse(c, http.StatusBadRequest, "Malformed request body")
		logrus.Errorf("Error parsing request body: %v", err)
		return
	}

	conf, err := config.ParseFile(configFile)
	if err != nil {
		logrus.Errorf("Unable to parse configuration file: %v", err)
		return
	}

	elastic.PostDataInSearch(conf, request.ID, request, c.Request.Context())
	logrus.Infof("Successfully pushed employee data to elasticsearch")
}

func fetchEmployeeData(c *gin.Context) {
	searchQuery := c.Request.URL.Query()
	var searchValue string
	response := &EmployeeInfo{}

	for _, value := range searchQuery {
		searchValue = strings.Join(value, "")
	}

	conf, err := config.ParseFile(configFile)
	if err != nil {
		logrus.Errorf("Unable to parse configuration file: %v", err)
		return
	}

	data := elastic.SearchDataInElastic(conf, searchValue, c.Request.Context())

	for _, parsedData := range data["hits"].(map[string]interface{})["hits"].([]interface{}) {
		empData, err := json.Marshal(parsedData.(map[string]interface{})["_source"])
		if err != nil {
			logrus.Errorf("Unable to marshal response JSON: %v", err)
			continue
		}

		if err := json.Unmarshal(empData, response); err != nil {
			logrus.Errorf("Unable to unmarshal response JSON: %v", err)
			continue
		}
	}

	c.JSON(http.StatusOK, response)
}

func fetchALLEmployeeData(c *gin.Context) {
	conf, err := config.ParseFile(configFile)
	if err != nil {
		logrus.Errorf("Unable to parse configuration file: %v", err)
		return
	}

	data := elastic.SearchALLDataInElastic(conf, c.Request.Context())
	var employeeInfo []EmployeeInfo

	for _, parsedData := range data["hits"].(map[string]interface{})["hits"].([]interface{}) {
		response := &EmployeeInfo{}

		empData, err := json.Marshal(parsedData.(map[string]interface{})["_source"])
		if err != nil {
			logrus.Errorf("Unable to marshal response JSON: %v", err)
			continue
		}

		if err := json.Unmarshal(empData, response); err != nil {
			logrus.Errorf("Unable to unmarshal response JSON: %v", err)
			continue
		}

		employeeInfo = append(employeeInfo, *response)
	}

	c.JSON(http.StatusOK, employeeInfo)
}

func fetchEmployeeRoles(c *gin.Context) {
	employeeInfo := fetchAllEmployeesInternal(c)
	result := make(map[string]int)

	for _, emp := range employeeInfo {
		result[emp.JobRole]++
	}

	c.JSON(http.StatusOK, result)
}

func fetchEmployeeLocation(c *gin.Context) {
	employeeInfo := fetchAllEmployeesInternal(c)
	result := make(map[string]int)

	for _, emp := range employeeInfo {
		result[emp.Location]++
	}

	c.JSON(http.StatusOK, result)
}

func fetchEmployeeStatus(c *gin.Context) {
	employeeInfo := fetchAllEmployeesInternal(c)
	result := make(map[string]int)

	for _, emp := range employeeInfo {
		result[emp.Status]++
	}

	c.JSON(http.StatusOK, result)
}

func healthCheck(c *gin.Context) {
	conf, err := config.ParseFile(configFile)
	if err != nil {
		logrus.Errorf("Unable to parse configuration file: %v", err)
		return
	}

	status, err := elastic.CheckElasticHealth(conf, c.Request.Context())
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "Elasticsearch is not running")
		return
	}

	if status {
		c.JSON(http.StatusOK, gin.H{
			"status":   "up",
			"database": "elasticsearch",
			"message":  "Elasticsearch is running",
		})
		return
	}

	errorResponse(c, http.StatusBadRequest, "Elasticsearch is not running")
}

func fetchAllEmployeesInternal(c *gin.Context) []EmployeeInfo {
	conf, err := config.ParseFile(configFile)
	if err != nil {
		logrus.Errorf("Unable to parse configuration file: %v", err)
		return nil
	}

	data := elastic.SearchALLDataInElastic(conf, c.Request.Context())
	var employees []EmployeeInfo

	for _, parsedData := range data["hits"].(map[string]interface{})["hits"].([]interface{}) {
		emp := &EmployeeInfo{}

		empData, err := json.Marshal(parsedData.(map[string]interface{})["_source"])
		if err != nil {
			continue
		}

		if err := json.Unmarshal(empData, emp); err != nil {
			continue
		}

		employees = append(employees, *emp)
	}

	return employees
}

func errorResponse(c *gin.Context, code int, errMsg string) {
	c.JSON(code, gin.H{
		"error": errMsg,
	})
}

