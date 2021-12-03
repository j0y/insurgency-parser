package main

import (
	"bufio"
	"fmt"
	insurgencylog "my.com/insurgency-log"
	"os"
)

// Usage:
//
// From file:
// go run main.go example.log
//
// From STDIN:
// cat example.log | go run main.go
//
// To File:
// go run main.go > out.txt
//
// Omit errors:
// go run main.go 2>/dev/null

type matchInfoStruct struct {
	Map       string `json:"map"`
	Rounds    uint8  `json:"rounds"`
	StartedAt uint64 `json:"started_at"`
	Duration  uint32 `json:"duration"`
	Won       bool   `json:"won"`
}
type playerStatsStruct struct {
	Kills       uint32            `json:"kills"`
	Deaths      uint32            `json:"deaths"`
	WeaponStats map[string]uint32 `json:"weapon_stats"`
}

func main() {
	var matchInfo matchInfoStruct
	playerStats := make(map[string]playerStatsStruct)

	var file *os.File
	var err error

	file, err = os.Open("example.log")
	/*if len(os.Args) < 2 {
		file = os.Stdin
	} else {
		file, err = os.Open(os.Args[1])
	}*/

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	r := bufio.NewReader(file)

	// read first line
	l, _, err := r.ReadLine()

	for err == nil {
		message, errParse := insurgencylog.Parse(string(l))
		if errParse != nil {
			// print parse errors to stderr
			fmt.Fprintf(os.Stderr, "ERROR: %s", insurgencylog.ToJSON(message))
		} else {
			// print to stdout
			//fmt.Fprintf(os.Stdout, "%s", insurgencylog.ToJSON(m))
		}

		switch m := message.(type) {
		case insurgencylog.LoadingMap:
			matchInfo.Map = m.Map
			matchInfo.StartedAt = uint64(m.Time.Unix())
		case insurgencylog.PlayerKill:
			if m.Attacker.SteamID == "BOT" && m.Victim.SteamID != "BOT" {
				stats, _ := playerStats[m.Victim.SteamID]
				stats.Deaths++
				playerStats[m.Victim.SteamID] = stats
			}
			if m.Victim.SteamID == "BOT" && m.Attacker.SteamID != "BOT" {
				stats, _ := playerStats[m.Attacker.SteamID]
				stats.Kills++

				if stats.WeaponStats == nil {
					stats.WeaponStats = make(map[string]uint32)
				}
				weaponStats, _ := stats.WeaponStats[m.Weapon]
				weaponStats++
				stats.WeaponStats[m.Weapon] = weaponStats

				playerStats[m.Attacker.SteamID] = stats
			}
		case insurgencylog.RoundWin:
			if m.Team == insurgencylog.TeamSecurity {
				matchInfo.Rounds++
			} else if m.Team == insurgencylog.TeamInsurgent {
				matchInfo.Won = true
				matchInfo.Duration = uint32(uint64(m.Time.Unix()) - matchInfo.StartedAt)
			}
		}

		// next line
		l, _, err = r.ReadLine()
	}

	fmt.Printf("%v", matchInfo)
	fmt.Printf("%v", playerStats)
}
