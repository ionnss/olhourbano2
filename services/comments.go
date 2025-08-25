package services

import (
	"database/sql"
	"fmt"
	"olhourbano2/models"
)

// CreateComment creates a new comment on a report
func CreateComment(db *sql.DB, reportID int, hashedCPF, content string) (*models.Comment, error) {
	// Validate content length
	if len(content) > 500 {
		return nil, fmt.Errorf("comment content exceeds 500 character limit")
	}

	// Insert the comment
	var comment models.Comment
	err := db.QueryRow(`
		INSERT INTO comments (report_id, hashed_cpf, content, created_at)
		VALUES ($1, $2, $3, NOW())
		RETURNING id, report_id, hashed_cpf, content, created_at
	`, reportID, hashedCPF, content).Scan(
		&comment.ID, &comment.ReportID, &comment.HashedCPF, &comment.Content, &comment.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("error creating comment: %w", err)
	}

	// Update comment count in reports table
	_, err = db.Exec(`
		UPDATE reports 
		SET comment_count = (
			SELECT COUNT(*) 
			FROM comments 
			WHERE report_id = $1
		)
		WHERE id = $1
	`, reportID)

	if err != nil {
		return nil, fmt.Errorf("error updating comment count: %w", err)
	}

	// Send email notification to report owner (async)
	go sendCommentNotificationEmail(db, reportID, hashedCPF, content)

	return &comment, nil
}

// sendCommentNotificationEmail sends email notification to report owner
func sendCommentNotificationEmail(db *sql.DB, reportID int, commenterHashedCPF, commentContent string) {
	// Get report owner's email
	var reportOwnerEmail string
	err := db.QueryRow(`
		SELECT email 
		FROM reports 
		WHERE id = $1
	`, reportID).Scan(&reportOwnerEmail)

	if err != nil {
		// Log error but don't fail the comment creation
		fmt.Printf("Error getting report owner email for notification: %v\n", err)
		return
	}

	// Don't send notification if no email or if commenter is the report owner
	if reportOwnerEmail == "" {
		return
	}

	// Check if commenter is the report owner (to avoid self-notifications)
	var reportOwnerHashedCPF string
	err = db.QueryRow(`
		SELECT hashed_cpf 
		FROM reports 
		WHERE id = $1
	`, reportID).Scan(&reportOwnerHashedCPF)

	if err != nil {
		fmt.Printf("Error getting report owner CPF for notification: %v\n", err)
		return
	}

	// Don't send notification if commenter is the report owner
	if commenterHashedCPF == reportOwnerHashedCPF {
		return
	}

	// Get commenter display name (first 8 characters of hashed CPF)
	commenterName := commenterHashedCPF
	if len(commenterHashedCPF) >= 8 {
		commenterName = commenterHashedCPF[:8]
	}

	// Truncate comment content for email (max 100 characters)
	emailCommentContent := commentContent
	if len(commentContent) > 100 {
		emailCommentContent = commentContent[:97] + "..."
	}

	// Send the notification email
	SendCommentNotificationEmail(reportOwnerEmail, reportID, commenterName, emailCommentContent)
}

// GetCommentsForReport retrieves comments for a specific report
func GetCommentsForReport(db *sql.DB, reportID int, sort string, limit int, offset int) ([]*models.CommentDisplay, error) {
	query := `
		SELECT c.id, c.report_id, c.content, c.created_at, c.hashed_cpf
		FROM comments c
		WHERE c.report_id = $1
		ORDER BY c.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := db.Query(query, reportID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error querying comments: %w", err)
	}
	defer rows.Close()

	var comments []*models.CommentDisplay
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(
			&comment.ID, &comment.ReportID, &comment.Content, &comment.CreatedAt, &comment.HashedCPF,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning comment: %w", err)
		}

		commentDisplay := &models.CommentDisplay{
			ID:               comment.ID,
			ReportID:         comment.ReportID,
			Content:          comment.Content,
			CreatedAt:        comment.CreatedAt,
			HashedCPFDisplay: comment.GetHashedCPFDisplay(),
		}

		comments = append(comments, commentDisplay)
	}

	return comments, nil
}

// GetCommentCountForReport returns the total number of comments for a report
func GetCommentCountForReport(db *sql.DB, reportID int) (int, error) {
	var count int
	err := db.QueryRow(`
		SELECT COUNT(*) 
		FROM comments 
		WHERE report_id = $1
	`, reportID).Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("error getting comment count: %w", err)
	}

	return count, nil
}

// UpdateAllReportCommentCounts updates comment counts for all reports
func UpdateAllReportCommentCounts(db *sql.DB) error {
	_, err := db.Exec(`
		UPDATE reports 
		SET comment_count = (
			SELECT COUNT(*) 
			FROM comments 
			WHERE comments.report_id = reports.id
		)
	`)
	return err
}
