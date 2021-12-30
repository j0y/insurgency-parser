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
	MedalObjectiveIWon         // Get 3 wins.
	MedalObjectiveImOnAStreak  // Get 3 wins in a row.
	MedalObjectiveTopFragger   // Get the most kills on your team 5 times.
	MedalObjectiveGoodTeammate // Get a kill/death ratio of over average in 5 matches.
	MedalObjectiveKnifeExpert
	MedalObjectivePistolExpert
	MedalObjectiveAKExpert
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
	MedalObjectiveAKExpert,
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
		case MedalObjectiveMostKillsCurrent:
			err := checkMostKills()
			if err != nil {
				log.Fatal(err)
			}
		}
	}
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
		}
		return err
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
				}
				return err
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
