package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	logger "github.com/mkm29/mpulsar/pkg/log"

	"github.com/mmcloughlin/geohash"
)

func geocode_http(w http.ResponseWriter, r *http.Request) {
	// Get query URL parameter string
	s := r.URL.Query().Get("q")
	err, l := geocode(s)
	if err != nil {
		logger.WithRequest(r).Error(err)
	}
	fmt.Printf("Location: %+v\n", l)
}

func geocode(s string) (error, Location) {
	/*
		Post request to Pelias API
		Parse response
		Create Location object from response
		Return Location
	*/

	/*
		Example cURL command: curl -X GET http://localhost:4000/v1/search\?text\=iah\&size\=1\&layers\=address | jq '.features | .[0] | .geometry."coordinates"'
		Get the first element of the features array, then get the geometry object, then get the coordinates array (this will be an array of 2 elements, latitude and longitude)

		Example response: see src/pelias_response.json

	*/

	url := fmt.Sprintf(PELIAS_URL, s)
	fmt.Printf("URL: %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		logger.Error(err)
		return err, Location{}
	}
	if err != nil {
		logger.Error(err)
		return err, Location{}
	}
	// Extract JSON from body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return err, Location{}
	}
	defer resp.Body.Close()

	// retrun if response Body is empty
	if len(body) == 0 {
		return nil, Location{}
	}

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		logger.Error(err)
		return err, Location{}
	}

	// Unmarshal JSON into Location object
	var l Location
	l.Latitude = data["features"].([]interface{})[0].(map[string]interface{})["geometry"].(map[string]interface{})["coordinates"].([]interface{})[0].(float64)
	l.Longitude = data["features"].([]interface{})[0].(map[string]interface{})["geometry"].(map[string]interface{})["coordinates"].([]interface{})[1].(float64)
	// extract properties from data
	properties := data["features"].([]interface{})[0].(map[string]interface{})["properties"].(map[string]interface{})

	// Populate Location object
	l.Confidence = properties["confidence"].(float64)
	l.Country = properties["country"].(string)
	l.CountryCode = properties["country_code"].(string)
	l.Label = properties["label"].(string)
	l.City = properties["locality"].(string)
	l.Region = properties["region"].(string)
	PostalCode := properties["postalcode"].(string)
	l.PostalCode, _ = strconv.Atoi(PostalCode)
	l.Address = properties["housenumber"].(string) + " " + properties["street"].(string)
	fmt.Printf("Location: \n%+v\n", l)

	// Get geohashes from Location object
	geo := geohash.Encode(l.Latitude, l.Longitude)
	for i := 5; i < 12; i++ {
		l.Geohashes = append(l.Geohashes, geo[:i])
	}

	return nil, l
}
