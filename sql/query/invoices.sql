-- name: GetAllInvoices :many
SELECT id, total, created
FROM invoices
ORDER BY created DESC;

-- name: GetInvoiceByID :one
SELECT id, total, created
FROM invoices
WHERE id = $1;

-- name: GetAllInvoicesWithReview :many
-- struct: InvoiceWithReview
SELECT
  i.id AS invoice_id,
  i.total AS invoice_total,
  i.created AS invoice_created,
  r.content,
  r.created AS review_created
FROM invoices i
LEFT JOIN reviews r ON r.id = i.id
ORDER BY i.created DESC;

-- name: GetInvoiceByIDWithReview :one
-- struct: InvoiceWithReview
SELECT
  i.id AS invoice_id,
  i.total AS invoice_total,
  i.created AS invoice_created,
  r.content,
  r.created AS review_created
FROM invoices i
LEFT JOIN reviews r ON r.id = i.id
WHERE i.id = $1;

-- name: CreateInvoice :one
INSERT INTO invoices (id, total)
VALUES ($1, $2)
RETURNING id, total, created;

-- name: DeleteInvoice :execrows
DELETE FROM invoices
WHERE id = $1;

-- name: CountInvoices :one
SELECT COUNT(*) FROM invoices;
