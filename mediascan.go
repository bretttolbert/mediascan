package main

import (
    "fmt"
    "log"
    "io/ioutil"
    "os"
    "strconv"
    "strings"
    "sort"
    "path/filepath"
    "github.com/dhowden/tag"
    "gopkg.in/yaml.v3"
)

type MediascanConf struct {
    MediaDir string
    MediaExts []string
    ExcludePath []string
    ExcludeTitle []string
    ExcludeArtist []string
    ExcludeAlbum []string
    ExcludeGenre []string
    SortBy string
    GroupBy string
}

type MediaFile struct {
    Path string
    Size int64
    Format string
    Title string
    Artist string
    Album string
    Genre string
    Year int
}

type MediaFileList struct {
    MediaFiles []MediaFile
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
    var mediaFileList MediaFileList
    countFailed := 0
    countSkipped := 0
    err := filepath.Walk(conf.MediaDir,
        func(path string, info os.FileInfo, err error) error {
            if err != nil {
                return err
            }
            if info.IsDir() {
                return nil
            }
            ext := filepath.Ext(path)
            if !stringInSlice(ext, conf.MediaExts) {
                return nil;
            }

            if containsAnyOf(path, conf.ExcludePath) {
                log.Printf("Skipping %s (ExcludePath)", path)
                countSkipped += 1
                return nil
            }

            f, err := os.Open(path)
            defer closeFile(f)
            check(err)

            tags, err2 := tag.ReadFrom(f)
            if err2 != nil {
                countFailed += 1
                log.Printf("ERROR: %v", err2)
                return nil
            }

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

            var m MediaFile
            m.Path = path
            m.Size = info.Size()
            m.Format = fmt.Sprintf("%s", tags.Format())
            m.Title = tags.Title()
            m.Artist = tags.Artist()
            m.Album = tags.Album()
            m.Genre = tags.Genre()
            m.Year = tags.Year()
            mediaFileList.MediaFiles = append(mediaFileList.MediaFiles, m)
            return nil
        })
    if err != nil {
        log.Printf("ERROR: %v", err)
    }

    if conf.SortBy == "year" {
        sort.SliceStable(mediaFileList.MediaFiles, func(i, j int) bool {
            return mediaFileList.MediaFiles[i].Year < mediaFileList.MediaFiles[j].Year
        })
    } else if conf.SortBy == "artist" {
        sort.SliceStable(mediaFileList.MediaFiles, func(i, j int) bool {
            return mediaFileList.MediaFiles[i].Artist < mediaFileList.MediaFiles[j].Artist
        })
    }


    log.Printf("Successfully loaded %d media files", len(mediaFileList.MediaFiles))
    if countFailed > 0 {
        log.Printf("Failed to load %d media files", countFailed)
    }
    if countSkipped > 0 {
        log.Printf("Skipped %d excluded media files", countSkipped)
    }

    if conf.GroupBy == "year" {
        var mediaFilePlaylistList MediaFilePlaylistList
        mediaFilePlaylistList.Playlists = make(map[string][]MediaFile)
        for _, m := range mediaFileList.MediaFiles {
            year := strconv.FormatInt(int64(m.Year), 10)
            _, ok := mediaFilePlaylistList.Playlists[year]
            if !ok {
                mediaFilePlaylistList.Playlists[year] = make([]MediaFile, 0)
            }
            mediaFilePlaylistList.Playlists[year] = append(mediaFilePlaylistList.Playlists[year], m)
        }
        yamlData, err := yaml.Marshal(&mediaFilePlaylistList)
        check(err);
        err2 := ioutil.WriteFile(outputYamlFilepath, yamlData, 0644)
        check(err2)
    } else {
        yamlData, err := yaml.Marshal(&mediaFileList)
        check(err);
        err2 := ioutil.WriteFile(outputYamlFilepath, yamlData, 0644)
        check(err2)
    }
}
