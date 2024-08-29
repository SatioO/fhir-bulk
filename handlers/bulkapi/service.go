package bulkapi

type BulkAPIService interface{}

type service struct {
	fhirJobRepo *repository
}

func NewBulkAPIService(fhirJobRepo *repository) BulkAPIService {
	return &service{fhirJobRepo}
}
