CREATE TABLE wishlists (
    id          VARCHAR(36)  NOT NULL,
    tenant_id   VARCHAR(36)  NOT NULL,
    customer_id VARCHAR(36)  NOT NULL,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id),
    CONSTRAINT fk_wishlists_customer FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
    CONSTRAINT uq_wishlists_tenant_customer UNIQUE (tenant_id, customer_id)
);

CREATE TABLE wishlist_items (
    id           VARCHAR(36) NOT NULL,
    wishlist_id  VARCHAR(36) NOT NULL,
    product_id   VARCHAR(36) NOT NULL,
    variant_id   VARCHAR(36) NULL,
    added_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id),
    CONSTRAINT fk_wishlist_items_wishlist FOREIGN KEY (wishlist_id) REFERENCES wishlists(id) ON DELETE CASCADE,
    CONSTRAINT uq_wishlist_items_product UNIQUE (wishlist_id, product_id, variant_id)
);

CREATE INDEX idx_wishlists_tenant ON wishlists(tenant_id);
CREATE INDEX idx_wishlists_customer ON wishlists(tenant_id, customer_id);
CREATE INDEX idx_wishlist_items_wishlist ON wishlist_items(wishlist_id);
