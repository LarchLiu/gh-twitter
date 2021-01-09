package model

import "testing"

func TestSmmsUpload(t *testing.T) {
	rsp, err := SmmsUpload("../../raw/twitter/images/FCBarcelona/Eq5P5dGW4AEDxdc.jpg", "Eq5P5dGW4AEDxdc.jpg")
	if err != nil {
		t.Error("SmmsUpload err: ", err)
	}
	_, err = SmmsHistory()
	if err != nil {
		t.Error("SmmsHistory err: ", err)
	}
	err = SmmsDelete(rsp.Hash)
	if err != nil {
		t.Error("SmmsDelete err: ", err)
	}
}
