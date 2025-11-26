package valueobject

import "github.com/startup-job-board/backend/internal/domain/entity"

type JobType entity.JobType

func (jt JobType) IsValid() bool {
	return entity.JobType(jt) == entity.JobTypeFullTime ||
		entity.JobType(jt) == entity.JobTypePartTime ||
		entity.JobType(jt) == entity.JobTypeContract ||
		entity.JobType(jt) == entity.JobTypeInternship
}



