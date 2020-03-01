package dkron

import (
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func Test_Job(t *testing.T) {

	store, err := NewDBStore(os.Getenv("DKRON_DB_HOST"), os.Getenv("DKRON_DB_NAME"), os.Getenv("DKRON_DB_USERNAME"), os.Getenv("DKRON_DB_PASSWORD"))
	if err != nil {
		t.Fatal(err)
	}
	err = store.db.AutoMigrate(&DBJob{}).Error
	if err != nil {
		t.Fatal(err)
	}

	var dbJob DBJob
	job := getJob()
	err = toDBValue(&job, &dbJob)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", dbJob)

	err = store.db.Create(&dbJob).Error
	if err != nil {
		t.Fatal(err)
	}

	var dbJob2 DBJob
	err = store.db.First(&dbJob2).Error
	if err != nil {
		t.Fatal(err)
	}

	var job2 Job
	err = fromDBValue(&dbJob2, &job2)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", job2)
}

func Test_Execution(t *testing.T) {
	store, err := NewDBStore(os.Getenv("DKRON_DB_HOST"), os.Getenv("DKRON_DB_NAME"), os.Getenv("DKRON_DB_USERNAME"), os.Getenv("DKRON_DB_PASSWORD"))
	if err != nil {
		t.Fatal(err)
	}
	err = store.db.AutoMigrate(&DBExecution{}).Error
	if err != nil {
		t.Fatal(err)
	}

	var dbExe DBExecution
	if err := toDBValue(getExecution(), &dbExe); err != nil {
		t.Fatal(err)
	}
	if err := store.db.Create(&dbExe).Error; err != nil {
		t.Fatal(err)
	}

	var dbExe2 DBExecution
	if err := store.db.Last(&dbExe2).Error; err != nil {
		t.Fatal(err)
	}

	var exe Execution
	if err := fromDBValue(dbExe2, &exe); err != nil {
		t.Fatal(err)
	}
	err = store.db.Unscoped().Debug().Where("TIMESTAMP(started_at) < TIMESTAMP(?)", exe.StartedAt).Delete(DBExecution{}).Error
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", exe)
}

func getJob() Job {
	return Job{
		Name:           "name",
		DisplayName:    "display name",
		Timezone:       "UTC",
		Schedule:       "@every 1s",
		Owner:          "user",
		OwnerEmail:     "username@outlook.com",
		SuccessCount:   2,
		ErrorCount:     2,
		Disabled:       true,
		Tags:           map[string]string{"tag": "v"},
		Metadata:       map[string]string{"meta": "v"},
		Retries:        4,
		DependentJobs:  []string{"job1", "job2"},
		ParentJob:      "job1",
		Processors:     map[string]PluginConfig{"p1": PluginConfig{}},
		Concurrency:    "true",
		Executor:       "http",
		Status:         "runing",
		ExecutorConfig: ExecutorPluginConfig{"d": "v"},
		Next:           time.Now(),
	}
}

func getExecution() Execution {
	return Execution{
		JobName:    "nameofjob",
		StartedAt:  time.Now(),
		FinishedAt: time.Now(),
		Success:    false,
		Output:     []byte{'a', 'b'},
		NodeName:   "",
		Group:      53,
		Attempt:    3,
	}
}
