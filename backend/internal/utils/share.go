package utils

const publicFileIDLength = 20

func PublicFileID(fileID string) string {
	if len(fileID) <= publicFileIDLength {
		return fileID
	}
	return fileID[:publicFileIDLength]
}
