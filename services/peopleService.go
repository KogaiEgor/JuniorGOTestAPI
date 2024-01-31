package services

import (
	"example/m/initializers"
	"example/m/models"
	"fmt"
	"strconv"
)

type PersonService struct {
	Name        string  `json:"name"`
	Surname     string  `json:"surname"`
	Patronymic  *string `json:"patronymic"`
	Age         int     `json:"age"`
	Sex         string  `json:"sex"`
	Nationality string  `json:"nationality"`
}

func NewPersonService() *PersonService {
	return &PersonService{}
}

// Finds person by id
func (s *PersonService) FindPerson(id string) (*models.Person, error) {
	initializers.Log.Infof("Finding person with ID: %s", id)
	var person models.Person

	if err := initializers.DB.First(&person, id).Error; err != nil {
		initializers.Log.WithError(err).Errorf("Error finding person with ID: %s", id)
		return nil, err
	}

	initializers.Log.Debug("Person found: ", person)
	return &person, nil
}

// create object instance in db
func (s *PersonService) CreatePerson(name, surname, patronymic, gender, nationality string, age int) (*models.Person, error) {
	initializers.Log.Infof("Creating new person: %s %s", name, surname)
	var patronymicPtr *string
	if patronymic != "" {
		patronymicPtr = &patronymic
	}

	man := &models.Person{
		Name:        name,
		Surname:     surname,
		Patronymic:  patronymicPtr,
		Age:         age,
		Gender:      gender,
		Nationality: nationality,
	}

	if err := initializers.DB.Create(man).Error; err != nil {
		initializers.Log.WithError(err).Error("Error creating new person")
		return nil, err
	}

	initializers.Log.Info("Person created successfully")
	return man, nil
}

// Updates instance in db
func (s *PersonService) UpdatePerson(m *models.Person, name, surname, patronymic, gender, nationality string, age int) error {
	initializers.Log.Infof("Updating person %s", name)
	var patronymicPtr *string
	if patronymic != "" {
		patronymicPtr = &patronymic
	}

	updateFields := &models.Person{
		Name:        name,
		Surname:     surname,
		Patronymic:  patronymicPtr,
		Age:         age,
		Gender:      gender,
		Nationality: nationality,
	}

	if err := initializers.DB.Model(m).Updates(updateFields).Error; err != nil {
		initializers.Log.WithError(err).Errorf("Error updating person %s", name)
		return err
	}

	initializers.Log.Info("Person updated successfully")
	return nil
}

// delete instance of given id
func (s *PersonService) DeletePerson(id string) error {
	i, _ := strconv.Atoi(id)

	initializers.Log.Infof("Deleting person with ID: %d", i)
	result := initializers.DB.Delete(&models.Person{}, i)

	if result.RowsAffected == 0 {
		err := fmt.Errorf("There is no person with id %s", id)
		initializers.Log.WithError(err).Error("Error deleting person")
		return err
	}

	initializers.Log.Info("Person deleted successfully")
	return nil
}

// gets people, filters optional
func (s *PersonService) GetPeople(name, surname, pageStr, pageSizeStr string) ([]models.Person, error) {
	initializers.Log.Infof("Fetching people with filters - Name: %s, Surname: %s, Page: %s, PageSize: %s", name, surname, pageStr, pageSizeStr)
	var people []models.Person

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	offset := (page - 1) * pageSize

	query := initializers.DB
	if name != "" {
		query = query.Where("name = ?", name)
	}
	if surname != "" {
		query = query.Where("surname = ?", surname)
	}

	if err := query.Offset(offset).Limit(pageSize).Find(&people).Error; err != nil {
		initializers.Log.WithError(err).Error("Error fetching people")
		return nil, err
	}

	initializers.Log.Debug("People fetched successfully")
	return people, nil
}
