package swizzle

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	_ "github.com/lib/pq" // postgresql driver

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func addServiceEndpoints(serverType string, config *Config, r *gin.Engine) {

	startDate := time.Now()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"source":  serverType,
			"message": "pong",
		})
	})

	r.GET("/echo", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"source":  serverType,
			"message": c.Query("message"),
		})
	})

	r.GET("/config", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"source": serverType,
			"config": config,
		})
	})

	r.GET("/env", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"source": serverType,
			"env":    os.Environ(),
		})
	})

	r.GET("/status", func(c *gin.Context) {
		status := fmt.Sprintf("%s server running on port %d up since %s (%s)",
			serverType,
			config.Port,
			startDate.Format(time.UnixDate),
			time.Now().Sub(startDate).String())

		c.JSON(200, gin.H{
			"source":  serverType,
			"message": status,
		})
	})

	r.GET("/s3", func(c *gin.Context) {

		content, err := readFileFromS3(config)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"source":  serverType,
			"bucket":  config.S3Config.Bucket,
			"key":     config.S3Config.File,
			"content": content,
		})
	})

	r.GET("/redis", func(c *gin.Context) {
		v, err := redisPing(config)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"source":    serverType,
			"redis_val": v,
		})
	})

	r.GET("/pgsql", func(c *gin.Context) {
		v, err := pgPing(config)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"source":    serverType,
			"pgsql_val": v,
		})
	})

}

func createAwsSession(config *Config) (*session.Session, error) {

	if config.S3Config.AccessKey != "" {
		os.Setenv("AWS_ACCESS_KEY", config.S3Config.AccessKey)
	}
	if config.S3Config.SecretKey != "" {
		os.Setenv("AWS_SECRET_KEY", config.S3Config.SecretKey)
	}

	var creds *credentials.Credentials

	if config.S3Config.UseMetadata {
		creds = credentials.NewChainCredentials(
			[]credentials.Provider{
				&ec2rolecreds.EC2RoleProvider{
					Client: ec2metadata.New(session.New()),
				},
				&credentials.SharedCredentialsProvider{},
				&credentials.EnvProvider{},
			})
	} else {
		creds = credentials.NewChainCredentials(
			[]credentials.Provider{
				&credentials.SharedCredentialsProvider{},
				&credentials.EnvProvider{},
			})
	}

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(config.S3Config.Region),
		Credentials: creds,
	})

	if err == nil {
		_, err = sess.Config.Credentials.Get()
	}

	if err != nil {
		return nil, fmt.Errorf("Invalid/missing AWS creds: %s", err.Error())
	}

	return sess, nil
}

func readFileFromS3(config *Config) (string, error) {

	sess, err := createAwsSession(config)
	if err != nil {
		return "", err
	}

	s3svc := s3.New(sess)
	getParams := &s3.GetObjectInput{
		Bucket: aws.String(config.S3Config.Bucket),
		Key:    aws.String(config.S3Config.File),
	}

	obj, err := s3svc.GetObject(getParams)
	if err != nil {
		return "", err
	}

	defer obj.Body.Close()

	data, err := ioutil.ReadAll(obj.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func redisPing(config *Config) (int64, error) {

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	res := client.Incr("swizzle_test")
	if res.Err() != nil {
		return 0, res.Err()
	}

	return res.Val(), nil
}

func pgPing(config *Config) (int64, error) {
	db, err := sql.Open("postgres", config.Pg.URI)
	if err != nil {
		return 0, err
	}

	err = db.Ping()
	if err != nil {
		return 0, err
	}

	row := db.QueryRow("SELECT count(*) FROM pg_stat_activity")

	var count int64
	if err := row.Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}
