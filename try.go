// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"net/http"
// 	"time"
// )

// const (
// 	url = "https://fantasy.premierleague.com/drf/entry/75082/event/4/picks"
// )

// type ResultStruct struct {
// 	Result []map[string]string
// }

// type MainElem struct {
// 	PicksElem []PicksElem `json:"picks"`
// }
// type PicksElem struct {
// 	Element       int64 `json:"element"`
// 	Position      int64 `json:"position"`
// 	IsCaptain     bool  `json:"is_captain"`
// 	IsViceCaptain bool  `json:"is_vice_captain"`
// 	Multiplier    int64 `json:"multiplier"`
// }

// func method1(resp *http.Response) {
// 	fmt.Println("method 1")
// 	// responseBody, err := ioutil.ReadAll(resp.Body)
// 	// if err != nil {
// 	// 	panic(err.Error())
// 	// }
// 	var data []map[string]interface{}
// 	//var jsonData ResultStruct
// 	//err = json.Unmarshal(responseBody, &jsonData)
// 	err := json.NewDecoder(resp.Body).Decode(&data)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(data)
// }

// func method2(resp *http.Response) {
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	s := new(MainElem)
// 	err = json.Unmarshal(body, &s)
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	fmt.Println(s.PicksElem[0].Element)
// }

// func method3() {
// 	respEx := `{
//    "active_chip":"3xc",
//    "automatic_subs":[  ],
//    "entry_history":{  },
//    "event":{
//       "id":4,
//       "name":"Gameweek 4",
//       "deadline_time":"2018-09-01T10:30:00Z",
//       "average_entry_score":44,
//       "finished":true,
//       "data_checked":true,
//       "highest_scoring_entry":2344578,
//       "deadline_time_epoch":1535797800,
//       "deadline_time_game_offset":3600,
//       "deadline_time_formatted":"01 Sep 11:30",
//       "highest_score":104,
//       "is_previous":false,
//       "is_current":false,
//       "is_next":false
//    },
//    "picks":[
//       {
//          "element":468,
//          "position":1,
//          "is_captain":false,
//          "is_vice_captain":false,
//          "multiplier":1
//       },
//       {
//          "element":113,
//          "position":2,
//          "is_captain":false,
//          "is_vice_captain":false,
//          "multiplier":1
//       },
//       {
//          "element":4,
//          "position":3,
//          "is_captain":false,
//          "is_vice_captain":false,
//          "multiplier":1
//       },
//       {
//          "element":332,
//          "position":4,
//          "is_captain":false,
//          "is_vice_captain":false,
//          "multiplier":1
//       },
//       {
//          "element":251,
//          "position":5,
//          "is_captain":false,
//          "is_vice_captain":true,
//          "multiplier":1
//       },
//       {
//          "element":164,
//          "position":6,
//          "is_captain":false,
//          "is_vice_captain":false,
//          "multiplier":1
//       },
//       {
//          "element":370,
//          "position":7,
//          "is_captain":false,
//          "is_vice_captain":false,
//          "multiplier":1
//       },
//       {
//          "element":149,
//          "position":8,
//          "is_captain":false,
//          "is_vice_captain":false,
//          "multiplier":1
//       },
//       {
//          "element":480,
//          "position":9,
//          "is_captain":false,
//          "is_vice_captain":false,
//          "multiplier":1
//       },
//       {
//          "element":280,
//          "position":10,
//          "is_captain":true,
//          "is_vice_captain":false,
//          "multiplier":3
//       },
//       {
//          "element":23,
//          "position":11,
//          "is_captain":false,
//          "is_vice_captain":false,
//          "multiplier":1
//       },
//       {
//          "element":136,
//          "position":12,
//          "is_captain":false,
//          "is_vice_captain":false,
//          "multiplier":1
//       },
//       {
//          "element":432,
//          "position":13,
//          "is_captain":false,
//          "is_vice_captain":false,
//          "multiplier":1
//       },
//       {
//          "element":28,
//          "position":14,
//          "is_captain":false,
//          "is_vice_captain":false,
//          "multiplier":1
//       },
//       {
//          "element":138,
//          "position":15,
//          "is_captain":false,
//          "is_vice_captain":false,
//          "multiplier":1
//       }
//    ]
// }`
// 	var s = new(MainElem)
// 	//var data map[string]interface{}
// 	err := json.Unmarshal([]byte(respEx), &s)
// 	//err := json.NewDecoder(respEx).Decode(&s)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(s)
// }

// func main() {
// 	fmt.Println("starting main program")

// 	var netClient = &http.Client{
// 		Timeout: time.Second * 10,
// 	}

// 	req, err := http.NewRequest(http.MethodGet, url, nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	req.Header.Set("User-Agent", "pg-tutorial")

// 	resp, getErr := netClient.Do(req)
// 	if getErr != nil {
// 		log.Fatal(getErr)
// 	}

// 	//resp, err := netClient.Get(url)
// 	if err != nil {
// 		panic(err)
// 	}
// 	//method1(resp)
// 	method2(resp)
// 	//method3()
// }
