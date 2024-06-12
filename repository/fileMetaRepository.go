package repository

import (
	"DisHub/common"
	"DisHub/common/db"
	"DisHub/common/utils"
	"DisHub/config"
	"errors"
	"gorm.io/gorm"
	"log"
)

type FileMetaRepository struct {
	db *gorm.DB
}

func NewFileMetaRepository(mysqlAddr string) *FileMetaRepository {
	if mysqlAddr == "" {
		mysqlAddr = config.GetLocalAddr()
	}
	d, err := db.NewConnect(mysqlAddr)
	if err != nil {
		log.Fatalln("connect to mysql err: ", err)
	}
	d.AutoMigrate(&FileMeta{})
	return &FileMetaRepository{
		db: d,
	}
}

func (fms *FileMetaRepository) Close() {
	sqlDB, err := fms.db.DB()
	if err != nil {
		log.Printf("close Mysql connect err:%v\n", err)
		return
	}
	sqlDB.Close()
}

// Insert 插入一条
func (fms *FileMetaRepository) Insert(meta *FileMeta) (bool, error) {
	result := fms.db.Create(meta)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return false, nil
		}
		return false, result.Error
	}
	return true, nil
}

// BatchInsert 插入一条
func (fms *FileMetaRepository) BatchInsert(metas *[]*FileMeta) (bool, error) {
	result := fms.db.CreateInBatches(metas, len(*metas))
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return false, nil
		}
		return false, result.Error
	}
	return true, nil
}

// Update 更新文件元数据
func (fms *FileMetaRepository) Update(meta *FileMeta) (bool, error) {
	result := fms.db.Save(meta)
	if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}

// GetFileMetaByHash 获得用户的单个文件的元数据
func (fms *FileMetaRepository) GetFileMetaByHash(uid int, hash string) (*FileMeta, error) {
	var meta FileMeta
	result := fms.db.Where("uid =? and hash = ? and is_del = 0", uid, hash).First(&meta)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &meta, nil
}

func (fms *FileMetaRepository) GetRootFolder(uid int) ([]FileMeta, error) {
	var metas []FileMeta
	result := fms.db.Where("uid = ? AND dir = ? and is_del = 0", uid, common.ROOTFOLDER).Find(&metas)
	if result.Error != nil {
		log.Printf("GetFileMetasByUserAndDir error: %v\n", result.Error)
		return nil, result.Error
	}
	return metas, nil
}

// GetFileMetasByUserAndDir 获得用户某个目录下的所有文件信息
func (fms *FileMetaRepository) GetFileMetasByUserAndDir(uid, dir int) ([]FileMeta, error) {
	var metas []FileMeta
	result := fms.db.Where("uid = ? AND dir = ? and is_del = 0", uid, dir).Order("name ASC").Find(&metas)
	if result.Error != nil {
		log.Printf("GetFileMetasByUserAndDir error: %v\n", result.Error)
		return nil, result.Error
	}
	return metas, nil
}

// GetDeleteFileMetasByUserAndDir 获得用户某个目录下的所有删除的文件信息
func (fms *FileMetaRepository) GetDeleteFileMetasByUserAndDir(uid, dir int) ([]FileMeta, error) {
	var metas []FileMeta
	query := fms.db.Where("uid = ? AND is_del = ?", uid, true)
	if dir >= 0 {
		query = query.Where("dir = ?", dir)
	}
	result := query.Order("fid ASC, name ASC").Find(&metas)
	if result.Error != nil {
		log.Printf("GetFileMetasByUserAndDir error: %v\n", result.Error)
		return nil, result.Error
	}
	return metas, nil
}

// GetFileMetasByUserAndFid 获得指定用户的文件信息
func (fms *FileMetaRepository) GetFileMetasByUserAndFid(uid, fid int) (*FileMeta, error) {
	var meta FileMeta
	result := fms.db.Where("uid = ? AND fid= ? AND is_del = false", uid, fid).First(&meta)
	if result.Error != nil {
		log.Printf("GetFileMetasByUserAndFid error: %v\n", result.Error)
		return nil, result.Error
	}
	return &meta, nil
}

// CheckUserFileOwnership 检查该用户是否拥有次文件夹
func (fms *FileMetaRepository) CheckUserFileOwnership(uid, fid int) (bool, error) {
	var count int64
	result := fms.db.Model(&FileMeta{}).Where("uid = ? AND fid = ? And content_type=?", uid, fid, common.FOLDER).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}
func (fms *FileMetaRepository) GetFileMetasByUserAndDirWithDB(db *gorm.DB, uid, dir int, isDel bool) ([]FileMeta, error) {
	var metas []FileMeta
	result := db.Where("uid = ? AND dir = ? and is_del = ?", uid, dir, isDel).Order("name ASC").Find(&metas)
	if result.Error != nil {
		log.Printf("GetFileMetasByUserAndDir error: %v\n", result.Error)
		return nil, result.Error
	}
	return metas, nil
}
func (fms *FileMetaRepository) GetAllFileMetasByUserAndDirWithDB(db *gorm.DB, uid, dir int) ([]FileMeta, error) {
	var metas []FileMeta
	result := db.Where("uid = ? AND dir = ? ", uid, dir).Order("name ASC").Find(&metas)
	if result.Error != nil {
		log.Printf("GetFileMetasByUserAndDir error: %v\n", result.Error)
		return nil, result.Error
	}
	return metas, nil
}
func (fms *FileMetaRepository) GetFileMetasByUserAndFidsWithDB(db *gorm.DB, uid int, fids []int, isDel bool) ([]FileMeta, error) {
	var metas []FileMeta
	result := db.Where("uid = ? AND fid in ? and is_del = ?", uid, fids, isDel).Order("name ASC").Find(&metas)
	if result.Error != nil {
		log.Printf("GetFileMetasByUserAndDir error: %v\n", result.Error)
		return nil, result.Error
	}
	return metas, nil
}

// DeleteFileMetaByHash 逻辑删除
func (fms *FileMetaRepository) DeleteFileMetaByHash(uid int, hash string) (bool, error) {
	result := fms.db.Model(&FileMeta{}).Where("uid = ? and hash = ? and is_del = 0", uid, hash).Updates(
		map[string]interface{}{
			"is_del":      1,
			"update_time": utils.GetNow(),
		})
	if result.Error != nil {
		return false, result.Error
	}
	if result.RowsAffected == 0 {
		return false, nil // 没有行受到影响，表示可能未找到匹配的记录
	}
	return true, nil
}

// LogicalDeleteFileMetasByFidsWithDB 逻辑删除
func (fms *FileMetaRepository) LogicalDeleteFileMetasByFidsWithDB(db *gorm.DB, uid int, fids []int) error {
	result := db.Model(&FileMeta{}).Where("uid = ? and is_del=false and fid in ?", uid, fids).Updates(
		map[string]interface{}{
			"is_del":      true,
			"update_time": utils.GetNow(),
		})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (fms *FileMetaRepository) CheckFileName(uid, dir int, name string) (bool, error) {
	// 检查是否存在相同目录下相同名称的文件
	var count int64
	result := fms.db.Model(&FileMeta{}).Where("uid = ? AND is_del = false AND dir = ? AND name = ?", uid, dir, name).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

func (fms *FileMetaRepository) CheckFileNameWithDB(db *gorm.DB, uid, dir int, name string) (bool, error) {
	// 检查是否存在相同目录下相同名称的文件
	var count int64
	result := db.Model(&FileMeta{}).Where("uid = ? AND is_del = false AND dir = ? AND name = ?", uid, dir, name).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

// LogicalRecoverFileMetasByFidsWithDB 逻辑恢复文件
func (fms *FileMetaRepository) LogicalRecoverFileMetasByFidsWithDB(db *gorm.DB, uid int, meta FileMeta, existed bool) error {
	var result *gorm.DB
	if !existed {
		// 不存在直接改
		result = db.Model(&FileMeta{}).Where("uid = ? and is_del=true and fid =?", uid, meta.Fid).Updates(
			map[string]interface{}{
				"is_del":      false,
				"update_time": utils.GetNow(),
			})
	} else {
		newName := meta.Name
		for {
			// 在原先名称后加上fid，更新name和is_del
			newName = utils.GenNewFileName(newName)
			ok, err := fms.CheckFileNameWithDB(db, uid, meta.Dir, newName)
			if err != nil {
				return err
			}
			// 不存在则改为新名称
			if !ok {
				result = db.Model(&FileMeta{}).Where("uid = ? and is_del=true and fid =?", uid, meta.Fid).Updates(
					map[string]interface{}{
						"name":        newName,
						"is_del":      false,
						"update_time": utils.GetNow(),
					})
				// 跳出循环
				break
			}
		}

	}
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (fms *FileMetaRepository) RemoveFileMetasByFidsWithDB(db *gorm.DB, uid int, fids []int) (bool, error) {
	// 在数据库中删除指定 fid 的文件元数据
	result := db.Where("uid = ? and is_del = true AND fid in ?", uid, fids).Delete(&FileMeta{})
	if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}
func (fms *FileMetaRepository) Begin() *gorm.DB {
	return fms.db.Begin()
}
