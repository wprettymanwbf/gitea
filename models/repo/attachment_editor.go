// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package repo

import (
	"context"

	"gitea.dev/models/db"
)

// LinkAttachmentsToRepoCode marks the given attachments (identified by UUID, scoped to
// repoID and uploaderID) as linked to repo code, so they survive the uploader-account-deletion
// cleanup that would otherwise reap them as abandoned uploads. UUIDs that don't match repoID
// or uploaderID, or that are already linked to an issue/release, are silently dropped, the
// same convention used by GetAttachmentsByUUIDs.
func LinkAttachmentsToRepoCode(ctx context.Context, repoID, uploaderID int64, uuids []string) error {
	if len(uuids) == 0 {
		return nil
	}

	_, err := db.GetEngine(ctx).
		Where("repo_id = ? AND uploader_id = ? AND issue_id = 0 AND release_id = 0", repoID, uploaderID).
		In("uuid", uuids).
		Cols("editor_linked").
		Update(&Attachment{EditorLinked: true})
	return err
}
