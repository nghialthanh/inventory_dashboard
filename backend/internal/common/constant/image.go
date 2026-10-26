package constant

// Primary image upload limits.
const (
	PrimaryImageMaxBytes = 5 << 20 // 5 MiB
	PrimaryImageField    = "file"
)

// Allowed primary-image content types.
const (
	ContentTypeJPEG = "image/jpeg"
	ContentTypePNG  = "image/png"
	ContentTypeWebP = "image/webp"
)
