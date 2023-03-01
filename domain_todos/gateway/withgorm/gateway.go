package withgorm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis"
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
	Redis   *redis.Client
}

// NewGateway ...
func NewGateway(log logger.Logger, appData gogen.ApplicationData, cfg *config.Config, redis *redis.Client) *Gateway {

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
		Redis:   redis,
	}
}

func (r *Gateway) SaveActivite(ctx context.Context, obj *entity.Activities) error {
	r.log.Info(ctx, "called")

	// Jika data tidak ada di cache, maka ambil data dari database
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

	// Simpan data ke dalam cache
	if err := r.Redis.Set(fmt.Sprintf("activities:%d", obj.ID), "1", 0).Err(); err != nil {
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

	if err := r.Redis.Set(fmt.Sprintf("todos:%d", obj.ID), "1", 0).Err(); err != nil {
		panic(err)
	}

	return nil
}

func (r *Gateway) UpdateActivitie(ctx context.Context, obj *entity.Activities) error {
	if err := r.Db.Model(entity.Activities{}).Where("id = ?", obj.ID).Update("title", obj.Title); err.Error != nil {
		return errorenum.SomethingError
	}
	if err := r.Redis.Del(fmt.Sprintf("activities:%d", obj.ID)).Err(); err != nil {
		return errorenum.SomethingError
	}
	return nil
}

func (r *Gateway) UpdateTodo(ctx context.Context, obj *entity.Todos) error {
	if err := r.Db.Save(obj).Error; err != nil {
		panic(err)
	}

	// Hapus data cache yang terkait dengan todo yang diupdate
	if err := r.Redis.Del(fmt.Sprintf("todos:%d", obj.ID)).Err(); err != nil {
		return errorenum.SomethingError
	}

	return nil
}

func (r *Gateway) GetAllActivite(ctx context.Context) []*entity.Activities {
	r.log.Info(ctx, "called")

	// Buat key untuk menyimpan data cache
	cacheKey := "activities:all"

	// Hapus data di cache
	err := r.Redis.Del(cacheKey).Err()
	if err != nil {
		return nil
	}

	// Ambil data dari database
	var Activities []*entity.Activities
	r.Db.Model(entity.Activities{}).Find(&Activities)

	// Simpan data di cache
	cacheData, err := json.Marshal(Activities)
	if err != nil {
		return nil
	} else {
		err = r.Redis.Set(cacheKey, cacheData, time.Hour).Err()
		if err != nil {
			return nil
		}
	}

	return Activities
}

func (r *Gateway) FindOneActivite(ctx context.Context, ID int64) (*entity.Activities, error) {
	r.log.Info(ctx, "called")
	var activity entity.Activities

	// Check if data is in cache
	cacheKey := fmt.Sprintf("activities:%d", ID)

	// If data is not in cache, query from database and save to cache
	if err := r.Db.Model(&entity.Activities{}).Where("id = ?", ID).Select("*").Scan(&activity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorenum.DataNotFound
		}
		return nil, err
	}

	dataBytes, err := json.Marshal(&activity.ID)
	if err != nil {
		panic(err)
	}
	if err := r.Redis.Set(cacheKey, dataBytes, 0).Err(); err != nil {
		panic(err)
	}

	return &activity, nil
}

func (r *Gateway) DeleteActivitie(ctx context.Context, ID int64) error {

	cacheKey := fmt.Sprintf("activities:%d", ID)

	// Hapus data dari cache jika ada
	if err := r.Redis.Del(cacheKey).Err(); err != nil {
		r.log.Error(ctx, err.Error())
	}

	var Activitie entity.Activities
	if err := r.Db.Where("id = ? ", ID).Delete(Activitie); err.RecordNotFound() {
		return errorenum.DataNotFound
	}

	return nil
}

func (r *Gateway) GetAllTodos(ctx context.Context, activity_group_id int64) []*entity.Todos {
	r.log.Info(ctx, "called")

	// Define cache key based on activity_group_id
	cacheKey := fmt.Sprintf("todos:%d", activity_group_id)
	// Hapus data di cache
	err := r.Redis.Del(cacheKey).Err()
	if err != nil {
		return nil
	}

	// Check if cache exists
	cacheData, err := r.Redis.Get(cacheKey).Result()
	if err == nil {
		// If cache exists, unmarshal the data and return it
		var Todos []*entity.Todos
		if err := json.Unmarshal([]byte(cacheData), &Todos); err != nil {
			r.log.Error(ctx, "error unmarshalling cache data", err)
		} else {
			return Todos
		}
	}

	// Hapus data dari cache jika ada
	if err := r.Redis.Del(cacheKey).Err(); err != nil {
		r.log.Error(ctx, err.Error())
	}
	// If cache doesn't exist or error occurred, query from DB
	var Todos []*entity.Todos
	if activity_group_id == 0 {
		r.Db.Model(&entity.Todos{}).Find(&Todos)
	} else {
		if err := r.Db.Model(&entity.Todos{}).Find(&Todos, "activity_group_id = ?", activity_group_id); err.RowsAffected == 0 {
			return Todos
		}
	}

	// Set cache data
	cacheDataBytes, err := json.Marshal(Todos)
	if err != nil {
		r.log.Error(ctx, "error marshalling data for cache", err)
	} else {
		if err := r.Redis.Set(cacheKey, cacheDataBytes, 0).Err(); err != nil {
			r.log.Error(ctx, "error setting cache data", err)
		}
	}

	return Todos
}

func (r *Gateway) FindOneTodos(ctx context.Context, id int64) (*entity.Todos, error) {
	var Todos entity.Todos
	cacheData, err := r.Redis.Get(fmt.Sprintf("todos:%d", id)).Result()
	if err == nil {
		// Jika data ada di cache, maka langsung mengembalikan data dari cache
		if err := json.Unmarshal([]byte(cacheData), &Todos); err != nil {
			r.log.Error(ctx, fmt.Sprintf("error unmarshaling cache data for todos with id %d: %v", id, err))
		} else {
			return &Todos, nil
		}
	}

	if err := r.Db.Model(&entity.Todos{}).Find(&Todos, "id = ?", id); err.RecordNotFound() {
		return &Todos, errorenum.DataNotFound
	}

	// Menyimpan data dari database ke cache selama 1 jam
	dataBytes, err := json.Marshal(Todos)
	if err != nil {
		r.log.Error(ctx, fmt.Sprintf("error marshaling data for todos with id %d: %v", id, err))
	} else {
		if err := r.Redis.Set(fmt.Sprintf("todos:%d", id), dataBytes, 0).Err(); err != nil {
			r.log.Error(ctx, fmt.Sprintf("error caching data for todos with id %d: %v", id, err))
		}
	}

	return &Todos, nil
}

func (r *Gateway) DeleteTodo(ctx context.Context, ID int64) error {
	// Hapus data di cache
	cacheKey := fmt.Sprintf("todos:%d", ID)

	// Hapus data dari cache jika ada
	if err := r.Redis.Del(cacheKey).Err(); err != nil {
		r.log.Error(ctx, err.Error())
	}

	var Todo entity.Todos
	if err := r.Db.Where("id = ?", ID).Delete(&Todo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorenum.DataNotFound
		}
		return err
	}
	return nil
}
