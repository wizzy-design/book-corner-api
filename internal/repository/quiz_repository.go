package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wizzy-design/book-corner-api/internal/models"
)

func GetQuizQuestions(pool *pgxpool.Pool) ([]models.QuizQuestion, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	questionQuery := `
		SELECT id, prompt, position
		FROM quiz_questions
		ORDER BY position ASC
	`

	// Fetch all questions ordered by position
	rows, err := pool.Query(ctx, questionQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	questions := []models.QuizQuestion{}
	questionIndexMap := make(map[int]int)

	for rows.Next() {
		var q models.QuizQuestion
		if err := rows.Scan(&q.ID, &q.Prompt, &q.Position);
		err != nil {
			return nil, err
		}

		// Initialize Options to an empty slice so JSON outputs [] instead of null
		q.Options = []models.QuizOption{}
		questionIndexMap[q.ID] = len(questions)
		questions = append(questions, q)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	optionQuery := `
		SELECT id, question_id, label
		FROM quiz_options
		ORDER BY id ASC
	`

	optRows, err := pool.Query(ctx, optionQuery)
	if err != nil {
		return nil, err
	}
	defer optRows.Close()

	for optRows.Next(){
		var opt models.QuizOption
		if err := optRows.Scan(&opt.ID, &opt.QuestionID, &opt.Label); err != nil {
			return nil, err
		}

		if idx, ok := questionIndexMap[opt.QuestionID]; ok {
			questions[idx].Options = append(questions[idx].Options, opt)
		}
	}

	if err := optRows.Err(); err != nil {
		return nil, err
	}

	return questions, nil
}