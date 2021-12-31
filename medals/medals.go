package medals

import (
	"database/sql"
	"errors"
	"log"
	"my.com/insurgency-parser/dbp"
)

const (
	MedalObjectiveMostKillsCurrent = iota + 1
	MedalObjectiveHighestKDCurrent
	MedalObjectiveIWon         // Get 5 wins.
	MedalObjectiveImOnAStreak  // Get 3 wins in a row.
	MedalObjectiveTopFragger   // Get the most kills on your team 5 times.
	MedalObjectiveGoodTeammate // Get a kill/death ratio of over average in 5 matches.
	MedalObjectiveKnifeExpert
	MedalObjectivePistolExpert
	MedalObjectiveBoltExpert
	MedalObjectiveRifleExpert
	MedalObjectiveExplosivesExpert
	MedalObjectiveOneManArmy // complete map alone
	MedalObjectiveDieHard    // don't die
	MedalObjective6MonthsMedal
	MedalObjective1YearMedal
	MedalObjectiveTwoYears
	MedalObjectiveThreeYears
	MedalObjectiveFourYears
	MedalObjectiveCount
)

var medals = []int{
	MedalObjectiveMostKillsCurrent,
	MedalObjectiveHighestKDCurrent,
	MedalObjectiveIWon,
	MedalObjectiveImOnAStreak,
	MedalObjectiveTopFragger,
	MedalObjectiveGoodTeammate,
	MedalObjectiveKnifeExpert,
	MedalObjectivePistolExpert,
	MedalObjectiveBoltExpert,
	MedalObjectiveRifleExpert,
	MedalObjectiveExplosivesExpert,
	MedalObjectiveOneManArmy,
	MedalObjectiveDieHard,
	MedalObjective6MonthsMedal,
	MedalObjective1YearMedal,
	MedalObjectiveTwoYears,
	MedalObjectiveThreeYears,
	MedalObjectiveFourYears,
	MedalObjectiveCount,
}

func UpdateMedals() {
	for _, medal := range medals {
		switch medal {
		/*case MedalObjectiveMostKillsCurrent:
			err := checkMostKills()
			if err != nil {
				log.Fatal(err)
			}
		case MedalObjectiveHighestKDCurrent:

		} */
		case MedalObjectiveIWon:
			err := checkIWon()
			if err != nil {
				log.Fatal(err)
			}
		case MedalObjectiveDieHard:
			err := checkDieHard()
			if err != nil {
				log.Fatal(err)
			}
		case MedalObjectiveKnifeExpert:
			err := checkKnifeExpert()
			if err != nil {
				log.Fatal(err)
			}
		case MedalObjectivePistolExpert:
			err := checkPistolExpert()
			if err != nil {
				log.Fatal(err)
			}
		case MedalObjectiveBoltExpert:
			err := checkBoltExpert()
			if err != nil {
				log.Fatal(err)
			}
		case MedalObjectiveExplosivesExpert:
			err := checkExplosivesExpert()
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func checkExplosivesExpert() error {
	explosExpertsQuery := `
SELECT users.id
from users
         LEFT JOIN user_medals um on users.id = um.user_id AND medal_id = $1
WHERE um.user_id IS NULL
  AND ((all_weapon_stats -> 'grenade_f1')::int + (all_weapon_stats -> 'grenade_ied')::int +
       (all_weapon_stats -> 'grenade_m67')::int + (all_weapon_stats -> 'grenade_c4')::int +
       (all_weapon_stats -> 'grenade_mk2')::int + (all_weapon_stats -> 'mortar_piat')::int +
       (all_weapon_stats -> 'rocket_rpg7')::int + (all_weapon_stats -> 'rocket_at4')::int +
       (all_weapon_stats -> 'grenade_m203_he')::int +
       (all_weapon_stats -> 'grenade_rifle_enfield')::int + (all_weapon_stats -> 'grenade_gp25_he')::int +
       (all_weapon_stats -> 'grenade_rifle_k98')::int + (all_weapon_stats -> 'grenade_gp25_lvg')::int + (all_weapon_stats -> 'm590')::int) >= 1000
`

	err := getIDAndAwardMedal(explosExpertsQuery, MedalObjectiveExplosivesExpert)
	if err != nil {
		return err
	}

	return nil
}

func checkBoltExpert() error {
	boltExpertsQuery := `
SELECT users.id
from users
         LEFT JOIN user_medals um on users.id = um.user_id AND medal_id = $1
WHERE um.user_id IS NULL
  AND ((all_weapon_stats -> 'mosin')::int + (all_weapon_stats -> 'k98')::int + (all_weapon_stats -> 'springfield')::int +
       (all_weapon_stats -> 'enfield')::int + (all_weapon_stats -> 'm1garand')::int) >= 1000
`

	err := getIDAndAwardMedal(boltExpertsQuery, MedalObjectiveBoltExpert)
	if err != nil {
		return err
	}

	return nil
}

func checkPistolExpert() error {
	pistolExpertsQuery := `
SELECT users.id
from users
         LEFT JOIN user_medals um on users.id = um.user_id AND medal_id = $1
WHERE um.user_id IS NULL
  AND ((all_weapon_stats -> 'deagle')::int + (all_weapon_stats -> 'sw500')::int + (all_weapon_stats -> 'makarov')::int +
       (all_weapon_stats -> 'model10')::int + (all_weapon_stats -> 'welrod')::int +
       (all_weapon_stats -> 'browninghp')::int + (all_weapon_stats -> 'sw1917')::int +
       (all_weapon_stats -> 'mr73')::int + (all_weapon_stats -> 'ots33')::int + (all_weapon_stats -> 'glock18')::int) >=
      1000
`

	err := getIDAndAwardMedal(pistolExpertsQuery, MedalObjectivePistolExpert)
	if err != nil {
		return err
	}

	return nil
}

func checkKnifeExpert() error {
	knifeExpertsQuery := `
SELECT users.id
from users
         LEFT JOIN user_medals um on users.id = um.user_id AND medal_id = $1
WHERE um.user_id IS NULL
  AND (all_weapon_stats -> 'gurkha')::int >= 100
`

	err := getIDAndAwardMedal(knifeExpertsQuery, MedalObjectiveKnifeExpert)
	if err != nil {
		return err
	}

	return nil
}

func getIDAndAwardMedal(query string, medal int) error {
	rows, err := dbp.DB.Query(query, medal)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	}
	defer rows.Close()

	userStats := make([]uint32, 0)
	for rows.Next() {
		var ID uint32
		err = rows.Scan(&ID)
		if err != nil {
			return err
		}

		userStats = append(userStats, ID)
	}

	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		return err
	}

	for _, ID := range userStats {
		insertQuery := `INSERT INTO user_medals (user_id, medal_id) VALUES ($1, $2)`
		_, err = dbp.DB.Exec(insertQuery, ID, medal)
		if err != nil {
			return err
		}
	}

	return nil
}

func checkDieHard() error {
	wonMatchesQuery := `
select id, MAX(max_kills)
from (
         SELECT users.id, mus.kills as max_kills
         from users
                  LEFT JOIN match_user_stats mus on users.id = mus.user_id
                  LEFT JOIN matches m on mus.match_id = m.id
         WHERE m.won = true
           AND mus.deaths = 0
           AND mus.kills > 20
         group by users.id, mus.kills) a
GROUP BY id
`
	rows, err := dbp.DB.Query(wonMatchesQuery)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	}
	defer rows.Close()

	type playerStatsStruct struct {
		ID       uint32 `json:"id"`
		MaxKills uint32 `json:"max_kills"`
	}

	userStats := make([]playerStatsStruct, 0)
	for rows.Next() {
		var player playerStatsStruct
		err = rows.Scan(&player.ID, &player.MaxKills)
		if err != nil {
			return err
		}

		userStats = append(userStats, player)
	}

	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		return err
	}

	for _, userStat := range userStats {
		// checking if medal is already awarded
		medalQuery := `SELECT value from user_medals where user_id = $1 AND medal_id = $2`

		var medalKills uint32
		err := dbp.DB.QueryRow(medalQuery, userStat.ID, MedalObjectiveDieHard).Scan(&medalKills)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				insertQuery := `INSERT INTO user_medals (user_id, medal_id, value) VALUES ($1, $2, $3)`
				_, err = dbp.DB.Exec(insertQuery, userStat.ID, MedalObjectiveDieHard, userStat.MaxKills)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		} else {
			if medalKills < userStat.MaxKills {
				updateQuery := `UPDATE user_medals SET value = $1 WHERE user_id = $2 AND medal_id = $3`

				_, err = dbp.DB.Exec(updateQuery, userStat.MaxKills, userStat.ID, MedalObjectiveDieHard)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func checkIWon() error {
	wonMatchesQuery := `
SELECT users.id, COUNT(*)
from users
         LEFT JOIN match_user_stats mus on users.id = mus.user_id
         LEFT JOIN matches m on mus.match_id = m.id
         LEFT JOIN user_medals um on users.id = um.user_id AND medal_id = $1
WHERE m.won = true
  AND um.user_id IS NULL
group by users.id
`
	rows, err := dbp.DB.Query(wonMatchesQuery, MedalObjectiveIWon)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	}
	defer rows.Close()

	type playerStatsStruct struct {
		ID    uint32 `json:"id"`
		Count uint32 `json:"count"`
	}

	userStats := make([]playerStatsStruct, 0)
	for rows.Next() {
		var player playerStatsStruct
		err = rows.Scan(&player.ID, &player.Count)
		if err != nil {
			return err
		}

		if player.Count >= 5 {
			userStats = append(userStats, player)
		}
	}

	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		return err
	}

	for _, userStat := range userStats {
		insertQuery := `INSERT INTO user_medals (user_id, medal_id) VALUES ($1, $2)`
		_, err = dbp.DB.Exec(insertQuery, userStat.ID, MedalObjectiveIWon)
		if err != nil {
			return err
		}
	}

	return nil
}

func checkMostKills() error {
	userQuery := `SELECT id, kills from users where ORDER BY kills DESC LIMIT 1`
	medalQuery := `SELECT user_id, value from user_medals where medal_id = $1 WHERE current = TRUE`
	insertQuery := `INSERT INTO user_medals (user_id, medal_id, value, current) VALUES ($1, $2, $3, $4)`

	var userID, kills int
	err := dbp.DB.QueryRow(userQuery).Scan(&userID, &kills)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	}

	var medalUserID, value int
	err = dbp.DB.QueryRow(medalQuery, MedalObjectiveMostKillsCurrent).Scan(&medalUserID, &value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_, err = dbp.DB.Exec(insertQuery, userID, MedalObjectiveMostKillsCurrent, kills, true)
			if err != nil {
				return err
			}
			return nil
		} else {
			return err
		}
	}

	if medalUserID == userID {
		if value < kills {
			updateQuery := `UPDATE user_medals SET value = $1 WHERE user_id = $2 AND medal_id = $3`

			_, err = dbp.DB.Exec(updateQuery, kills, userID, MedalObjectiveMostKillsCurrent)
			if err != nil {
				return err
			}
		}
	} else {
		if value < kills {
			//check if user has medal
			oldMedalQuery := `SELECT value from user_medals where user_id = $1 AND medal_id = $2`

			var oldValue int
			err = dbp.DB.QueryRow(oldMedalQuery, userID, MedalObjectiveMostKillsCurrent).Scan(&oldValue)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					//adding new medal
					_, err = dbp.DB.Exec(insertQuery, userID, MedalObjectiveMostKillsCurrent, kills, true)
					if err != nil {
						return err
					}
				} else {
					return err
				}
			} else {
				//updating old medal
				updateQuery := `UPDATE user_medals SET value = $1, current = TRUE WHERE user_id = $2 AND medal_id = $3 AND current = FALSE`

				_, err = dbp.DB.Exec(updateQuery, kills, userID, MedalObjectiveMostKillsCurrent)
				if err != nil {
					return err
				}
			}

			updateSecondPlaceQuery := `UPDATE user_medals SET current = FALSE WHERE user_id = $1 AND medal_id = $2 AND current = TRUE`

			_, err = dbp.DB.Exec(updateSecondPlaceQuery, medalUserID, MedalObjectiveMostKillsCurrent)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
