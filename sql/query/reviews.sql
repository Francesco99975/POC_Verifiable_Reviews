-- name: GetAllReviews :many
SELECT id, content, created
FROM reviews
ORDER BY created DESC;

-- name: GetReviewByID :one
SELECT id, content, created
FROM reviews
WHERE id = $1;

-- name: CreateReview :one
INSERT INTO reviews (id, content)
VALUES ($1, $2)
RETURNING id, content, created;

-- name: DeleteReview :execrows
DELETE FROM reviews
WHERE id = $1;

-- name: CountReviews :one
SELECT COUNT(*) FROM reviews;
