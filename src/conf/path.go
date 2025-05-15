package conf

import "strings"

func (conf Conf) parsePath(pth string) (p Path) {
	p.FullPath = pth
	p.Path = pth
	p.Exists, p.Errors = conf.exists(p.FullPath)
	p.IsFolder = nil
	p.IsLocal = conf.isLocalPath(p.FullPath)
	p.IsHealthy, p.Errors = conf.isHealthy(p)
	if p.IsLocal {
		p.IsFolder = conf.isFolder(p.FullPath)
		p.IsEmpty, _ = conf.isEmpty(p.FullPath)
	} else {
		p.IsLocal = false
		arr := strings.Split(p.FullPath, ":")
		p.Machine = arr[0]
		p.Path = arr[1]
	}
	return
}

func (pth Path) isEmpty() (b bool) {
	b = false
	switch val := pth.IsEmpty.(type) {
	case bool:
		if val {
			b = true
		}
	}
	return
}

// func (pth Path) isFolder() (b bool) {
// 	b = false
// 	switch val := pth.IsFolder.(type) {
// 	case bool:
// 		if val {
// 			b = true
// 		}
// 	}
// 	return
// }

func (conf Conf) isLocalPath(path string) bool {
	return !strings.Contains(path, ":")
}
