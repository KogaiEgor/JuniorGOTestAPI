package services

import (
	"encoding/json"
	"example/m/initializers"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Body struct {
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic"`
	Age         int    `json:"age,omitempty"`
	Gender      string `json:"gender,omitempty"`
	Nationality string `json:"nationality,omitempty"`
}

// Ответы от API
type AgifyResponse struct {
	Age int `json:"age"`
}

type GenderizeResponse struct {
	Gender string `json:"gender"`
}

type NationalizeResponse struct {
	Country []struct {
		CountryID string `json:"country_id"`
	} `json:"country"`
}

// Add information to the object from external API
func CompleteObject(person *Body) error {
	initializers.Log.Infof("Starting to complete object for: %s", person.Name)
	// Получение возраста
	agifyResp, err := http.Get(fmt.Sprintf("https://api.agify.io/?name=%s", person.Name))
	if err != nil {
		initializers.Log.WithError(err).Error("Failed to request age from Agify")
		return err
	}
	defer agifyResp.Body.Close()
	agifyBody, _ := ioutil.ReadAll(agifyResp.Body)
	var agifyData AgifyResponse
	json.Unmarshal(agifyBody, &agifyData)
	person.Age = agifyData.Age
	initializers.Log.Debugf("Received age from Agify: %d", person.Age)

	// Получение пола
	genderizeResp, err := http.Get(fmt.Sprintf("https://api.genderize.io/?name=%s", person.Name))
	if err != nil {
		initializers.Log.WithError(err).Error("Failed to request gender from Genderize")
		return err
	}
	defer genderizeResp.Body.Close()
	genderizeBody, _ := ioutil.ReadAll(genderizeResp.Body)
	var genderizeData GenderizeResponse
	json.Unmarshal(genderizeBody, &genderizeData)
	person.Gender = genderizeData.Gender
	initializers.Log.Debugf("Received gender from Genderize: %s", person.Gender)

	// Получение национальности
	nationalizeResp, err := http.Get(fmt.Sprintf("https://api.nationalize.io/?name=%s", person.Name))
	if err != nil {
		initializers.Log.WithError(err).Error("Failed to request nationality from Nationalize")
		return err
	}
	defer nationalizeResp.Body.Close()
	nationalizeBody, _ := ioutil.ReadAll(nationalizeResp.Body)
	var nationalizeData NationalizeResponse
	json.Unmarshal(nationalizeBody, &nationalizeData)
	if len(nationalizeData.Country) > 0 {
		person.Nationality = nationalizeData.Country[0].CountryID
		initializers.Log.Debugf("Received nationality from Nationalize: %s", person.Nationality)
	}

	initializers.Log.Info("Object completed successfully")
	return nil
}
