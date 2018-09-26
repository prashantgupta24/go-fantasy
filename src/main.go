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
	allPlayersURL   = "https://fantasy.premierleague.com/drf/bootstrap-static"
	participantsURL = "https://fantasy.premierleague.com/drf/leagues-classic-standings/313?phase=1&le-page=1&ls-page=1"
)

type fantasyMain struct {
	httpClient         *http.Client
	playerMap          map[int64]string
	leagueParticipants []int64
	playerOccurances   map[string]int
}

type ParticipantTeamInfo struct {
	TeamPlayers []TeamPlayers `json:"picks"`
}
type TeamPlayers struct {
	Element int64 `json:"element"`
}
type AllPlayers struct {
	Players []Players `json:"elements"`
}
type Players struct {
	Id      int64  `json:"id"`
	WebName string `json:"web_name"`
}
type LeagueParticipants struct {
	LeagueStandings LeagueStandings `json:"standings"`
}
type LeagueStandings struct {
	LeagueResults []LeagueResults `json:"results"`
}
type LeagueResults struct {
	Entry int64 `json:"entry"`
}

func makeRequest(fantasyMain *fantasyMain, URL string) []byte {
	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "pg-fpl")

	resp, err := fantasyMain.httpClient.Do(req)
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
	return body
}

func getTeamInfoForParticipant(participantNumber int64, fantasyMain *fantasyMain) {

	teamURL := teamURL + strconv.FormatInt(participantNumber, 10) + "/event/6/picks"

	response := makeRequest(fantasyMain, teamURL)

	ParticipantTeamInfo := new(ParticipantTeamInfo)
	err := json.Unmarshal(response, &ParticipantTeamInfo)
	if err != nil {
		panic(err.Error())
	}

	for _, player := range ParticipantTeamInfo.TeamPlayers {
		fantasyMain.playerOccurances[fantasyMain.playerMap[player.Element]]++
	}

}

func getPlayerMapping(fantasyMain *fantasyMain) {

	response := makeRequest(fantasyMain, allPlayersURL)

	allPlayers := new(AllPlayers)
	err := json.Unmarshal(response, &allPlayers)
	if err != nil {
		panic(err.Error())
	}

	for _, player := range allPlayers.Players {
		fantasyMain.playerMap[player.Id] = player.WebName
	}

	fmt.Println("Fetched " + strconv.Itoa(len(fantasyMain.playerMap)) + " players")

}

func getParticipantsInLeague(fantasyMain *fantasyMain) {

	response := makeRequest(fantasyMain, participantsURL)

	leagueParticipants := new(LeagueParticipants)
	err := json.Unmarshal(response, &leagueParticipants)
	if err != nil {
		panic(err.Error())
	}

	for _, participant := range leagueParticipants.LeagueStandings.LeagueResults {
		fantasyMain.leagueParticipants = append(fantasyMain.leagueParticipants, participant.Entry)
	}

	fmt.Println("Fetched " + strconv.Itoa(len(fantasyMain.leagueParticipants)) + " participants")
}

func writeToFile(fantasyMain *fantasyMain) {

	file, err := os.Create("result.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	for key, value := range fantasyMain.playerOccurances {
		//fmt.Println("player: ", key, "Used: ", value)
		s := []string{string(key), strconv.Itoa(value)}
		err := writer.Write(s)
		if err != nil {
			panic(err)
		}
	}
}
func main() {
	fmt.Println("starting main program")

	var httpClient = &http.Client{
		Timeout: time.Second * 10,
	}

	fantasyMain := &fantasyMain{
		httpClient:       httpClient,
		playerMap:        make(map[int64]string),
		playerOccurances: make(map[string]int),
	}

	getPlayerMapping(fantasyMain)

	getParticipantsInLeague(fantasyMain)
	for _, participant := range fantasyMain.leagueParticipants[0:5] {
		getTeamInfoForParticipant(participant, fantasyMain)
	}

	writeToFile(fantasyMain)

}
