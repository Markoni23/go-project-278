package service_test

import (
	"context"
	"errors"
	"testing"

	"markoni23/url-shortener/internal/domain"
	"markoni23/url-shortener/internal/repository"
	"markoni23/url-shortener/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockLinkRepository struct {
	mock.Mock
}

func (m *MockLinkRepository) FindAll(ctx context.Context) ([]domain.Link, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Link), args.Error(1)
}

func (m *MockLinkRepository) Create(ctx context.Context, dto repository.CreateLinkDTO) (domain.Link, error) {
	args := m.Called(ctx, dto)
	return args.Get(0).(domain.Link), args.Error(1)
}

func (m *MockLinkRepository) Get(ctx context.Context, id int64) (domain.Link, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.Link), args.Error(1)
}

func (m *MockLinkRepository) Update(ctx context.Context, id int64, dto repository.UpdateLinkDTO) (domain.Link, error) {
	args := m.Called(ctx, id, dto)
	return args.Get(0).(domain.Link), args.Error(1)
}

func (m *MockLinkRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestLinkService_GetLinksList(t *testing.T) {
	tests := []struct {
		name          string
		mockSetup     func(*MockLinkRepository)
		expectedLen   int
		expectedError error
	}{
		{
			name: "success - get empty list",
			mockSetup: func(m *MockLinkRepository) {
				m.On("FindAll", mock.Anything).Return([]domain.Link{}, nil)
			},
			expectedLen:   0,
			expectedError: nil,
		},
		{
			name: "success - get links list",
			mockSetup: func(m *MockLinkRepository) {
				links := []domain.Link{
					{ID: 1, OriginalUrl: "https://google.com", ShortName: "google", ShortUrl: "http://localhost:8080/google"},
					{ID: 2, OriginalUrl: "https://github.com", ShortName: "github", ShortUrl: "http://localhost:8080/github"},
				}
				m.On("FindAll", mock.Anything).Return(links, nil)
			},
			expectedLen:   2,
			expectedError: nil,
		},
		{
			name: "error - repository error",
			mockSetup: func(m *MockLinkRepository) {
				m.On("FindAll", mock.Anything).Return([]domain.Link{}, errors.New("database error"))
			},
			expectedLen:   0,
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockLinkRepository)
			tt.mockSetup(mockRepo)

			linkService := service.NewLinkService("http://localhost:8080", mockRepo)

			links, err := linkService.GetLinksList(context.Background())

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Len(t, links, tt.expectedLen)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestLinkService_CreateLink(t *testing.T) {
	basePath := "http://localhost:8080"

	t.Run("create with custom short name", func(t *testing.T) {
		mockRepo := new(MockLinkRepository)
		linkService := service.NewLinkService(basePath, mockRepo)

		originalUrl := "https://example.com"
		customName := "custom"

		mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(dto repository.CreateLinkDTO) bool {
			return *dto.OriginalUrl == originalUrl &&
				*dto.ShortName == customName &&
				*dto.ShortUrl != ""
		})).
			Return(domain.Link{ID: 1, OriginalUrl: originalUrl, ShortName: customName}, nil)

		link, err := linkService.CreateLink(context.Background(), originalUrl, customName)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), link.ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("create with generated short name", func(t *testing.T) {
		mockRepo := new(MockLinkRepository)
		linkService := service.NewLinkService(basePath, mockRepo)

		originalUrl := "https://example.com"

		mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(dto repository.CreateLinkDTO) bool {
			return *dto.OriginalUrl == originalUrl &&
				*dto.ShortName != "" &&
				*dto.ShortUrl != ""
		})).Return(domain.Link{}, nil)

		_, err := linkService.CreateLink(context.Background(), originalUrl, "")

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}
