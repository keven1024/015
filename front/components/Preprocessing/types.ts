export type FileHandleKey = 'file-share' | 'file-image-compress' | 'file-image-convert'
export type FileShareHandleProps = { type: FileHandleKey; config: Record<string, any> }

export type TextHandleKey = 'text-share' | 'text-translate'
export type TextShareHandleProps = { type: TextHandleKey; config: Record<string, any> }
