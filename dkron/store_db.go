package dkron

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"sort"
	"time"

	"github.com/distribworks/dkron/v2/ntime"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

// MysqlStore is the local implementation of the Storage interface.
type MysqlStore struct {
	db *gorm.DB
}

// NewDBStore return storage impl by mysql
func NewDBStore(host, dbName, username, password string) (*MysqlStore, error) {
	dsn := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, dbName)
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &MysqlStore{
		db: db,
	}, nil
}

// SetJob stores a job in the storage
func (s *MysqlStore) SetJob(job *Job, copyDependentJobs bool) error {
	if err := job.Validate(); err != nil {
		return err
	}

	var ej Job
	// Abort if parent not found before committing job to the store
	if job.ParentJob != "" {
		if j, _ := s.GetJob(job.ParentJob, nil); j == nil {
			return ErrParentJobNotFound
		}
	}

	//
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var dbJob DBJob

		err := tx.Where(&DBJob{Name: job.Name}).First(&dbJob).Error
		if err != nil {
			if gorm.IsRecordNotFoundError(err) {
				goto save
			}
			return err
		}

		if err := fromDBValue(&dbJob, &ej); err != nil {
			return err
		}
		// When the job runs, these status vars are updated
		// otherwise use the ones that are stored
		if ej.LastError.After(job.LastError) {
			job.LastError = ej.LastError
		}
		if ej.LastSuccess.After(job.LastSuccess) {
			job.LastSuccess = ej.LastSuccess
		}
		if ej.SuccessCount > job.SuccessCount {
			job.SuccessCount = ej.SuccessCount
		}
		if ej.ErrorCount > job.ErrorCount {
			job.ErrorCount = ej.ErrorCount
		}
		if len(ej.DependentJobs) != 0 && copyDependentJobs {
			job.DependentJobs = ej.DependentJobs
		}
		if job.Schedule != ej.Schedule {
			next, err := job.GetNext()
			if err != nil {
				return err
			}
			job.Next = next
		}

	save:
		if err := toDBValue(job, &dbJob); err != nil {
			return err
		}
		tx.Save(&dbJob)

		return nil
	})

	if err != nil {
		return err
	}

	// If the parent job changed update the parents of the old (if any) and new jobs
	if job.ParentJob != ej.ParentJob {
		if err := s.removeFromParent(&ej); err != nil {
			return err
		}
		if err := s.addToParent(job); err != nil {
			return err
		}
	}

	return nil
}

// DeleteJob deletes the given job from the store, along with
// all its executions and references to it.
func (s *MysqlStore) DeleteJob(name string) (*Job, error) {
	var job Job
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var savedDBJob DBJob
		err := tx.Where(&DBJob{Name: name}).First(&savedDBJob).Error
		if err != nil {
			if gorm.IsRecordNotFoundError(err) {
				return nil
			}
			return err
		}

		if err := fromDBValue(&savedDBJob, &job); err != nil {
			return err
		}

		// delete executions
		if err := s.DeleteExecutions(name); err != nil {
			return err
		}

		return tx.Delete(&savedDBJob).Error
	})
	if err != nil {
		return nil, err
	}
	return &job, nil
}

// GetJobs returns all jobs
func (s *MysqlStore) GetJobs(options *JobOptions) ([]*Job, error) {
	jobs := make([]*Job, 0)
	dbJobs := make([]DBJob, 0)
	if err := s.db.Find(&dbJobs).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}

	for _, dbJob := range dbJobs {
		var job Job
		if err := fromDBValue(&dbJob, &job); err == nil {
			jobs = append(jobs, &job)
		} else {
			// do log
		}
	}
	return jobs, nil
}

// GetJob finds and return a Job from the store
func (s *MysqlStore) GetJob(name string, options *JobOptions) (*Job, error) {
	var dbJob DBJob
	if err := s.db.Where(&DBJob{Name: name}).First(&dbJob).Error; err != nil {
		return nil, err
	}
	var job Job
	if err := fromDBValue(&dbJob, &job); err != nil {
		return nil, err
	}
	return &job, nil
}

// Removes the given job from its parent.
// Does nothing if nil is passed as child.
func (s *MysqlStore) removeFromParent(child *Job) error {
	// Do nothing if no job was given or job has no parent
	if child == nil || child.ParentJob == "" {
		return nil
	}

	parent, err := s.GetJob(child.ParentJob, nil)
	if err != nil {
		return err
	}

	// Remove all occurrences from the parent, not just one.
	// Due to an old bug (in v1), a parent can have the same child more than once.
	djs := []string{}
	for _, djn := range parent.DependentJobs {
		if djn != child.Name {
			djs = append(djs, djn)
		}
	}
	parent.DependentJobs = djs
	data, err := json.Marshal(parent.DependentJobs)
	if err != nil {
		return err
	}

	err = s.db.Model(&DBJob{}).Where(DBJob{Name: parent.Name}).UpdateColumn(DBJob{DependentJobs: string(data)}).Error
	return err
}

// Adds the given job to its parent.
func (s *MysqlStore) addToParent(child *Job) error {
	// Do nothing if job has no parent
	if child.ParentJob == "" {
		return nil
	}

	parent, err := s.GetJob(child.ParentJob, nil)
	if err != nil {
		return err
	}
	parent.DependentJobs = append(parent.DependentJobs, child.Name)
	data, err := json.Marshal(parent.DependentJobs)
	if err != nil {
		return err
	}

	err = s.db.Model(&DBJob{}).Where(DBJob{Name: parent.Name}).UpdateColumn(DBJob{DependentJobs: string(data)}).Error
	return err
}

// SetExecution Save a new execution and returns the key of the new saved item or an error.
func (s *MysqlStore) SetExecution(execution *Execution) (string, error) {
	log.WithFields(logrus.Fields{
		"job":       execution.JobName,
		"execution": execution.Key(),
		"finished":  execution.FinishedAt.String(),
	}).Debug("store: Execution")

	var dbExecution DBExecution
	// avoid write twice, just save result on SetExecutionDone
	// if err := toDBValue(execution, &dbExecution); err != nil {
	// 	return "", err
	// }

	// if err := s.db.Create(&dbExecution).Error; err != nil {
	// 	log.WithError(err).WithFields(logrus.Fields{
	// 		"job":       execution.JobName,
	// 		"execution": execution.Key(),
	// 	}).Debug("store: Failed to Create")
	// 	return "", err
	// }

	var count int
	s.db.Model(&dbExecution).Where("job_name = ?", execution.JobName).Count(&count)

	if count > MaxExecutions {
		err := s.db.Where(&DBExecution{JobName: execution.JobName}).Order("TIMESTAMP(started_at) DESC").Offset(MaxExecutions - 2).First(&dbExecution).Error
		if err == nil {
			s.db.Unscoped().Where("TIMESTAMP(started_at) < TIMESTAMP(?)", dbExecution.StartedAt).Delete(DBExecution{})
		}
	}

	return "", nil
}

// DeleteExecutions removes all executions of a job
func (s *MysqlStore) DeleteExecutions(jobName string) error {
	return s.db.Unscoped().Where("job_name = ?", jobName).Delete(DBExecution{}).Error
}

// SetExecutionDone saves the execution and updates the job with the corresponding
// results
func (s *MysqlStore) SetExecutionDone(execution *Execution) (bool, error) {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var dbJob DBJob
		err := tx.Where(&DBJob{Name: execution.JobName}).First(&dbJob).Error
		if gorm.IsRecordNotFoundError(err) {
			log.Warning(ErrExecutionDoneForDeletedJob)
			return ErrExecutionDoneForDeletedJob
		}

		var dbExecution DBExecution
		startAt := execution.StartedAt.Local() // use Local mode
		if err := tx.Where(DBExecution{
			JobName:   execution.JobName,
			StartedAt: &startAt,
			NodeName:  execution.NodeName,
		}).FirstOrCreate(&dbExecution).Error; err != nil {
			return err
		}

		if err := toDBValue(execution, &dbExecution); err != nil {
			return err
		}

		if err := tx.Save(&dbExecution).Error; err != nil {
			return err
		}

		execTime := ntime.NullableTime{}
		execTime.Set(execution.FinishedAt)
		data, err := execTime.MarshalJSON()
		if err != nil {
			return err
		}
		if execution.Success {
			dbJob.LastSuccess = string(data)
			dbJob.SuccessCount++
		} else {
			dbJob.LastError = string(data)
			dbJob.ErrorCount++
		}

		return tx.Save(&dbJob).Error
	})
	if err != nil {
		log.WithError(err).Error("store: Error in SetExecutionDone")
		return false, err
	}
	return true, nil
}

// GetExecutions returns the exections given a Job name.
func (s *MysqlStore) GetExecutions(jobName string) ([]*Execution, error) {
	executions := make([]*Execution, 0)
	dbExecutions := make([]DBExecution, 0)

	err := s.db.Where("job_name = ?", jobName).Find(&dbExecutions).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}

	for _, dbExecution := range dbExecutions {
		var execution Execution
		if err := fromDBValue(&dbExecution, &execution); err != nil {
			// log
		} else {
			executions = append(executions, &execution)
		}
	}

	return executions, nil
}

// GetGroupedExecutions returns executions for a job grouped and with an ordered index
// to facilitate access.
func (s *MysqlStore) GetGroupedExecutions(jobName string) (map[int64][]*Execution, []int64, error) {
	execs, err := s.GetExecutions(jobName)
	if err != nil {
		return nil, nil, err
	}
	groups := make(map[int64][]*Execution)
	for _, exec := range execs {
		groups[exec.Group] = append(groups[exec.Group], exec)
	}

	// Build a separate data structure to show in order
	var byGroup int64arr
	for key := range groups {
		byGroup = append(byGroup, key)
	}
	sort.Sort(sort.Reverse(byGroup))

	return groups, byGroup, nil
}

// GetExecutionGroup returns all executions in the same group of a given execution
func (s *MysqlStore) GetExecutionGroup(execution *Execution) ([]*Execution, error) {
	res, err := s.GetExecutions(execution.JobName)
	if err != nil {
		return nil, err
	}

	var executions []*Execution
	for _, ex := range res {
		if ex.Group == execution.Group {
			executions = append(executions, ex)
		}
	}
	return executions, nil
}

// GetLastExecutionGroup get last execution group given the Job name.
func (s *MysqlStore) GetLastExecutionGroup(jobName string) ([]*Execution, error) {
	executions, byGroup, err := s.GetGroupedExecutions(jobName)
	if err != nil {
		return nil, err
	}

	if len(executions) > 0 && len(byGroup) > 0 {
		return executions[byGroup[0]], nil
	}

	return nil, nil
}

// Shutdown close the DB store
func (s *MysqlStore) Shutdown() error {
	return s.db.Close()
}

// Snapshot no need in mysql store
func (s *MysqlStore) Snapshot(w io.WriteCloser) error {
	return nil
}

// Restore no need in mysql store
func (s *MysqlStore) Restore(r io.ReadCloser) error {
	return nil
}

// DBJob for Job store in mysql
type DBJob struct {
	Name           string     `json:"name"`
	DisplayName    string     `json:"displayname"`
	Timezone       string     `json:"timezone"`
	Schedule       string     `json:"schedule"`
	Owner          string     `json:"owner"`
	OwnerEmail     string     `json:"owner_email"`
	SuccessCount   int        `json:"success_count"`
	ErrorCount     int        `json:"error_count"`
	LastSuccess    string     `json:"last_success"` // ntime.NullableTime to json string
	LastError      string     `json:"last_error"`   // ntime.NullableTime to json string
	Disabled       bool       `json:"disabled"`
	Tags           string     `json:"tags"`     // map to json string
	Metadata       string     `json:"metadata"` // map to json string
	Retries        uint       `json:"retries"`
	DependentJobs  string     `json:"dependent_jobs"` // []string to json string
	ParentJob      string     `json:"parent_job"`
	Processors     string     `json:"processors"` // map[string]PluginConfig to json string
	Concurrency    string     `json:"concurrency"`
	Executor       string     `json:"executor"`
	ExecutorConfig string     `json:"executor_config"` // ExecutorPluginConfig to json string
	Status         string     `json:"status"`
	Next           *time.Time `json:"next"`

	ID        uint `gorm:"primary_key"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

// DBExecution for Execution store in mysql
type DBExecution struct {
	JobName    string     `json:"job_name,omitempty"`
	StartedAt  *time.Time `json:"started_at,omitempty"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`
	Success    bool       `json:"success"`
	Output     string     `json:"output,omitempty"`
	NodeName   string     `json:"node_name,omitempty"`
	Group      int64      `json:"group,omitempty"`
	Attempt    uint       `json:"attempt,omitempty"`

	ID        uint `gorm:"primary_key"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

func fromDBValue(dbValue interface{}, dest interface{}) error {
	destValue := reflect.ValueOf(dest)
	if destValue.Kind() != reflect.Ptr || destValue.IsNil() {
		return errors.New("dest must be ptr not nil")
	}

	srcValue := reflect.ValueOf(dbValue)
	if srcValue.Kind() != reflect.Ptr || srcValue.IsNil() {
		return errors.New("dbValue must be can be ptr not nil")
	}
	srcValue = srcValue.Elem()
	srcType := srcValue.Type()

	destValue = destValue.Elem()
	destType := destValue.Type()

	if srcType.Kind() != reflect.Struct || destType.Kind() != reflect.Struct {
		return errors.New("dbValue and dest must be be struct")
	}

	for i := 0; i < destType.NumField(); i++ {
		destFieldType := destType.Field(i)
		destFieldValue := destValue.Field(i)

		if !destFieldValue.CanSet() {
			continue
		}

		if tmpFieldType, ok := srcType.FieldByName(destFieldType.Name); ok {
			tmpFieldValue := srcValue.FieldByName(destFieldType.Name)

			if dt := destFieldType.Type; dt == tmpFieldType.Type {
				destFieldValue.Set(tmpFieldValue)
			} else if tmpFieldType.Type.Kind() == reflect.String {
				if destFieldType.Type == typeOfBytes {
					destFieldValue.SetBytes([]byte(tmpFieldValue.String()))
					continue
				}
				switch destFieldType.Type.Kind() {
				case reflect.Struct, reflect.Map, reflect.Slice:
					o := destFieldValue.Addr().Interface()
					if err := json.Unmarshal([]byte(tmpFieldValue.String()), o); err != nil {
						// log
					}
				}
			} else if tmpFieldType.Type.Kind() == reflect.Ptr {
				if destFieldType.Type == tmpFieldType.Type.Elem() && !tmpFieldValue.IsNil() {
					destFieldValue.Set(tmpFieldValue.Elem())
				}
			}

		}

	}
	return nil
}

func toDBValue(src interface{}, dbValue interface{}) error {
	destValue := reflect.ValueOf(dbValue)
	if destValue.Kind() != reflect.Ptr || destValue.IsNil() {
		return errors.New("dbValue must be ptr not nil")
	}

	srcValue := reflect.ValueOf(src)
	if srcValue.Kind() != reflect.Ptr || srcValue.IsNil() {
		return errors.New("src must be can be ptr not nil")
	}
	srcValue = srcValue.Elem()
	srcType := srcValue.Type()

	destValue = destValue.Elem()
	destType := destValue.Type()

	if srcType.Kind() != reflect.Struct || destType.Kind() != reflect.Struct {
		return errors.New("src and dbValue must be be struct")
	}

	for i := 0; i < destType.NumField(); i++ {
		destFieldType := destType.Field(i)
		destFieldValue := destValue.Field(i)

		if !destFieldValue.CanSet() {
			continue
		}

		if tmpFieldType, ok := srcType.FieldByName(destFieldType.Name); ok {
			tmpFieldValue := srcValue.FieldByName(destFieldType.Name)

			if dt := destFieldType.Type; dt == tmpFieldType.Type {
				destFieldValue.Set(tmpFieldValue)
			} else if dt.Kind() == reflect.String {
				if tmpFieldType.Type == typeOfBytes {
					destFieldValue.SetString(string(tmpFieldValue.Interface().([]byte)))
					continue
				}
				data, err := json.Marshal(tmpFieldValue.Addr().Interface())
				if err != nil {
					return err
				}
				destFieldValue.SetString(string(data))
			} else if dt.Kind() == reflect.Ptr {
				if tmpFieldType.Type == dt.Elem() && !tmpFieldValue.IsZero() {
					tmp := reflect.New(tmpFieldType.Type)
					tmp.Elem().Set(tmpFieldValue)
					destFieldValue.Set(tmp)
				}
			}
		}
	}

	return nil
}

var typeOfBytes = reflect.TypeOf([]byte(nil))
