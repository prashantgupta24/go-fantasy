package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	//teamURL         = "https://fantasy.premierleague.com/drf/entry/{{.}}/event/4/picks"
	teamURL         = "https://fantasy.premierleague.com/drf/entry/"
	playersURL      = "https://fantasy.premierleague.com/drf/bootstrap-static"
	participantsURL = "https://fantasy.premierleague.com/drf/leagues-classic-standings/313?phase=1&le-page=1&ls-page=1"
)

type fplMain struct {
	playerMap          map[int64]string
	leagueParticipants []int64
	playerOccurances   map[string]int
}
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

type Participants struct {
	LeagueStandings LeagueStandings `json:"standings"`
}

type LeagueStandings struct {
	LeagueResults []LeagueResults `json:"results"`
}

type LeagueResults struct {
	Entry int64 `json:"entry"`
}

func getTeamInfo(participantNumber int64, fplMain *fplMain) {

	// t := template.New("Participant template")

	// t, err := t.Parse(teamURL)
	// if err != nil {
	// 	log.Fatal("Parse: ", err)
	// 	return
	// }
	// // err = t.Execute(os.Stdout, participantNumber)
	// // if err != nil {
	// // 	log.Fatal("Execute: ", err)
	// // 	return
	// // }

	// fmt.Println("url is ", t)

	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}

	teamURL := teamURL + strconv.FormatInt(participantNumber, 10) + "/event/6/picks"
	req, err := http.NewRequest(http.MethodGet, teamURL, nil)
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

	for _, player := range s.TeamPicksElem {
		fplMain.playerOccurances[fplMain.playerMap[player.Element]]++
		//fmt.Println(fplMain.playerMap[player.Element])
	}

}

func getPlayerMapping(fplMain *fplMain) {
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest(http.MethodGet, playersURL, nil)
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
	for _, player := range s.PlayerElem {
		fplMain.playerMap[player.Id] = player.WebName
	}

	fmt.Println("Fetched " + strconv.Itoa(len(fplMain.playerMap)) + " players")
	//fmt.Println(fplMain.playerMap[332])

}

func getParticipantsInLeague(fplMain *fplMain) {
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest(http.MethodGet, participantsURL, nil)
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
	s := new(Participants)
	err = json.Unmarshal(body, &s)
	if err != nil {
		panic(err.Error())
	}
	for _, participant := range s.LeagueStandings.LeagueResults {
		fplMain.leagueParticipants = append(fplMain.leagueParticipants, participant.Entry)
	}

	fmt.Println("Fetched " + strconv.Itoa(len(fplMain.leagueParticipants)) + " participants")
}

func main() {
	file, err := os.Create("result.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	//var data = [][]string{{"abc", 1}, {"Line2", "golangcode.com"}}
	// data := make(map[string]int)

	// data["a"] = 1
	// data["b"] = 4
	// data["c"] = 2

	// for key, value := range data {
	// 	s := []string{key, strconv.Itoa(value)}
	// 	err := writer.Write(s)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }
	fmt.Println("starting main program")
	fplMain := &fplMain{
		playerMap:        make(map[int64]string),
		playerOccurances: make(map[string]int),
	}

	getPlayerMapping(fplMain)

	getParticipantsInLeague(fplMain)
	for _, participant := range fplMain.leagueParticipants[0:5] {
		//fmt.Println("Team ", i)
		getTeamInfo(participant, fplMain)
	}

	//sort.(fplMain.playerOccurances)
	for key, value := range fplMain.playerOccurances {
		//fmt.Println("player: ", key, "Used: ", value)
		s := []string{string(key), strconv.Itoa(value)}
		err := writer.Write(s)
		if err != nil {
			panic(err)
		}
	}
}
