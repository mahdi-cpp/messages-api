package phasset

import (
	"time"

	"github.com/mahdi-cpp/iris-tools/update"
)

type UpdateOptions struct {
	AssetIds []string `json:"assetIds,omitempty"` // Asset Ids

	FileSize string `json:"fileSize"`
	FileType string `json:"fileType"`
	MimeType string `json:"mimeType"`

	CameraMake  *string `json:"cameraMake,omitempty"`
	CameraModel *string `json:"cameraModel,omitempty"`

	IsCamera        *bool
	IsFavorite      *bool
	IsScreenshot    *bool
	IsHidden        *bool
	NotInOnePHAsset *bool

	Albums       *[]string `json:"albums,omitempty"`       // Full album replacement
	AddAlbums    []string  `json:"addAlbums,omitempty"`    // PHAssets to add
	RemoveAlbums []string  `json:"removeAlbums,omitempty"` // PHAssets to remove

	Trips       *[]string `json:"trips,omitempty"`       // Full trip replacement
	AddTrips    []string  `json:"addTrips,omitempty"`    // Trips to add
	RemoveTrips []string  `json:"removeTrips,omitempty"` // Trips to remove

	Persons       *[]string `json:"persons,omitempty"`       // Full Person replacement
	AddPersons    []string  `json:"addPersons,omitempty"`    // Persons to add
	RemovePersons []string  `json:"removePersons,omitempty"` // Persons to remove
}

// Initialize updater
var metadataUpdater = update.NewUpdater[PHAsset, UpdateOptions]()

func init() {

	// Configure scalar field updates
	metadataUpdater.AddScalarUpdater(func(a *PHAsset, u UpdateOptions) {
		if u.FileType != "" {
			a.FileInfo.FileType = u.FileType
		}
	})

	metadataUpdater.AddScalarUpdater(func(a *PHAsset, u UpdateOptions) {
		if u.FileSize != "" {
			a.FileInfo.FileSize = u.FileSize
		}
	})

	// Add other scalar fields similarly...

	// Configure collection operations
	metadataUpdater.AddCollectionUpdater(func(a *PHAsset, u UpdateOptions) {
		op := update.CollectionUpdateOp[string]{
			FullReplace: u.Albums,
			Add:         u.AddAlbums,
			Remove:      u.RemoveAlbums,
		}
		a.Albums = update.ApplyCollectionUpdate(a.Albums, op)
	})

	// Set modification timestamp
	metadataUpdater.AddPostUpdateHook(func(a *PHAsset) {
		a.UpdatedAt = time.Now()
	})
}

func Update(p *PHAsset, update UpdateOptions) *PHAsset {
	metadataUpdater.Apply(p, update)
	return p
}

// IsEmpty checks if the Place struct contains zero values for all its fields.
func (l Location) IsEmpty() bool {
	return l.Latitude == 0.0 &&
		l.Longitude == 0.0 &&
		l.City == "" &&
		l.Country == ""
}
