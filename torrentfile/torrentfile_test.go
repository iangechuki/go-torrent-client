package torrentfile

import (
	"bytes"
	"os"
	"testing"

	"github.com/jackpal/bencode-go"
)

func TestOpenTorrentFile(t *testing.T) {
	// Create a mocj bencodeTorrent struct
	mockBencodeTorrent := bencodeTorrent{
		Announce: "http://tracker.example.com",
		Info: bencodeInfo{
			Pieces:      "abcdef1234567890abcdef1234567890abcdef12",
			PieceLength: 524288,
			Length:      1048576,
			Name:        "example.txt",
		},
	}
	tempFile, err := os.CreateTemp("", "test.torrent")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())
	var buf bytes.Buffer
	err = bencode.Marshal(&buf, mockBencodeTorrent)
	if err != nil {
		t.Fatalf("Failed to marshal mock torrent: %v", err)
	}

	_, err = tempFile.Write(buf.Bytes())
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()
	t.Run("valid torrent file", func(t *testing.T) {
		bto, err := OpenTorrentFile(tempFile.Name())
		if err != nil {
			t.Fatalf("Failed to open torrent file: %v", err)
		}
		t.Logf("parsed torrent file: %+v", bto)
		if bto.Announce != mockBencodeTorrent.Announce {
			t.Errorf("Expected announce to be %s, got %s", mockBencodeTorrent.Announce, bto.Announce)
		}
		if bto.Info.Name != mockBencodeTorrent.Info.Name {
			t.Errorf("Expected name to be %s, got %s", mockBencodeTorrent.Info.Name, bto.Info.Name)
		}
		if bto.Info.Length != mockBencodeTorrent.Info.Length {
			t.Errorf("Expected length to be %d, got %d", mockBencodeTorrent.Info.Length, bto.Info.Length)
		}
	})
	t.Run("missing torrent file", func(t *testing.T) {
		_, err := OpenTorrentFile("missing.torrent")
		if err == nil {
			t.Errorf("Expected error for missing file, got nil")
		} else {
			t.Logf("Received expected error for missing file: %v", err)
		}
	})
	t.Run("Invalid torrent file format", func(t *testing.T) {
		invalidFile, err := os.CreateTemp("", "invalid.torrent")
		if err != nil {
			t.Fatalf("Failed to create temp invalid file: %v", err)
		}
		defer os.Remove(invalidFile.Name())
		_, err = invalidFile.Write([]byte("invalid bencode data"))
		if err != nil {
			t.Fatalf("Failed to write to temp invalid file: %v", err)
		}
		invalidFile.Close()
		_, err = OpenTorrentFile(invalidFile.Name())
		if err == nil {
			t.Errorf("Expected error for invalid file, got nil")
		} else {
			t.Logf("Received expected error for invalid file: %v", err)
		}
	})
}
func TestSplitPieceHashes(t *testing.T) {
	t.Run("valid pieces string", func(t *testing.T) {
		info := bencodeInfo{
			Pieces:      "abcdef1234567890abcdef1234567890abcdef12",
			PieceLength: 524288,
		}
		hashes, err := info.splitPieceHashes()

		if err != nil {
			t.Fatalf("Expected no error,got %v", err)
		}
		if len(hashes) != 2 {
			t.Errorf("Expected 2 hashes, got %d", len(hashes))
		}
		// if err != nil {
	})
	t.Run("malformed pieces string", func(t *testing.T) {
		info := bencodeInfo{
			Pieces: "abcdef1234567890abcdef1234567890abc",
		}
		_, err := info.splitPieceHashes()
		if err == nil {
			t.Errorf("Expected error for malformed pieces string,got nil")
		} else {
			t.Logf("Received expected error for malformed pieces string: %v", err)
		}
	})
}
func TestToTorrentFile(t *testing.T) {
	t.Run("valid bencodeTorrent", func(t *testing.T) {
		bto := bencodeTorrent{
			Announce: "http://tracker.example.com",
			Info: bencodeInfo{
				Pieces:      "abcdef1234567890abcdef1234567890abcdef12",
				PieceLength: 524288,
				Length:      1048576,
				Name:        "example.txt",
			},
		}
		torrentFile, err := bto.toTorrentFile()
		if err != nil {
			t.Fatalf("Expected no error,got %v", err)
		}
		t.Logf("Converted to TorrentFile: %+v", torrentFile)
		if torrentFile.Name != bto.Info.Name {
			t.Errorf("Expected name %q,got %q", bto.Info.Name, torrentFile.Name)
		}
	})
}
