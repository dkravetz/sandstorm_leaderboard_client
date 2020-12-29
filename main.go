package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/sahilm/fuzzy"
	"log"
	"net/http"
	"os"
	"time"
)

func main(){
	nickNamePtr := flag.String("Nickname", "", "Steam nickname to search")
	flag.Parse()

	var players Players
	GetPlayers("https://leaderboard.sandstorm.game/api/v0/Players/getrankedplayers", &players)

	if *nickNamePtr == ""{
		fmt.Println("Type 'exit' or 'quit' to exit the application.")
		fmt.Print("Type in the nickname you're looking for:")
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan(){
			if scanner.Text() == "quit" || scanner.Text() == "exit"{
				os.Exit(0)
			}
			c := make(chan string)
			go searchPlayers(&players, scanner.Text(), c)
			for i := range c {
				fmt.Println(i)
			}
			fmt.Print("Type in the nickname you're looking for: ")
		}
	} else {
		if results := fuzzy.FindFrom(*nickNamePtr, &players); results != nil {
			for _, item := range results {
				fmt.Println(players[item.Index])
			}
		} else {
			fmt.Println("Sorry, can't find that player")
		}
	}
}

func GetPlayers(url string, players interface{})  error{
	client := http.Client{Timeout: 10 * time.Second}
	res, err := client.Get(url)
	if err != nil {
		log.Fatal("Couldn't request URL. ", err)
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		return err
	}
	err = json.NewDecoder(res.Body).Decode(&players)
	if err != nil{
		return err
	}
	return nil
}

func searchPlayers(players *Players, query string, c chan string){
	if results := fuzzy.FindFrom(query, players); results != nil {
		for _, item := range results {
			c <- fmt.Sprintf("%s", (*players)[item.Index])
		}
		close(c)
	} else {
		c <- "Sorry, can't find that player"
		close(c)
	}
}