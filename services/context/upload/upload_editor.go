// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package upload

import (
	"strings"

	repo_model "gitea.dev/models/repo"
	"gitea.dev/modules/reqctx"
	"gitea.dev/modules/setting"
)

// AddUploadContextForRepoEditorAttachment renders the dropzone template values for the file
// editor's "attach files" dropzone, pointing it at the editor-attachment endpoints and reusing
// the repo upload settings (setting.Repository.Upload.*) for limits, consistent with the
// "Upload files" feature.
func AddUploadContextForRepoEditorAttachment(ctx reqctx.RequestContext, repo *repo_model.Repository) {
	ctxData, repoLink := ctx.GetData(), repo.Link()
	ctxData["UploadUrl"] = repoLink + "/editor/attachments"
	ctxData["UploadRemoveUrl"] = repoLink + "/editor/attachments/remove"
	ctxData["UploadLinkUrl"] = repoLink + "/editor/attachments"
	ctxData["UploadAccepts"] = strings.ReplaceAll(setting.Repository.Upload.AllowedTypes, "|", ",")
	ctxData["UploadMaxFiles"] = setting.Repository.Upload.MaxFiles
	ctxData["UploadMaxSize"] = setting.Repository.Upload.FileMaxSize
}
