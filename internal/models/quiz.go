package models

type QuizOption struct {
	ID         int    `json:"id" db:"id"`
	QuestionID int    `json:"question_id" db:"question_id"`
	Label      string `json:"label" db:"label"`
}

type QuizQuestion struct {
	ID       int          `json:"id" db:"id"`
	Prompt   string       `json:"prompt" db:"prompt"`
	Position int          `json:"position" db:"position"`
	Options  []QuizOption `json:"options"`
}

type QuizAnswer struct {
	QuestionID int `json:"question_id" binding:"required"`
	OptionID   int `json:"option_id" binding:"required"`
}

type QuizRecommendRequest struct {
	Answers []QuizAnswer `json:"answers" binding:"required,min=1"`
}