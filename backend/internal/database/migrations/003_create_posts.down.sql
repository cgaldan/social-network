DROP INDEX IF EXISTS idx_comments_created_at;
DROP INDEX IF EXISTS idx_comments_user_id;
DROP INDEX IF EXISTS idx_comments_post_id;

DROP INDEX IF EXISTS idx_post_audiences_user_id;
DROP INDEX IF EXISTS idx_post_audiences_post_id;

DROP INDEX IF EXISTS idx_posts_privacy_level;
DROP INDEX IF EXISTS idx_posts_created_at;
DROP INDEX IF EXISTS idx_posts_category;
DROP INDEX IF EXISTS idx_posts_user_id;

DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS post_audiences;
DROP TABLE IF EXISTS posts;
