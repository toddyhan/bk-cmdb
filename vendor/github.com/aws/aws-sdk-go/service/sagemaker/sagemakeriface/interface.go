// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

// Package sagemakeriface provides an interface to enable mocking the Amazon SageMaker Service service client
// for testing your code.
//
// It is important to note that this interface will have breaking changes
// when the service model is updated and adds new API operations, paginators,
// and waiters.
package sagemakeriface

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/sagemaker"
)

// SageMakerAPI provides an interface to enable mocking the
// sagemaker.SageMaker service client's API operation,
// paginators, and waiters. This make unit testing your code that calls out
// to the SDK's service client's calls easier.
//
// The best way to use this interface is so the SDK's service client's calls
// can be stubbed out for unit testing your code with the SDK without needing
// to inject custom request handlers into the SDK's request pipeline.
//
//    // myFunc uses an SDK service client to make a request to
//    // Amazon SageMaker Service.
//    func myFunc(svc sagemakeriface.SageMakerAPI) bool {
//        // Make svc.AddTags request
//    }
//
//    func main() {
//        sess := session.New()
//        svc := sagemaker.New(sess)
//
//        myFunc(svc)
//    }
//
// In your _test.go file:
//
//    // Define a mock struct to be used in your unit tests of myFunc.
//    type mockSageMakerClient struct {
//        sagemakeriface.SageMakerAPI
//    }
//    func (m *mockSageMakerClient) AddTags(input *sagemaker.AddTagsInput) (*sagemaker.AddTagsOutput, error) {
//        // mock response/functionality
//    }
//
//    func TestMyFunc(t *testing.T) {
//        // Setup Test
//        mockSvc := &mockSageMakerClient{}
//
//        myfunc(mockSvc)
//
//        // Verify myFunc's functionality
//    }
//
// It is important to note that this interface will have breaking changes
// when the service model is updated and adds new API operations, paginators,
// and waiters. Its suggested to use the pattern above for testing, or using
// tooling to generate mocks to satisfy the interfaces.
type SageMakerAPI interface {
	AddTags(*sagemaker.AddTagsInput) (*sagemaker.AddTagsOutput, error)
	AddTagsWithContext(aws.Context, *sagemaker.AddTagsInput, ...request.Option) (*sagemaker.AddTagsOutput, error)
	AddTagsRequest(*sagemaker.AddTagsInput) (*request.Request, *sagemaker.AddTagsOutput)

	CreateAlgorithm(*sagemaker.CreateAlgorithmInput) (*sagemaker.CreateAlgorithmOutput, error)
	CreateAlgorithmWithContext(aws.Context, *sagemaker.CreateAlgorithmInput, ...request.Option) (*sagemaker.CreateAlgorithmOutput, error)
	CreateAlgorithmRequest(*sagemaker.CreateAlgorithmInput) (*request.Request, *sagemaker.CreateAlgorithmOutput)

	CreateCodeRepository(*sagemaker.CreateCodeRepositoryInput) (*sagemaker.CreateCodeRepositoryOutput, error)
	CreateCodeRepositoryWithContext(aws.Context, *sagemaker.CreateCodeRepositoryInput, ...request.Option) (*sagemaker.CreateCodeRepositoryOutput, error)
	CreateCodeRepositoryRequest(*sagemaker.CreateCodeRepositoryInput) (*request.Request, *sagemaker.CreateCodeRepositoryOutput)

	CreateCompilationJob(*sagemaker.CreateCompilationJobInput) (*sagemaker.CreateCompilationJobOutput, error)
	CreateCompilationJobWithContext(aws.Context, *sagemaker.CreateCompilationJobInput, ...request.Option) (*sagemaker.CreateCompilationJobOutput, error)
	CreateCompilationJobRequest(*sagemaker.CreateCompilationJobInput) (*request.Request, *sagemaker.CreateCompilationJobOutput)

	CreateEndpoint(*sagemaker.CreateEndpointInput) (*sagemaker.CreateEndpointOutput, error)
	CreateEndpointWithContext(aws.Context, *sagemaker.CreateEndpointInput, ...request.Option) (*sagemaker.CreateEndpointOutput, error)
	CreateEndpointRequest(*sagemaker.CreateEndpointInput) (*request.Request, *sagemaker.CreateEndpointOutput)

	CreateEndpointConfig(*sagemaker.CreateEndpointConfigInput) (*sagemaker.CreateEndpointConfigOutput, error)
	CreateEndpointConfigWithContext(aws.Context, *sagemaker.CreateEndpointConfigInput, ...request.Option) (*sagemaker.CreateEndpointConfigOutput, error)
	CreateEndpointConfigRequest(*sagemaker.CreateEndpointConfigInput) (*request.Request, *sagemaker.CreateEndpointConfigOutput)

	CreateHyperParameterTuningJob(*sagemaker.CreateHyperParameterTuningJobInput) (*sagemaker.CreateHyperParameterTuningJobOutput, error)
	CreateHyperParameterTuningJobWithContext(aws.Context, *sagemaker.CreateHyperParameterTuningJobInput, ...request.Option) (*sagemaker.CreateHyperParameterTuningJobOutput, error)
	CreateHyperParameterTuningJobRequest(*sagemaker.CreateHyperParameterTuningJobInput) (*request.Request, *sagemaker.CreateHyperParameterTuningJobOutput)

	CreateLabelingJob(*sagemaker.CreateLabelingJobInput) (*sagemaker.CreateLabelingJobOutput, error)
	CreateLabelingJobWithContext(aws.Context, *sagemaker.CreateLabelingJobInput, ...request.Option) (*sagemaker.CreateLabelingJobOutput, error)
	CreateLabelingJobRequest(*sagemaker.CreateLabelingJobInput) (*request.Request, *sagemaker.CreateLabelingJobOutput)

	CreateModel(*sagemaker.CreateModelInput) (*sagemaker.CreateModelOutput, error)
	CreateModelWithContext(aws.Context, *sagemaker.CreateModelInput, ...request.Option) (*sagemaker.CreateModelOutput, error)
	CreateModelRequest(*sagemaker.CreateModelInput) (*request.Request, *sagemaker.CreateModelOutput)

	CreateModelPackage(*sagemaker.CreateModelPackageInput) (*sagemaker.CreateModelPackageOutput, error)
	CreateModelPackageWithContext(aws.Context, *sagemaker.CreateModelPackageInput, ...request.Option) (*sagemaker.CreateModelPackageOutput, error)
	CreateModelPackageRequest(*sagemaker.CreateModelPackageInput) (*request.Request, *sagemaker.CreateModelPackageOutput)

	CreateNotebookInstance(*sagemaker.CreateNotebookInstanceInput) (*sagemaker.CreateNotebookInstanceOutput, error)
	CreateNotebookInstanceWithContext(aws.Context, *sagemaker.CreateNotebookInstanceInput, ...request.Option) (*sagemaker.CreateNotebookInstanceOutput, error)
	CreateNotebookInstanceRequest(*sagemaker.CreateNotebookInstanceInput) (*request.Request, *sagemaker.CreateNotebookInstanceOutput)

	CreateNotebookInstanceLifecycleConfig(*sagemaker.CreateNotebookInstanceLifecycleConfigInput) (*sagemaker.CreateNotebookInstanceLifecycleConfigOutput, error)
	CreateNotebookInstanceLifecycleConfigWithContext(aws.Context, *sagemaker.CreateNotebookInstanceLifecycleConfigInput, ...request.Option) (*sagemaker.CreateNotebookInstanceLifecycleConfigOutput, error)
	CreateNotebookInstanceLifecycleConfigRequest(*sagemaker.CreateNotebookInstanceLifecycleConfigInput) (*request.Request, *sagemaker.CreateNotebookInstanceLifecycleConfigOutput)

	CreatePresignedNotebookInstanceUrl(*sagemaker.CreatePresignedNotebookInstanceUrlInput) (*sagemaker.CreatePresignedNotebookInstanceUrlOutput, error)
	CreatePresignedNotebookInstanceUrlWithContext(aws.Context, *sagemaker.CreatePresignedNotebookInstanceUrlInput, ...request.Option) (*sagemaker.CreatePresignedNotebookInstanceUrlOutput, error)
	CreatePresignedNotebookInstanceUrlRequest(*sagemaker.CreatePresignedNotebookInstanceUrlInput) (*request.Request, *sagemaker.CreatePresignedNotebookInstanceUrlOutput)

	CreateTrainingJob(*sagemaker.CreateTrainingJobInput) (*sagemaker.CreateTrainingJobOutput, error)
	CreateTrainingJobWithContext(aws.Context, *sagemaker.CreateTrainingJobInput, ...request.Option) (*sagemaker.CreateTrainingJobOutput, error)
	CreateTrainingJobRequest(*sagemaker.CreateTrainingJobInput) (*request.Request, *sagemaker.CreateTrainingJobOutput)

	CreateTransformJob(*sagemaker.CreateTransformJobInput) (*sagemaker.CreateTransformJobOutput, error)
	CreateTransformJobWithContext(aws.Context, *sagemaker.CreateTransformJobInput, ...request.Option) (*sagemaker.CreateTransformJobOutput, error)
	CreateTransformJobRequest(*sagemaker.CreateTransformJobInput) (*request.Request, *sagemaker.CreateTransformJobOutput)

	CreateWorkteam(*sagemaker.CreateWorkteamInput) (*sagemaker.CreateWorkteamOutput, error)
	CreateWorkteamWithContext(aws.Context, *sagemaker.CreateWorkteamInput, ...request.Option) (*sagemaker.CreateWorkteamOutput, error)
	CreateWorkteamRequest(*sagemaker.CreateWorkteamInput) (*request.Request, *sagemaker.CreateWorkteamOutput)

	DeleteAlgorithm(*sagemaker.DeleteAlgorithmInput) (*sagemaker.DeleteAlgorithmOutput, error)
	DeleteAlgorithmWithContext(aws.Context, *sagemaker.DeleteAlgorithmInput, ...request.Option) (*sagemaker.DeleteAlgorithmOutput, error)
	DeleteAlgorithmRequest(*sagemaker.DeleteAlgorithmInput) (*request.Request, *sagemaker.DeleteAlgorithmOutput)

	DeleteCodeRepository(*sagemaker.DeleteCodeRepositoryInput) (*sagemaker.DeleteCodeRepositoryOutput, error)
	DeleteCodeRepositoryWithContext(aws.Context, *sagemaker.DeleteCodeRepositoryInput, ...request.Option) (*sagemaker.DeleteCodeRepositoryOutput, error)
	DeleteCodeRepositoryRequest(*sagemaker.DeleteCodeRepositoryInput) (*request.Request, *sagemaker.DeleteCodeRepositoryOutput)

	DeleteEndpoint(*sagemaker.DeleteEndpointInput) (*sagemaker.DeleteEndpointOutput, error)
	DeleteEndpointWithContext(aws.Context, *sagemaker.DeleteEndpointInput, ...request.Option) (*sagemaker.DeleteEndpointOutput, error)
	DeleteEndpointRequest(*sagemaker.DeleteEndpointInput) (*request.Request, *sagemaker.DeleteEndpointOutput)

	DeleteEndpointConfig(*sagemaker.DeleteEndpointConfigInput) (*sagemaker.DeleteEndpointConfigOutput, error)
	DeleteEndpointConfigWithContext(aws.Context, *sagemaker.DeleteEndpointConfigInput, ...request.Option) (*sagemaker.DeleteEndpointConfigOutput, error)
	DeleteEndpointConfigRequest(*sagemaker.DeleteEndpointConfigInput) (*request.Request, *sagemaker.DeleteEndpointConfigOutput)

	DeleteModel(*sagemaker.DeleteModelInput) (*sagemaker.DeleteModelOutput, error)
	DeleteModelWithContext(aws.Context, *sagemaker.DeleteModelInput, ...request.Option) (*sagemaker.DeleteModelOutput, error)
	DeleteModelRequest(*sagemaker.DeleteModelInput) (*request.Request, *sagemaker.DeleteModelOutput)

	DeleteModelPackage(*sagemaker.DeleteModelPackageInput) (*sagemaker.DeleteModelPackageOutput, error)
	DeleteModelPackageWithContext(aws.Context, *sagemaker.DeleteModelPackageInput, ...request.Option) (*sagemaker.DeleteModelPackageOutput, error)
	DeleteModelPackageRequest(*sagemaker.DeleteModelPackageInput) (*request.Request, *sagemaker.DeleteModelPackageOutput)

	DeleteNotebookInstance(*sagemaker.DeleteNotebookInstanceInput) (*sagemaker.DeleteNotebookInstanceOutput, error)
	DeleteNotebookInstanceWithContext(aws.Context, *sagemaker.DeleteNotebookInstanceInput, ...request.Option) (*sagemaker.DeleteNotebookInstanceOutput, error)
	DeleteNotebookInstanceRequest(*sagemaker.DeleteNotebookInstanceInput) (*request.Request, *sagemaker.DeleteNotebookInstanceOutput)

	DeleteNotebookInstanceLifecycleConfig(*sagemaker.DeleteNotebookInstanceLifecycleConfigInput) (*sagemaker.DeleteNotebookInstanceLifecycleConfigOutput, error)
	DeleteNotebookInstanceLifecycleConfigWithContext(aws.Context, *sagemaker.DeleteNotebookInstanceLifecycleConfigInput, ...request.Option) (*sagemaker.DeleteNotebookInstanceLifecycleConfigOutput, error)
	DeleteNotebookInstanceLifecycleConfigRequest(*sagemaker.DeleteNotebookInstanceLifecycleConfigInput) (*request.Request, *sagemaker.DeleteNotebookInstanceLifecycleConfigOutput)

	DeleteTags(*sagemaker.DeleteTagsInput) (*sagemaker.DeleteTagsOutput, error)
	DeleteTagsWithContext(aws.Context, *sagemaker.DeleteTagsInput, ...request.Option) (*sagemaker.DeleteTagsOutput, error)
	DeleteTagsRequest(*sagemaker.DeleteTagsInput) (*request.Request, *sagemaker.DeleteTagsOutput)

	DeleteWorkteam(*sagemaker.DeleteWorkteamInput) (*sagemaker.DeleteWorkteamOutput, error)
	DeleteWorkteamWithContext(aws.Context, *sagemaker.DeleteWorkteamInput, ...request.Option) (*sagemaker.DeleteWorkteamOutput, error)
	DeleteWorkteamRequest(*sagemaker.DeleteWorkteamInput) (*request.Request, *sagemaker.DeleteWorkteamOutput)

	DescribeAlgorithm(*sagemaker.DescribeAlgorithmInput) (*sagemaker.DescribeAlgorithmOutput, error)
	DescribeAlgorithmWithContext(aws.Context, *sagemaker.DescribeAlgorithmInput, ...request.Option) (*sagemaker.DescribeAlgorithmOutput, error)
	DescribeAlgorithmRequest(*sagemaker.DescribeAlgorithmInput) (*request.Request, *sagemaker.DescribeAlgorithmOutput)

	DescribeCodeRepository(*sagemaker.DescribeCodeRepositoryInput) (*sagemaker.DescribeCodeRepositoryOutput, error)
	DescribeCodeRepositoryWithContext(aws.Context, *sagemaker.DescribeCodeRepositoryInput, ...request.Option) (*sagemaker.DescribeCodeRepositoryOutput, error)
	DescribeCodeRepositoryRequest(*sagemaker.DescribeCodeRepositoryInput) (*request.Request, *sagemaker.DescribeCodeRepositoryOutput)

	DescribeCompilationJob(*sagemaker.DescribeCompilationJobInput) (*sagemaker.DescribeCompilationJobOutput, error)
	DescribeCompilationJobWithContext(aws.Context, *sagemaker.DescribeCompilationJobInput, ...request.Option) (*sagemaker.DescribeCompilationJobOutput, error)
	DescribeCompilationJobRequest(*sagemaker.DescribeCompilationJobInput) (*request.Request, *sagemaker.DescribeCompilationJobOutput)

	DescribeEndpoint(*sagemaker.DescribeEndpointInput) (*sagemaker.DescribeEndpointOutput, error)
	DescribeEndpointWithContext(aws.Context, *sagemaker.DescribeEndpointInput, ...request.Option) (*sagemaker.DescribeEndpointOutput, error)
	DescribeEndpointRequest(*sagemaker.DescribeEndpointInput) (*request.Request, *sagemaker.DescribeEndpointOutput)

	DescribeEndpointConfig(*sagemaker.DescribeEndpointConfigInput) (*sagemaker.DescribeEndpointConfigOutput, error)
	DescribeEndpointConfigWithContext(aws.Context, *sagemaker.DescribeEndpointConfigInput, ...request.Option) (*sagemaker.DescribeEndpointConfigOutput, error)
	DescribeEndpointConfigRequest(*sagemaker.DescribeEndpointConfigInput) (*request.Request, *sagemaker.DescribeEndpointConfigOutput)

	DescribeHyperParameterTuningJob(*sagemaker.DescribeHyperParameterTuningJobInput) (*sagemaker.DescribeHyperParameterTuningJobOutput, error)
	DescribeHyperParameterTuningJobWithContext(aws.Context, *sagemaker.DescribeHyperParameterTuningJobInput, ...request.Option) (*sagemaker.DescribeHyperParameterTuningJobOutput, error)
	DescribeHyperParameterTuningJobRequest(*sagemaker.DescribeHyperParameterTuningJobInput) (*request.Request, *sagemaker.DescribeHyperParameterTuningJobOutput)

	DescribeLabelingJob(*sagemaker.DescribeLabelingJobInput) (*sagemaker.DescribeLabelingJobOutput, error)
	DescribeLabelingJobWithContext(aws.Context, *sagemaker.DescribeLabelingJobInput, ...request.Option) (*sagemaker.DescribeLabelingJobOutput, error)
	DescribeLabelingJobRequest(*sagemaker.DescribeLabelingJobInput) (*request.Request, *sagemaker.DescribeLabelingJobOutput)

	DescribeModel(*sagemaker.DescribeModelInput) (*sagemaker.DescribeModelOutput, error)
	DescribeModelWithContext(aws.Context, *sagemaker.DescribeModelInput, ...request.Option) (*sagemaker.DescribeModelOutput, error)
	DescribeModelRequest(*sagemaker.DescribeModelInput) (*request.Request, *sagemaker.DescribeModelOutput)

	DescribeModelPackage(*sagemaker.DescribeModelPackageInput) (*sagemaker.DescribeModelPackageOutput, error)
	DescribeModelPackageWithContext(aws.Context, *sagemaker.DescribeModelPackageInput, ...request.Option) (*sagemaker.DescribeModelPackageOutput, error)
	DescribeModelPackageRequest(*sagemaker.DescribeModelPackageInput) (*request.Request, *sagemaker.DescribeModelPackageOutput)

	DescribeNotebookInstance(*sagemaker.DescribeNotebookInstanceInput) (*sagemaker.DescribeNotebookInstanceOutput, error)
	DescribeNotebookInstanceWithContext(aws.Context, *sagemaker.DescribeNotebookInstanceInput, ...request.Option) (*sagemaker.DescribeNotebookInstanceOutput, error)
	DescribeNotebookInstanceRequest(*sagemaker.DescribeNotebookInstanceInput) (*request.Request, *sagemaker.DescribeNotebookInstanceOutput)

	DescribeNotebookInstanceLifecycleConfig(*sagemaker.DescribeNotebookInstanceLifecycleConfigInput) (*sagemaker.DescribeNotebookInstanceLifecycleConfigOutput, error)
	DescribeNotebookInstanceLifecycleConfigWithContext(aws.Context, *sagemaker.DescribeNotebookInstanceLifecycleConfigInput, ...request.Option) (*sagemaker.DescribeNotebookInstanceLifecycleConfigOutput, error)
	DescribeNotebookInstanceLifecycleConfigRequest(*sagemaker.DescribeNotebookInstanceLifecycleConfigInput) (*request.Request, *sagemaker.DescribeNotebookInstanceLifecycleConfigOutput)

	DescribeSubscribedWorkteam(*sagemaker.DescribeSubscribedWorkteamInput) (*sagemaker.DescribeSubscribedWorkteamOutput, error)
	DescribeSubscribedWorkteamWithContext(aws.Context, *sagemaker.DescribeSubscribedWorkteamInput, ...request.Option) (*sagemaker.DescribeSubscribedWorkteamOutput, error)
	DescribeSubscribedWorkteamRequest(*sagemaker.DescribeSubscribedWorkteamInput) (*request.Request, *sagemaker.DescribeSubscribedWorkteamOutput)

	DescribeTrainingJob(*sagemaker.DescribeTrainingJobInput) (*sagemaker.DescribeTrainingJobOutput, error)
	DescribeTrainingJobWithContext(aws.Context, *sagemaker.DescribeTrainingJobInput, ...request.Option) (*sagemaker.DescribeTrainingJobOutput, error)
	DescribeTrainingJobRequest(*sagemaker.DescribeTrainingJobInput) (*request.Request, *sagemaker.DescribeTrainingJobOutput)

	DescribeTransformJob(*sagemaker.DescribeTransformJobInput) (*sagemaker.DescribeTransformJobOutput, error)
	DescribeTransformJobWithContext(aws.Context, *sagemaker.DescribeTransformJobInput, ...request.Option) (*sagemaker.DescribeTransformJobOutput, error)
	DescribeTransformJobRequest(*sagemaker.DescribeTransformJobInput) (*request.Request, *sagemaker.DescribeTransformJobOutput)

	DescribeWorkteam(*sagemaker.DescribeWorkteamInput) (*sagemaker.DescribeWorkteamOutput, error)
	DescribeWorkteamWithContext(aws.Context, *sagemaker.DescribeWorkteamInput, ...request.Option) (*sagemaker.DescribeWorkteamOutput, error)
	DescribeWorkteamRequest(*sagemaker.DescribeWorkteamInput) (*request.Request, *sagemaker.DescribeWorkteamOutput)

	GetSearchSuggestions(*sagemaker.GetSearchSuggestionsInput) (*sagemaker.GetSearchSuggestionsOutput, error)
	GetSearchSuggestionsWithContext(aws.Context, *sagemaker.GetSearchSuggestionsInput, ...request.Option) (*sagemaker.GetSearchSuggestionsOutput, error)
	GetSearchSuggestionsRequest(*sagemaker.GetSearchSuggestionsInput) (*request.Request, *sagemaker.GetSearchSuggestionsOutput)

	ListAlgorithms(*sagemaker.ListAlgorithmsInput) (*sagemaker.ListAlgorithmsOutput, error)
	ListAlgorithmsWithContext(aws.Context, *sagemaker.ListAlgorithmsInput, ...request.Option) (*sagemaker.ListAlgorithmsOutput, error)
	ListAlgorithmsRequest(*sagemaker.ListAlgorithmsInput) (*request.Request, *sagemaker.ListAlgorithmsOutput)

	ListCodeRepositories(*sagemaker.ListCodeRepositoriesInput) (*sagemaker.ListCodeRepositoriesOutput, error)
	ListCodeRepositoriesWithContext(aws.Context, *sagemaker.ListCodeRepositoriesInput, ...request.Option) (*sagemaker.ListCodeRepositoriesOutput, error)
	ListCodeRepositoriesRequest(*sagemaker.ListCodeRepositoriesInput) (*request.Request, *sagemaker.ListCodeRepositoriesOutput)

	ListCompilationJobs(*sagemaker.ListCompilationJobsInput) (*sagemaker.ListCompilationJobsOutput, error)
	ListCompilationJobsWithContext(aws.Context, *sagemaker.ListCompilationJobsInput, ...request.Option) (*sagemaker.ListCompilationJobsOutput, error)
	ListCompilationJobsRequest(*sagemaker.ListCompilationJobsInput) (*request.Request, *sagemaker.ListCompilationJobsOutput)

	ListCompilationJobsPages(*sagemaker.ListCompilationJobsInput, func(*sagemaker.ListCompilationJobsOutput, bool) bool) error
	ListCompilationJobsPagesWithContext(aws.Context, *sagemaker.ListCompilationJobsInput, func(*sagemaker.ListCompilationJobsOutput, bool) bool, ...request.Option) error

	ListEndpointConfigs(*sagemaker.ListEndpointConfigsInput) (*sagemaker.ListEndpointConfigsOutput, error)
	ListEndpointConfigsWithContext(aws.Context, *sagemaker.ListEndpointConfigsInput, ...request.Option) (*sagemaker.ListEndpointConfigsOutput, error)
	ListEndpointConfigsRequest(*sagemaker.ListEndpointConfigsInput) (*request.Request, *sagemaker.ListEndpointConfigsOutput)

	ListEndpointConfigsPages(*sagemaker.ListEndpointConfigsInput, func(*sagemaker.ListEndpointConfigsOutput, bool) bool) error
	ListEndpointConfigsPagesWithContext(aws.Context, *sagemaker.ListEndpointConfigsInput, func(*sagemaker.ListEndpointConfigsOutput, bool) bool, ...request.Option) error

	ListEndpoints(*sagemaker.ListEndpointsInput) (*sagemaker.ListEndpointsOutput, error)
	ListEndpointsWithContext(aws.Context, *sagemaker.ListEndpointsInput, ...request.Option) (*sagemaker.ListEndpointsOutput, error)
	ListEndpointsRequest(*sagemaker.ListEndpointsInput) (*request.Request, *sagemaker.ListEndpointsOutput)

	ListEndpointsPages(*sagemaker.ListEndpointsInput, func(*sagemaker.ListEndpointsOutput, bool) bool) error
	ListEndpointsPagesWithContext(aws.Context, *sagemaker.ListEndpointsInput, func(*sagemaker.ListEndpointsOutput, bool) bool, ...request.Option) error

	ListHyperParameterTuningJobs(*sagemaker.ListHyperParameterTuningJobsInput) (*sagemaker.ListHyperParameterTuningJobsOutput, error)
	ListHyperParameterTuningJobsWithContext(aws.Context, *sagemaker.ListHyperParameterTuningJobsInput, ...request.Option) (*sagemaker.ListHyperParameterTuningJobsOutput, error)
	ListHyperParameterTuningJobsRequest(*sagemaker.ListHyperParameterTuningJobsInput) (*request.Request, *sagemaker.ListHyperParameterTuningJobsOutput)

	ListHyperParameterTuningJobsPages(*sagemaker.ListHyperParameterTuningJobsInput, func(*sagemaker.ListHyperParameterTuningJobsOutput, bool) bool) error
	ListHyperParameterTuningJobsPagesWithContext(aws.Context, *sagemaker.ListHyperParameterTuningJobsInput, func(*sagemaker.ListHyperParameterTuningJobsOutput, bool) bool, ...request.Option) error

	ListLabelingJobs(*sagemaker.ListLabelingJobsInput) (*sagemaker.ListLabelingJobsOutput, error)
	ListLabelingJobsWithContext(aws.Context, *sagemaker.ListLabelingJobsInput, ...request.Option) (*sagemaker.ListLabelingJobsOutput, error)
	ListLabelingJobsRequest(*sagemaker.ListLabelingJobsInput) (*request.Request, *sagemaker.ListLabelingJobsOutput)

	ListLabelingJobsPages(*sagemaker.ListLabelingJobsInput, func(*sagemaker.ListLabelingJobsOutput, bool) bool) error
	ListLabelingJobsPagesWithContext(aws.Context, *sagemaker.ListLabelingJobsInput, func(*sagemaker.ListLabelingJobsOutput, bool) bool, ...request.Option) error

	ListLabelingJobsForWorkteam(*sagemaker.ListLabelingJobsForWorkteamInput) (*sagemaker.ListLabelingJobsForWorkteamOutput, error)
	ListLabelingJobsForWorkteamWithContext(aws.Context, *sagemaker.ListLabelingJobsForWorkteamInput, ...request.Option) (*sagemaker.ListLabelingJobsForWorkteamOutput, error)
	ListLabelingJobsForWorkteamRequest(*sagemaker.ListLabelingJobsForWorkteamInput) (*request.Request, *sagemaker.ListLabelingJobsForWorkteamOutput)

	ListLabelingJobsForWorkteamPages(*sagemaker.ListLabelingJobsForWorkteamInput, func(*sagemaker.ListLabelingJobsForWorkteamOutput, bool) bool) error
	ListLabelingJobsForWorkteamPagesWithContext(aws.Context, *sagemaker.ListLabelingJobsForWorkteamInput, func(*sagemaker.ListLabelingJobsForWorkteamOutput, bool) bool, ...request.Option) error

	ListModelPackages(*sagemaker.ListModelPackagesInput) (*sagemaker.ListModelPackagesOutput, error)
	ListModelPackagesWithContext(aws.Context, *sagemaker.ListModelPackagesInput, ...request.Option) (*sagemaker.ListModelPackagesOutput, error)
	ListModelPackagesRequest(*sagemaker.ListModelPackagesInput) (*request.Request, *sagemaker.ListModelPackagesOutput)

	ListModels(*sagemaker.ListModelsInput) (*sagemaker.ListModelsOutput, error)
	ListModelsWithContext(aws.Context, *sagemaker.ListModelsInput, ...request.Option) (*sagemaker.ListModelsOutput, error)
	ListModelsRequest(*sagemaker.ListModelsInput) (*request.Request, *sagemaker.ListModelsOutput)

	ListModelsPages(*sagemaker.ListModelsInput, func(*sagemaker.ListModelsOutput, bool) bool) error
	ListModelsPagesWithContext(aws.Context, *sagemaker.ListModelsInput, func(*sagemaker.ListModelsOutput, bool) bool, ...request.Option) error

	ListNotebookInstanceLifecycleConfigs(*sagemaker.ListNotebookInstanceLifecycleConfigsInput) (*sagemaker.ListNotebookInstanceLifecycleConfigsOutput, error)
	ListNotebookInstanceLifecycleConfigsWithContext(aws.Context, *sagemaker.ListNotebookInstanceLifecycleConfigsInput, ...request.Option) (*sagemaker.ListNotebookInstanceLifecycleConfigsOutput, error)
	ListNotebookInstanceLifecycleConfigsRequest(*sagemaker.ListNotebookInstanceLifecycleConfigsInput) (*request.Request, *sagemaker.ListNotebookInstanceLifecycleConfigsOutput)

	ListNotebookInstanceLifecycleConfigsPages(*sagemaker.ListNotebookInstanceLifecycleConfigsInput, func(*sagemaker.ListNotebookInstanceLifecycleConfigsOutput, bool) bool) error
	ListNotebookInstanceLifecycleConfigsPagesWithContext(aws.Context, *sagemaker.ListNotebookInstanceLifecycleConfigsInput, func(*sagemaker.ListNotebookInstanceLifecycleConfigsOutput, bool) bool, ...request.Option) error

	ListNotebookInstances(*sagemaker.ListNotebookInstancesInput) (*sagemaker.ListNotebookInstancesOutput, error)
	ListNotebookInstancesWithContext(aws.Context, *sagemaker.ListNotebookInstancesInput, ...request.Option) (*sagemaker.ListNotebookInstancesOutput, error)
	ListNotebookInstancesRequest(*sagemaker.ListNotebookInstancesInput) (*request.Request, *sagemaker.ListNotebookInstancesOutput)

	ListNotebookInstancesPages(*sagemaker.ListNotebookInstancesInput, func(*sagemaker.ListNotebookInstancesOutput, bool) bool) error
	ListNotebookInstancesPagesWithContext(aws.Context, *sagemaker.ListNotebookInstancesInput, func(*sagemaker.ListNotebookInstancesOutput, bool) bool, ...request.Option) error

	ListSubscribedWorkteams(*sagemaker.ListSubscribedWorkteamsInput) (*sagemaker.ListSubscribedWorkteamsOutput, error)
	ListSubscribedWorkteamsWithContext(aws.Context, *sagemaker.ListSubscribedWorkteamsInput, ...request.Option) (*sagemaker.ListSubscribedWorkteamsOutput, error)
	ListSubscribedWorkteamsRequest(*sagemaker.ListSubscribedWorkteamsInput) (*request.Request, *sagemaker.ListSubscribedWorkteamsOutput)

	ListSubscribedWorkteamsPages(*sagemaker.ListSubscribedWorkteamsInput, func(*sagemaker.ListSubscribedWorkteamsOutput, bool) bool) error
	ListSubscribedWorkteamsPagesWithContext(aws.Context, *sagemaker.ListSubscribedWorkteamsInput, func(*sagemaker.ListSubscribedWorkteamsOutput, bool) bool, ...request.Option) error

	ListTags(*sagemaker.ListTagsInput) (*sagemaker.ListTagsOutput, error)
	ListTagsWithContext(aws.Context, *sagemaker.ListTagsInput, ...request.Option) (*sagemaker.ListTagsOutput, error)
	ListTagsRequest(*sagemaker.ListTagsInput) (*request.Request, *sagemaker.ListTagsOutput)

	ListTagsPages(*sagemaker.ListTagsInput, func(*sagemaker.ListTagsOutput, bool) bool) error
	ListTagsPagesWithContext(aws.Context, *sagemaker.ListTagsInput, func(*sagemaker.ListTagsOutput, bool) bool, ...request.Option) error

	ListTrainingJobs(*sagemaker.ListTrainingJobsInput) (*sagemaker.ListTrainingJobsOutput, error)
	ListTrainingJobsWithContext(aws.Context, *sagemaker.ListTrainingJobsInput, ...request.Option) (*sagemaker.ListTrainingJobsOutput, error)
	ListTrainingJobsRequest(*sagemaker.ListTrainingJobsInput) (*request.Request, *sagemaker.ListTrainingJobsOutput)

	ListTrainingJobsPages(*sagemaker.ListTrainingJobsInput, func(*sagemaker.ListTrainingJobsOutput, bool) bool) error
	ListTrainingJobsPagesWithContext(aws.Context, *sagemaker.ListTrainingJobsInput, func(*sagemaker.ListTrainingJobsOutput, bool) bool, ...request.Option) error

	ListTrainingJobsForHyperParameterTuningJob(*sagemaker.ListTrainingJobsForHyperParameterTuningJobInput) (*sagemaker.ListTrainingJobsForHyperParameterTuningJobOutput, error)
	ListTrainingJobsForHyperParameterTuningJobWithContext(aws.Context, *sagemaker.ListTrainingJobsForHyperParameterTuningJobInput, ...request.Option) (*sagemaker.ListTrainingJobsForHyperParameterTuningJobOutput, error)
	ListTrainingJobsForHyperParameterTuningJobRequest(*sagemaker.ListTrainingJobsForHyperParameterTuningJobInput) (*request.Request, *sagemaker.ListTrainingJobsForHyperParameterTuningJobOutput)

	ListTrainingJobsForHyperParameterTuningJobPages(*sagemaker.ListTrainingJobsForHyperParameterTuningJobInput, func(*sagemaker.ListTrainingJobsForHyperParameterTuningJobOutput, bool) bool) error
	ListTrainingJobsForHyperParameterTuningJobPagesWithContext(aws.Context, *sagemaker.ListTrainingJobsForHyperParameterTuningJobInput, func(*sagemaker.ListTrainingJobsForHyperParameterTuningJobOutput, bool) bool, ...request.Option) error

	ListTransformJobs(*sagemaker.ListTransformJobsInput) (*sagemaker.ListTransformJobsOutput, error)
	ListTransformJobsWithContext(aws.Context, *sagemaker.ListTransformJobsInput, ...request.Option) (*sagemaker.ListTransformJobsOutput, error)
	ListTransformJobsRequest(*sagemaker.ListTransformJobsInput) (*request.Request, *sagemaker.ListTransformJobsOutput)

	ListTransformJobsPages(*sagemaker.ListTransformJobsInput, func(*sagemaker.ListTransformJobsOutput, bool) bool) error
	ListTransformJobsPagesWithContext(aws.Context, *sagemaker.ListTransformJobsInput, func(*sagemaker.ListTransformJobsOutput, bool) bool, ...request.Option) error

	ListWorkteams(*sagemaker.ListWorkteamsInput) (*sagemaker.ListWorkteamsOutput, error)
	ListWorkteamsWithContext(aws.Context, *sagemaker.ListWorkteamsInput, ...request.Option) (*sagemaker.ListWorkteamsOutput, error)
	ListWorkteamsRequest(*sagemaker.ListWorkteamsInput) (*request.Request, *sagemaker.ListWorkteamsOutput)

	ListWorkteamsPages(*sagemaker.ListWorkteamsInput, func(*sagemaker.ListWorkteamsOutput, bool) bool) error
	ListWorkteamsPagesWithContext(aws.Context, *sagemaker.ListWorkteamsInput, func(*sagemaker.ListWorkteamsOutput, bool) bool, ...request.Option) error

	RenderUiTemplate(*sagemaker.RenderUiTemplateInput) (*sagemaker.RenderUiTemplateOutput, error)
	RenderUiTemplateWithContext(aws.Context, *sagemaker.RenderUiTemplateInput, ...request.Option) (*sagemaker.RenderUiTemplateOutput, error)
	RenderUiTemplateRequest(*sagemaker.RenderUiTemplateInput) (*request.Request, *sagemaker.RenderUiTemplateOutput)

	Search(*sagemaker.SearchInput) (*sagemaker.SearchOutput, error)
	SearchWithContext(aws.Context, *sagemaker.SearchInput, ...request.Option) (*sagemaker.SearchOutput, error)
	SearchRequest(*sagemaker.SearchInput) (*request.Request, *sagemaker.SearchOutput)

	SearchPages(*sagemaker.SearchInput, func(*sagemaker.SearchOutput, bool) bool) error
	SearchPagesWithContext(aws.Context, *sagemaker.SearchInput, func(*sagemaker.SearchOutput, bool) bool, ...request.Option) error

	StartNotebookInstance(*sagemaker.StartNotebookInstanceInput) (*sagemaker.StartNotebookInstanceOutput, error)
	StartNotebookInstanceWithContext(aws.Context, *sagemaker.StartNotebookInstanceInput, ...request.Option) (*sagemaker.StartNotebookInstanceOutput, error)
	StartNotebookInstanceRequest(*sagemaker.StartNotebookInstanceInput) (*request.Request, *sagemaker.StartNotebookInstanceOutput)

	StopCompilationJob(*sagemaker.StopCompilationJobInput) (*sagemaker.StopCompilationJobOutput, error)
	StopCompilationJobWithContext(aws.Context, *sagemaker.StopCompilationJobInput, ...request.Option) (*sagemaker.StopCompilationJobOutput, error)
	StopCompilationJobRequest(*sagemaker.StopCompilationJobInput) (*request.Request, *sagemaker.StopCompilationJobOutput)

	StopHyperParameterTuningJob(*sagemaker.StopHyperParameterTuningJobInput) (*sagemaker.StopHyperParameterTuningJobOutput, error)
	StopHyperParameterTuningJobWithContext(aws.Context, *sagemaker.StopHyperParameterTuningJobInput, ...request.Option) (*sagemaker.StopHyperParameterTuningJobOutput, error)
	StopHyperParameterTuningJobRequest(*sagemaker.StopHyperParameterTuningJobInput) (*request.Request, *sagemaker.StopHyperParameterTuningJobOutput)

	StopLabelingJob(*sagemaker.StopLabelingJobInput) (*sagemaker.StopLabelingJobOutput, error)
	StopLabelingJobWithContext(aws.Context, *sagemaker.StopLabelingJobInput, ...request.Option) (*sagemaker.StopLabelingJobOutput, error)
	StopLabelingJobRequest(*sagemaker.StopLabelingJobInput) (*request.Request, *sagemaker.StopLabelingJobOutput)

	StopNotebookInstance(*sagemaker.StopNotebookInstanceInput) (*sagemaker.StopNotebookInstanceOutput, error)
	StopNotebookInstanceWithContext(aws.Context, *sagemaker.StopNotebookInstanceInput, ...request.Option) (*sagemaker.StopNotebookInstanceOutput, error)
	StopNotebookInstanceRequest(*sagemaker.StopNotebookInstanceInput) (*request.Request, *sagemaker.StopNotebookInstanceOutput)

	StopTrainingJob(*sagemaker.StopTrainingJobInput) (*sagemaker.StopTrainingJobOutput, error)
	StopTrainingJobWithContext(aws.Context, *sagemaker.StopTrainingJobInput, ...request.Option) (*sagemaker.StopTrainingJobOutput, error)
	StopTrainingJobRequest(*sagemaker.StopTrainingJobInput) (*request.Request, *sagemaker.StopTrainingJobOutput)

	StopTransformJob(*sagemaker.StopTransformJobInput) (*sagemaker.StopTransformJobOutput, error)
	StopTransformJobWithContext(aws.Context, *sagemaker.StopTransformJobInput, ...request.Option) (*sagemaker.StopTransformJobOutput, error)
	StopTransformJobRequest(*sagemaker.StopTransformJobInput) (*request.Request, *sagemaker.StopTransformJobOutput)

	UpdateCodeRepository(*sagemaker.UpdateCodeRepositoryInput) (*sagemaker.UpdateCodeRepositoryOutput, error)
	UpdateCodeRepositoryWithContext(aws.Context, *sagemaker.UpdateCodeRepositoryInput, ...request.Option) (*sagemaker.UpdateCodeRepositoryOutput, error)
	UpdateCodeRepositoryRequest(*sagemaker.UpdateCodeRepositoryInput) (*request.Request, *sagemaker.UpdateCodeRepositoryOutput)

	UpdateEndpoint(*sagemaker.UpdateEndpointInput) (*sagemaker.UpdateEndpointOutput, error)
	UpdateEndpointWithContext(aws.Context, *sagemaker.UpdateEndpointInput, ...request.Option) (*sagemaker.UpdateEndpointOutput, error)
	UpdateEndpointRequest(*sagemaker.UpdateEndpointInput) (*request.Request, *sagemaker.UpdateEndpointOutput)

	UpdateEndpointWeightsAndCapacities(*sagemaker.UpdateEndpointWeightsAndCapacitiesInput) (*sagemaker.UpdateEndpointWeightsAndCapacitiesOutput, error)
	UpdateEndpointWeightsAndCapacitiesWithContext(aws.Context, *sagemaker.UpdateEndpointWeightsAndCapacitiesInput, ...request.Option) (*sagemaker.UpdateEndpointWeightsAndCapacitiesOutput, error)
	UpdateEndpointWeightsAndCapacitiesRequest(*sagemaker.UpdateEndpointWeightsAndCapacitiesInput) (*request.Request, *sagemaker.UpdateEndpointWeightsAndCapacitiesOutput)

	UpdateNotebookInstance(*sagemaker.UpdateNotebookInstanceInput) (*sagemaker.UpdateNotebookInstanceOutput, error)
	UpdateNotebookInstanceWithContext(aws.Context, *sagemaker.UpdateNotebookInstanceInput, ...request.Option) (*sagemaker.UpdateNotebookInstanceOutput, error)
	UpdateNotebookInstanceRequest(*sagemaker.UpdateNotebookInstanceInput) (*request.Request, *sagemaker.UpdateNotebookInstanceOutput)

	UpdateNotebookInstanceLifecycleConfig(*sagemaker.UpdateNotebookInstanceLifecycleConfigInput) (*sagemaker.UpdateNotebookInstanceLifecycleConfigOutput, error)
	UpdateNotebookInstanceLifecycleConfigWithContext(aws.Context, *sagemaker.UpdateNotebookInstanceLifecycleConfigInput, ...request.Option) (*sagemaker.UpdateNotebookInstanceLifecycleConfigOutput, error)
	UpdateNotebookInstanceLifecycleConfigRequest(*sagemaker.UpdateNotebookInstanceLifecycleConfigInput) (*request.Request, *sagemaker.UpdateNotebookInstanceLifecycleConfigOutput)

	UpdateWorkteam(*sagemaker.UpdateWorkteamInput) (*sagemaker.UpdateWorkteamOutput, error)
	UpdateWorkteamWithContext(aws.Context, *sagemaker.UpdateWorkteamInput, ...request.Option) (*sagemaker.UpdateWorkteamOutput, error)
	UpdateWorkteamRequest(*sagemaker.UpdateWorkteamInput) (*request.Request, *sagemaker.UpdateWorkteamOutput)

	WaitUntilEndpointDeleted(*sagemaker.DescribeEndpointInput) error
	WaitUntilEndpointDeletedWithContext(aws.Context, *sagemaker.DescribeEndpointInput, ...request.WaiterOption) error

	WaitUntilEndpointInService(*sagemaker.DescribeEndpointInput) error
	WaitUntilEndpointInServiceWithContext(aws.Context, *sagemaker.DescribeEndpointInput, ...request.WaiterOption) error

	WaitUntilNotebookInstanceDeleted(*sagemaker.DescribeNotebookInstanceInput) error
	WaitUntilNotebookInstanceDeletedWithContext(aws.Context, *sagemaker.DescribeNotebookInstanceInput, ...request.WaiterOption) error

	WaitUntilNotebookInstanceInService(*sagemaker.DescribeNotebookInstanceInput) error
	WaitUntilNotebookInstanceInServiceWithContext(aws.Context, *sagemaker.DescribeNotebookInstanceInput, ...request.WaiterOption) error

	WaitUntilNotebookInstanceStopped(*sagemaker.DescribeNotebookInstanceInput) error
	WaitUntilNotebookInstanceStoppedWithContext(aws.Context, *sagemaker.DescribeNotebookInstanceInput, ...request.WaiterOption) error

	WaitUntilTrainingJobCompletedOrStopped(*sagemaker.DescribeTrainingJobInput) error
	WaitUntilTrainingJobCompletedOrStoppedWithContext(aws.Context, *sagemaker.DescribeTrainingJobInput, ...request.WaiterOption) error

	WaitUntilTransformJobCompletedOrStopped(*sagemaker.DescribeTransformJobInput) error
	WaitUntilTransformJobCompletedOrStoppedWithContext(aws.Context, *sagemaker.DescribeTransformJobInput, ...request.WaiterOption) error
}

var _ SageMakerAPI = (*sagemaker.SageMaker)(nil)
