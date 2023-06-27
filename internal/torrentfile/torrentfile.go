package torrentfile

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"math/rand"
	"os"

	"github.com/jackpal/bencode-go"
	"github.com/khachi-at/torrent-client/internal/p2p"
)

// Port to listen on.
const Port uint16 = 6881

// TorrentFile encodes the metadata from a .torrent file.
type TorrentFile struct {
	PieceLength int
	Length      int
	Announce    string
	Name        string
	InfoHash    [20]byte
	PieceHashes [][20]byte
}

func (t *TorrentFile) DownloadToFile(path string) error {
	var peerID [20]byte
	if _, err := rand.Read(peerID[:]); err != nil {
		return err
	}
	peers, err := t.requestPeers(peerID, Port)
	if err != nil {
		return err
	}
	torrent := p2p.Torrent{
		Peers:       peers,
		PeerID:      peerID,
		InfoHash:    t.InfoHash,
		PieceHashes: t.PieceHashes,
		PieceLength: t.PieceLength,
		Length:      t.Length,
		Name:        t.Name,
	}
	buf, err := torrent.Download()
	if err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(buf)
	return err
}

func Open(path string) (TorrentFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return TorrentFile{}, err
	}
	defer file.Close()

	t := bencodeTorrent{}
	if err := bencode.Unmarshal(file, &t); err != nil {
		return TorrentFile{}, err
	}
	return t.toTorrentFile()
}

type bencodeInfo struct {
	Name        string `bencode:"name"`
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Length      int    `bencode:"length"`
}

type bencodeTorrent struct {
	Announce string      `bencode:"announce"`
	Info     bencodeInfo `bencode:"info"`
}

func (i *bencodeInfo) hash() ([20]byte, error) {
	var buf bytes.Buffer
	if err := bencode.Marshal(&buf, *i); err != nil {
		return [20]byte{}, err
	}
	return sha1.Sum(buf.Bytes()), nil
}

func (i *bencodeInfo) splitPieceHashes() ([][20]byte, error) {
	buf := []byte(i.Pieces)
	bufLen := len(buf)
	hashLen := 20
	if bufLen%hashLen != 0 {
		err := fmt.Errorf("received malformed pieces of length %d", len(buf))
		return nil, err
	}
	hashesNum := bufLen / hashLen
	hashes := make([][20]byte, hashesNum)
	for i := 0; i < hashesNum; i++ {
		copy(hashes[i][:], buf[i*hashLen:(i+1)*hashLen])
	}
	return hashes, nil
}

func (t *bencodeTorrent) toTorrentFile() (TorrentFile, error) {
	hash, err := t.Info.hash()
	if err != nil {
		return TorrentFile{}, err
	}
	pHashes, err := t.Info.splitPieceHashes()
	if err != nil {
		return TorrentFile{}, err
	}
	ret := TorrentFile{
		Announce:    t.Announce,
		InfoHash:    hash,
		PieceHashes: pHashes,
		PieceLength: t.Info.PieceLength,
		Length:      t.Info.Length,
		Name:        t.Info.Name,
	}
	return ret, nil
}
