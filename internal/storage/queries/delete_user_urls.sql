UPDATE urls
SET is_deleted = true
WHERE user_id = $1
    and short_url = in ANY($2)
    and is_deleted is false;