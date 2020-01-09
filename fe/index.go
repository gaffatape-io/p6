package fe

import (
	"html/template"
	"k8s.io/klog"
	"net/http"
	"sync"
)

type IndexPage struct {
	CacheTemplate bool
	mtx           sync.Mutex
	template      *template.Template
}

func (i *IndexPage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	i.mtx.Lock()
	if !i.CacheTemplate || i.template == nil {
		var err error
		const path = "fe/index.html"
		i.template, err = template.ParseFiles(path)
		if err != nil {
			klog.Errorf("failed to parse html template:%q err:%+v", path, err)
			const code = http.StatusInternalServerError
			txt := http.StatusText(code)
			http.Error(w, txt, code)
			return
		}
		klog.Info("parsed html template:%q", path)
	}
	i.mtx.Unlock()
}
