package core

import "testing"

func TestJobRequest_gen(t *testing.T) {
	job := &JobRequest{
		Pattern: "location/[20200901-20200904]/tokens/<token1,token2,token3>.tar.gz",
	}
	job.gen()
}
