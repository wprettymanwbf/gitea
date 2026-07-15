import {initDropzone, DropzoneCustomEventUploadDone, generateMarkdownLinkForAttachment} from './dropzone.ts';
import {imageInfo} from '../utils/image.ts';
import type {EditorView} from '@codemirror/view';

type CodeEditor = {view: EditorView};

let placeholderIdCounter = 0;

function getPastedImages(e: ClipboardEvent): Array<File> {
  const images: Array<File> = [];
  for (const item of e.clipboardData?.items ?? []) {
    if (item.type?.startsWith('image/')) {
      const file = item.getAsFile();
      if (file) images.push(file);
    }
  }
  return images;
}

function insertAtCursor(view: EditorView, text: string) {
  view.dispatch(view.state.replaceSelection(text));
  view.focus();
}

function replacePlaceholder(view: EditorView, placeholder: string, text: string) {
  const idx = view.state.doc.toString().indexOf(placeholder);
  if (idx < 0) {
    insertAtCursor(view, text);
    return;
  }
  view.dispatch({changes: {from: idx, to: idx + placeholder.length, insert: text}});
}

async function addFilesToDropzone(view: EditorView, dzInst: any, files: Array<File> | FileList) {
  for (const file of files) {
    // A unique placeholder lets us replace exactly this file's text once its upload completes.
    const placeholder = `![Uploading ${file.name}…#${placeholderIdCounter++}]()`;
    (file as any)._giteaEditorPlaceholder = placeholder;
    const {width, dppx} = await imageInfo(file);
    (file as any)._giteaEditorImageInfo = {width, dppx};
    insertAtCursor(view, placeholder);
    dzInst.addFile(file);
  }
}

// initEditorUpload wires the file-editor dropzone to the code editor so that dragging, pasting
// or picking a file uploads it as a repo attachment and inserts a markdown link (referencing
// the attachment by URL) at the cursor. The attachment becomes permanently linked to the repo
// once the edit is committed (see repo_model.LinkAttachmentsToRepoCode on the backend).
export async function initEditorUpload(editor: CodeEditor, container: HTMLElement) {
  const dropzoneEl = container.querySelector<HTMLElement>('.editor-upload-dropzone .dropzone');
  if (!dropzoneEl) return;

  const view = editor.view;
  const dzInst = await initDropzone(dropzoneEl);

  // Insert a link when an upload finishes. This covers the dropzone's own file picker as well as
  // drag-and-drop and paste (which pre-insert a placeholder that we replace here).
  dzInst.on(DropzoneCustomEventUploadDone, ({file}: {file: any}) => {
    const link = generateMarkdownLinkForAttachment(file, file._giteaEditorImageInfo ?? {});
    if (file._giteaEditorPlaceholder) {
      replacePlaceholder(view, file._giteaEditorPlaceholder, link);
    } else {
      insertAtCursor(view, link);
    }
  });

  view.dom.addEventListener('paste', (e: ClipboardEvent) => {
    const images = getPastedImages(e);
    if (!images.length) return;
    e.preventDefault();
    addFilesToDropzone(view, dzInst, images);
  });

  view.dom.addEventListener('drop', (e: DragEvent) => {
    if (!e.dataTransfer?.files.length) return;
    e.preventDefault();
    addFilesToDropzone(view, dzInst, e.dataTransfer.files);
  });
}
