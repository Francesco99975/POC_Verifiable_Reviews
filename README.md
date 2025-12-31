# Verifiable True Reviews - Proof of Concept

A simple web application that demonstrates a system for **verifiably authentic customer reviews** by tying each review directly to a real purchase.

## Core Concept

Traditional review systems suffer from fake reviews because anyone can post without proving they actually bought the product. This proof-of-concept solves that by **associating exactly one review per invoice**.

- A customer can only leave a review if they have a valid invoice ID from a completed purchase.
- Each invoice can have **at most one review** (enforced at the database level).
- No purchase = no review. This eliminates fake, incentivized, or competitor sabotage reviews.

## How It Works

1. A customer completes a purchase → an `invoice` record is created with a unique UUID.
2. After purchase, the customer receives their invoice ID (e.g., via email or order confirmation page).
3. The customer visits the review page and enters their invoice ID.
4. The app verifies the invoice exists and checks if a review has already been submitted for it.
5. If valid and no prior review exists, the customer can submit their review.
6. The review is stored with the **same ID** as the invoice, enforced by a foreign key constraint.

This design guarantees **one review per purchase** and makes fake reviews practically impossible without access to a real invoice.

## Database Schema

The authenticity is enforced directly in the database using a clever schema design:

```sql
CREATE TABLE IF NOT EXISTS invoices (
  id UUID NOT NULL,
  total INT NOT NULL,
  created TIMESTAMP NOT NULL DEFAULT NOW(),
  updated TIMESTAMP NOT NULL DEFAULT NOW(),
  PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS reviews (
  id UUID NOT NULL,
  content TEXT NOT NULL,
  created TIMESTAMP NOT NULL DEFAULT NOW(),
  updated TIMESTAMP NOT NULL DEFAULT NOW(),
  PRIMARY KEY(id),
  CONSTRAINT fk_reviews_invoice
    FOREIGN KEY (id)
    REFERENCES invoices (id)
    ON DELETE CASCADE
);
