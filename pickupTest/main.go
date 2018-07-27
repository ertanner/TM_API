package main

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"bytes"
	"time"
	"strings"
	"strconv"
)

func main()  {
	start := time.Now()
	ch := make(chan string)

	for i:= 0; i< 10000; i++ {
		fmt.Println(i)
		go getQuote(i, ch)
		time.Sleep( 1 * time.Second)
	}
	for i:= 0; i< 10000; i++ {
		fmt.Println(<-ch)
	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}


func getQuote(count int, ch chan<-string) {
	startQTime := time.Now()
	type AutoGenerated struct {
		DyltPickupReqs struct {
			DyltPickupReq struct {
				AccountNumber          string `json:"accountNumber"`
				UserName               string `json:"userName"`
				Password               string `json:"password"`
				ShipmentID             string `json:"shipmentID"`
				BillTerms              string `json:"billTerms"`
				ServiceType            string `json:"serviceType"`
				ShipperName            string `json:"shipperName"`
				ShipperAddress1        string `json:"shipperAddress1"`
				ShipperCity            string `json:"shipperCity"`
				ShipperState           string `json:"shipperState"`
				ShipperZip             string `json:"shipperZip"`
				ShipperContactName     string `json:"shipperContactName"`
				ShipperContactNumber   string `json:"shipperContactNumber"`
				ConsigneeName          string `json:"consigneeName"`
				ConsigneeAddress1      string `json:"consigneeAddress1"`
				ConsigneeAddress2      string `json:"consigneeAddress2"`
				ConsigneeCity          string `json:"consigneeCity"`
				ConsigneeState         string `json:"consigneeState"`
				ConsigneeZip           string `json:"consigneeZip"`
				ConsigneeContactName   string `json:"consigneeContactName"`
				ConsigneeContactNumber string `json:"consigneeContactNumber"`
				PickupStartDate        string `json:"pickupStartDate"`
				PickupStartTime        string `json:"pickupStartTime"`
				PickupEndDate          string `json:"pickupEndDate"`
				PickupEndTime          string `json:"pickupEndTime"`
				Items                  struct {
					Item struct {
						Weight      string `json:"weight"`
						Pcs         string `json:"pcs"`
						Description string `json:"description"`
						ActualClass string `json:"actualClass"`
					} `json:"item"`
				} `json:"items"`
				ShipReferences struct {
					ShipReference struct {
						ReferenceType   string `json:"referenceType"`
						ReferenceNumber string `json:"referenceNumber"`
					} `json:"shipReference"`
				} `json:"shipReferences"`
				Accessorials struct {
					Accessorial struct {
						AccName string `json:"accName"`
						AccID   string `json:"accId"`
					} `json:"accessorial"`
				} `json:"accessorials"`
			} `json:"dyltPickupReq"`
		} `json:"dyltPickupReqs"`
	}

	jsonData, err := ioutil.ReadFile("./pickup_test.json")
	if err  != nil {
		panic(err.Error())
	}

	var jsonFile AutoGenerated
	err = json.Unmarshal( jsonData, &jsonFile)
	if err  != nil {
		panic(err.Error())
	}

	jsonValue, _ := json.Marshal(jsonFile)

	token := getToken()
	request, _ := http.NewRequest("POST", "https://test-api.dylt.com/pickup", bytes.NewBuffer(jsonValue))
	request.Header.Set("Content-Type", "application/json")
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

























