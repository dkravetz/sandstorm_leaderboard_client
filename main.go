package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/sahilm/fuzzy"
	"log"
	"os"
	"sync"
)

func main(){
	nickNamePtr := flag.String("Nickname", "", "Steam nickname to search")
	flag.Parse()

	var players Players
	var wg sync.WaitGroup
	fmt.Println("Loading player list, please wait...")
	err := GetPlayers(&players, &wg)
	wg.Wait()
	fmt.Println("Player list loaded.")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	if *nickNamePtr == ""{
		fmt.Println("Type 'exit' or 'quit' to exit the application.")
		fmt.Println("Type in the nickname you're looking for:")
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
			fmt.Println("Type in the nickname you're looking for:")
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