CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS pod_monthly_metrics (
  id       uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  pod_name TEXT NOT NULL,
  pod_version TEXT NOT NULL,
  downloads integer,
  apps integer,
  tests integer,
  extensions integer,
  watch integer,
  tries integer,
  rollup_datestamp TEXT, -- e.g 2016-07
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (NOW()),
  updated_at TIMESTAMP WITH TIME ZONE,

  UNIQUE(pod_name, pod_version, rollup_datestamp)
);

