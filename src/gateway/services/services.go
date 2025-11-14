package services

import (

)

type Gateway interface {

}

type Services struct {
	Gateway
}

func NewServices() *Services {
	return &Services{
		Gateway: NewGatewayService(),
	}
}
