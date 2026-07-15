// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package v1_27

import (
	"testing"

	"gitea.dev/models/migrations/migrationtest"

	"github.com/stretchr/testify/require"
)

func Test_AddEditorLinkedToAttachment(t *testing.T) {
	type Attachment struct {
		ID         int64  `xorm:"pk autoincr"`
		UUID       string `xorm:"uuid UNIQUE"`
		RepoID     int64  `xorm:"INDEX"`
		UploaderID int64  `xorm:"INDEX DEFAULT 0"`
	}

	x, deferable := migrationtest.PrepareTestEnv(t, 0, new(Attachment))
	defer deferable()

	_, err := x.Insert(&Attachment{UUID: "uuid-1", RepoID: 1, UploaderID: 2})
	require.NoError(t, err)

	require.NoError(t, AddEditorLinkedToAttachment(x))

	type AttachmentAfter struct {
		ID           int64  `xorm:"pk autoincr"`
		UUID         string `xorm:"uuid UNIQUE"`
		RepoID       int64  `xorm:"INDEX"`
		UploaderID   int64  `xorm:"INDEX DEFAULT 0"`
		EditorLinked bool   `xorm:"NOT NULL DEFAULT false"`
	}

	var a AttachmentAfter
	has, err := x.Table("attachment").Where("uuid = ?", "uuid-1").Get(&a)
	require.NoError(t, err)
	require.True(t, has)
	require.False(t, a.EditorLinked)
}
