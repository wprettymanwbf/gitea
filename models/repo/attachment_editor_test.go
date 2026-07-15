// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package repo_test

import (
	"testing"

	"gitea.dev/models/db"
	repo_model "gitea.dev/models/repo"
	"gitea.dev/models/unittest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLinkAttachmentsToRepoCode(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	own := &repo_model.Attachment{UUID: "editor-test-own", RepoID: 1, UploaderID: 2, Name: "img.png"}
	require.NoError(t, db.Insert(t.Context(), own))
	otherRepo := &repo_model.Attachment{UUID: "editor-test-other-repo", RepoID: 2, UploaderID: 2, Name: "img.png"}
	require.NoError(t, db.Insert(t.Context(), otherRepo))
	otherUploader := &repo_model.Attachment{UUID: "editor-test-other-uploader", RepoID: 1, UploaderID: 4, Name: "img.png"}
	require.NoError(t, db.Insert(t.Context(), otherUploader))

	err := repo_model.LinkAttachmentsToRepoCode(t.Context(), 1, 2, []string{
		"editor-test-own", "editor-test-other-repo", "editor-test-other-uploader", "editor-test-missing",
	})
	require.NoError(t, err)

	linked, err := repo_model.GetAttachmentByUUID(t.Context(), "editor-test-own")
	require.NoError(t, err)
	assert.True(t, linked.EditorLinked)

	// UUIDs that don't match repoID/uploaderID must not be linked
	notLinked, err := repo_model.GetAttachmentByUUID(t.Context(), "editor-test-other-repo")
	require.NoError(t, err)
	assert.False(t, notLinked.EditorLinked)

	notLinked, err = repo_model.GetAttachmentByUUID(t.Context(), "editor-test-other-uploader")
	require.NoError(t, err)
	assert.False(t, notLinked.EditorLinked)

	// a linked attachment must no longer be returned as "unlinked" (survives account-deletion cleanup)
	unlinked, err := repo_model.GetUnlinkedAttachmentsByUserID(t.Context(), 2)
	require.NoError(t, err)
	for _, a := range unlinked {
		assert.NotEqual(t, "editor-test-own", a.UUID)
	}
}
