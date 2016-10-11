package awsapitools

import (
	"encoding/json"
	"net/http"

	"github.com/notonthehighstreet/awsnathealth/errhandling"
)

//GetInstanceJSONUserData retrives the UserData form aws metadata service end converst the json to a
func GetInstanceJSONUserData(url, key string) string {
	//Catch and log panic events
	var err error
	defer errhandling.CatchPanic(&err, "GetInstanceJSONUserData")

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var mymap map[string]string
	json.NewDecoder(resp.Body).Decode(&mymap)
	return mymap[key]
}
