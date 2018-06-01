package swizzle

// Config the swizzle configuration struct
type Config struct {
	Port       int
	PathPrefix string // location where the index page should be served, e.g. "" for simply "/"

	S3Config struct {
		UseMetadata bool
		Bucket      string
		File        string
		Region      string
		AccessKey   string
		SecretKey   string
	}

	Redis struct {
		Host string
		Port int
	}

	Pg struct {
		URI string
	}
}
