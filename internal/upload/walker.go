package upload

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/facebookgo/symwalk"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/log"
)

const defaultFileListPath = "/content/drive/MyDrive/gphotos-uploader/config-folder/filelist.txt"

// ScanFolder return the list of Items{} to be uploaded. It scans the folder and skip
// non allowed files (includePatterns & excludePattens).
func (job *UploadFolderJob) ScanFolderFromList(logger log.Logger) ([]FileItem, error) {
	var result []FileItem

	fileList, err := os.Open(defaultFileListPath)
	if err != nil {
		return nil, err
	}
	defer fileList.Close()

	scanner := bufio.NewScanner(fileList)
	for scanner.Scan() {
		filePath := strings.TrimSpace(scanner.Text())
		if filePath == "" {
			continue
		}

		relativePath := RelativePath(job.SourceFolder, filePath)

		// Skip already uploaded files
		if job.FileTracker.IsUploaded(filePath) {
			logger.Debugf("Skipping already uploaded file '%s'.", filePath)
			continue
		}

		albumName := job.albumName(relativePath)
		logger.Debugf("Adding file '%s' to the upload list for album '%s'.", filePath, albumName)

		result = append(result, FileItem{
			Path:      filePath,
			AlbumName: albumName,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// RelativePath returns a path relative to the base.
// If a relative path could not be calculated or it contains ' ../`,
// returns the original path.
func RelativePath(base string, path string) string {
	rp, err := filepath.Rel(base, path)
	if err != nil {
		return path
	}
	if strings.HasPrefix(rp, "../") {
		return path
	}
	return rp
}
