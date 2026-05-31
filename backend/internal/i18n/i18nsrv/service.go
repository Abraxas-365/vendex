package i18nsrv

import (
	"context"
	"time"

	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/i18n"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/google/uuid"
)

// Service implements the business logic for the i18n domain.
type Service struct {
	repo i18n.Repository
}

// New creates a new i18n Service.
func New(repo i18n.Repository) *Service {
	return &Service{repo: repo}
}

// validateLocale checks if the locale is in the supported locales map.
func validateLocale(locale string) error {
	if _, ok := i18n.SupportedLocales[locale]; !ok {
		return i18n.ErrInvalidLocale
	}
	return nil
}

// validateEntityType checks if the entity type is valid.
func validateEntityType(entityType string) error {
	if !i18n.ValidEntityTypes[entityType] {
		return i18n.ErrInvalidEntityType
	}
	return nil
}

// validateField checks if the field name is valid.
func validateField(field string) error {
	if !i18n.ValidFields[field] {
		return i18n.ErrInvalidField
	}
	return nil
}

// SetTranslation sets a single translated field value for an entity+locale+field.
func (s *Service) SetTranslation(ctx context.Context, tenantID kernel.TenantID, entityType, entityID, locale, field, value string) error {
	if err := validateEntityType(entityType); err != nil {
		return err
	}
	if err := validateLocale(locale); err != nil {
		return err
	}
	if err := validateField(field); err != nil {
		return err
	}
	if entityID == "" {
		return errx.New("entity_id is required", errx.TypeValidation)
	}

	now := time.Now()
	t := &i18n.Translation{
		ID:         kernel.NewTranslationID(uuid.NewString()),
		TenantID:   tenantID,
		EntityType: entityType,
		EntityID:   entityID,
		Locale:     locale,
		Field:      field,
		Value:      value,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := s.repo.Upsert(ctx, t); err != nil {
		return errx.Wrap(err, "upserting translation", errx.TypeInternal)
	}
	return nil
}

// SetTranslations bulk-sets multiple field translations for an entity+locale.
func (s *Service) SetTranslations(ctx context.Context, tenantID kernel.TenantID, entityType, entityID, locale string, fields map[string]string) error {
	if err := validateEntityType(entityType); err != nil {
		return err
	}
	if err := validateLocale(locale); err != nil {
		return err
	}
	if entityID == "" {
		return errx.New("entity_id is required", errx.TypeValidation)
	}

	now := time.Now()
	for field, value := range fields {
		if err := validateField(field); err != nil {
			return err
		}
		t := &i18n.Translation{
			ID:         kernel.NewTranslationID(uuid.NewString()),
			TenantID:   tenantID,
			EntityType: entityType,
			EntityID:   entityID,
			Locale:     locale,
			Field:      field,
			Value:      value,
			CreatedAt:  now,
			UpdatedAt:  now,
		}
		if err := s.repo.Upsert(ctx, t); err != nil {
			return errx.Wrap(err, "upserting translation", errx.TypeInternal)
		}
	}
	return nil
}

// GetTranslations returns the translation bundle (all fields) for an entity+locale.
func (s *Service) GetTranslations(ctx context.Context, tenantID kernel.TenantID, entityType, entityID, locale string) (*i18n.TranslationBundle, error) {
	if err := validateEntityType(entityType); err != nil {
		return nil, err
	}
	if err := validateLocale(locale); err != nil {
		return nil, err
	}

	bundle, err := s.repo.GetBundle(ctx, tenantID, entityType, entityID, locale)
	if err != nil {
		return nil, errx.Wrap(err, "fetching translation bundle", errx.TypeInternal)
	}
	return bundle, nil
}

// ListLocales returns all locales that have at least one translation for the entity.
func (s *Service) ListLocales(ctx context.Context, tenantID kernel.TenantID, entityType, entityID string) ([]string, error) {
	if err := validateEntityType(entityType); err != nil {
		return nil, err
	}

	locales, err := s.repo.ListLocales(ctx, tenantID, entityType, entityID)
	if err != nil {
		return nil, errx.Wrap(err, "listing locales", errx.TypeInternal)
	}
	return locales, nil
}

// DeleteTranslation removes a single translated field for an entity+locale+field.
func (s *Service) DeleteTranslation(ctx context.Context, tenantID kernel.TenantID, entityType, entityID, locale, field string) error {
	if err := validateEntityType(entityType); err != nil {
		return err
	}
	if err := validateLocale(locale); err != nil {
		return err
	}
	if err := validateField(field); err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, tenantID, entityType, entityID, locale, field); err != nil {
		return errx.Wrap(err, "deleting translation", errx.TypeInternal)
	}
	return nil
}

// DeleteAllTranslations removes all translations for an entity (all locales, all fields).
func (s *Service) DeleteAllTranslations(ctx context.Context, tenantID kernel.TenantID, entityType, entityID string) error {
	if err := validateEntityType(entityType); err != nil {
		return err
	}

	if err := s.repo.DeleteAll(ctx, tenantID, entityType, entityID); err != nil {
		return errx.Wrap(err, "deleting all translations for entity", errx.TypeInternal)
	}
	return nil
}

// ListSupportedLocales returns the static map of supported locales.
func (s *Service) ListSupportedLocales() map[string]i18n.LocaleInfo {
	return i18n.SupportedLocales
}
