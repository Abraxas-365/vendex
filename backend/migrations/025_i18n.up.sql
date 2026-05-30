CREATE TABLE translations (
    id          TEXT        NOT NULL,
    tenant_id   TEXT        NOT NULL,
    entity_type TEXT        NOT NULL,
    entity_id   TEXT        NOT NULL,
    locale      TEXT        NOT NULL,
    field       TEXT        NOT NULL,
    value       TEXT        NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id),
    UNIQUE (tenant_id, entity_type, entity_id, locale, field)
);

CREATE INDEX idx_translations_entity ON translations(tenant_id, entity_type, entity_id);
CREATE INDEX idx_translations_locale ON translations(tenant_id, locale);
