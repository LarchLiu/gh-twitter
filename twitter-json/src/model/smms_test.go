package model

import "testing"

func TestSmmsUpload(t *testing.T) {
	_, err := SmmsUpload("../../raw/twitter/images/FCBarcelona/Eq5P5dGW4AEDxdc.jpg", "Eq5P5dGW4AEDxdc.jpg")
	if err != nil {
		t.Error("SmmsUpload err: ", err)
	}
}
