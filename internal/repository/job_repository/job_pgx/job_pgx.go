package jobPGX

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/playture/backend/internal/entity"
	"github.com/playture/backend/internal/infrastructure/postgresql"
	jobRepository "github.com/playture/backend/internal/repository/job_repository"
)

const (
	insertJob = `
		INSERT INTO jobs (
			user_email, user_name, input_image_url, input_image_s3_key, style,
			status, veo_video_url, veo_video_s3_key, veo_duration,
			que_job_id, que_job_status, final_video_url, final_video_s3_key,
			final_video_duration, final_video_size, signed_url, signed_url_expiry,
			email_sent, email_sent_at, error_message, error_stack, retry_count,
			ip_address, user_agent, started_at, completed_at, total_processing_time,
			converted_to_order, order_id, content_moderated, content_moderation_result,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9,
			$10, $11, $12, $13,
			$14, $15, $16, $17,
			$18, $19, $20, $21, $22,
			$23, $24, $25, $26, $27,
			$28, $29, $30, $31,
			$32, $33
		) RETURNING id
	`

	deleteJob = `DELETE FROM jobs WHERE id = $1`
)

type JobPgx struct {
	logger   *slog.Logger
	postgres *postgresql.Postgres
}

func NewJobPgx(
	logger *slog.Logger,
	postgres *postgresql.Postgres,
) *JobPgx {
	return &JobPgx{
		logger:   logger.With("layer", "JobRepository"),
		postgres: postgres,
	}
}

func (j *JobPgx) Create(
	ctx context.Context,
	job *entity.Job,
	tx pgx.Tx,
) (string, error) {
	var id string
	query := insertJob

	args := []interface{}{
		job.UserEmail, job.UserName, job.InputImageURL, job.InputImageS3Key, job.Style,
		job.Status, job.VeoVideoURL, job.VeoVideoS3Key, job.VeoDuration,
		job.QueJobID, job.QueJobStatus, job.FinalVideoURL, job.FinalVideoS3Key,
		job.FinalVideoDuration, job.FinalVideoSize, job.SignedURL, job.SignedURLExpiry,
		job.EmailSent, job.EmailSentAt, job.ErrorMessage, job.ErrorStack, job.RetryCount,
		job.IPAddress, job.UserAgent, job.StartedAt, job.CompletedAt, job.TotalProcessingTime,
		job.ConvertedToOrder, job.OrderID, job.ContentModerated, job.ContentModerationResult,
		job.CreatedAt, job.UpdatedAt,
	}

	var err error
	if tx != nil {
		err = tx.QueryRow(ctx, query, args...).Scan(&id)
	} else {
		err = j.postgres.PrimaryConn.QueryRow(ctx, query, args...).Scan(&id)
	}
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			j.logger.Error("Create failed", "pgErr", pgErr.Message)
		}
		return "", err
	}
	return id, nil
}

func (j *JobPgx) Delete(
	ctx context.Context,
	id string,
	tx pgx.Tx,
) error {
	var err error
	if tx != nil {
		_, err = tx.Exec(ctx, deleteJob, id)
	} else {
		_, err = j.postgres.PrimaryConn.Exec(ctx, deleteJob, id)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return jobRepository.ErrJobNotFound
		}
		return err
	}
	return nil
}

func (j *JobPgx) FindByField(
	ctx context.Context,
	field string,
	value interface{},
	tx pgx.Tx,
) (*entity.Job, error) {
	query := fmt.Sprintf(`SELECT 
		id, user_email, user_name, input_image_url, input_image_s3_key, style,
		status, veo_video_url, veo_video_s3_key, veo_duration,
		que_job_id, que_job_status, final_video_url, final_video_s3_key,
		final_video_duration, final_video_size, signed_url, signed_url_expiry,
		email_sent, email_sent_at, error_message, error_stack, retry_count,
		ip_address, user_agent, started_at, completed_at, total_processing_time,
		converted_to_order, order_id, content_moderated, content_moderation_result,
		created_at, updated_at
		FROM jobs WHERE %s = $1 LIMIT 1`, field)

	var row pgx.Row
	if tx != nil {
		row = tx.QueryRow(ctx, query, value)
	} else {
		row = j.postgres.PrimaryConn.QueryRow(ctx, query, value)
	}

	job := &entity.Job{}
	err := row.Scan(
		&job.ID, &job.UserEmail, &job.UserName, &job.InputImageURL, &job.InputImageS3Key, &job.Style,
		&job.Status, &job.VeoVideoURL, &job.VeoVideoS3Key, &job.VeoDuration,
		&job.QueJobID, &job.QueJobStatus, &job.FinalVideoURL, &job.FinalVideoS3Key,
		&job.FinalVideoDuration, &job.FinalVideoSize, &job.SignedURL, &job.SignedURLExpiry,
		&job.EmailSent, &job.EmailSentAt, &job.ErrorMessage, &job.ErrorStack, &job.RetryCount,
		&job.IPAddress, &job.UserAgent, &job.StartedAt, &job.CompletedAt, &job.TotalProcessingTime,
		&job.ConvertedToOrder, &job.OrderID, &job.ContentModerated, &job.ContentModerationResult,
		&job.CreatedAt, &job.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, jobRepository.ErrJobNotFound
		}
		return nil, err
	}
	return job, nil
}

func (j *JobPgx) Update(
	ctx context.Context,
	job *entity.Job,
	tx pgx.Tx,
) error {
	// Collect fields dynamically
	fields := []string{}
	args := []interface{}{}
	i := 1

	addField := func(name string, value interface{}) {
		if value != nil && value != "" {
			fields = append(fields, fmt.Sprintf("%s=$%d", name, i))
			args = append(args, value)
			i++
		}
	}

	addField("status", job.Status)
	addField("veo_video_url", job.VeoVideoURL)
	addField("veo_video_s3_key", job.VeoVideoS3Key)
	addField("final_video_url", job.FinalVideoURL)
	addField("final_video_s3_key", job.FinalVideoS3Key)
	addField("final_video_duration", job.FinalVideoDuration)
	addField("final_video_size", job.FinalVideoSize)
	addField("error_message", job.ErrorMessage)
	addField("error_stack", job.ErrorStack)
	addField("updated_at", job.UpdatedAt)

	if len(fields) == 0 {
		return nil
	}

	query := fmt.Sprintf("UPDATE jobs SET %s WHERE id=$%d", strings.Join(fields, ", "), i)
	args = append(args, job.ID)

	var err error
	if tx != nil {
		_, err = tx.Exec(ctx, query, args...)
	} else {
		_, err = j.postgres.PrimaryConn.Exec(ctx, query, args...)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return jobRepository.ErrJobNotFound
		}
		return err
	}
	return nil
}

func (j *JobPgx) List(ctx context.Context,
	statuses []entity.JobStatus,
	orderBy string,
	ascending bool,
	limit,
	page int,
	tx pgx.Tx,
) ([]*entity.Job, error) {
	lg := j.logger.With("method", "List")

	// Default order
	if orderBy == "" {
		orderBy = "created_at"
	}
	direction := "ASC"
	if !ascending {
		direction = "DESC"
	}

	// Build WHERE clause
	where := ""
	args := []interface{}{}
	if len(statuses) > 0 {
		placeholders := []string{}
		for i, st := range statuses {
			placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
			args = append(args, st)
		}
		where = fmt.Sprintf("WHERE status IN (%s)", strings.Join(placeholders, ", "))
	}

	// Pagination
	offset := (page - 1) * limit
	args = append(args, limit, offset)

	query := fmt.Sprintf(`
		SELECT 
			id, user_email, user_name, input_image_url, input_image_s3_key, style,
			status, veo_video_url, veo_video_s3_key, veo_duration,
			que_job_id, que_job_status, final_video_url, final_video_s3_key,
			final_video_duration, final_video_size, signed_url, signed_url_expiry,
			email_sent, email_sent_at, error_message, error_stack, retry_count,
			ip_address, user_agent, started_at, completed_at, total_processing_time,
			converted_to_order, order_id, content_moderated, content_moderation_result,
			created_at, updated_at
		FROM jobs
		%s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d
	`, where, orderBy, direction, len(args)-1, len(args))

	var rows pgx.Rows
	var err error
	if tx != nil {
		rows, err = tx.Query(ctx, query, args...)
	} else {
		rows, err = j.postgres.PrimaryConn.Query(ctx, query, args...)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []*entity.Job
	for rows.Next() {
		job := &entity.Job{}
		err = rows.Scan(
			&job.ID, &job.UserEmail, &job.UserName, &job.InputImageURL, &job.InputImageS3Key, &job.Style,
			&job.Status, &job.VeoVideoURL, &job.VeoVideoS3Key, &job.VeoDuration,
			&job.QueJobID, &job.QueJobStatus, &job.FinalVideoURL, &job.FinalVideoS3Key,
			&job.FinalVideoDuration, &job.FinalVideoSize, &job.SignedURL, &job.SignedURLExpiry,
			&job.EmailSent, &job.EmailSentAt, &job.ErrorMessage, &job.ErrorStack, &job.RetryCount,
			&job.IPAddress, &job.UserAgent, &job.StartedAt, &job.CompletedAt, &job.TotalProcessingTime,
			&job.ConvertedToOrder, &job.OrderID, &job.ContentModerated, &job.ContentModerationResult,
			&job.CreatedAt, &job.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}

	if len(jobs) == 0 {
		return nil, jobRepository.ErrJobNotFound
	}

	lg.Info("fetched jobs", "count", len(jobs))
	return jobs, nil
}
