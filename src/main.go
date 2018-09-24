package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	url        = "https://fantasy.premierleague.com/drf/entry/75082/event/4/picks"
	playersUrl = "https://fantasy.premierleague.com/drf/bootstrap-static"
)

type TeamMainElem struct {
	TeamPicksElem []TeamPicksElem `json:"picks"`
}
type TeamPicksElem struct {
	Element       int64 `json:"element"`
	Position      int64 `json:"position"`
	IsCaptain     bool  `json:"is_captain"`
	IsViceCaptain bool  `json:"is_vice_captain"`
	Multiplier    int64 `json:"multiplier"`
}
type PlayerMainElem struct {
	PlayerElem []PlayerElem `json:"elements"`
}
type PlayerElem struct {
	Id      int64  `json:"id"`
	WebName string `json:"web_name"`
}

func getTeamInfo() {
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "pg-fpl")

	resp, err := netClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}
	s := new(TeamMainElem)
	err = json.Unmarshal(body, &s)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(s.TeamPicksElem[0].Element)
}

func getPlayerMapping() {
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest(http.MethodGet, playersUrl, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "pg-fpl")

	resp, err := netClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}
	s := new(PlayerMainElem)
	err = json.Unmarshal(body, &s)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(s.PlayerElem[1].WebName)

}

func main() {
	fmt.Println("starting main program")

	//getTeamInfo()
	getPlayerMapping()
}
