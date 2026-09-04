package dto

type CropTransformation struct {
	Height uint32 `json:"height"`
	Width uint32 `json:"width"`
}

type ScaleTransformation struct {
	Height uint32 `json:"height"`
	Width uint32 `json:"width"`
}

type ConvertTransformation struct {
	Format string `json:"format"`
}

type CompressTransformation struct {
	Compress bool `json:"compress"`
}