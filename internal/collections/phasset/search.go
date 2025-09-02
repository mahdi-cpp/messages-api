package phasset

import (
	"strings"

	"github.com/mahdi-cpp/iris-tools/search"
)

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

func BuildPHAssetSearchCriteria(with *SearchOptions) search.SearchCriteria[*PHAsset] {

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

		// Boolean flags
		if with.IsCamera != nil && c.IsCamera != *with.IsCamera {
			return false
		}
		if with.IsFavorite != nil && c.IsFavorite != *with.IsFavorite {
			return false
		}
		if with.IsScreenshot != nil && c.IsScreenshot != *with.IsScreenshot {
			return false
		}
		if with.IsHidden != nil && c.IsHidden != *with.IsHidden {
			return false
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

		if len(with.Trips) > 0 {
			found := false
			for _, adminID := range with.Trips {
				if search.StringInSlice(adminID, c.Trips) {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}

		if len(with.Persons) > 0 {
			found := false
			for _, blockID := range with.Persons {
				if search.StringInSlice(blockID, c.Persons) {
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
	criteria := BuildPHAssetSearchCriteria(with)

	// Execute search_manager
	results := search.Search(chats, criteria)

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
