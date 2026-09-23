package producttesting

import (
	"context"
	"errors"
	"testing"

	"telegram-service-platform/entity/smmentity"
	"telegram-service-platform/params/productparams"
)

func TestService_GetDistinctPlatforms_CacheHit(
	t *testing.T,
) {
	cache := &fakeCatalogCache{
		platforms: []smmentity.Platform{
			{Name: "telegram"},
		},
		platformsFound: true,
	}

	repository := &fakeProductRepository{}

	service := newProductServiceWithCatalogCache(
		repository,
		&fakeSMMAdapter{},
		cache,
	)

	response, err := service.GetDistinctPlatforms(
		context.Background(),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(response.Platforms) != 1 {
		t.Fatalf(
			"expected 1 cached platform, got %d",
			len(response.Platforms),
		)
	}

	if repository.platformCalls != 0 {
		t.Fatalf(
			"expected database not to be called on cache hit, got %d calls",
			repository.platformCalls,
		)
	}
}

func TestService_GetDistinctPlatforms_CacheReadErrorFallsBackToDB(
	t *testing.T,
) {
	cache := &fakeCatalogCache{
		getPlatformsErr: errors.New("redis unavailable"),
	}

	repository := &fakeProductRepository{
		platforms: []smmentity.Platform{
			{Name: "telegram"},
		},
	}

	service := newProductServiceWithCatalogCache(
		repository,
		&fakeSMMAdapter{},
		cache,
	)

	response, err := service.GetDistinctPlatforms(
		context.Background(),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(response.Platforms) != 1 {
		t.Fatalf(
			"expected DB platform to be returned, got %d",
			len(response.Platforms),
		)
	}

	if repository.platformCalls != 1 {
		t.Fatalf(
			"expected one database call, got %d",
			repository.platformCalls,
		)
	}
}

func TestService_GetDistinctPlatforms_CacheWriteErrorDoesNotFail(
	t *testing.T,
) {
	cache := &fakeCatalogCache{
		setPlatformsErr: errors.New("redis write failed"),
	}

	repository := &fakeProductRepository{
		platforms: []smmentity.Platform{
			{Name: "telegram"},
		},
	}

	service := newProductServiceWithCatalogCache(
		repository,
		&fakeSMMAdapter{},
		cache,
	)

	response, err := service.GetDistinctPlatforms(
		context.Background(),
	)
	if err != nil {
		t.Fatalf(
			"cache write failure must not fail business operation: %v",
			err,
		)
	}

	if len(response.Platforms) != 1 {
		t.Fatalf(
			"expected DB platform to be returned, got %d",
			len(response.Platforms),
		)
	}
}

func TestService_GetDistinctCategoriesByPlatform_CacheFailureFallsBackToDB(
	t *testing.T,
) {
	cache := &fakeCatalogCache{
		getCategoriesErr: errors.New("redis unavailable"),
	}

	repository := &fakeProductRepository{
		categories: []smmentity.Category{
			{Name: "members"},
		},
	}

	service := newProductServiceWithCatalogCache(
		repository,
		&fakeSMMAdapter{},
		cache,
	)

	response, err := service.GetDistinctCategoriesByPlatform(
		context.Background(),
		productparams.GetDistinctCategoriesByPlatformRequest{
			Platform: smmentity.TelegramPlatform,
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(response.Categories) != 1 {
		t.Fatalf(
			"expected 1 category, got %d",
			len(response.Categories),
		)
	}

	if repository.categoryCalls != 1 {
		t.Fatalf(
			"expected one database call, got %d",
			repository.categoryCalls,
		)
	}
}

func TestService_CreateSMMMapping_InvalidatesCatalogCache(
	t *testing.T,
) {
	cache := &fakeCatalogCache{}

	repository := &fakeProductRepository{}

	service := newProductServiceWithCatalogCache(
		repository,
		&fakeSMMAdapter{},
		cache,
	)

	_, err := service.CreateSMMMapping(
		context.Background(),
		productparams.CreateSMMMappingRequest{
			SmmServiceId: 100,
			Name:         "Telegram Members",
			Platform:     smmentity.TelegramPlatform,
			Category:     smmentity.MemberCategory,
			IsActive:     true,
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(cache.invalidatedPlatforms) != 1 {
		t.Fatalf(
			"expected one invalidated platform, got %d",
			len(cache.invalidatedPlatforms),
		)
	}

	if cache.invalidatedPlatforms[0] !=
		smmentity.TelegramPlatform.String() {
		t.Fatalf(
			"expected telegram cache invalidation, got %s",
			cache.invalidatedPlatforms[0],
		)
	}
}

func TestService_CreateSMMMapping_CacheInvalidationFailureDoesNotFailBusinessOperation(
	t *testing.T,
) {
	cache := &fakeCatalogCache{
		invalidateErr: errors.New("redis unavailable"),
	}

	repository := &fakeProductRepository{}

	service := newProductServiceWithCatalogCache(
		repository,
		&fakeSMMAdapter{},
		cache,
	)

	_, err := service.CreateSMMMapping(
		context.Background(),
		productparams.CreateSMMMappingRequest{
			SmmServiceId: 100,
			Name:         "Telegram Members",
			Platform:     smmentity.TelegramPlatform,
			Category:     smmentity.MemberCategory,
			IsActive:     true,
		},
	)
	if err != nil {
		t.Fatalf(
			"cache invalidation failure must not fail DB mutation: %v",
			err,
		)
	}

	if repository.createdMapping == nil {
		t.Fatal("expected mapping to be persisted")
	}
}

func TestService_UpdateSMMMapping_InvalidatesOldAndNewPlatforms(
	t *testing.T,
) {
	cache := &fakeCatalogCache{}

	repository := &fakeProductRepository{
		existingMapping: &smmentity.SmmMapping{
			Id:       10,
			Platform: smmentity.TelegramPlatform,
			Category: smmentity.MemberCategory,
		},
	}

	service := newProductServiceWithCatalogCache(
		repository,
		&fakeSMMAdapter{},
		cache,
	)

	_, err := service.UpdateSMMMapping(
		context.Background(),
		productparams.UpdateSMMMappingRequest{
			Id:           10,
			SmmServiceId: 200,
			Name:         "Instagram Followers",
			Platform:     smmentity.InstagramPlatform,
			Category:     smmentity.FollowerCategory,
			IsActive:     true,
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(cache.invalidatedPlatforms) != 2 {
		t.Fatalf(
			"expected 2 invalidations, got %d",
			len(cache.invalidatedPlatforms),
		)
	}

	if cache.invalidatedPlatforms[0] !=
		smmentity.TelegramPlatform.String() {
		t.Fatalf(
			"expected old platform telegram, got %s",
			cache.invalidatedPlatforms[0],
		)
	}

	if cache.invalidatedPlatforms[1] !=
		smmentity.InstagramPlatform.String() {
		t.Fatalf(
			"expected new platform instagram, got %s",
			cache.invalidatedPlatforms[1],
		)
	}
}
func TestService_GetDistinctPlatforms_WithoutCatalogCache(
	t *testing.T,
) {
	repository := &fakeProductRepository{
		platforms: []smmentity.Platform{
			{
				Name: "telegram",
			},
		},
	}

	service := newProductService(
		repository,
		&fakeSMMAdapter{},
	)

	response, err := service.GetDistinctPlatforms(
		context.Background(),
	)
	if err != nil {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}

	if len(response.Platforms) != 1 {
		t.Fatalf(
			"expected 1 platform, got %d",
			len(response.Platforms),
		)
	}

	if repository.platformCalls != 1 {
		t.Fatalf(
			"expected one database call, got %d",
			repository.platformCalls,
		)
	}
}

func TestService_InvalidateCatalogCache_WithoutCatalogCache(
	t *testing.T,
) {
	service := newProductService(
		&fakeProductRepository{},
		&fakeSMMAdapter{},
	)

	err := service.InvalidateCatalogCache(
		context.Background(),
		smmentity.TelegramPlatform,
	)
	if err != nil {
		t.Fatalf(
			"expected nil error without catalog cache, got %v",
			err,
		)
	}
}
