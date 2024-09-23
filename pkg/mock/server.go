package mock

func New() (Router, Mocker) {
	repository := newRepository()
	service := newService(repository)
	mocker := &mocker{
		service: service,
	}
	router := newRouter(service)
	return router, mocker
}

func internalNew(service Service) Mocker {
	return &mocker{
		service: service,
	}
}
