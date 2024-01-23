package dao

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"main/db"
	"main/model"
)

var Uploadhist uploadhist

type uploadhist struct {
}

// 插入上传数据
func (u *uploadhist) UploadData(upload_data *model.Upload_History) error {
	tx := db.GORM.Create(upload_data)
	if tx.Error != nil {
		return errors.New("存入上传数据失败: " + tx.Error.Error())
	}
	return nil
}

// 更新文件上传状态
func (u *uploadhist) UpdateUploadDataStatus(code, status string) error {
	fmt.Println("将" + code + "状态改为：" + status)
	clu := &model.Upload_History{}
	tx := db.GORM.Model(clu).Where("code = ?", code).Update(model.Upload_History{Status: status})
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		return errors.New("更新上传状态失败，" + tx.Error.Error())
	}
	return nil
}

// 查询每个集群的所有文件记录
func (u *uploadhist) GetUploadHistory(clusterName, fiterName string, page, limit int) ([]model.Upload_History, int, error) {
	fmt.Println("filterName:", fiterName)
	//定义分页的起始位置
	startSet := (page - 1) * limit
	var uploadH []model.Upload_History

	var tx *gorm.DB
	var total int
	if fiterName == "" {
		//先查出符合条件的总数量
		tx = db.GORM.Model(uploadH).Where("cluster_name = ?", clusterName).Find(&uploadH)
		if tx.Error != nil && tx.Error.Error() != "record not found" {
			return nil, 0, errors.New("查询每个集群的所有文件记录失败，" + tx.Error.Error())
		}
		fmt.Println("total= ", len(uploadH))
		total = len(uploadH)
		uploadH = nil
	}

	tx = db.GORM.Model(uploadH).Where("cluster_name = ?", clusterName).Where("file like ?", "%"+fiterName+"%").
		Limit(limit).
		Offset(startSet).
		Order("id desc").
		Find(&uploadH)
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		return nil, 0, errors.New("查询每个集群的所有文件记录失败，" + tx.Error.Error())
	}
	if fiterName != "" {
		total = len(uploadH)
	}
	return uploadH, total, nil
}

// 删除指定集群的所有上传记录
func (u *uploadhist) DeleteAllFilesInfo(name string) error {
	uploadu := new(model.Upload_History)
	tx := db.GORM.Model(uploadu).Where("cluster_name = ?", name).Delete(&uploadu)
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		return errors.New("删除指定集群的所有上传记录失败，" + tx.Error.Error())
	}
	return nil
}
