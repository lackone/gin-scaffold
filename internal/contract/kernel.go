package contract

import "net/http"

const KernelKey = "app:kernel"

type Kernel interface {
	HttpEngine() http.Handler
}
