package main

import (
	"cms/version_control"
)

func main() {
	m := version_control.Metadata{
		Title:       "test",
		Timestamp:   "test",
		Description: "test",
	}

	data := make(map[string]string)

	data["test"] = "test"
	data["test 2"] = "test 2"

	s := version_control.Snapshot{
		Mdata:   m,
		Muuid:   "test",
		Content: data,
	}

	data_2 := make(map[string]string)

	data_2["test"] = "test"

	d := version_control.Delta{
		Mdata:     m,
		Muuid:     "test-2",
		BaseUuid:  "test",
		Changes:   data_2,
		Deletions: []string{"test 2"},
	}

	version_control.DumpVersion(s)
	version_control.DumpVersion(d)
	version_control.Restore("test-2")
	//Restore("test")

}
