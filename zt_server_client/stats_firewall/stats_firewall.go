package main

import (
	"fmt"
	"net/http"

	"github.com/project-1/config"
	"github.com/project-1/dataframe"

	"github.com/project-1/stats"
)

func main() {
	df := dataframe.ReadCSV(config.FILE) // Create a dataframe to be used
	df.DropCol(0)
	http.Handle("/", http.FileServer(http.Dir("web"))) // Use the index.html for the landing page
	http.HandleFunc("/player", func(w http.ResponseWriter, r *http.Request) {
		var player = r.FormValue("player_name")      // Take the response from the player name entry
		fmt.Fprint(w, stats.PlayerStats(df, player)) // Call the PlayerStats function with the df and player name
	})
	http.HandleFunc("/sorted_stats", func(w http.ResponseWriter, r *http.Request) {
		var stat = r.FormValue("stat") // Take the response from the stat selection
		fmt.Fprintln(w, stat, "\n")    // Print out the stat at the top of the page
		// fmt.Fprint(w, stats.StatLeader(df, stat)) // Print out the sorted and formatted results
	})
	fmt.Println("Listening on port 4000")
	http.ListenAndServe(":4000", nil)
}
