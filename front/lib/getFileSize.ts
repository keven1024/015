import { filesize } from 'filesize'

const getFileSize = (size: number) => {
    return filesize(size, { standard: 'jedec' })
}

export default getFileSize
