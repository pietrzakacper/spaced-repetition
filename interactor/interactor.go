package interactor

import "controller"

type Interactor interface {
	Start(c controller.Controller)
}
