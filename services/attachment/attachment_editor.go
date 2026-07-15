// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package attachment

import (
	"context"

	repo_model "gitea.dev/models/repo"
	"gitea.dev/modules/setting"
)

// UploadAttachmentForEditor uploads a file attached via the repo file editor's dropzone (e.g.
// an image embedded in a markdown file), using the repo upload settings (shared with the
// "Upload files" feature) for allowed types and size limits.
func UploadAttachmentForEditor(ctx context.Context, file *UploaderFile, attach *repo_model.Attachment) (*repo_model.Attachment, error) {
	return uploadAttachment(ctx, file, setting.Repository.Upload.AllowedTypes, setting.Repository.Upload.FileMaxSize<<20, attach)
}
