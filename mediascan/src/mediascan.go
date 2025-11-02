package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/dhowden/tag"
	"github.com/tcolgate/mp3"
	"gopkg.in/yaml.v3"
)

type MediascanConf struct {
	MediaDirs          []string
	MediaExts          []string
	ExcludePaths       []string
	ExcludeTitle       []string
	ExcludeArtist      []string
	ExcludeAlbumArtist []string
	ExcludeAlbum       []string
	ExcludeGenre       []string
	SortBy             string
	GroupBy            string
	GetMp3Duration     bool
}

type MediaFile struct {
	Path        string
	Size        int64
	Format      string
	Title       string
	Artist      string
	AlbumArtist string
	Album       string
	Genre       string
	Year        int
	Duration    float64
	ModTime     time.Time
}

type MediaFiles struct {
	Files []MediaFile
}

type MediaFilePlaylistList struct {
	Playlists map[string][]MediaFile
}

func check(e error) {
	if e != nil {
		log.Fatalf("ERROR: %v", e)
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

/**
 * Returns true if the given string (a) contains any of the substrings
 * in the given slice of strings (list) (case-insensitive)
 */
func containsAnyOf(a string, list []string) bool {
	if len(a) > 0 {
		a = strings.ToLower(a)
		for _, b := range list {
			if strings.Contains(a, strings.ToLower(b)) {
				return true
			}
		}
	}
	return false
}

func closeFile(f *os.File) {
	err := f.Close()
	if err != nil {
		log.Printf("ERROR: Failed to close file %v\n", err)
		os.Exit(1)
	}
}

func loadConf(configYamlFilepath string) (conf MediascanConf) {
	yfile, err := ioutil.ReadFile(configYamlFilepath)
	check(err)
	err2 := yaml.Unmarshal(yfile, &conf)
	check(err2)
	return
}

func getMp3Duration(path string) (duration float64) {
	t := 0.0
	r, err := os.Open(path)
	if err != nil {
		log.Printf("getMp3Duration error: %v", err)
		return 0.0
	}
	d := mp3.NewDecoder(r)
	var f mp3.Frame
	skipped := 0
	for {
		if err := d.Decode(&f, &skipped); err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			return 0.0
		}
		t = t + f.Duration().Seconds()
	}
	return math.Round(t*100) / 100
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Error: Invalid arguments")
		fmt.Println("Usage: go run mediascan.go {config yaml filepath} {output yaml filepath}")
		fmt.Println("Example: go run mediascan.go conf.yaml files.yaml")
		os.Exit(1)
	}
	configYamlFilepath := os.Args[1]
	outputYamlFilepath := os.Args[2]
	conf := loadConf(configYamlFilepath)
	var files MediaFiles
	countLoadFailed := 0
	countTagsFailed := 0
	countSkipped := 0
	for _, mediaDir := range conf.MediaDirs {
		err := filepath.Walk(mediaDir,
			func(path string, info os.FileInfo, err error) error {
				var m MediaFile
				m.Path = path
				m.Size = info.Size()
				m.ModTime = info.ModTime()
				m.Duration = 0.0
				m.Format = ""
				ext := filepath.Ext(info.Name())
				name := strings.TrimSuffix(info.Name(), ext)
				m.Title = name

				log.Printf("Reading filepath %s", path)
				if err != nil {
					return err
				}
				if info.IsDir() {
					return nil
				}
				if !stringInSlice(ext, conf.MediaExts) {
					return nil
				}

				if containsAnyOf(path, conf.ExcludePaths) {
					log.Printf("Skipping %s (ExcludePaths)", path)
					countSkipped += 1
					return nil
				}

				f, err := os.Open(path)
				defer closeFile(f)
				check(err)

				tags, err2 := tag.ReadFrom(f)
				if err2 != nil {
					countTagsFailed += 1
					//log.Printf("ERROR reading tags: %v", err2)
				} else {
					if containsAnyOf(tags.Title(), conf.ExcludeTitle) {
						log.Printf("Skipping %s (ExcludeTitle %s)", path, tags.Title())
						countSkipped += 1
						return nil
					}
					if containsAnyOf(tags.Artist(), conf.ExcludeArtist) {
						log.Printf("Skipping %s (ExcludeArtist %s)", path, tags.Artist())
						countSkipped += 1
						return nil
					}
					if containsAnyOf(tags.AlbumArtist(), conf.ExcludeArtist) {
						log.Printf("Skipping %s (ExcludeAlbumArtist %s)", path, tags.AlbumArtist())
						countSkipped += 1
						return nil
					}
					if containsAnyOf(tags.Album(), conf.ExcludeAlbum) {
						log.Printf("Skipping %s (ExcludeAlbum %s)", path, tags.Album())
						countSkipped += 1
						return nil
					}
					if containsAnyOf(tags.Genre(), conf.ExcludeGenre) {
						log.Printf("Skipping %s (ExcludeGenre %s)", path, tags.Genre())
						countSkipped += 1
						return nil
					}

					m.Format = string(tags.Format())
					m.Title = tags.Title()
					m.Artist = tags.Artist()
					m.AlbumArtist = tags.AlbumArtist()
					m.Album = tags.Album()
					m.Genre = tags.Genre()
					m.Year = tags.Year()
				}

				if conf.GetMp3Duration && ext == ".mp3" {
					m.Duration = getMp3Duration(path)
				}

				files.Files = append(files.Files, m)

				return nil
			})
		if err != nil {
			log.Printf("ERROR loading file: %v", err)
			countLoadFailed += 1
		}
	}

	if conf.SortBy == "year" {
		sort.SliceStable(files.Files, func(i, j int) bool {
			return files.Files[i].Year < files.Files[j].Year
		})
	} else if conf.SortBy == "artist" {
		sort.SliceStable(files.Files, func(i, j int) bool {
			return files.Files[i].Artist < files.Files[j].Artist
		})
	}

	log.Printf("Successfully loaded %d media files", len(files.Files))
	if countLoadFailed > 0 {
		log.Printf("Failed to load for %d media files", countLoadFailed)
	}
	if countTagsFailed > 0 {
		log.Printf("Failed to load tags for %d media files", countTagsFailed)
	}
	if countSkipped > 0 {
		log.Printf("Skipped %d excluded media files", countSkipped)
	}

	if conf.GroupBy == "year" {
		var mediaFilePlaylistList MediaFilePlaylistList
		mediaFilePlaylistList.Playlists = make(map[string][]MediaFile)
		for _, m := range files.Files {
			year := strconv.FormatInt(int64(m.Year), 10)
			_, ok := mediaFilePlaylistList.Playlists[year]
			if !ok {
				mediaFilePlaylistList.Playlists[year] = make([]MediaFile, 0)
			}
			mediaFilePlaylistList.Playlists[year] = append(mediaFilePlaylistList.Playlists[year], m)
		}
		yamlData, err := yaml.Marshal(&mediaFilePlaylistList)
		check(err)
		err2 := os.WriteFile(outputYamlFilepath, yamlData, 0644)
		check(err2)
	} else {
		yamlData, err := yaml.Marshal(&files)
		check(err)
		err2 := os.WriteFile(outputYamlFilepath, yamlData, 0644)
		check(err2)
	}
}
