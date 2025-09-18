package jobPGX

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/playture/backend/utils"
	"log/slog"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/playture/backend/internal/entity"
	"github.com/playture/backend/internal/infrastructure/postgresql"
	jobRepository "github.com/playture/backend/internal/repository/job_repository"
)

const (
	CreateQuery = `
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
		) RETURNING id`

	deleteQuery = `DELETE FROM jobs WHERE id = $1`

	findByFieldQuery = `SELECT 
		id, user_email, user_name, input_image_url, input_image_s3_key, style,
		status, veo_video_url, veo_video_s3_key, veo_duration,
		que_job_id, que_job_status, final_video_url, final_video_s3_key,
		final_video_duration, final_video_size, signed_url, signed_url_expiry,
		email_sent, email_sent_at, error_message, error_stack, retry_count,
		ip_address, user_agent, started_at, completed_at, total_processing_time,
		converted_to_order, order_id, content_moderated, content_moderation_result,
		created_at, updated_at
		FROM jobs WHERE %s = $1 LIMIT 1`

	updateQuery = `
		UPDATE jobs SET
			user_email=$2, user_name=$3, input_image_url=$4, input_image_s3_key=$5, style=$6,
			status=$7, veo_video_url=$8, veo_video_s3_key=$9, veo_duration=$10,
			que_job_id=$11, que_job_status=$12, final_video_url=$13, final_video_s3_key=$14,
			final_video_duration=$15, final_video_size=$16, signed_url=$17, signed_url_expiry=$18,
			email_sent=$19, email_sent_at=$20, error_message=$21, error_stack=$22, retry_count=$23,
			ip_address=$24, user_agent=$25, started_at=$26, completed_at=$27, total_processing_time=$28,
			converted_to_order=$29, order_id=$30, content_moderated=$31, content_moderation_result=$32,
			updated_at=$33
		WHERE id=$1`

	listQuery = `
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
		LIMIT $%d OFFSET $%d`
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
	lg := j.logger.With("method", "Create")
	var id string

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
		err = tx.QueryRow(ctx, CreateQuery, args...).Scan(&id)
	} else {
		err = j.postgres.PrimaryConn.QueryRow(ctx, CreateQuery, args...).Scan(&id)
	}
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			lg.Error("Create failed", "pgErr", pgErr.Message)
		}
		return "", utils.WrapError("create job", err)
	}
	return id, nil
}

func (j *JobPgx) Delete(
	ctx context.Context,
	id string,
	tx pgx.Tx,
) error {
	lg := j.logger.With("method", "Delete")
	var err error
	if tx != nil {
		_, err = tx.Exec(ctx, deleteQuery, id)
	} else {
		_, err = j.postgres.PrimaryConn.Exec(ctx, deleteQuery, id)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return jobRepository.ErrJobNotFound
		}
		lg.Error("Delete failed", "id", id, "err", err)
		return utils.WrapError("delete job", err)
	}
	return nil
}

func (j *JobPgx) FindByField(
	ctx context.Context,
	field string,
	value interface{},
	tx pgx.Tx,
) (*entity.Job, error) {
	lg := j.logger.With("method", "FindByField")
	query := fmt.Sprintf(findByFieldQuery, field)

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
			return nil, utils.WrapError("FindByFiled", jobRepository.ErrJobNotFound)
		}
		lg.Error("FindByField failed", "err", err)
		return nil, utils.WrapError("FindByField job", err)
	}
	return job, nil
}
func (j *JobPgx) Update(
	ctx context.Context,
	job *entity.Job,
	tx pgx.Tx,
) error {
	lg := j.logger.With("method", "Update")

	args := []interface{}{
		job.ID, job.UserEmail, job.UserName, job.InputImageURL, job.InputImageS3Key, job.Style,
		job.Status, job.VeoVideoURL, job.VeoVideoS3Key, job.VeoDuration,
		job.QueJobID, job.QueJobStatus, job.FinalVideoURL, job.FinalVideoS3Key,
		job.FinalVideoDuration, job.FinalVideoSize, job.SignedURL, job.SignedURLExpiry,
		job.EmailSent, job.EmailSentAt, job.ErrorMessage, job.ErrorStack, job.RetryCount,
		job.IPAddress, job.UserAgent, job.StartedAt, job.CompletedAt, job.TotalProcessingTime,
		job.ConvertedToOrder, job.OrderID, job.ContentModerated, job.ContentModerationResult,
		job.UpdatedAt,
	}

	var (
		cmd pgconn.CommandTag
		err error
	)

	if tx != nil {
		cmd, err = tx.Exec(ctx, updateQuery, args...)
	} else {
		cmd, err = j.postgres.PrimaryConn.Exec(ctx, updateQuery, args...)
	}
	if err != nil {
		lg.Error("Update failed", "id", job.ID, "err", err)
		return utils.WrapError("update job", err)
	}
	if cmd.RowsAffected() == 0 {
		return utils.WrapError("update job", jobRepository.ErrJobNotFound)
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

	if orderBy == "" {
		orderBy = "created_at"
	}
	direction := "ASC"
	if !ascending {
		direction = "DESC"
	}

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

	offset := (page - 1) * limit
	args = append(args, limit, offset)

	query := fmt.Sprintf(listQuery, where, orderBy, direction, len(args)-1, len(args))

	var rows pgx.Rows
	var err error
	if tx != nil {
		rows, err = tx.Query(ctx, query, args...)
	} else {
		rows, err = j.postgres.PrimaryConn.Query(ctx, query, args...)
	}
	if err != nil {
		return nil, utils.WrapError("failed to Query against database", err)
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
			return nil, utils.WrapError("failed to scan row", err)
		}
		jobs = append(jobs, job)
	}

	if len(jobs) == 0 {
		return nil, utils.WrapError("list", jobRepository.ErrJobNotFound)
	}

	lg.Info("fetched jobs", "count", len(jobs))
	return jobs, nil
}
