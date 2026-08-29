package share

import "testing"

func TestJsonShareInfoToDomainConvertsLegacyFileData(t *testing.T) {
	shareInfo, err := JsonShareInfoToDomain(`{"type":"file","data":"file-id","file_name":"example.txt"}`)
	if err != nil {
		t.Fatal(err)
	}
	if shareInfo.Data != "file-id" {
		t.Fatalf("Data = %q, want %q", shareInfo.Data, "file-id")
	}
	if len(shareInfo.Files) != 1 {
		t.Fatalf("len(Files) = %d, want 1", len(shareInfo.Files))
	}
	if shareInfo.Files[0].Id != "file-id" || shareInfo.Files[0].FileName != "example.txt" {
		t.Fatalf("Files[0] = %#v", shareInfo.Files[0])
	}
}

func TestJsonShareInfoToDomainKeepsCurrentFiles(t *testing.T) {
	shareInfo, err := JsonShareInfoToDomain(`{"type":"file","data":"legacy-id","file_name":"legacy.txt","files":[{"id":"current-id","file_name":"current.txt"}]}`)
	if err != nil {
		t.Fatal(err)
	}
	if len(shareInfo.Files) != 1 || shareInfo.Files[0].Id != "current-id" {
		t.Fatalf("Files = %#v", shareInfo.Files)
	}
}

func TestJsonShareInfoToDomainDoesNotConvertFileWithoutData(t *testing.T) {
	shareInfo, err := JsonShareInfoToDomain(`{"type":"file","file_name":"example.txt"}`)
	if err != nil {
		t.Fatal(err)
	}
	if len(shareInfo.Files) != 0 {
		t.Fatalf("Files = %#v, want empty", shareInfo.Files)
	}
}

func TestJsonShareInfoToDomainConvertsLegacyTextData(t *testing.T) {
	shareInfo, err := JsonShareInfoToDomain(`{"type":"text","data":"legacy text"}`)
	if err != nil {
		t.Fatal(err)
	}
	if shareInfo.Text != "legacy text" {
		t.Fatalf("Text = %q, want %q", shareInfo.Text, "legacy text")
	}
}
