package controllers

import (
	"example/m/initializers"
	"example/m/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

var personService = services.NewPersonService()

// Endpoint for creating a new object
func PeopleCreate(c *gin.Context) {
	initializers.Log.Info("Create Person request received")

	var body services.Body

	if err := c.Bind(&body); err != nil {
		initializers.Log.WithError(err).Error("Error binding request body for update")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong data format"})
		return
	}

	if body.Name == "" || body.Surname == "" {
		initializers.Log.Errorln("Name and surname are requierd, instance wasn't created")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name and Surname are required"})
		return
	}

	// Complete data with gender, age and nationality
	initializers.Log.Debug("Calling CompleteObject with: ", body)
	err := services.CompleteObject(&body)
	if err != nil {
		initializers.Log.WithError(err).Error("Error completing object")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "External API error"})
		return
	}

	// Create person
	person, er := personService.CreatePerson(body.Name, body.Surname, body.Patronymic, body.Gender, body.Nationality, body.Age)

	if er != nil {
		initializers.Log.WithError(err).Error("Error creating person")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong todo format"})
		return
	}

	initializers.Log.Info("Person created successfully: ", person)
	c.JSON(http.StatusCreated, gin.H{
		"person": person,
	})
}

func PeopleDelete(c *gin.Context) {
	// Get param
	id := c.Param("id")
	initializers.Log.Info("Request to delete person with ID: ", id)

	// Delete person using the service
	if err := personService.DeletePerson(id); err != nil {
		initializers.Log.WithError(err).Error("Error deleting person with ID: ", id)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Respond with status
	initializers.Log.Info("Person deleted successfully with ID: ", id)
	c.Status(http.StatusNoContent)
}

func PeopleUpdate(c *gin.Context) {
	//Get param
	id := c.Param("id")
	initializers.Log.Info("Request to update person with ID: ", id)

	//Get person
	person, err := personService.FindPerson(id)

	if err != nil {
		initializers.Log.WithError(err).Error("Person not found with ID: ", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "Person doesn't exist"})
		return
	}

	//Get data
	var body services.Body

	if err := c.Bind(&body); err != nil {
		initializers.Log.WithError(err).Error("Error binding request body for update")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong data format"})
		return
	}

	// Update data with new params
	if body.Name != person.Name {
		initializers.Log.Debug("Calling CompleteObject with: ", body)
		err = services.CompleteObject(&body)
		if err != nil {
			initializers.Log.WithError(err).Error("Error completing object")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "External API error"})
			return
		}
	}

	//Update person
	if err := personService.UpdatePerson(person, body.Name, body.Surname, body.Patronymic, body.Gender, body.Nationality, body.Age); err != nil {
		initializers.Log.WithError(err).Error("Error updating person with ID: ", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//Respond with updated person
	initializers.Log.Info("Person updated successfully with ID: ", id)
	c.JSON(http.StatusOK, gin.H{
		"person": person,
	})
}

func PeopleGet(c *gin.Context) {
	// Get param
	name := c.Query("name")
	surname := c.Query("surname")
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "10")
	initializers.Log.Infof("Fetching people with filters - Name: %s, Surname: %s, Page: %s, PageSize: %s", name, surname, page, pageSize)

	// Get people
	people, err := personService.GetPeople(name, surname, page, pageSize)

	if err != nil {
		initializers.Log.WithError(err).Error("Error fetching people")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	initializers.Log.Debug("People fetched successfully")
	c.JSON(http.StatusOK, gin.H{
		"people": people,
	})
}
