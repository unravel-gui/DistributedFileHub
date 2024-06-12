package response

import "DisHub/repository"

type FolderInfo struct {
	Folder string `json:"folder"`
	Fid    int    `json:"fid"`
}

func NewFolderInfo(folder *repository.FileMeta) *FolderInfo {
	return &FolderInfo{
		Folder: folder.Name,
		Fid:    folder.Fid,
	}
}

type FidsResp struct {
	JwtToken    string      `json:"jwt_token"`
	HomeFolder  *FolderInfo `json:"home_folder"`
	VideoFolder *FolderInfo `json:"video_folder"`
	ImageFolder *FolderInfo `json:"image_folder"`
}
