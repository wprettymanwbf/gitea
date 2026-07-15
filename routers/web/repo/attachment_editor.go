// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package repo

import (
	"net/http"

	repo_model "gitea.dev/models/repo"
	"gitea.dev/modules/log"
	"gitea.dev/services/attachment"
	"gitea.dev/services/context"
	"gitea.dev/services/context/upload"
)

// EditorUploadAttachment handles file uploads from the file editor's dropzone (e.g. images
// embedded in a markdown file). The uploaded attachment is stored outside git and referenced
// by URL; it is only permanently linked to the repo once the edit referencing it is committed
// (see repo_model.LinkAttachmentsToRepoCode in EditFilePost).
func EditorUploadAttachment(ctx *context.Context) {
	file, header, err := ctx.Req.FormFile("file")
	if err != nil {
		ctx.ServerError("FormFile", err)
		return
	}
	defer file.Close()

	uploaderFile := attachment.NewLimitedUploaderKnownSize(file, header.Size)
	attach, err := attachment.UploadAttachmentForEditor(ctx, uploaderFile, &repo_model.Attachment{
		Name:       header.Filename,
		UploaderID: ctx.Doer.ID,
		RepoID:     ctx.Repo.Repository.ID,
	})
	if err != nil {
		if upload.IsErrFileTypeForbidden(err) {
			ctx.HTTPError(http.StatusBadRequest, err.Error())
			return
		}
		ctx.ServerError("EditorUploadAttachment", err)
		return
	}

	log.Trace("New editor attachment uploaded: %s", attach.UUID)
	ctx.JSON(http.StatusOK, map[string]string{
		"uuid": attach.UUID,
	})
}
