export type FileHandleKey = 'file-share' | 'file-image-compress'
export type FileShareHandleProps = { type: FileHandleKey; config: Record<string, any> }

export type TextHandleKey = 'text-share'
export type TextShareHandleProps = { type: TextHandleKey; config: Record<string, any> }
