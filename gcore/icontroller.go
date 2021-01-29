package gcore

type IController interface {
	Bean
	RouterMap() RoutesInfo
}

