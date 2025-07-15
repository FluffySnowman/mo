package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	// "log"
	// _ "log"
	_ "log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"
	Bold    = "\033[1m"
	Dim     = "\033[2m"
)

type Config struct {
	Endpoint  string
	Region    string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
	Recursive bool
}

type Mo struct {
	client *minio.Client
	config *Config
}

func main() {
	if len(os.Args) < 2 {
		printHelpShit()
		os.Exit(1)
	}

	// straight up doing help here cos of all the references + its gonna exit
	// doing --help and nothing else so ye
	if os.Args[1] == "help" || os.Args[1] == "-help" || os.Args[1] == "--help" || os.Args[1] == "-h" || os.Args[1] == "--h" {
		printHelpShit()
		os.Exit(0)
		// compilers do be yappin
		// return nil;
	}
	

	config, err := parseConfig()
	if err != nil {
		printError("configuration error: %v", err)
		os.Exit(1)
	}

	// cfg, err := parseConfig()
	// if err != nil {
	// 	printError("config parse failed: %v", err)
	// 	os.Exit(1)
	// }

	mo, err := NewMo(config)
	if err != nil {
		printError("failed to initialize client: %v", err)
		os.Exit(1)
	}

	if err := mo.__EXEC_CMD(); err != nil {
		printError("command failed: %v", err)
		os.Exit(1)
	}
}

func NewMo(config *Config) (*Mo, error) {
	printInfo("connecting to %s...", config.Endpoint)
	
	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKey, config.SecretKey, ""),
		Secure: config.UseSSL,
		Region: config.Region,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	// _, err = client.ListBuckets(ctx)
	// if err != nil {
	// 	return nil, fmt.Errorf("connection test failed: %v", err)
	// }

	// check if bucket exists first
	exists, err := client.BucketExists(ctx, config.Bucket)
	if err != nil {
		return nil, fmt.Errorf("connection test failed: %v", err)
	}
	if !exists && config.Bucket != "" {
		printWarning("bucket '%s' does not exist", config.Bucket)
	}

	printSuccess("successfully connected to %s", config.Endpoint)
	return &Mo{client: client, config: config}, nil
}

func parseConfig() (*Config, error) {
	config := &Config{
		UseSSL: true,
	}

	// load config file if it exists
	if err := loadConfigFile(config); err != nil {
		printWarning("could not load mo.conf: %v", err)
	}

	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-endpoint":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("missing value for -endpoint")
			}
			config.Endpoint = args[i+1]
			i++
		case "-region":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("missing value for -region")
			}
			config.Region = args[i+1]
			i++
		case "-key":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("missing value for -key")
			}
			config.AccessKey = args[i+1]
			i++
		case "-secret":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("missing value for -secret")
			}
			config.SecretKey = args[i+1]
			i++
		case "-bucket":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("missing value for -bucket")
			}
			config.Bucket = args[i+1]
			i++
		case "-insecure":
			config.UseSSL = false
		case "-r":
			config.Recursive = true
		// case "-R":
		// 	config.Recursive = true
		}
	}

	// override with env vars if exis t 
	if endpoint := os.Getenv("MO_ENDPOINT"); endpoint != "" {
		config.Endpoint = endpoint
	}
	if region := os.Getenv("MO_REGION"); region != "" {
		config.Region = region
	}
	if key := os.Getenv("MO_ACCESS_KEY"); key != "" {
		config.AccessKey = key
	}
	if secret := os.Getenv("MO_SECRET_KEY"); secret != "" {
		config.SecretKey = secret
	}
	if bucket := os.Getenv("MO_BUCKET"); bucket != "" {
		config.Bucket = bucket
	}

	// validate required stuff
	if config.Endpoint == "" {
		return nil, fmt.Errorf("endpoint is required")
	}
	if config.AccessKey == "" {
		return nil, fmt.Errorf("access key is required")
	}
	if config.SecretKey == "" {
		return nil, fmt.Errorf("secret key is required")
	}

	// cleaning/sanitising(?)  idk what the word would be for this lol
	if strings.HasPrefix(config.Endpoint, "http://") {
		config.Endpoint = strings.TrimPrefix(config.Endpoint, "http://")
		config.UseSSL = false
	} else if strings.HasPrefix(config.Endpoint, "https://") {
		config.Endpoint = strings.TrimPrefix(config.Endpoint, "https://")
		config.UseSSL = true
	}

	return config, nil
}

func loadConfigFile(config *Config) error {
	file, err := os.Open("mo.conf")
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "endpoint":
			config.Endpoint = value
		case "region":
			config.Region = value
		case "access_key":
			config.AccessKey = value
		case "secret_key":
			config.SecretKey = value
		case "bucket":
			config.Bucket = value
		}
	}

	return scanner.Err()
}

func (m *Mo) __EXEC_CMD() error {
	args := os.Args[1:]
	
	var command string
	var commandArgs []string

	// find the actual cmd in args
	for i, arg := range args {
		if !strings.HasPrefix(arg, "-") && 
		   (i == 0 || !strings.HasPrefix(args[i-1], "-") || 
		    args[i-1] == "-insecure" || args[i-1] == "-r") {
			command = arg
			commandArgs = args[i+1:]
			break
		}
	}

	switch command {
	case "ls":
		return m.listObjects(commandArgs)
	case "buckets":
		return m.listBuckets()
	case "cp":
		return m.copyObject(commandArgs)
	case "rm":
		return m.removeObject(commandArgs)
	case "mv":
		return m.moveObject(commandArgs)
	case "stat":
		return m.statObject(commandArgs)
	case "mb":
		return m.makeBucket(commandArgs)
	case "rb":
		return m.removeBucket(commandArgs)
	default:
		return fmt.Errorf("unknown command: %s. use --help for all availableu cmds/opts", command)
	}
}

func (m *Mo) listBuckets() error {
	printInfo("listing buckets...")
	
	ctx := context.Background()
	buckets, err := m.client.ListBuckets(ctx)
	if err != nil {
		return fmt.Errorf("failed to list buckets: %v", err)
	}

	if len(buckets) == 0 {
		printInfo("no buckets found")
		return nil
	}

	printSuccess("found %d bucket(s):", len(buckets))
	for _, bucket := range buckets {
		fmt.Printf("  %s%-20s%s %s%s%s\n", 
			Cyan, bucket.Name, Reset,
			Dim, bucket.CreationDate.Format("2006-01-02 15:04:05"), Reset)
	}

	return nil
}

func (m *Mo) listObjects(args []string) error {
	if m.config.Bucket == "" {
		return fmt.Errorf("bucket is required for ls command")
	}

	prefix := ""
	if len(args) > 0 {
		prefix = args[0]
	}

	// remove leading slash if present
	if strings.HasPrefix(prefix, "/") {
		prefix = prefix[1:]
	}

	printInfo("listing objects in bucket '%s' with prefix '%s'...", m.config.Bucket, prefix)

	ctx := context.Background()
	
	// opts := minio.ListObjectsOptions{
	// 	Prefix:    prefix,
	// 	Recursive: false,
	// }

	objectCh := m.client.ListObjects(ctx, m.config.Bucket, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	count := 0
	for object := range objectCh {
		if object.Err != nil {
			return fmt.Errorf("error listing objects: %v", object.Err)
		}

		count++
		size := formatSize(object.Size)
		fmt.Printf("  %s%-50s%s %s%8s%s %s%s%s\n",
			Yellow, object.Key, Reset,
			Green, size, Reset,
			Dim, object.LastModified.Format("2006-01-02 15:04:05"), Reset)
	}

	if count == 0 {
		printInfo("no objects found")
	} else {
		printSuccess("listed %d object(s)", count)
	}

	return nil
}

func (m *Mo) copyObject(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("cp requires source and destination arguments")
	}

	src := args[0]
	dst := args[1]

	// if strings.HasPrefix(src, "./") || !strings.Contains(src, "/") {
	// 	return m.uploadFile(src, dst)
	// } else if strings.HasPrefix(dst, "./") || !strings.Contains(dst, "/") {
	// 	return m.downloadFile(src, dst)
	// }

	if isLocalPath(src) && !isLocalPath(dst) {
		if m.config.Recursive {
			return m.uploadRecursive(src, dst)
		}
		return m.uploadFile(src, dst)
	} else if !isLocalPath(src) && isLocalPath(dst) {
		return m.downloadFile(src, dst)
	} else {
		return m.copyRemoteObject(src, dst)
	}
}

func (m *Mo) uploadRecursive(localPath, remotePath string) error {
	if m.config.Bucket == "" {
		return fmt.Errorf("bucket is ReQUIReD for upload")
	}

	// check if local is dir
	stat, err := os.Stat(localPath)
	if err != nil {
		return fmt.Errorf("failed to stat local path: %v", err)
	}

	if stat.IsDir() {
		return m.uploadDirectory(localPath, remotePath)
	} else {
		return m.uploadFile(localPath, remotePath)
	}
}

func (m *Mo) uploadDirectory(localDir, remotePrefix string) error {
	printInfo("uploading directory %s to %s/%s...", localDir, m.config.Bucket, remotePrefix)

	// going out for an excursion
	err := filepath.Walk(localDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(localDir, path)
		if err != nil {
			return err
		}

		// cos s3 is retarded
		remotePath := remotePrefix + "/" + strings.Replace(relPath, "\\", "/", -1)
		
		file, err := os.Open(path)
		if err != nil {
			printError("failed to open %s: %v", path, err)
			return nil // do other files rn
		}
		defer file.Close()

		printInfo("uploading %s -> %s", path, remotePath)

		ctx := context.Background()
		
		contentType := "application/octet-stream"
		if strings.HasSuffix(path, ".png") {
			contentType = "image/png"
		} else if strings.HasSuffix(path, ".jpg") || strings.HasSuffix(path, ".jpeg") {
			contentType = "image/jpeg"
		} else if strings.HasSuffix(path, ".txt") {
			contentType = "text/plain"
		}

		_, err = m.client.PutObject(ctx, m.config.Bucket, remotePath, file, info.Size(), minio.PutObjectOptions{
			ContentType: contentType,
		})
		if err != nil {
			printError("upload failed for %s: %v", path, err)
			return nil // do other files rn
		}

		printSuccess("uploaded %s (%s)", path, formatSize(info.Size()))
		return nil
	})

	if err != nil {
		return fmt.Errorf("directory walk failed: %v", err)
	}

	return nil
}

func isLocalPath(path string) bool {
	return strings.HasPrefix(path, "./") || strings.HasPrefix(path, "/") || 
		   strings.HasPrefix(path, "~/") || !strings.Contains(path, "/") ||
		   (len(path) > 1 && path[1] == ':')
}

func (m *Mo) uploadFile(localPath, remotePath string) error {
	if m.config.Bucket == "" {
		return fmt.Errorf("bucket is required for upload")
	}

	printInfo("uploading %s to %s/%s...", localPath, m.config.Bucket, remotePath)

	file, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("failed to open local file: %v", err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat local file: %v", err)
	}

	ctx := context.Background()
	
	contentType := "application/octet-stream"
	if strings.HasSuffix(localPath, ".png") {
		contentType = "image/png"
	} else if strings.HasSuffix(localPath, ".jpg") || strings.HasSuffix(localPath, ".jpeg") {
		contentType = "image/jpeg"
	} else if strings.HasSuffix(localPath, ".txt") {
		contentType = "text/plain"
	}

	_, err = m.client.PutObject(ctx, m.config.Bucket, remotePath, file, stat.Size(), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("upload failed: %v", err)
	}

	printSuccess("successfully uploaded %s (%s)", localPath, formatSize(stat.Size()))
	return nil
}

func (m *Mo) downloadFile(remotePath, localPath string) error {
	if m.config.Bucket == "" {
		return fmt.Errorf("bucket required for download retard lol")
	}

	printInfo("downloading %s/%s to %s...", m.config.Bucket, remotePath, localPath)

	ctx := context.Background()
	object, err := m.client.GetObject(ctx, m.config.Bucket, remotePath, minio.GetObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to get object: %v", err)
	}
	defer object.Close()

	// dir := filepath.Dir(localPath)
	// if err := os.MkdirAll(dir, 0755); err != nil {
	// 	return fmt.Errorf("failed to create directory: %v", err)
	// }

	localFile, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create local file: %v", err)
	}
	defer localFile.Close()

	written, err := io.Copy(localFile, object)
	if err != nil {
		return fmt.Errorf("download failed: %v", err)
	}

	printSuccess("successfully downloaded %s (%s)", remotePath, formatSize(written))
	return nil
}

func (m *Mo) copyRemoteObject(src, dst string) error {
	if m.config.Bucket == "" {
		return fmt.Errorf("bucket r equired for copy")
	}

	printInfo("copying %s to %s...", src, dst)

	ctx := context.Background()
	_, err := m.client.CopyObject(ctx, minio.CopyDestOptions{
		Bucket: m.config.Bucket,
		Object: dst,
	}, minio.CopySrcOptions{
		Bucket: m.config.Bucket,
		Object: src,
	})
	if err != nil {
		return fmt.Errorf("copy failed: %v", err)
	}

	printSuccess("successfully copied %s to %s", src, dst)
	return nil
}

func (m *Mo) removeObject(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("rm requires object path argument")
	}

	if m.config.Bucket == "" {
		return fmt.Errorf("bucket is required for rm command")
	}

	objectPath := args[0]
	printInfo("removing %s/%s...", m.config.Bucket, objectPath)

	ctx := context.Background()
	err := m.client.RemoveObject(ctx, m.config.Bucket, objectPath, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("remove failed: %v", err)
	}

	printSuccess("successfully removed %s", objectPath)
	return nil
}

func (m *Mo) moveObject(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("mv requires src & dst arguments")
	}

	src := args[0]
	dst := args[1]

	if m.config.Bucket == "" {
		return fmt.Errorf("bucket is required dumbass")
	}

	printInfo("moving %s to %s...", src, dst)

	ctx := context.Background()
	
	_, err := m.client.CopyObject(ctx, minio.CopyDestOptions{
		Bucket: m.config.Bucket,
		Object: dst,
	}, minio.CopySrcOptions{
		Bucket: m.config.Bucket,
		Object: src,
	})
	if err != nil {
		return fmt.Errorf("FAILED TO MOVE: FAILED DURING COPY: %v", err)
	}

	err = m.client.RemoveObject(ctx, m.config.Bucket, src, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("move failed while removing src: %v", err)
	}

	printSuccess("successfully moved %s to %s", src, dst)
	return nil
}

func (m *Mo) statObject(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("stat requires object path argument")
	}

	if m.config.Bucket == "" {
		return fmt.Errorf("bucket is required for stat command")
	}

	objectPath := args[0]
	printInfo("getting metadata for %s/%s...", m.config.Bucket, objectPath)

	ctx := context.Background()
	objInfo, err := m.client.StatObject(ctx, m.config.Bucket, objectPath, minio.StatObjectOptions{})
	if err != nil {
		return fmt.Errorf("stat failed: %v", err)
	}

	printSuccess("object metadata:")
	fmt.Printf("  %s%sname:%s        %s%s%s\n", Bold, Green, Reset, Yellow, objInfo.Key, Reset)
	fmt.Printf("  %s%ssize:%s        %s%s%s\n", Bold, Green, Reset, Cyan, formatSize(objInfo.Size), Reset)
	fmt.Printf("  %s%setag:%s        %s%s%s\n", Bold, Green, Reset, Magenta, objInfo.ETag, Reset)
	fmt.Printf("  %s%scontent-type:%s %s%s%s\n", Bold, Green, Reset, Blue, objInfo.ContentType, Reset)
	fmt.Printf("  %s%slast modified:%s %s%s%s\n", Bold, Green, Reset, Dim, objInfo.LastModified.Format("2006-01-02 15:04:05 MST"), Reset)

	if len(objInfo.UserMetadata) > 0 {
		fmt.Printf("  %s%suser metadata:%s\n", Bold, Green, Reset)
		for key, value := range objInfo.UserMetadata {
			fmt.Printf("    %s%s:%s %s%s%s\n", Yellow, key, Reset, White, value, Reset)
		}
	}

	return nil
}

func (m *Mo) makeBucket(args []string) error {
	bucketName := ""
	if len(args) > 0 {
		bucketName = args[0]
	} else if m.config.Bucket != "" {
		bucketName = m.config.Bucket
	} else {
		return fmt.Errorf("bucket name is required for mb command")
	}

	fmt.Printf("%swarnign:%s about to create bucket '%s%s%s'. proceed? (%syes please%s / %sno%s): ", 
		Yellow, Reset, Yellow, bucketName, Reset, Green, Reset, Red, Reset)
	
	var response string
	fmt.Scanln(&response)
	if response != "yes" && response != "yes please" {
		printInfo("operation cancelled")
		return nil
	}

	// confirming this shit due to the chance of being victim to a DDoW:
	// Distributed Denial of Wallet

	fmt.Printf("\n%splease confirm credentials:%s\n", Bold, Reset)
	fmt.Printf("endpoint: %s%s%s\n", Cyan, m.config.Endpoint, Reset)
	fmt.Printf("region: %s%s%s\n", Cyan, m.config.Region, Reset)
	fmt.Printf("access key: %s%s%s\n", Cyan, m.config.AccessKey, Reset)
	fmt.Printf("secret key: %s%s...%s\n", Cyan, m.config.SecretKey[:min(len(m.config.SecretKey), 8)], Reset)
	
	fmt.Printf("\n%sfinal confirmation:%s create bucket with these settings? (%syes please kthnxbye%s / %sno%s): ", 
		Red, Reset, Green, Reset, Red, Reset)
	
	fmt.Scanln(&response)
	if response != "yes please kthnxbye" {
		printInfo("operation cancelled")
		return nil
	}

	printInfo("creating bucket %s...", bucketName)

	ctx := context.Background()
	err := m.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{
		Region: m.config.Region,
	})
	if err != nil {
		return fmt.Errorf("make bucket failed: %v", err)
	}

	printSuccess("successfully created bucket %s", bucketName)
	return nil
}

func (m *Mo) removeBucket(args []string) error {
	bucketName := ""
	if len(args) > 0 {
		bucketName = args[0]
	} else if m.config.Bucket != "" {
		bucketName = m.config.Bucket
	} else {
		return fmt.Errorf("bucket name is required for rb command")
	}

	fmt.Printf("%s%sDANGER:%s about to DELETE bucket '%s%s%s'. this is irreversible! proceed? (%syes please%s / %sno%s): ", 
		Bold, Red, Reset, Red, bucketName, Reset, Green, Reset, Red, Reset)
	
	var response string
	fmt.Scanln(&response)
	if response != "yes" && response != "yes please" {
		printInfo("operation cancelled")
		return nil
	}

	fmt.Printf("\n%splease confirm credentials:%s\n", Bold, Reset)
	fmt.Printf("endpoint: %s%s%s\n", Cyan, m.config.Endpoint, Reset)
	fmt.Printf("region: %s%s%s\n", Cyan, m.config.Region, Reset)
	fmt.Printf("access key: %s%s%s\n", Cyan, m.config.AccessKey, Reset)
	fmt.Printf("secret key: %s%s...%s\n", Cyan, m.config.SecretKey[:min(len(m.config.SecretKey), 8)], Reset)
	
	fmt.Printf("\n%sfinal confirmation:%s DELETE bucket with these settings? (%syes please kthnxbye%s / %sno%s): ", 
		Red, Reset, Green, Reset, Red, Reset)
	
	fmt.Scanln(&response)
	if response != "yes please kthnxbye" {
		printInfo("operation cancelled")
		return nil
	}

	printInfo("removing bucket %s...", bucketName)

	ctx := context.Background()
	err := m.client.RemoveBucket(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("remove bucket failed: %v", err)
	}

	printSuccess("successfully removed bucket %s", bucketName)
	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func printInfo(format string, args ...interface{}) {
	fmt.Printf("%s[info]%s     %s\n", Blue, Reset, fmt.Sprintf(format, args...))
}

func printSuccess(format string, args ...interface{}) {
	fmt.Printf("%s[success]%s  %s\n", Green, Reset, fmt.Sprintf(format, args...))
}

func printWarning(format string, args ...interface{}) {
	fmt.Printf("%s[warning]%s  %s\n", Yellow, Reset, fmt.Sprintf(format, args...))
}

func printError(format string, args ...interface{}) {
	fmt.Printf("%s[error]%s    %s\n", Red, Reset, fmt.Sprintf(format, args...))
}

func printHelpShit() {
	fmt.Printf(`mo - s3-compatible object storage manager


usage:
    mo [options] command [args...]


opts:
    -endpoint URL     s3 endpoint url
    -region REGION    s3 region (default: us-east-1)
    -key KEY          access key id
    -secret SECRET    secret access key
    -bucket BUCKET    default bucket name
    -insecure         use http instead of https
    -r                recursive copy (for cp command)


cmds:
    buckets           list all buckets
    ls [prefix]       list objects in bucket
    cp SOURCE DEST    copy/upload/download files
    rm OBJECT         remove(delet) object
    mv SOURCE DEST    move/rename object
    stat OBJECT       get object metadata
    mb [bucket]       make bucket
    rb [bucket]       remove bucket


doc/example(s(?)):
    mo -endpoint s3.amazonaws.com -key YOUARERETARDEDLOLHEH -secret wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY -bucket somebucket ls
    mo -endpoint localhost:9000 -key minioadmin -secret minioadmin -bucket test cp ./file.txt remote/file.txt
    mo buckets
    mo ls myfolder/
    mo cp ./localfile.txt remote/path/file.txt
    mo cp remote/file.txt ./localfile.txt
    mo -r cp ./mydirectory/ remote/backup/
    mo rm remote/file.txt
    mo mv old/path.txt new/path.txt
    mo stat remote/file.txt


configuration file if you can't handle typing credentials unlike gigachads:
    create a file 'mo.conf' in the $CWD/$PWD w/e 
    
    endpoint s3.amazonaws.com
    region us-east-1
    access_key AKIAIOSFODNN7EXAMPLE
    secret_key wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
    bucket mybucket


environment variables (overrides previouus opts):
    MO_ENDPOINT, MO_REGION, MO_ACCESS_KEY, MO_SECRET_KEY, MO_BUCKET

`)
}
