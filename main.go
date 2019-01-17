package main

import (
	"encoding/json"
	"fmt"
	"github.com/bradfitz/slice"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type Response struct {
	Status  string `json:"status"`
	Code    int    `json:"code "`
	Data    *Data  `json:"data"`
	Message string `json:"message"`
}

type Data struct {
	Team *Team `json:"team"`
}

type Team struct {
	Name    string    `json:"name"`
	Players []*Player `json:"players"`
}

type Player struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Age  string `json:"age"`
}

type AllPlayers struct {
	Id    string
	Name  string
	Age   string
	Teams []string
}

func main() {
	teamIds := []int{2, 5, 6, 9, 21, 26, 34, 45, 61, 96} // hard code because of the specific teams
	var url string
	var response *Response
	players := make(map[int][]AllPlayers)

	for _, teamId := range teamIds {
		url = fmt.Sprintf("https://.../api/teams/en/%d.json", teamId)
		req, err := http.NewRequest("GET", url, nil)
		checkErr(err)

		req.Header.Add("content-type", "application/json")
		res, err := http.DefaultClient.Do(req)
		checkErr(err)

		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		checkErr(err)

		err = json.Unmarshal(body, &response)
		checkErr(err)

		addPlayers(response, players)
	}

	sortedPlayers := sortPlayers(players)
	printPlayers(sortedPlayers)
}



func addPlayers(response *Response, players map[int][]AllPlayers){
	for _, player := range response.Data.Team.Players {
		id, _ := strconv.Atoi(player.Id)
		if _, is := players[id]; is { // check if player is already inserted, only add his team
			players[id][0].Teams = append(players[id][0].Teams, response.Data.Team.Name)
		} else {
			var p AllPlayers
			p.Name = player.Name
			p.Age = player.Age
			p.Teams = append(p.Teams, response.Data.Team.Name)
			players[id] = append(players[id], p)
		}
	}
}

func sortPlayers(players map[int][]AllPlayers) []AllPlayers{
	var sorted []AllPlayers
	for _, x := range players {
		sorted = append(sorted, AllPlayers{x[0].Id, x[0].Name, x[0].Age, x[0].Teams}) // create slice from map
	}

	slice.Sort(sorted[:], func(i, j int) bool { // sort slice
		return sorted[i].Name < sorted[j].Name
	})

	return sorted
}

func printPlayers(players []AllPlayers) {
	i :=1
	for _, x := range players {
		final := fmt.Sprintf("%d. %s; %s;",i, x.Name, x.Age)
		for _, team := range x.Teams {
			final += fmt.Sprintf(" %s,", team)
		}
		final = strings.TrimRight(final, ",")
		fmt.Println(final)
		i++
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

