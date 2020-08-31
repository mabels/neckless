package kvpearl

import "testing"

func TestEmptyMatchTag(t *testing.T) {
	ma := MergeArgs{}
	if mkvp, found := ma.Match("xx", Tags{}); !found && mkvp == nil {
		t.Error("should be true")
	}
}

func TestEmptyMatchEmptyTags(t *testing.T) {
	ma := MergeArgs{}
	ma["test1"] = &KVParsed{Tags: Tags{}}
	ma["test2"] = &KVParsed{Tags: Tags{}}
	if mkvp, found := ma.Match("xx", Tags{}); found && mkvp == nil {
		t.Error("should be true")
	}
	if mkvp, found := ma.Match("test1", Tags{}); !found && mkvp == nil {
		t.Error("should be true")
	}
	if mkvp, found := ma.Match("test2", Tags{}); !found && mkvp == nil {
		t.Error("should be true")
	}
}

func TestEmptyMatchTags(t *testing.T) {
	ma := MergeArgs{}
	test1Tag := Tags{}
	test1Tag["t1T1"] = 1
	test1Tag["t1T2"] = 1
	ma["test1"] = &KVParsed{Tags: test1Tag}
	test2Tag := Tags{}
	test2Tag["t2T1"] = 1
	test2Tag["t2T2"] = 1
	ma["test2"] = &KVParsed{Tags: test2Tag}
	if ma.Match("xx", Tags{}) {
		t.Error("should be true")
	}
	if ma.Match("test1", Tags{}) {
		t.Error("should be true")
	}
	if ma.Match("test2", Tags{}) {
		t.Error("should be true")
	}

	unMatch := Tags{}
	unMatch["unMatch"] = 1
	if ma.Match("test1", unMatch) {
		t.Error("should be true")
	}
	if ma.Match("test2", unMatch) {
		t.Error("should be true")
	}

	matcht1T1 := Tags{}
	matcht1T1["t1T1"] = 1
	if !ma.Match("test1", matcht1T1) {
		t.Error("should be true")
	}
	if ma.Match("test2", matcht1T1) {
		t.Error("should be true")
	}

	matcht1T2 := Tags{}
	matcht1T2["t1T2"] = 1
	if !ma.Match("test1", matcht1T2) {
		t.Error("should be true")
	}
	if ma.Match("test2", matcht1T2) {
		t.Error("should be true")
	}

}

func TestNoTags(t *testing.T) {
	ma := MergeArgs{}
	test1 := "test1"
	ma[test1] = &KVParsed{Key: &test1, Tags: Tags{}}
	test2 := "test2"
	ma[test2] = &KVParsed{Key: &test2, Tags: Tags{}}
	if ma.Match("xxx", Tags{}) {
		t.Error("Should fail")
	}
	if !ma.Match("test1", Tags{}) {
		t.Error("Should match")
	}
	if !ma.Match("test2", Tags{}) {
		t.Error("Should match")
	}
}
