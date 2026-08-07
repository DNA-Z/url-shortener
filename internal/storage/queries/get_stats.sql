SELECT
    COUNT(*) AS urls_count,
    COUNT(DISTINCT user_id) AS users_count
FROM urls
WHERE is_deleted = false;