package handler

var dynamicProxies = map[string]ProxyRequests{}

// AddProxies add proxies to the dynamic collection
func AddProxies(file string, requests ProxyRequests) {
	dynamicProxies[file] = requests
}

// RemoveProxies removes proxies from dynamic collection
func RemoveProxies(file string) {
	delete(dynamicProxies, file)
}

// GetDynamicProxy returns Proxy if there is a matching one
func GetDynamicProxy(path string, method string) *Proxy {
	for _, requests := range dynamicProxies {
		for _, request := range requests {
			if request.Path == path && request.Method == method {
				return &request.Proxy
			}
		}
	}
	return nil
}
