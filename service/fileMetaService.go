package service

import (
	"DisHub/common"
	"DisHub/common/utils"
	"DisHub/config"
	"DisHub/repository"
	"sort"
)

type FileMetaService struct {
	r *repository.FileMetaRepository
}

var G_FileMeta FileMetaService

func NewFileMetaService(mysqlAddr string) *FileMetaService {
	if mysqlAddr == "" {
		mysqlAddr = config.GetMySQLAddr()
	}
	r := repository.NewFileMetaRepository(mysqlAddr)
	return &FileMetaService{
		r: r,
	}
}

func (fms *FileMetaService) Load() {
	mysqlAddr := config.GetMySQLAddr()
	r := repository.NewFileMetaRepository(mysqlAddr)
	fms.r = r
}

func (fms *FileMetaService) Close() {
	fms.r.Close()
}

func (fms *FileMetaService) PutFileMeta(meta *repository.FileMeta) (bool, error) {
	return fms.r.Insert(meta)
}

func (fms *FileMetaService) PutFileMetaForce(meta *repository.FileMeta) (bool, error) {
	newName := meta.Name
	for {
		ok, err := fms.r.CheckFileName(meta.Uid, meta.Dir, newName)
		if err != nil {
			return false, err
		}
		// 不存在则插入
		if !ok {
			meta.Name = newName
			return fms.r.Insert(meta)
		}
		newName = utils.GenNewFileName(newName)
	}
}

func (fms *FileMetaService) BatchPutFileMeta(meta *[]*repository.FileMeta) (bool, error) {
	return fms.r.BatchInsert(meta)
}

// GetFileMetaByHash 获得用户的单个文件的元数据
func (fms *FileMetaService) GetFileMetaByHash(uid int, hash string) (*repository.FileMeta, error) {
	return fms.r.GetFileMetaByHash(uid, hash)
}

func (fms *FileMetaService) GetRootFolder(uid int) ([]repository.FileMeta, error) {
	return fms.r.GetRootFolder(uid)
}

// GetFileMetasByUserAndDir 获得用户某个目录下的所有文件信息
func (fms *FileMetaService) GetFileMetasByUserAndDir(uid, dir int) ([]repository.FileMeta, error) {
	metas, err := fms.r.GetFileMetasByUserAndDir(uid, dir)
	if err != nil {
		return nil, err
	}
	// 排序使得目录的优先级更高
	// 目录类型优先
	metas = sortFileFirstFolder(metas)
	return metas, nil
}

// GetDeleteFileMetasByUserAndDir 获得用户某个目录下的所有文件信息
func (fms *FileMetaService) GetDeleteFileMetasByUserAndDir(uid, dir int) ([]repository.FileMeta, error) {
	metas, err := fms.r.GetDeleteFileMetasByUserAndDir(uid, dir)
	if err != nil {
		return nil, err
	}
	fmap := make(map[int]struct{})
	currentFile := make([]repository.FileMeta, 0)
	for _, meta := range metas {
		// 文件夹类型加入hmap中
		if meta.ContentType == common.FOLDER {
			fmap[meta.Fid] = struct{}{}
		}
		// 如果父目录已经被删除则不加入到结果集中
		if _, ok := fmap[meta.Dir]; ok {
			continue
		}
		// 加入结果集
		currentFile = append(currentFile, meta)
	}
	currentFile = sortFileFirstFolder(currentFile)
	return currentFile, nil
}

func (fms *FileMetaService) CheckUserFileOwnership(uid, dir int) (bool, error) {
	return fms.r.CheckUserFileOwnership(uid, dir)
}

func (fms *FileMetaService) CheckUserFileExisted(uid, dir int, name string) (bool, error) {
	return fms.r.CheckFileName(uid, dir, name)
}

func (fms *FileMetaService) UpdateFileMeta(uid int, meta *repository.FileMeta) (bool, error) {
	fm, err := fms.r.GetFileMetasByUserAndFid(uid, meta.Fid)
	if err != nil {
		return false, err
	}
	fm.UpdateFileMeta(meta)
	return fms.r.Update(meta)
}

func (fms *FileMetaService) DeleteFileMetaByHash(uid int, hash string) (bool, error) {
	return fms.r.DeleteFileMetaByHash(uid, hash)
}

func (fms *FileMetaService) LogicalDeleteFileMetasByDir(uid int, delFids []int) error {
	fids := make([]int, 0)
	dirs := make([]int, 0)
	dirs = append(dirs, delFids...)
	// 开启事务
	tx := fms.r.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for len(dirs) > 0 {
		size := len(dirs)
		for i := 0; i < size; i++ {
			fids = append(fids, dirs[i])
			fileMetas, err := fms.r.GetFileMetasByUserAndDirWithDB(tx, uid, dirs[i], false)
			// 出错需要回滚
			if err != nil {
				tx.Rollback()
				return err
			}
			for _, meta := range fileMetas {
				if meta.ContentType == common.FOLDER {
					dirs = append(dirs, meta.Fid)
				}
				fids = append(fids, meta.Fid)
			}
		}
		dirs = dirs[size:]
		err := fms.r.LogicalDeleteFileMetasByFidsWithDB(tx, uid, fids)
		// 出错需要回滚
		if err != nil {
			tx.Rollback()
			return err
		}
		fids = make([]int, 0)
	}
	return tx.Commit().Error
}

func (fms *FileMetaService) LogicalRecoverFileMetasByDir(uid int, delFids []int) error {

	// 开启事务
	tx := fms.r.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	dirs := make([]int, 0)

	var delWithFileMetas func(fileMetas []repository.FileMeta) error
	delWithFileMetas = func(fileMetas []repository.FileMeta) error {
		for _, meta := range fileMetas {
			// 记录恢复的目录恢复深层子文件
			if meta.ContentType == common.FOLDER {
				dirs = append(dirs, meta.Fid)
			}
			// 判断恢复后是否会重复
			ok, err := fms.r.CheckFileNameWithDB(tx, uid, meta.Dir, meta.Name)
			if err != nil {
				tx.Rollback()
				return err
			}
			err = fms.r.LogicalRecoverFileMetasByFidsWithDB(tx, uid, meta, ok)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
		return nil
	}

	fileMetas, err := fms.r.GetFileMetasByUserAndFidsWithDB(tx, uid, delFids, true)
	if err != nil {
		tx.Rollback()
		return err
	}
	delWithFileMetas(fileMetas)
	for len(dirs) > 0 {
		size := len(dirs)
		for i := 0; i < size; i++ {
			// 寻找已经被删除的文件元数据
			fileMetas, err := fms.r.GetFileMetasByUserAndDirWithDB(tx, uid, dirs[i], true)
			// 出错需要回滚
			if err != nil {
				tx.Rollback()
				return err
			}
			delWithFileMetas(fileMetas)
		}
		dirs = dirs[size:]
	}
	return tx.Commit().Error
}

func (fms *FileMetaService) RemoveFileMetasByDir(uid int, delFids []int) error {
	fids := make([]int, 0)
	dirs := make([]int, 0)
	dirs = append(dirs, delFids...)
	// 开启事务
	tx := fms.r.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for len(dirs) > 0 {
		size := len(dirs)
		for i := 0; i < size; i++ {
			fids = append(fids, dirs[i])
			fileMetas, err := fms.r.GetAllFileMetasByUserAndDirWithDB(tx, uid, dirs[i])
			// 出错需要回滚
			if err != nil {
				tx.Rollback()
				return err
			}
			for _, meta := range fileMetas {
				if meta.ContentType == common.FOLDER {
					dirs = append(dirs, meta.Fid)
				}
				fids = append(fids, meta.Fid)
			}
		}
		dirs = dirs[size:]
		_, err := fms.r.RemoveFileMetasByFidsWithDB(tx, uid, fids)
		// 出错需要回滚
		if err != nil {
			tx.Rollback()
			return err
		}
		fids = make([]int, 0)
	}
	return tx.Commit().Error
}

// 目录优先
func sortFileFirstFolder(metas []repository.FileMeta) []repository.FileMeta {
	sort.SliceStable(metas, func(i, j int) bool {
		if (metas)[i].ContentType == common.FOLDER && metas[j].ContentType != common.FOLDER {
			return true // 目录优先级最高
		} else if metas[i].ContentType != common.FOLDER && metas[j].ContentType == common.FOLDER {
			return false
		} else {
			return metas[i].Name < metas[j].Name // 按文件名排序
		}
	})
	return metas
}
