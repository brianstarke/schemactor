CREATE TABLE reviews (
    id bigserial,
    product_id bigint NOT NULL,
    user_id bigint NOT NULL,
    rating integer NOT NULL,
    title varchar(200),
    comment text,
    created_at timestamptz DEFAULT now() NOT NULL,
    updated_at timestamptz DEFAULT now() NOT NULL,
    verified_purchase boolean,
    helpful_count integer,
    seller_response text,
    seller_responded_at timestamptz,
    image_urls text[],
    PRIMARY KEY (id),
    FOREIGN KEY (product_id) REFERENCES products (id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE INDEX idx_reviews_product ON reviews (product_id);

CREATE INDEX idx_reviews_user ON reviews (user_id);

CREATE INDEX idx_reviews_rating ON reviews (rating);

CREATE INDEX idx_reviews_verified ON reviews (verified_purchase);
