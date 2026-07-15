// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package v1_27

import (
	"gitea.dev/models/db"

	"xorm.io/xorm"
)

// AddEditorLinkedToAttachment distinguishes an attachment linked to repo code (embedded in a
// committed file via the file editor) from an abandoned/orphaned upload. Unlike issue or
// release attachments, editor attachments have no IssueID/ReleaseID to link against, so
// without this marker they would be indistinguishable from an orphaned upload: restricted to
// uploader-only viewing and reaped when the uploader's account is deleted.
func AddEditorLinkedToAttachment(x db.EngineMigration) error {
	type Attachment struct {
		EditorLinked bool `xorm:"NOT NULL DEFAULT false"`
	}

	_, err := x.SyncWithOptions(xorm.SyncOptions{
		IgnoreConstrains: true,
		IgnoreIndices:    true,
	}, new(Attachment))
	return err
}
