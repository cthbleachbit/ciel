package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"ciel-driver"
)

func cleanNormal(c *ciel.Container) error {
	return clean(c, []string{
		`^/dev`,
		`^/efi`,
		`^/etc`,
		`^/run`,
		`^/usr`,
		`^/var/lib/dpkg`,
		`^/var/log/journal$`,
		`^/root`,
		`^/home`,
		`/\.updated$`,
	}, []string{
		`$^`, // => Nothing
	}, func(path string, info os.FileInfo, err error) error {
		if err := os.RemoveAll(path); err != nil {
			log.Println(path, err.Error())
		}
		return nil
	})
}

func cleanFactoryReset(c *ciel.Container) error {
	return clean(c, []string{
		`^/dev`,
		`^/efi`,
		`^/etc`,
		`^/run`,
		`^/usr`,
		`^/var/lib/dpkg`,
		`^/var/log/journal$`,
		`^/root`,
		`^/home`,
		`/\.updated$`,
	}, []string{
		`^/etc/.*-$`,
		`^/etc/machine-id`,
		`^/etc/ssh/ssh_host_.*`,
		`^/var/lib/dpkg/[^/]*-old`,
	}, func(path string, info os.FileInfo, err error) error {
		if err := os.RemoveAll(path); err != nil {
			log.Println(path, err.Error())
		}
		return nil
	})
}

func clean(c *ciel.Container, re []string, reN []string, fn filepath.WalkFunc) error {
	relst := []string{}
	for _, reitem := range re {
		relst = append(relst, "("+reitem+")")
	}

	relstN := []string{}
	for _, reitem := range reN {
		relstN = append(relstN, "("+reitem+")")
	}

	regex := regexp.MustCompile("(" + strings.Join(relst, "|") + ")")
	regexN := regexp.MustCompile("(" + strings.Join(relstN, "|") + ")")

	dpkgfiles := dpkgPackageFiles(c)
	if dpkgfiles == nil {
		return errors.New("no file in dpkg")
	}
	c.Shutdown()
	c.Fs.Mount()
	if !c.Fs.IsMounted() {
		return errors.New("cannot mount file system")
	}
	target := c.Fs.TargetDir()
	filepath.Walk(target, wrapWalkFunc(target, func(path string, info os.FileInfo, err error) error {
		if _, indpkg := dpkgfiles[path]; indpkg {
			return nil
		}
		if !(regex.MatchString(path) && !regexN.MatchString(path)) {
			return fn(filepath.Join(target, path), info, err)
		}
		return nil
	}))

	return nil
}

func wrapWalkFunc(root string, fn filepath.WalkFunc) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return err
		}
		if path == root {
			return nil
		}
		rel, _ := filepath.Rel(root, path)
		rel = "/" + rel
		return fn(rel, info, err)
	}
}
