package updater

import (
	"github.com/hashicorp/go-version"
	"github.com/levigross/grequests"
	"time"
)

type Metadata struct {
	Version      string
	Date         time.Time
	ChangeLogs   []string
	Announcement string
	Force        bool
	SetupPackage string
}

type Updater struct {
	url     string
	version string
}

func NewUpdater(url string, version string) *Updater {
	return &Updater{
		url:     url,
		version: version,
	}
}

func (u *Updater) GetMetadata() (*Metadata, error) {
	md := Metadata{}
	resp, err := grequests.Get(u.url, nil)
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&md); err != nil {
		return nil, err
	}
	return &md, nil
}

func (u *Updater) IsNeedUpdate() bool {
	md, err := u.GetMetadata()
	if err != nil {
		return false
	}
	v1, err := version.NewVersion(u.version)
	if err != nil {
		return false
	}
	v2, err := version.NewVersion(md.Version)
	if err != nil {
		return false
	}
	if md.Force || v1.LessThan(v2) {
		return true
	}
	return false
}

func (u *Updater) Update(md Metadata) error {
	if _, err := grequests.Get(md.SetupPackage, nil); err != nil {
		return err
	}
	return nil
}
