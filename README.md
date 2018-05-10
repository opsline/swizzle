# Swizzle

Swizzle is a simple "hello world" type tool for testing cloud-native infrastructure deployments.

# Configuration

The following environment config vars are available:

Variable            | Default | Description
--------------------|---------|-------------------------
`SWIZZLE_PORT`      | 8080    | Port number to listen on
`SWIZZLE_S3_BUCKET` | com.opsline.swizzle | S3 bucket to read from
`SWIZZLE_S3_FILE`   | hello.txt | S3 filename to read
`SWIZZLE_AWS_REGION`| us-east-1 | AWS Region
`SWIZZLE_AWS_ACCESS_KEY`|  | AWS Access Key (reads env/instance creds when empty)
`SWIZZLE_AWS_SECRET_KEY`|  | AWS Secret Key (reads env/instance creds when empty)
`SWIZZLE_AWS_USE_METADATA`| false | Load AWS creds from instance metadata
`SWIZZLE_REDIS_HOST`| localhost | Redis hostname
`SWIZZLE_REDIS_PORT`| 6379 | Redis port number
`SWIZZLE_PG_URI`| postgres://postgres:postgres@localhost/postgres?sslmode=disable | PostgreSQL connection URI
