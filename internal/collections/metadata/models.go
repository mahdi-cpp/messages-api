package metadata

import (
	"time"

	"github.com/google/uuid"
)

type Metadata struct {
	ID               uuid.UUID              `json:"id"`
	FileInfo         FileInfo               `json:"fileInfo,omitempty"`
	Image            ImageInfo              `json:"image,omitempty"`
	Camera           CameraInfo             `json:"camera,omitempty"`
	Video            VideoInfo              `json:"video,omitempty"`
	Location         Location               `json:"location,omitempty"`
	DateTimeOriginal time.Time              `json:"dateTimeOriginal,omitempty"`
	RawData          map[string]interface{} `json:"-"` // Raw EXIF data for debugging
}

type FileInfo struct {
	BaseURL  string `json:"baseURL"`
	FileSize int    `json:"fileSize"`
	FileType string `json:"fileType"`
	MimeType string `json:"mimeType"`
}

type ImageInfo struct {
	Width           int     `json:"width,omitempty"`
	Height          int     `json:"height,omitempty"`
	Megapixels      float64 `json:"megapixels,omitempty"`
	Orientation     string  `json:"orientation,omitempty"`
	ColorSpace      string  `json:"colorSpace,omitempty"`
	EncodingProcess string  `json:"encodingProcess,omitempty"`
}

type CameraInfo struct {
	Make            string  `json:"make,omitempty"`
	Model           string  `json:"model,omitempty"`
	Software        string  `json:"software,omitempty"`
	ExposureTime    string  `json:"exposureTime,omitempty"`
	FNumber         float64 `json:"fNumber,omitempty"`
	ISO             int     `json:"iso,omitempty"`
	FocalLength     string  `json:"focalLength,omitempty"`
	FocalLength35mm string  `json:"focalLength35mm,omitempty"`
	Flash           string  `json:"flash,omitempty"`
	LightSource     string  `json:"lightSource,omitempty"`
	ExposureMode    string  `json:"exposureMode,omitempty"`
	WhiteBalance    string  `json:"whiteBalance,omitempty"`
}

type VideoInfo struct {
	MediaDuration      string  `json:"mediaDuration,omitempty"`
	Width              int     `json:"width,omitempty"`
	Height             int     `json:"height,omitempty"`
	VideoFrameRate     float64 `json:"videoFrameRate,omitempty"`
	AvgBitrate         string  `json:"avgBitrate,omitempty"`
	Encoder            string  `json:"encoder,omitempty"`
	Rotation           int     `json:"rotation,omitempty"`
	AudioFormat        string  `json:"audioFormat,omitempty"`
	AudioChannels      int     `json:"audioChannels,omitempty"`
	AudioSampleRate    int     `json:"audioSampleRate,omitempty"`
	AudioBitsPerSample int     `json:"audioBitsPerSample,omitempty"`
}

type Location struct {
	Latitude   float64 `json:"latitude,omitempty"`
	Longitude  float64 `json:"longitude,omitempty"`
	Country    string  `json:"country,omitempty"`
	Province   string  `json:"province,omitempty"`
	County     string  `json:"county,omitempty"`
	City       string  `json:"city,omitempty"`
	Village    string  `json:"village,omitempty"`
	Electronic int     `json:"electronic,omitempty"`
}
