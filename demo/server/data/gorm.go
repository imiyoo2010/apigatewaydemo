package data

import (
	"apigatewaydemo/demo/server/model"
	"fmt"
	"github.com/jinzhu/gorm"
)

type Sqlite struct {

	db 	*gorm.DB
}

func NewSqlite() *Sqlite {

	s := new(Sqlite)

	s.Init()

	return s
}


func (s *Sqlite) Init() {

	db, err := gorm.Open("sqlite3","abc.db")

	if err != nil {
		return
	}


	//如果不设置这个参数，gorm会在表名后加个s
	//db.SingularTable(true)

	db.AutoMigrate(&model.Route{},&model.Upstream{},&model.Version{})

	s.db = db

	s.AddVersion("api_map")
	//s.UpdateVersion("api_map")
}


func (s *Sqlite) GetRouter() []model.Route {

	var rs []model.Route

	s.db.Find(&rs).Limit(100)

	return rs

}

func (s *Sqlite) AddRouter(r *model.Route) int {

	s.db.Create(r)

	return r.ID

}


func (s *Sqlite) GetUpstreamp() []model.Upstream{

	var us []model.Upstream

	s.db.Find(&us).Limit(100)

	return us

}

func (s *Sqlite) AddUpstreamp(u *model.Upstream) int{

	s.db.Create(u)

	return u.ID

}

func (s *Sqlite) AddVersion(filename string) int {

	var (
		vs model.Version

		 filename_count int
	)

	s.db.Model(&vs).Where("file_name=?",filename).Count(&filename_count)

	fmt.Println(filename_count)

	if filename_count < 1 {

		fmt.Println(filename)
		vs.FileName = filename

		s.db.Create(&vs)

		return vs.ID

	}

	return -1
}

func (s *Sqlite) GetVersion(filename string) int{

	var vs model.Version

	s.db.Where("file_name=?",filename).Find(&vs)

	return vs.VersionNum

}

func (s *Sqlite) UpdateVersion(filename string) {

	var vs model.Version

	s.db.Model(&vs).Where("file_name=?",filename).UpdateColumn("version_num", gorm.Expr("version_num + ?", 1))

}



func (s *Sqlite) Close() {
	defer s.db.Close()
}