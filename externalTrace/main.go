package main

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"time"
	"strings"
	"strconv"
)

func main()  {
	ch := make(chan string)
	start := time.Now()
	proBill := 2500

	for i:= 0; i< 10000; i++ {
		fmt.Println(i)
		if (proBill == 3500){
			proBill = 2500
		}
		proBill ++
		go getQuote(i, strconv.Itoa(proBill), ch)
		time.Sleep( 2500 * time.Millisecond)
		fmt.Println(time.Since(start))
	}
	for i:= 0; i< 10000; i++ {
		fmt.Println(<-ch)
	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func getQuote(count int, proBill string, ch chan<-string) {
	startQTime := time.Now()

	token := getToken()
	request, _ := http.NewRequest("GET", "https://test-api.dylt.com/externalTrace/" + proBill, nil)
	request.Header.Set("Authorization", "Bearer "+ token)
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
	}
	secs := time.Since(startQTime).Seconds()
	fmt.Println("Fuction Time: " + strconv.FormatFloat(secs, 'f', -1, 64))
	//fmt.Println("response Status:", response.Status)
	body, _ := ioutil.ReadAll(response.Body)
	ch <- fmt.Sprintf("%.2f elapsed with response length: %d ", secs, len(body))
}

func getToken() string{
	//Consumer Key: x5Vxusddiy2pYqwpZytwxqkG0lW7Z6a5
	//Consumer Secret: ThzO25vxF0RDuA2U
	body := strings.NewReader(`client_secret=ThzO25vxF0RDuA2U&grant_type=client_credentials&client_id=x5Vxusddiy2pYqwpZytwxqkG0lW7Z6a5`)
	req, err := http.NewRequest("POST", "https://api.dylt.com/oauth/client_credential/accesstoken?grant_type=client_credentials", body)
	if err != nil {
		// handle err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	token, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}
	var data map[string]string
	json.Unmarshal(token, &data)
	//fmt.Println(data)
	//fmt.Println(data["access_token"])
	return data["access_token"]
}


























