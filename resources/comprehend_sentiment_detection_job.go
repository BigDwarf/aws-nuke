package resources

import (
	"github.com/BigDwarf/aws-nuke/v2/pkg/types"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/comprehend"
)

func init() {
	register("ComprehendSentimentDetectionJob", ListComprehendSentimentDetectionJobs)
}

func ListComprehendSentimentDetectionJobs(sess *session.Session) ([]Resource, error) {
	svc := comprehend.New(sess)

	params := &comprehend.ListSentimentDetectionJobsInput{}
	resources := make([]Resource, 0)

	for {
		resp, err := svc.ListSentimentDetectionJobs(params)
		if err != nil {
			return nil, err
		}
		for _, sentimentDetectionJob := range resp.SentimentDetectionJobPropertiesList {
			if *sentimentDetectionJob.JobStatus == "STOPPED" ||
				*sentimentDetectionJob.JobStatus == "FAILED" {
				// if the job has already been stopped, do not try to delete it again
				continue
			}
			resources = append(resources, &ComprehendSentimentDetectionJob{
				svc:                   svc,
				sentimentDetectionJob: sentimentDetectionJob,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

type ComprehendSentimentDetectionJob struct {
	svc                   *comprehend.Comprehend
	sentimentDetectionJob *comprehend.SentimentDetectionJobProperties
}

func (ce *ComprehendSentimentDetectionJob) Remove() error {
	_, err := ce.svc.StopSentimentDetectionJob(&comprehend.StopSentimentDetectionJobInput{
		JobId: ce.sentimentDetectionJob.JobId,
	})
	return err
}

func (ce *ComprehendSentimentDetectionJob) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("JobName", ce.sentimentDetectionJob.JobName)
	properties.Set("JobId", ce.sentimentDetectionJob.JobId)

	return properties
}

func (ce *ComprehendSentimentDetectionJob) String() string {
	if ce.sentimentDetectionJob.JobName == nil {
		return "Unnamed job"
	} else {
		return *ce.sentimentDetectionJob.JobName
	}
}
