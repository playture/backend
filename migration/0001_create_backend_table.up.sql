CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE jobs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_email TEXT NOT NULL,
    user_name TEXT NOT NULL,
    input_image_url TEXT NOT NULL,
    input_image_s3_key TEXT NOT NULL,
    style TEXT DEFAULT 'default',

    -- Processing status 
    status SMALLINT NOT NULL DEFAULT 1,

    -- Google Veo 3 output
    veo_video_url TEXT,
    veo_video_s3_key TEXT,
    veo_duration INT,

    -- Dataclay QUE integration
    que_job_id TEXT,
    que_job_status TEXT,

    -- Final output
    final_video_url TEXT,
    final_video_s3_key TEXT,
    final_video_duration INT,
    final_video_size BIGINT,

    -- CloudFront signed URL
    signed_url TEXT,
    signed_url_expiry BIGINT,

    -- Email delivery
    email_sent BOOLEAN DEFAULT FALSE,
    email_sent_at BIGINT,

    -- Error tracking
    error_message TEXT,
    error_stack TEXT,
    retry_count INT DEFAULT 0,

    -- Metadata
    ip_address TEXT,
    user_agent TEXT,

    -- Timing
    started_at BIGINT,
    completed_at BIGINT,
    total_processing_time BIGINT,

    -- Convert to order
    converted_to_order BOOLEAN DEFAULT FALSE,
    order_id UUID,

    -- Content moderation
    content_moderated BOOLEAN DEFAULT FALSE,
    content_moderation_result JSONB,

    -- Metadata
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL
);

