import { noop } from 'lodash-es'
import { md5 } from 'js-md5'

interface CalcFileHashProps {
    file: File
    onProgress?: (current: number) => void
    chunkSize?: number
}

const calcFileHash = async (props: CalcFileHashProps) => {
    const { file, onProgress = noop, chunkSize = 100 } = props || {}
    const blob = await file.arrayBuffer()
    const hash = md5(blob)
    return hash
    // const finalChunkSize = chunkSize * 1024 * 1024;
    // const chunks = Math.ceil(file.size / finalChunkSize);
    // const spark = new SparkMD5.ArrayBuffer(); // 使用 SparkMD5 增量计算哈希
    // const fileReader = new FileReader();

    // const readChunk = (start: number): Promise<ArrayBuffer> => {
    //     return new Promise((resolve, reject) => {
    //         const chunk = file.slice(start, Math.min(start + finalChunkSize, file.size));
    //         fileReader.onload = (e) => resolve(e.target?.result as ArrayBuffer);
    //         fileReader.onerror = reject;
    //         fileReader.readAsArrayBuffer(chunk);
    //     });
    // };

    // try {
    //     const progressCallback = (current: number) => {
    //         const percentage = Math.round((current / chunks) * 100);
    //         onProgress(percentage);
    //     };

    //     for (let i = 0; i < chunks; i++) {
    //         const chunk = await readChunk(i * chunkSize);
    //         spark.append(chunk);
    //         progressCallback(i + 1);
    //     }

    //     return spark.end();
    // } catch (error) {
    //     throw error;
    // }
}

export default calcFileHash
