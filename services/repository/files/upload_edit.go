// Copyright 2025 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package files

import (
	"context"
	"fmt"
	"path"

	repo_model "gitea.dev/models/repo"
)

// UploadsToChangeRepoFiles resolves the given staged upload UUIDs into "upload"
// ChangeRepoFile operations placed under dirPath. It returns the operations along
// with the resolved Upload rows so the caller can DeleteUploads after committing.
// This mirrors the loop in UploadRepoFiles but is reusable from callers (such as the
// file editor) that commit uploads together with other file changes.
func UploadsToChangeRepoFiles(ctx context.Context, dirPath string, uuids []string) ([]*ChangeRepoFile, []*repo_model.Upload, error) {
	if len(uuids) == 0 {
		return nil, nil, nil
	}

	uploads, err := repo_model.GetUploadsByUUIDs(ctx, uuids)
	if err != nil {
		return nil, nil, fmt.Errorf("GetUploadsByUUIDs [uuids: %v]: %w", uuids, err)
	}

	files := make([]*ChangeRepoFile, 0, len(uploads))
	for _, upload := range uploads {
		files = append(files, &ChangeRepoFile{
			Operation:     "upload",
			TreePath:      path.Join(dirPath, upload.Name),
			ContentReader: &lazyLocalFileReader{localFilename: upload.LocalPath()},
		})
	}
	return files, uploads, nil
}
