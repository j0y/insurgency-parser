package main

import (
	"bufio"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/MrWaggel/gosteamconv"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	insurgencylog "my.com/insurgency-log"
	"os"
	"regexp"
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
	Ip        string `json:"ip"`
}

type weaponStatsStruct map[string]uint32

type playerStatsStruct struct {
	Kills       uint32            `json:"kills"`
	Deaths      uint32            `json:"deaths"`
	WeaponStats weaponStatsStruct `json:"weapon_stats"`
}

// Value Returns the JSON-encoded representation
func (a weaponStatsStruct) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	var matchInfo matchInfoStruct
	playerStats := make(map[string]playerStatsStruct)

	var file *os.File

	filename := "91.105.183.199_1638288928.log"
	file, err = os.Open(filename)

	re := regexp.MustCompile(`^[0-9,.]*`)

	ip := re.FindString(filename)
	if len(ip) == 0 {
		log.Fatal("ip not found")
	}
	matchInfo.Ip = ip

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
			_, err = fmt.Fprintf(os.Stderr, "ERROR: %s", insurgencylog.ToJSON(message))
			if err != nil {
				fmt.Println(insurgencylog.ToJSON(message))
				fmt.Println(err)
			}
		}

		switch m := message.(type) {
		case insurgencylog.LoadingMap:
			matchInfo.Map = m.Map
			matchInfo.StartedAt = uint64(m.Time.Unix())
		case insurgencylog.PlayerKill:
			if m.Attacker.SteamID == insurgencylog.PlayerBot && m.Victim.SteamID != insurgencylog.PlayerBot {
				stats, _ := playerStats[m.Victim.SteamID]
				stats.Deaths++
				playerStats[m.Victim.SteamID] = stats
			}
			if m.Victim.SteamID == insurgencylog.PlayerBot && m.Attacker.SteamID != insurgencylog.PlayerBot {
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

	db, err := sql.Open("postgres", os.Getenv("PSQL_CONN"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	sqlStatement := `
INSERT INTO matches (ip, started_at, map, rounds, duration, won)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT(ip, started_at, map) DO UPDATE SET rounds = $4, duration = $5, won = $6
RETURNING id`
	var matchID uint32
	err = db.QueryRow(sqlStatement, matchInfo.Ip, matchInfo.StartedAt, matchInfo.Map, matchInfo.Rounds, matchInfo.Duration, matchInfo.Won).Scan(&matchID)
	if err != nil {
		log.Fatal(err)
	}

	for s, statsStruct := range playerStats {
		userID, err := gosteamconv.SteamStringToInt32(s)
		if err != nil {
			log.Fatal(err)
		}

		err = checkOrCreateUser(db, userID)
		if err != nil {
			log.Fatal(err)
		}

		err = insertUserStats(db, matchID, userID, statsStruct)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Finished processing match ", matchID)
}

func checkOrCreateUser(db *sql.DB, userID int) error {
	userQuery := `SELECT 1 from users where id = $1`
	insertQuery := `INSERT INTO users (id) VALUES ($1)`

	var dummy int
	err := db.QueryRow(userQuery, userID).Scan(&dummy)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_, err = db.Exec(insertQuery, userID)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}

func insertUserStats(db *sql.DB, matchID uint32, userID int, stats playerStatsStruct) error {
	insertQuery := `INSERT INTO match_user_stats (match_id, user_id, kills, deaths, weapon_stats) 
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT(match_id, user_id) DO UPDATE SET kills = $3, deaths = $4, weapon_stats = $5;`

	_, err := db.Exec(insertQuery, matchID, userID, stats.Kills, stats.Deaths, stats.WeaponStats)
	if err != nil {
		return err
	}

	return nil
}
