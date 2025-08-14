# gcsls - Google Cloud Storage List with Wildcards

A command-line tool for listing Google Cloud Storage (GCS) objects with advanced wildcard pattern matching support, including recursive `**` patterns.

## Features

- 🔍 **Wildcard Support**: Use `*`, `?`, and `**` patterns to match objects
- 🚀 **Efficient Listing**: Optimizes GCS API calls by using prefixes when possible
- 🔐 **ADC Authentication**: Uses Google Cloud Application Default Credentials
- 📁 **Recursive Matching**: Support for `**` to match across directory levels
- 🎯 **Precise Filtering**: Client-side pattern matching for exact results

## Installation

### Prerequisites

- Go 1.16 or later
- Google Cloud SDK (gcloud) installed and configured
- Access to Google Cloud Storage buckets

### Build from Source

```bash
git clone https://github.com/biolog71/gcsls.git
cd gcsls
go build -o gcsls main.go
```

### Install with Go

```bash
go install github.com/biolog71/gcsls@latest
```

## Authentication

Before using gcsls, ensure you're authenticated with Google Cloud:

```bash
# Authenticate with your user account
gcloud auth application-default login

# OR set up service account credentials
export GOOGLE_APPLICATION_CREDENTIALS="/path/to/service-account-key.json"
```

## Usage

```bash
gcsls "gs://bucket-name/pattern"
```

### Basic Examples

```bash
# List all objects in a bucket
gcsls "gs://my-bucket"

# List all CSV files in a specific folder
gcsls "gs://my-bucket/data/*.csv"

# List all log files recursively
gcsls "gs://my-bucket/logs/**/*.log"

# Find specific files in any subdirectory
gcsls "gs://my-bucket/**/config.json"

# List files with specific naming pattern
gcsls "gs://my-bucket/backup-*.tar.gz"
```

### Advanced Pattern Examples

```bash
# Files starting with "data" and ending with any extension
gcsls "gs://my-bucket/folder/data*.*"

# All files in directories that start with "202"
gcsls "gs://my-bucket/202*/**"

# Specific file patterns across multiple directory levels
gcsls "gs://my-bucket/**/logs/**/*.txt"

# Files with single character wildcards
gcsls "gs://my-bucket/file?.log"
```

## Wildcard Patterns

| Pattern | Description | Example |
|---------|-------------|---------|
| `*` | Matches any sequence of characters (except `/`) | `*.csv` matches `data.csv` |
| `**` | Matches any sequence including `/` (recursive) | `**/logs` matches `a/b/logs` |
| `?` | Matches exactly one character | `file?.txt` matches `file1.txt` |
| `[abc]` | Matches any character in brackets | `file[123].txt` matches `file2.txt` |

## Output Format

The tool outputs matching GCS paths in the format:
```
gs://bucket-name/path/to/object
```

If no objects match the pattern, it displays:
```
No objects found matching the pattern.
```

## Dependencies

- [cloud.google.com/go/storage](https://pkg.go.dev/cloud.google.com/go/storage) - Google Cloud Storage client library
- [github.com/bmatcuk/doublestar/v4](https://pkg.go.dev/github.com/bmatcuk/doublestar/v4) - Advanced glob pattern matching with `**` support

## Error Handling

The tool provides clear error messages for common issues:

- **Invalid GCS path**: Path must start with `gs://`
- **Missing bucket name**: Bucket name is required
- **Authentication errors**: Check your GCloud authentication
- **Invalid patterns**: Malformed glob patterns will be reported
- **Access denied**: Ensure you have permissions to list objects in the bucket

## Performance Considerations

- The tool optimizes GCS API calls by extracting prefixes from patterns
- For patterns like `logs/**/*.txt`, only objects with prefix `logs/` are fetched
- Client-side filtering ensures exact pattern matching
- Large buckets with broad patterns may take longer to process

## Examples in Practice

### Data Analysis Workflows
```bash
# Find all CSV files for processing
gcsls "gs://data-lake/raw/**/*.csv"

# Locate specific date-partitioned files
gcsls "gs://analytics/events/2024/01/**/*.parquet"
```

### Log Management
```bash
# Find error logs across all services
gcsls "gs://app-logs/**/error*.log"

# Locate logs for a specific date pattern
gcsls "gs://logs/*/2024-01-??.log"
```

### Backup Verification
```bash
# Check for backup files
gcsls "gs://backups/**/backup-*.tar.gz"

# Verify daily snapshots
gcsls "gs://snapshots/*/snapshot-$(date +%Y-%m-%d)*"
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

If you encounter any issues or have questions:

1. Check the [Issues](https://github.com/biolog71/gcsls/issues) page
2. Create a new issue with detailed information about your problem
3. Include the command you ran and the error message received

## Changelog

### v1.0.0
- Initial release
- Basic wildcard pattern matching
- Recursive `**` support
- GCS integration with ADC authentication