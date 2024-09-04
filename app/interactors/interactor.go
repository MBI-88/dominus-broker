package interactors

import "dominus/app/domain/entities"


type interactor struct {
	m *entities.MessageRest // Rest struct
	// grpc struct
}



// RestToGrp sends information from  rest protocol receiver to  grpc protocol client
func (i interactor) RestToGrp() {

}

// RestToRest sends information from rest protocol receiver to rest protocol client
func (i interactor) RestToRest() {

}

// GrpcToRest sends information from grpc protocol receiver to rest protocol client
func (i interactor) GrpcToRest() {

}

// UpdateNode updates topics and subcribers in the node
func (i interactor) UpdateNode() {

}




type InteractorInt interface {
	RestToGrp() 
	RestToRest()
	GrpcToRest()
	UpdateNode()
}


func NewInteractor() InteractorInt {
	return &interactor{}
}