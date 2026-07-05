package database

-- Enable UUID extension for generating secure, unique identifiers
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- =========================================================================
-- 1. Users Table (Production-ready with UUID and Timestamps)
-- =========================================================================
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),          -- Secure, non-guessable primary key
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,       -- TIMESTAMPTZ automatically handles timezone offsets
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- =========================================================================
-- 2. Smart Tasks Table (Optimized for AI-extracted structured data)
-- =========================================================================
CREATE TABLE tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    raw_text TEXT NOT NULL,                                 -- Original unstructured text input from the user
    title VARCHAR(255) NOT NULL,                            -- AI-extracted task title
    description TEXT,                                       -- AI-extracted task description
    category VARCHAR(50),                                   -- Category assigned by AI (e.g., Work, Personal, Study)
    priority VARCHAR(20) NOT NULL DEFAULT 'medium',          -- Priority level assigned by AI (low, medium, high)
    due_date TIMESTAMPTZ,                                   -- AI-extracted deadline/due date
    sub_tasks JSONB DEFAULT '[]'::jsonb,                    -- Actionable breakdown steps stored efficiently as binary JSON
    is_completed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

    -- Data integrity constraint to ensure priority values are strictly validated
    CONSTRAINT check_priority CHECK (priority IN ('low', 'medium', 'high'))
);

-- =========================================================================
-- 3. Professional Indexing Strategy for Maximum Query Performance
-- =========================================================================

-- Foreign Key Index: Crucial for optimizing lookups of a specific user's tasks
CREATE INDEX idx_tasks_user_id ON tasks(user_id);

-- Composite Index: Optimizes the highly frequent query "Get all active/completed tasks for User X"
CREATE INDEX idx_tasks_user_status ON tasks(user_id, is_completed);

-- Partial Index: Speeds up sorting active tasks by due date (nearest deadline first)
-- Excludes completed tasks since their chronological order is rarely queried on active dashboards
CREATE INDEX idx_tasks_active_due_date ON tasks(due_date) WHERE is_completed = FALSE;

-- GIN (Generalized Inverted Index): Enables ultra-fast querying and nested filtering inside the JSONB sub_tasks array
CREATE INDEX idx_tasks_sub_tasks_gin ON tasks USING gin (sub_tasks);
