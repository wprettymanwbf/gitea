import {initDropzone, DropzoneCustomEventUploadDone} from './dropzone.ts';
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

// Uploaded files are committed alongside the edit into the same directory as the edited file,
// so the inserted link is the (relative) filename rather than an /attachments/ URL.
function generateRelativeMarkdownLink(file: {name: string, type?: string}): string {
  const encodedName = encodeURIComponent(file.name);
  if (file.type?.startsWith('image/')) {
    const alt = file.name.slice(0, file.name.lastIndexOf('.')) || file.name;
    return `![${alt}](${encodedName})`;
  }
  return `[${file.name}](${encodedName})`;
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

function addFilesToDropzone(view: EditorView, dzInst: any, files: Array<File> | FileList) {
  for (const file of files) {
    // A unique placeholder lets us replace exactly this file's text once its upload completes.
    const placeholder = `![Uploading ${file.name}…#${placeholderIdCounter++}]()`;
    (file as any)._giteaEditorPlaceholder = placeholder;
    insertAtCursor(view, placeholder);
    dzInst.addFile(file);
  }
}

// initEditorUpload wires the file-editor dropzone to the code editor so that dragging, pasting
// or picking a file uploads it (staged via the repo upload endpoint) and inserts a relative
// markdown link at the cursor. The uploaded file is committed together with the edit.
export async function initEditorUpload(editor: CodeEditor, container: HTMLElement) {
  const dropzoneEl = container.querySelector<HTMLElement>('.editor-upload-dropzone .dropzone');
  if (!dropzoneEl) return;

  const view = editor.view;
  const dzInst = await initDropzone(dropzoneEl);

  // Insert a link when an upload finishes. This covers the dropzone's own file picker as well as
  // drag-and-drop and paste (which pre-insert a placeholder that we replace here).
  dzInst.on(DropzoneCustomEventUploadDone, ({file}: {file: any}) => {
    const link = generateRelativeMarkdownLink(file);
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
