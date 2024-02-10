package golanggorm

import (
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func OpenConnection() *gorm.DB {

	dsn := "host=localhost user=postgres password=123456 dbname=databaseTest port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	return db

}

var db = OpenConnection()

func TestOpenConnection(t *testing.T) {
	assert.NotNil(t, db)
}

func TestExecuteSQL(t *testing.T) {
	err := db.Exec("insert into sample(name) values ($1)", "Widi").Error
	assert.Nil(t, err)

	err = db.Exec("insert into sample(name) values ($1)", "Budi").Error
	assert.Nil(t, err)

	err = db.Exec("insert into sample(name) values ($1)", "Joko").Error
	assert.Nil(t, err)

	err = db.Exec("insert into sample(name) values ($1)", "Rully").Error
	assert.Nil(t, err)
}

type Sample struct {
	Id   uuid.UUID
	Name string
}

func TestRawSQL(t *testing.T) {
	var sample Sample
	err := db.Raw("select id, name from sample where id = $1", "c0f00736-c826-11ee-bfa7-db3b6b38fe04").Scan(&sample).Error
	assert.Nil(t, err)
	assert.Equal(t, "Widi", sample.Name)

	var samples []Sample
	err = db.Raw("select id, name from sample").Scan(&samples).Error
	assert.Nil(t, err)
	assert.Equal(t, 4, len(samples))
}

func TestSqlRow(t *testing.T) {
	rows, err := db.Raw("select id, name from sample").Rows()
	assert.Nil(t, err)
	defer rows.Close()

	var samples []Sample
	for rows.Next() {
		var id uuid.UUID
		var name string

		err := rows.Scan(&id, &name)
		assert.Nil(t, err)

		samples = append(samples, Sample{
			Id:   id,
			Name: name,
		})
	}
	assert.Equal(t, 4, len(samples))
}

func TestScanRow(t *testing.T) {
	rows, err := db.Raw("select id, name from sample").Rows()
	assert.Nil(t, err)
	defer rows.Close()

	var samples []Sample
	for rows.Next() {
		err := db.ScanRows(rows, &samples)
		assert.Nil(t, err)
	}
	assert.Equal(t, 4, len(samples))
}

func TestCreateUser(t *testing.T) {
	user := User{
		ID:       uuid.New(),
		Password: "rahasia",
		Name: Name{
			FirstName:  "Tessa",
			MiddleName: "Widi",
			LastName:   "Nugroho",
		},
		Information: "ini akan di ignore",
	}

	response := db.Create(&user)
	assert.Nil(t, response.Error)
	assert.Equal(t, int64(1), response.RowsAffected)
}

func TestBatchInsert(t *testing.T) {
	var users []User
	for i := 2; i < 10; i++ {
		users = append(users, User{
			ID:       uuid.New(),
			Password: "rahasia",
			Name: Name{
				FirstName: "User " + strconv.Itoa(i),
			},
		})
	}

	result := db.Create(&users)
	assert.Nil(t, result.Error)
	assert.Equal(t, 8, int(result.RowsAffected))
}
