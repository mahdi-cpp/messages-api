package phasset

import (
	"strings"
	"time"

	"github.com/mahdi-cpp/iris-tools/search"
)

type SearchOptions struct {
	ID     string
	UserID string

	TextQuery string

	FileSize string `json:"fileSize"`
	FileType string `json:"fileType"`
	MimeType string `json:"mimeType"`

	PixelWidth  int
	PixelHeight int

	CameraMake  string
	CameraModel string

	IsCamera        *bool
	IsFavorite      *bool
	IsScreenshot    *bool
	IsHidden        *bool
	IsLandscape     *bool
	NotInOnePHAsset *bool

	HideScreenshot *bool `json:"hideScreenshot"`

	Albums  []string
	Trips   []string
	Persons []string

	NearPoint    []float64 `json:"nearPoint"`    // [latitude, longitude]
	WithinRadius float64   `json:"withinRadius"` // in kilometers
	BoundingBox  []float64 `json:"boundingBox"`  // [minLat, minLon, maxLat, maxLon]

	// Date filters
	CreatedAfter  *time.Time `json:"createdAfter,omitempty"`
	CreatedBefore *time.Time `json:"createdBefore,omitempty"`
	ActiveAfter   *time.Time `json:"activeAfter,omitempty"`

	// Pagination
	Offset int `json:"offset,omitempty"`
	Limit  int `json:"limit,omitempty"`

	// Sorting
	SortBy    string `json:"sortBy,omitempty"`    // "title", "created", "members", "lastActivity"
	SortOrder string `json:"sortOrder,omitempty"` // "asc" or "desc"
}

const MaxLimit = 1000

var LessFunks = map[string]search.LessFunction[*PHAsset]{
	"id":        func(a, b *PHAsset) bool { return a.ID < b.ID },
	"createdAt": func(a, b *PHAsset) bool { return a.CreatedAt.Before(b.CreatedAt) },
	"updatedAt": func(a, b *PHAsset) bool { return a.UpdatedAt.Before(b.UpdatedAt) },
}

func GetLessFunc(sortBy, sortOrder string) search.LessFunction[*PHAsset] {

	fn, exists := LessFunks[sortBy]
	if !exists {
		return nil
	}

	if sortOrder == "end" {
		return func(a, b *PHAsset) bool { return !fn(a, b) }
	}
	return fn
}

func BuildPHAssetCriteria(with *SearchOptions) search.Criteria[*PHAsset] {

	return func(c *PHAsset) bool {

		// ID filter
		//if with.ID != "" && c.ID != with.ID {
		//	return false
		//}

		// Title search_manager (case-insensitive)
		if with.TextQuery != "" {
			query := strings.ToLower(with.FileType)
			title := strings.ToLower(c.FileInfo.FileType)
			if !strings.Contains(title, query) {
				return false
			}
		}

		// Collection membership filters
		if len(with.Albums) > 0 {
			found := false
			for _, memberID := range with.Albums {
				if search.StringInSlice(memberID, c.Albums) {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}

		// Date filters
		if with.CreatedAfter != nil && c.CreatedAt.Before(*with.CreatedAfter) {
			return false
		}
		if with.CreatedBefore != nil && c.CreatedAt.After(*with.CreatedBefore) {
			return false
		}

		return true
	}
}

func Search(chats []*PHAsset, with *SearchOptions) []*PHAsset {

	// Build criteria
	criteria := BuildPHAssetCriteria(with)

	// Execute search_manager
	results := search.Find(chats, criteria)

	// Sort results if needed
	if with.SortBy != "" {
		lessFn := GetLessFunc(with.SortBy, with.SortOrder)
		if lessFn != nil {
			search.SortIndexedItems(results, lessFn)
		}
	}

	// Extract final assets
	final := make([]*PHAsset, len(results))
	for i, item := range results {
		final[i] = item.Value
	}

	if with.Limit == 0 { // if not set default is MAX_LIMIT
		with.Limit = MaxLimit
	}

	// Apply pagination
	start := with.Offset
	end := start + with.Limit
	if end > len(final) {
		end = len(final)
	}
	return final[start:end]
}
