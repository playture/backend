package entity

import (
	"github.com/google/uuid"
)

type JobStatus uint8

const (
	JobStatusReceived      JobStatus = 1
	JobStatusProcessing    JobStatus = 2
	JobStatusVeoGenerating JobStatus = 3
	JobStatusVeoCompleted  JobStatus = 4
	JobStatusQueProcessing JobStatus = 5
	JobStatusRendering     JobStatus = 6
	JobStatusCompleted     JobStatus = 7
	JobStatusFailed        JobStatus = 8
	JobStatusCancelled     JobStatus = 9
)

func (j JobStatus) String() string {
	switch j {
	case JobStatusReceived:
		return "RECEIVED"
	case JobStatusProcessing:
		return "PROCESSING"
	case JobStatusVeoGenerating:
		return "VEO-GENERATING"
	case JobStatusVeoCompleted:
		return "VEO-COMPLETED"
	case JobStatusQueProcessing:
		return "QUE-PROCESSING"
	case JobStatusRendering:
		return "RENDERING"
	case JobStatusCompleted:
		return "COMPLETED"
	case JobStatusFailed:
		return "FAILED"
	case JobStatusCancelled:
		return "CANCELLED"
	default:
		return "UNKNOWN"
	}
}

type Job struct {
	ID                      uuid.UUID   `json:"id" bson:"_id"`
	UserEmail               string      `json:"userEmail" bson:"userEmail"`
	UserName                string      `json:"userName" bson:"userName"`
	InputImageURL           string      `json:"inputImageUrl" bson:"inputImageUrl"`
	InputImageS3Key         string      `json:"inputImageS3Key" bson:"inputImageS3Key"`
	Style                   string      `json:"style" bson:"style"`
	Status                  JobStatus   `json:"status" bson:"status"`
	VeoVideoURL             string      `json:"veoVideoUrl,omitempty" bson:"veoVideoUrl,omitempty"`
	VeoVideoS3Key           string      `json:"veoVideoS3Key,omitempty" bson:"veoVideoS3Key,omitempty"`
	VeoDuration             int         `json:"veoDuration,omitempty" bson:"veoDuration,omitempty"`
	QueJobID                string      `json:"queJobId,omitempty" bson:"queJobId,omitempty"`
	QueJobStatus            string      `json:"queJobStatus,omitempty" bson:"queJobStatus,omitempty"`
	FinalVideoURL           string      `json:"finalVideoUrl,omitempty" bson:"finalVideoUrl,omitempty"`
	FinalVideoS3Key         string      `json:"finalVideoS3Key,omitempty" bson:"finalVideoS3Key,omitempty"`
	FinalVideoDuration      int         `json:"finalVideoDuration,omitempty" bson:"finalVideoDuration,omitempty"`
	FinalVideoSize          int64       `json:"finalVideoSize,omitempty" bson:"finalVideoSize,omitempty"`
	SignedURL               string      `json:"signedUrl,omitempty" bson:"signedUrl,omitempty"`
	SignedURLExpiry         int64       `json:"signedUrlExpiry,omitempty" bson:"signedUrlExpiry,omitempty"`
	EmailSent               bool        `json:"emailSent" bson:"emailSent"`
	EmailSentAt             int64       `json:"emailSentAt,omitempty" bson:"emailSentAt,omitempty"`
	ErrorMessage            string      `json:"errorMessage,omitempty" bson:"errorMessage,omitempty"`
	ErrorStack              string      `json:"errorStack,omitempty" bson:"errorStack,omitempty"`
	RetryCount              int         `json:"retryCount" bson:"retryCount"`
	IPAddress               string      `json:"ipAddress,omitempty" bson:"ipAddress,omitempty"`
	UserAgent               string      `json:"userAgent,omitempty" bson:"userAgent,omitempty"`
	StartedAt               int64       `json:"startedAt,omitempty" bson:"startedAt,omitempty"`
	CompletedAt             int64       `json:"completedAt,omitempty" bson:"completedAt,omitempty"`
	TotalProcessingTime     int64       `json:"totalProcessingTime,omitempty" bson:"totalProcessingTime,omitempty"`
	ConvertedToOrder        bool        `json:"convertedToOrder" bson:"convertedToOrder"`
	OrderID                 *uuid.UUID  `json:"orderId,omitempty" bson:"orderId,omitempty"`
	ContentModerated        bool        `json:"contentModerated" bson:"contentModerated"`
	ContentModerationResult interface{} `json:"contentModerationResult,omitempty" bson:"contentModerationResult,omitempty"`
	CreatedAt               int64       `json:"createdAt" bson:"createdAt"`
	UpdatedAt               int64       `json:"updatedAt" bson:"updatedAt"`
}
