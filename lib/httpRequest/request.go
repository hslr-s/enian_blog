package httpRequest

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func HttpPostJson(url string, params map[string]interface{}) ([]byte, error) {
	jsonByte, _ := json.Marshal(params)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonByte))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	// statuscode := resp.StatusCode
	// hea := resp.Header
	return ioutil.ReadAll(resp.Body)

	// fmt.Println(string(body))
	// fmt.Println(statuscode)
	// fmt.Println(hea)
}
