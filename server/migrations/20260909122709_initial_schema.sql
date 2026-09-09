-- +goose Up

CREATE TABLE universities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    country TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (name, country)
);

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL,
    university_id UUID NOT NULL REFERENCES universities(id),

    name TEXT NOT NULL
        CHECK (char_length(name) BETWEEN 1 AND 99),

    username TEXT NOT NULL
        CHECK (char_length(username) BETWEEN 1 AND 29),

    role TEXT NOT NULL DEFAULT 'user'
        CHECK (role IN ('user', 'admin')),

    is_email_verified BOOLEAN NOT NULL DEFAULT false,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX users_email_lower_unique
    ON users (lower(email));

CREATE UNIQUE INDEX users_username_lower_unique
    ON users (lower(username));

CREATE TABLE modules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    university_id UUID NOT NULL REFERENCES universities(id),
    code TEXT NOT NULL,
    title TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (university_id, code)
);

CREATE TABLE submissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE mappings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    home_module_id UUID NOT NULL REFERENCES modules(id),
    partner_module_id UUID NOT NULL REFERENCES modules(id),

    academic_year SMALLINT NOT NULL
        CHECK (academic_year BETWEEN 1990 AND 2100),

    semester SMALLINT NOT NULL
        CHECK (semester IN (1, 2)),

    partner_module_link TEXT,
    partner_module_desc TEXT,

    university_decision TEXT NOT NULL
        CHECK (university_decision IN ('approved', 'rejected')),

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CHECK (home_module_id <> partner_module_id),

    UNIQUE (
        home_module_id,
        partner_module_id,
        academic_year,
        semester,
        university_decision
    )
);

CREATE TABLE submission_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    submission_id UUID NOT NULL REFERENCES submissions(id),
    home_module_id UUID NOT NULL REFERENCES modules(id),
    partner_module_id UUID NOT NULL REFERENCES modules(id),

    academic_year SMALLINT NOT NULL
        CHECK (academic_year BETWEEN 1990 AND 2100),

    semester SMALLINT NOT NULL
        CHECK (semester IN (1, 2)),

    partner_module_link TEXT,
    partner_module_desc TEXT,

    university_decision TEXT NOT NULL
        CHECK (university_decision IN ('approved', 'rejected')),

    review_status TEXT NOT NULL DEFAULT 'pending'
        CHECK (review_status IN ('pending', 'approved', 'rejected')),

    mapping_id UUID REFERENCES mappings(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CHECK (home_module_id <> partner_module_id),

    CHECK (
        (review_status = 'approved' AND mapping_id IS NOT NULL)
        OR
        (review_status IN ('pending', 'rejected') AND mapping_id IS NULL)
    )
);

CREATE TABLE submission_evidence (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    submission_id UUID NOT NULL REFERENCES submissions(id),
    object_key TEXT NOT NULL UNIQUE,
    original_filename TEXT NOT NULL,

    file_type TEXT NOT NULL CHECK (
        file_type IN (
            'image/png',
            'image/jpeg',
            'image/webp',
            'application/pdf'
        )
    ),

    size_bytes BIGINT NOT NULL CHECK (size_bytes > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    submission_item_id UUID NOT NULL REFERENCES submission_items(id),
    reviewer_id UUID NOT NULL REFERENCES users(id),

    decision TEXT NOT NULL
        CHECK (decision IN ('approved', 'rejected')),

    reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX submissions_user_created_idx
    ON submissions (user_id, created_at DESC);

CREATE INDEX submission_items_submission_idx
    ON submission_items (submission_id);

CREATE INDEX submission_items_mapping_idx
    ON submission_items (mapping_id)
    WHERE mapping_id IS NOT NULL;

CREATE INDEX submission_items_pending_idx
    ON submission_items (created_at)
    WHERE review_status = 'pending';

CREATE INDEX submission_evidence_submission_idx
    ON submission_evidence (submission_id);

CREATE INDEX mappings_partner_module_idx
    ON mappings (partner_module_id);

-- +goose Down

DROP TABLE reviews;
DROP TABLE submission_evidence;
DROP TABLE submission_items;
DROP TABLE mappings;
DROP TABLE submissions;
DROP TABLE modules;
DROP TABLE users;
DROP TABLE universities;
