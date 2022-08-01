package fcore

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

type CdnStructureList struct {
	Path string                       `json:"path"`
	List map[string]*CdnStructureItem `json:"list"`
}

type CdnStructureItem struct {
	Depth int    `json:"depth"` // Dosya derinliği: kaç tane iç içe klasör olacak: eğer 0 ise direkt hash kodu klasör adı oluyor
	Path  string `json:"path"`  // Dosya yolu bu yol daha sonra dirs ile birleşerek tam yola dönüşecek
}

func NewCdnStructureList(path string) *CdnStructureList {
	return &CdnStructureList{
		Path: path,
		List: make(map[string]*CdnStructureItem),
	}
}

func (csl *CdnStructureList) AddItem(cdnType string, depth int, path string) {
	csl.List[cdnType] = &CdnStructureItem{
		Depth: depth,
		Path:  path,
	}
}

func (csl *CdnStructureList) RemoveItem(cdnType string) {
	delete(csl.List, cdnType)
}

func (csl *CdnStructureList) Exists(cdnType string) bool {
	_, exists := csl.List[cdnType]

	return exists
}

func (csl *CdnStructureList) GetMD5(name string) string {
	// Generate MD5 Hash
	hash := md5.Sum([]byte(name))

	return hex.EncodeToString(hash[:])
}

func (csl *CdnStructureList) GetDirs(hash string, depth int) []string {

	var dirs []string

	for i := 0; i < depth; i++ {
		start := i * 2
		end := start + 2
		dirs = append(dirs, hash[start:end])
	}

	return dirs
}

func (csl *CdnStructureList) GetAbsolutePath(path string) string {
	return strings.Replace(path, csl.Path, "", 1)
}

func (csl *CdnStructureList) GetRelativePath(path string) string {
	return csl.Path + path
}

func (csl *CdnStructureList) GetCdnUrlFromName(cdnType, folderName, fileName string) string {
	// Hash
	hash := csl.GetMD5(folderName)

	// Collect Dirs Depth
	var dirs []string

	dirs = append(dirs, csl.Path, csl.List[cdnType].Path)
	dirs = append(dirs, csl.GetDirs(hash, csl.List[cdnType].Depth)...)

	// if only generate filePath, send empty file name
	if fileName != "" {
		dirs = append(dirs, fileName)
	}

	return strings.Join(dirs, "/")
}

func (csl *CdnStructureList) GetCdnUrlFromHash(cdnType, folderHash, fileName string) string {
	// Collect Dirs Depth
	var dirs []string

	dirs = append(dirs, csl.Path, csl.List[cdnType].Path)
	dirs = append(dirs, csl.GetDirs(folderHash, csl.List[cdnType].Depth)...)

	// if only generate filePath, send empty file name
	if fileName != "" {
		dirs = append(dirs, fileName)
	}

	return strings.Join(dirs, "/")
}
