package bulkapi

type BulkAPIService interface{}

type service struct {
	fhirJobRepo *FHIRJobRepo
}

func NewBulkAPIService(fhirJobRepo *FHIRJobRepo) BulkAPIService {
	return &service{fhirJobRepo}
}
