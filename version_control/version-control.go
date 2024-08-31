package version_control

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Metadata struct {
	Title       string
	Timestamp   string
	Description string
}

type Version interface {
	IsSnapshot() bool
	UUID() string
}

type Snapshot struct {
	Mdata   Metadata
	Muuid   string
	Content map[string]string
}

type Delta struct {
	Mdata     Metadata
	Muuid     string
	BaseUuid  string
	Changes   map[string]string
	Deletions []string
}

func (s Snapshot) IsSnapshot() bool {
	return true
}

func (d Delta) IsSnapshot() bool {
	return false
}

func (s Snapshot) UUID() string {
	return s.Muuid
}

func (d Delta) UUID() string {
	return d.Muuid
}

func setSnapshot(snapshot *Snapshot) error {

	err := filepath.WalkDir("content/", func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			if err := os.Remove(path); err != nil {
				return fmt.Errorf("failed to remove file %s: %w", path, err)
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	for fname := range snapshot.Content {
		f, err := os.OpenFile("content/"+fname, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)

		if err != nil {
			return err
		}

		defer f.Close()

		_, err = f.WriteString(snapshot.Content[fname])

		if err != nil {
			return err
		}
	}

	return nil
}

func applyDelta(v_delta Version, base_snapshot Snapshot) Snapshot {
	if v_delta.IsSnapshot() {
		return base_snapshot
	}

	delta := v_delta.(Delta)

	for changed_file := range delta.Changes {
		base_snapshot.Content[changed_file] = delta.Changes[changed_file]
	}

	for _, deletion := range delta.Deletions {
		delete(base_snapshot.Content, deletion)
	}

	return base_snapshot
}

func compress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	_, err := w.Write(data)

	if err != nil {
		return nil, err
	}

	err = w.Close()

	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func getCompressedForVersion(v Version) ([]byte, error) {
	var data []byte
	var err error

	data, err = json.Marshal(v)

	if err != nil {
		return nil, err
	}

	compressedData, err := compress(data)

	if err != nil {
		return nil, err
	}

	return compressedData, nil
}

func DumpVersion(v Version) error {
	fname := "versions/" + v.UUID()

	file, err := os.Create(fname)

	if err != nil {
		return err
	}

	defer file.Close()

	compressedData, err := getCompressedForVersion(v)

	if err != nil {
		panic(err)
	}

	_, err = file.Write(compressedData)

	if err != nil {
		panic(err)
	}

	return nil

}

func decompress(compressed []byte) ([]byte, error) {
	b := bytes.NewReader(compressed)
	r, err := zlib.NewReader(b)

	if err != nil {
		return nil, err
	}

	defer r.Close()

	var out bytes.Buffer
	_, err = io.Copy(&out, r)
	if err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

func getObjectFrom(uuid string) ([]byte, error) {
	data, err := os.ReadFile("versions/" + uuid)
	if err != nil {
		return nil, err
	}

	decompressedData, err := decompress(data)

	if err != nil {
		return nil, err
	}

	return decompressedData, nil
}

func lookup(uuid string) (Version, error) {
	var v_json map[string]interface{}

	jsonData, err := getObjectFrom(uuid)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonData, &v_json)

	if err != nil {
		return nil, err
	}

	if _, ok := v_json["BaseUuid"]; ok {
		var d Delta
		err = json.Unmarshal(jsonData, &d)

		if err != nil {
			return nil, err
		}

		return d, nil
	} else {
		var s Snapshot
		err = json.Unmarshal(jsonData, &s)

		if err != nil {
			return nil, err
		}

		return s, nil
	}
}

func reconstructSnapshot(v Version) Snapshot {
	if v.IsSnapshot() {
		return v.(Snapshot)
	} else if delta_v, ok := v.(Delta); ok {
		parent_v, err := lookup(delta_v.BaseUuid)

		if err != nil {
			panic("Error looking up uuid")
		}

		parent_snap := reconstructSnapshot(parent_v)
		cur_snap := applyDelta(delta_v, parent_snap)
		return cur_snap
	}

	return Snapshot{}
}

func Restore(uuid string) {
	v, err := lookup(uuid)

	if err != nil {
		panic("Error looking up uuid")
	}

	snapshot := reconstructSnapshot(v)
	setSnapshot(&snapshot)
}
