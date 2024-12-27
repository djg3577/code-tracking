package repository

import (
	"database/sql"
	"fmt"
	"strconv"
)

type LeaderBoardRepository struct {
	DB *sql.DB
}

type UserScore struct {
	Username string `json:"username"`
	Score    int    `json:"score"`
}

func (r *LeaderBoardRepository) GetTopScores(days string) ([]UserScore, error) {
	query := `
		SELECT username, SUM(total_duration) as total_score
		FROM Users
		JOIN activities ON activities.user_id = Users.id
	`
	if days != "" { 
		daysInt, err := strconv.Atoi(days)
		if err != nil {
			return nil, fmt.Errorf("invalid days parameter: %v", err)
		}
		query += fmt.Sprintf(`
				WHERE activities.date >= CURRENT_DATE - INTERVAL '%d days'
		`, daysInt)
	}
// testing push
	query += `
				GROUP BY username
				ORDER BY total_score DESC
	`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topScores []UserScore
	for rows.Next() {
		var score UserScore
		if err := rows.Scan(&score.Username, &score.Score); err != nil {
			return nil, err
		}
		topScores = append(topScores, score)
	}

	return topScores, nil
}
