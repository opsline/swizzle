package swizzle

// Config the swizzle configuration struct
type Config struct {
	Port int

	S3Config struct {
		UseMetadata bool
		Bucket      string
		File        string
		Region      string
		AccessKey   string
		SecretKey   string
	}
}
