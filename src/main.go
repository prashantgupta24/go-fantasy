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
	teamURL         = "https://fantasy.premierleague.com/drf/entry/%v/event/%v/picks"
	allPlayersURL   = "https://fantasy.premierleague.com/drf/bootstrap-static"
	participantsURL = "https://fantasy.premierleague.com/drf/leagues-classic-standings/%v?phase=1&le-page=1&ls-page=1"
	csvFileName     = "data/result-%v-%v.csv"
)

type fantasyMain struct {
	httpClient         *http.Client
	playerMap          map[int64]string
	leagueParticipants []int64
	playerOccurances   []map[string]int
}

/* Structure of JSON

picks
    0
    element	260
    1
    element	247
*/
type ParticipantTeamInfo struct {
	TeamPlayers []TeamPlayers `json:"picks"`
}
type TeamPlayers struct {
	Element int64 `json:"element"`
}

/* Structure of JSON

elements
    0
    id	1
    photo	"11334.jpg"
    web_name	"Cech"
    team_code	3
    status	"i"
    code	11334
    first_name	"Petr"
    second_name	"Cech"
    squad_number	1

    1
    id	2
    photo	"80201.jpg"
    web_name	"Leno"
    team_code	3
    status	"a"
    code	80201
    first_name	"Bernd"
    second_name	"Leno"
    squad_number	19
*/
type AllPlayers struct {
	Players []Players `json:"elements"`
}
type Players struct {
	Id      int64  `json:"id"`
	WebName string `json:"web_name"`
}

/* Structure of JSON

standings
    has_next	true
    number	1
    results
        0
        id	13987896
        rank	1
        last_rank	1
        rank_sort	1
        total	575
        entry	2557010

        1
        id	13148025
        rank	2
        last_rank	5
        rank_sort	2
        total	572
        entry	2415205
*/
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

func getTeamInfoForParticipant(participantNumber int64, gameweek int, playerOccurance *map[string]int, fantasyMain *fantasyMain) error {
	teamURL := fmt.Sprintf(teamURL, participantNumber, gameweek)

	response := makeRequest(fantasyMain, teamURL)
	ParticipantTeamInfo := new(ParticipantTeamInfo)
	err := json.Unmarshal(response, &ParticipantTeamInfo)
	if err != nil {
		return err
	}

	for _, player := range ParticipantTeamInfo.TeamPlayers {
		(*playerOccurance)[fantasyMain.playerMap[player.Element]]++
	}
	return nil
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

func getParticipantsInLeague(fantasyMain *fantasyMain, leagueCode int) {
	participantsURL := fmt.Sprintf(participantsURL, leagueCode)

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

func writeToFile(fantasyMain *fantasyMain, leagueCode int) {

	fileName := fmt.Sprintf(csvFileName, time.Now().Format("2006-01-02"), leagueCode)
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	//Headers
	var record []string
	record = append(record, "Player")
	for gameweekNum := range fantasyMain.playerOccurances {
		record = append(record, fmt.Sprintf("Gameweek %v", gameweekNum+1))
	}
	err = writer.Write(record)
	if err != nil {
		panic(err)
	}

	allPlayers := fantasyMain.playerOccurances[len(fantasyMain.playerOccurances)-1]

	for player := range allPlayers {

		var record []string
		record = append(record, string(player))

		for _, playerOccurances := range fantasyMain.playerOccurances {
			record = append(record, strconv.Itoa(playerOccurances[player]))
		}

		//fmt.Println("player: ", player, "Used: ", occurances)

		err := writer.Write(record)
		if err != nil {
			panic(err)
		}
	}
}
func main() {
	fmt.Println("starting main program")

	start := time.Now()
	defer func() {
		fmt.Printf("Took %v to fetch all gameweeks data!\n", time.Since(start))
	}()

	var httpClient = &http.Client{
		Timeout: time.Second * 10,
	}

	fantasyMain := &fantasyMain{
		httpClient: httpClient,
		playerMap:  make(map[int64]string),
	}

	gameweekMax := 38
	leagueCode := 313

	getPlayerMapping(fantasyMain)
	getParticipantsInLeague(fantasyMain, leagueCode)

	for gameweek := 1; gameweek <= gameweekMax; gameweek++ {
		playerOccuranceForGameweek := make(map[string]int)
		var err error
		fmt.Println("Fetching gameweek ", gameweek)
		for _, participant := range fantasyMain.leagueParticipants[0:10] {
			err = getTeamInfoForParticipant(participant, gameweek, &playerOccuranceForGameweek, fantasyMain)
			if err != nil {
				break
			}
		}
		if err != nil {
			break
		}
		fantasyMain.playerOccurances = append(fantasyMain.playerOccurances, playerOccuranceForGameweek)
	}

	writeToFile(fantasyMain, leagueCode)

}
