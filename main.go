package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/bmatcuk/doublestar/v4"
	"google.golang.org/api/iterator"
)

// showHelp displays the usage information for the tool.
func showHelp() {
	fmt.Printf("gcsls - List Google Cloud Storage objects with wildcard support\n\n")
	fmt.Printf("USAGE:\n")
	fmt.Printf("  %s [OPTIONS] \"gs://bucket/object-pattern\"\n\n", os.Args[0])
	fmt.Printf("OPTIONS:\n")
	fmt.Printf("  -h, --help    Show this help message and exit\n\n")
	fmt.Printf("EXAMPLES:\n")
	fmt.Printf("  %s \"gs://my-bucket/logs/**/*.log\"\n", os.Args[0])
	fmt.Printf("  %s \"gs://my-bucket/data/*.csv\"\n", os.Args[0])
	fmt.Printf("  %s \"gs://my-bucket/folder/**/data.txt\"\n", os.Args[0])
	fmt.Printf("  %s \"gs://my-bucket/\"\n\n", os.Args[0])
	fmt.Printf("DESCRIPTION:\n")
	fmt.Printf("  This tool lists objects in Google Cloud Storage that match a given pattern.\n")
	fmt.Printf("  It supports glob patterns including:\n")
	fmt.Printf("    *     - matches any sequence of characters (except /)\n")
	fmt.Printf("    **    - matches any sequence of characters (including /)\n")
	fmt.Printf("    ?     - matches any single character\n")
	fmt.Printf("    [abc] - matches any character in the set\n\n")
	fmt.Printf("AUTHENTICATION:\n")
	fmt.Printf("  Ensure you have authenticated with Google Cloud:\n")
	fmt.Printf("    gcloud auth application-default login\n")
}

// main is the entry point of the program.
// It expects exactly one command-line argument: a GCS path like gs://bucket-name/prefix.
// Example Usage:
// go run . "gs://my-bucket/some-folder/*.csv"
// go run . "gs://my-bucket/some-folder/**/data.txt"
func main() {
	// Check for help flags first
	if len(os.Args) == 2 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		showHelp()
		os.Exit(0)
	}

	// Check for the correct number of command-line arguments.
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s \"gs://bucket/object-pattern\"\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Example: %s \"gs://my-bucket/logs/**/*.log\"\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Use -h or --help for more information.\n")
		os.Exit(1)
	}

	gcsPath := os.Args[1]

	// The context is used to manage the lifecycle of API requests.
	ctx := context.Background()

	// Call the core logic function and handle any errors.
	if err := listObjectsWithWildcard(ctx, gcsPath); err != nil {
		log.Fatalf("Failed to list objects: %v", err)
	}
}

// listObjectsWithWildcard lists objects in GCS that match a given path with wildcards.
func listObjectsWithWildcard(ctx context.Context, gcsPath string) error {
	// --- 1. Parse the GCS Path ---
	// The path must start with "gs://".
	if !strings.HasPrefix(gcsPath, "gs://") {
		return fmt.Errorf("invalid GCS path: must start with gs://")
	}

	// Remove the "gs://" prefix to work with the bucket and object path.
	pathWithoutScheme := strings.TrimPrefix(gcsPath, "gs://")

	// Split the path into bucket name and the object pattern.
	parts := strings.SplitN(pathWithoutScheme, "/", 2)
	if len(parts) == 0 || parts[0] == "" {
		return fmt.Errorf("invalid GCS path: bucket name is missing")
	}
	bucketName := parts[0]
	objectPattern := ""
	if len(parts) > 1 {
		objectPattern = parts[1]
	}

	// If the pattern is empty, it means we should list everything in the bucket.
	// We'll use the "**" wildcard for this, which matches everything recursively.
	if objectPattern == "" {
		objectPattern = "**"
	}

	// --- 2. Initialize GCS Client ---
	// This uses Application Default Credentials (ADC) to authenticate.
	// Ensure you have authenticated via `gcloud auth application-default login`
	// or that the environment is configured with a service account.
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create GCS client: %w", err)
	}
	defer client.Close()

	bucket := client.Bucket(bucketName)

	// --- 3. Determine Prefix for API Query ---
	// To make the GCS API call more efficient, we find the part of the pattern
	// before any wildcards. This reduces the number of objects we have to
	// process client-side.
	prefix := getPrefixFromPattern(objectPattern)

	query := &storage.Query{
		Prefix: prefix,
	}

	// --- 4. Iterate and Filter ---
	fmt.Printf("Listing objects in gs://%s matching pattern: %s\n", bucketName, objectPattern)

	it := bucket.Objects(ctx, query)
	found := false
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			// End of the results.
			break
		}
		if err != nil {
			return fmt.Errorf("failed to iterate objects: %w", err)
		}

		// Client-side filtering using the doublestar library, which supports "**".
		matched, err := doublestar.Match(objectPattern, attrs.Name)
		if err != nil {
			return fmt.Errorf("invalid glob pattern '%s': %w", objectPattern, err)
		}

		if matched {
			fmt.Printf("gs://%s/%s\n", bucketName, attrs.Name)
			found = true
		}
	}

	if !found {
		fmt.Println("No objects found matching the pattern.")
	}

	return nil
}

// getPrefixFromPattern extracts the part of a string before the first wildcard character.
// Wildcards are considered to be '*', '?', and '['.
func getPrefixFromPattern(pattern string) string {
	wildcardIndex := strings.IndexAny(pattern, "*?[")
	if wildcardIndex == -1 {
		// No wildcards, the whole pattern is a prefix.
		return pattern
	}
	// Return the substring up to the first wildcard.
	return pattern[:wildcardIndex]
}
