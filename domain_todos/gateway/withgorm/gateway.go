package withgorm

import (
	"context"
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"vikishptra/domain_todos/model/entity"
	"vikishptra/domain_todos/model/errorenum"
	"vikishptra/shared/gogen"
	"vikishptra/shared/infrastructure/config"
	"vikishptra/shared/infrastructure/logger"
)

type Gateway struct {
	log     logger.Logger
	appData gogen.ApplicationData
	config  *config.Config
	Db      *gorm.DB
}

// NewGateway ...
func NewGateway(log logger.Logger, appData gogen.ApplicationData, cfg *config.Config) *Gateway {

	// err := godotenv.Load(".env")
	// if err != nil {
	// 	panic(err)
	// }
	dbUser := os.Getenv("MYSQL_USER")
	dbPassword := os.Getenv("MYSQL_PASSWORD")
	dbHost := os.Getenv("MYSQL_HOST")
	dbPort := os.Getenv("MYSQL_PORT")
	database := os.Getenv("MYSQL_DBNAME")

	// DbHost := os.Getenv("DB_HOST")
	// DbUser := os.Getenv("DB_USER")
	// DbPassword := os.Getenv("DB_PASSWORD")
	// DbName := os.Getenv("DB_NAME")
	// DbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbPort, database)

	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)

	Db, err := gorm.Open("mysql", dsn)

	if err != nil {
		panic(err)
	}

	err = Db.AutoMigrate(entity.Activities{}, entity.Todos{}).Error
	Db.Model(entity.Todos{}).AddForeignKey("activity_group_id", "activities(id)", "CASCADE", "CASCADE")

	if err != nil {
		panic(err)
	}
	return &Gateway{
		log:     log,
		appData: appData,
		config:  cfg,
		Db:      Db,
	}
}

func (r *Gateway) SaveActivite(ctx context.Context, obj *entity.Activities) error {
	r.log.Info(ctx, "called")
	var MaxStatus int64
	if err := r.Db.Model(&entity.Activities{}).Select("MAX(id)").Row().Scan(&MaxStatus); err != nil {
		MaxStatus = 1
		if MaxStatus < 0 {
			return errorenum.SomethingError
		}
		MaxStatus = 0
	}
	obj.ID = MaxStatus + 1
	MaxStatus = obj.ID
	if err := r.Db.Save(obj).Error; err != nil {
		panic(err)
	}

	return nil
}

func (r *Gateway) SaveTodo(ctx context.Context, obj *entity.Todos) error {
	r.log.Info(ctx, "called")
	var MaxStatus int64
	if err := r.Db.Model(&entity.Todos{}).Select("MAX(id)").Row().Scan(&MaxStatus); err != nil {
		MaxStatus = 1
		if MaxStatus < 0 {
			return errorenum.SomethingError
		}
		MaxStatus = 0
	}
	obj.ID = MaxStatus + 1
	MaxStatus = obj.ID
	if err := r.Db.Save(obj).Error; err != nil {
		panic(err)
	}

	return nil
}

func (r *Gateway) UpdateActivitie(ctx context.Context, obj *entity.Activities) error {
	if err := r.Db.Model(entity.Activities{}).Where("id = ?", obj.ID).Update("title", obj.Title); err.Error != nil {
		return errorenum.SomethingError
	}
	return nil
}

func (r *Gateway) UpdateTodo(ctx context.Context, obj *entity.Todos) error {
	if err := r.Db.Save(obj).Error; err != nil {
		panic(err)
	}
	return nil
}

func (r *Gateway) GetAllActivite(ctx context.Context) []*entity.Activities {
	r.log.Info(ctx, "called")
	var Activities []*entity.Activities
	r.Db.Model(entity.Activities{}).Find(&Activities)
	return Activities
}
func (r *Gateway) FindOneActivite(ctx context.Context, ID int64) (*entity.Activities, error) {
	r.log.Info(ctx, "called")
	var Activities entity.Activities
	if err := r.Db.Model(&Activities).Where("id = ?", ID).Select("*").Scan(&Activities); err.RowsAffected == 0 {
		return nil, errorenum.DataNotFound
	}
	return &Activities, nil
}

func (r *Gateway) DeleteActivitie(ctx context.Context, ID int64) error {
	var Activitie entity.Activities
	if err := r.Db.Where("id = ? ", ID).Delete(Activitie); err.RecordNotFound() {
		return errorenum.DataNotFound
	}

	return nil
}

func (r *Gateway) GetAllTodos(ctx context.Context, activity_group_id int64) []*entity.Todos {
	var Todos []*entity.Todos

	if activity_group_id == 0 {
		r.Db.Model(&entity.Todos{}).Find(&Todos)
		return Todos
	}
	if err := r.Db.Model(&entity.Todos{}).Find(&Todos, "activity_group_id = ?", activity_group_id); err.RowsAffected == 0 {
		return Todos
	}
	return Todos
}

func (r *Gateway) FindOneTodos(ctx context.Context, id int64) (*entity.Todos, error) {
	var Todos entity.Todos
	if err := r.Db.Model(&entity.Todos{}).Find(&Todos, "id = ?", id); err.RecordNotFound() {
		return &Todos, errorenum.DataNotFound
	}
	return &Todos, nil
}

func (r *Gateway) DeleteTodo(ctx context.Context, ID int64) error {
	var Todo entity.Todos
	if err := r.Db.Where("id = ? ", ID).Delete(Todo); err.RecordNotFound() {
		return errorenum.DataNotFound
	}

	return nil
}
