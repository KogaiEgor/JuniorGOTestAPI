package main

import (
	"bytes"
	"encoding/json"
	"example/m/controllers"
	"example/m/initializers"
	"example/m/models"
	"example/m/services"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("TEST_DB_NAME"),
		os.Getenv("DB_PORT"),
	)
	fmt.Println("DSN:", dsn)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database")
	}
}

func InitializeTestDB(db *gorm.DB) error {
	if err := db.AutoMigrate(&models.Person{}); err != nil {
		return err
	}

	testPeople := []models.Person{
		{
			Name:        "Ivan",
			Surname:     "Ivanov",
			Age:         30,
			Gender:      "male",
			Nationality: "RU",
		},
		{
			Name:        "Ivan",
			Surname:     "Sergeev",
			Age:         30,
			Gender:      "male",
			Nationality: "RU",
		},
		{
			Name:        "Maria",
			Surname:     "Petrova",
			Age:         25,
			Gender:      "female",
			Nationality: "RU",
		},
		{
			Name:        "Lee",
			Surname:     "Cheng",
			Age:         66,
			Gender:      "female",
			Nationality: "KR",
		},
		{
			Name:        "John",
			Surname:     "Doe",
			Age:         73,
			Gender:      "male",
			Nationality: "IE",
		},
	}

	for _, todo := range testPeople {
		if result := db.Create(&todo); result.Error != nil {
			return result.Error
		}
	}

	return nil
}

func ClearTestDB(db *gorm.DB) {
	db.Exec("TRUNCATE TABLE people RESTART IDENTITY CASCADE")
}

func TestMain(m *testing.M) {
	if _, err := os.Stat("../.env"); err == nil {
		if err := godotenv.Load("../.env"); err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	ConnectToDB()

	if err := InitializeTestDB(DB); err != nil {
		log.Fatalf("Failed to initialize test todo: %s", err)
	}

	initializers.DB = DB

	code := m.Run()

	ClearTestDB(DB)

	os.Exit(code)
}

func TestPeopleGet(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/get", controllers.PeopleGet)

	// check name and surname as params
	req, _ := http.NewRequest("GET", "/get?name=Ivan&surname=Ivanov&page=1&pageSize=10", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected status code to be 200 OK")

	var response map[string][]models.Person
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	people, exists := response["people"]
	assert.True(t, exists)
	assert.NotEmpty(t, people)
	assert.Len(t, people, 1)
	// check get without params
	req, _ = http.NewRequest("GET", "/get", nil)
	w = httptest.NewRecorder()

	r.ServeHTTP(w, req)

	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err, "Should unmarshal response without error")

	people, exists = response["people"]
	assert.Equal(t, http.StatusOK, w.Code, "Expected status code to be 200 OK")
	assert.Len(t, people, 5)

	// check with name param
	req, _ = http.NewRequest("GET", "/get?name=Ivan", nil)
	w = httptest.NewRecorder()

	r.ServeHTTP(w, req)

	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err, "Should unmarshal response without error")

	people, exists = response["people"]
	assert.Equal(t, http.StatusOK, w.Code, "Expected status code to be 200 OK")
	assert.Len(t, people, 2)

	// check pagination
	req, _ = http.NewRequest("GET", "/get?page=2&pageSize=3", nil)
	w = httptest.NewRecorder()

	r.ServeHTTP(w, req)

	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err, "Should unmarshal response without error")

	people, exists = response["people"]
	assert.Equal(t, http.StatusOK, w.Code, "Expected status code to be 200 OK")
	assert.Len(t, people, 2)
}

func TestPersonCreate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/create", controllers.PeopleCreate)

	body := services.Body{
		Name:       "Dmitriy",
		Surname:    "TestSurname",
		Patronymic: "TestPatronymic",
	}
	bodyJSON, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/create", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code, "Expected status code to be 201 Created")

	var response map[string]*models.Person
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err, "Should unmarshal response without error")

	assert.Equal(t, "male", response["person"].Gender)
	assert.Equal(t, "UA", response["person"].Nationality)
	assert.Equal(t, 43, response["person"].Age)
}

func TestPersonCreateWrongData(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.POST("/todos", func(c *gin.Context) {
		controllers.PeopleCreate(c)
	})

	body := services.Body{
		Name: "TestName",
	}

	jsonData, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/todos", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeletePerson(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.DELETE("/delete/:id", controllers.PeopleDelete)

	testPerson := models.Person{
		Name:    "Oleg",
		Surname: "Olegov",
	}
	initializers.DB.Create(&testPerson)

	req, _ := http.NewRequest("DELETE", "/delete/"+fmt.Sprint(testPerson.ID), nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code, "Expected status code to be 204 No Content")

	var deletedPerson models.Person
	result := initializers.DB.First(&deletedPerson, testPerson.ID)
	assert.Error(t, result.Error, "Person should be deleted")
}

func TestPeopleUpdate(t *testing.T) {
	// Настройка тестового окружения
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.PUT("/update/:id", controllers.PeopleUpdate)

	// Создаем и сохраняем тестовую запись
	testPerson := models.Person{
		Name:    "Timur",
		Surname: "Orig",
	}
	initializers.DB.Create(&testPerson)

	updateData := services.Body{
		Name:    "Arthur",
		Surname: "New",
	}
	updateDataJSON, _ := json.Marshal(updateData)

	req, _ := http.NewRequest("PUT", fmt.Sprintf("/update/%d", testPerson.ID), bytes.NewBuffer(updateDataJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected status code to be 200 OK")

	var updatedPerson models.Person
	initializers.DB.First(&updatedPerson, testPerson.ID)
	assert.Equal(t, "Arthur", updatedPerson.Name, "Name should be updated")
	assert.Equal(t, "New", updatedPerson.Surname, "Surname should be updated")
	assert.Equal(t, 71, updatedPerson.Age)
	assert.Equal(t, "male", updatedPerson.Gender)
	assert.Equal(t, "GH", updatedPerson.Nationality)
}
