package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	logger "github.com/mkm29/mpulsar/pkg/logging"
	"github.com/mmcloughlin/geohash"
)

func geocodeHTTP(c *gin.Context) {
	// Get query URL parameter string
	s := c.Request.URL.Query().Get("q")
	err, l := geocode(s)
	if err != nil {
		logger.WithRequest(c.Request).Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	logger.Log("INFO", fmt.Sprintf("Location: %+v", l))
}

// Get HTTP response from URL
func getHTTP(url string) ([]byte, error) {
	logger.Log("INFO", fmt.Sprintf("GET %s", url))
	resp, err := http.Get(url)
	if err != nil {
		logger.Log("ERROR", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Log("ERROR", err)
		return nil, err
	}
	return body, nil
}

func geocode(s string) (*Location, error) {
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

	url := fmt.Sprintf(peliasURL, s)
	// call get_http to get the response
	body, err := getHTTP(url)
	if err != nil {
		return nil, err
	}

	// retrun if response Body is empty
	if len(body) == 0 {
		// create Error object
		err := fmt.Errorf("Empty response body")
		logger.Log("ERROR", err)
		return nil, err
	}

	// Unmarshal response body into PeliasResponse struct
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		logger.Log("ERROR", err)
		return nil, err
	}

	// Unmarshal JSON into Location object
	l := &Location{}
	l.populate(data)

	/* Geohash
	 *
	 * Get geohashes from Location object
	 *
	 * should probably be published to a different topic for processing (adding Geohashes)
	 * The geohash metadata will be stored in a separate Cassandra table
	 *
	 */
	l.addGeohash()

	return l, nil
}

func (l *Location) addGeohash() {
	geo := geohash.Encode(l.Latitude, l.Longitude)
	for i := 5; i < 12; i++ {
		l.Geohashes = append(l.Geohashes, geo[:i])
	}
	// publish to topic

}

func (l *Location) populate(data map[string]interface{}) {
	/*
		Populate
	*/
	logger.Log("INFO", "Populating Location object")
	coords := data["features"].([]interface{})[0].(map[string]interface{})["geometry"].(map[string]interface{})["coordinates"].([]interface{})
	l.Latitude = coords[0].(float64)
	l.Longitude = coords[1].(float64)
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
	logger.Log("INFO", fmt.Sprintf("Location object populated: %+v", l))
}
