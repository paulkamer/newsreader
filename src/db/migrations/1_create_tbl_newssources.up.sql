CREATE TABLE newssources (
  id GUID PRIMARY KEY,
  title TEXT NOT NULL,
  url TEXT NOT NULL,
  update_priority TEXT CHECK( update_priority IN ('URGENT','HIGH','MED', 'LOW') ) NOT NULL DEFAULT 'MED',
  feed_type TEXT CHECK( feed_type IN ('rss','atom') ) NOT NULL,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME
);