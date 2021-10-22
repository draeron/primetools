package enginedj

import (
	"database/sql"
)

/*
CREATE TABLE Track (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  playOrder INTEGER,
  length INTEGER,
  bpm INTEGER,
  year INTEGER,
  path TEXT,
  filename TEXT,
  bitrate INTEGER,
  bpmAnalyzed REAL,
  albumArtId INTEGER,
  fileBytes INTEGER,
  title TEXT,
  artist TEXT,
  album TEXT,
  genre TEXT,
  comment TEXT,
  label TEXT,
  composer TEXT,
  remixer TEXT,
  key INTEGER,
  rating INTEGER,
  albumArt TEXT,
  timeLastPlayed DATETIME,
  isPlayed BOOLEAN,
  fileType TEXT,
  isAnalyzed BOOLEAN,
  dateCreated DATETIME,
  dateAdded DATETIME,
  isAvailable BOOLEAN,
  isMetadataOfPackedTrackChanged BOOLEAN,
  isPerfomanceDataOfPackedTrackChanged BOOLEAN,
  playedIndicator INTEGER,
  isMetadataImported BOOLEAN,
  pdbImportKey INTEGER,
  streamingSource TEXT,
  uri TEXT,
  isBeatGridLocked BOOLEAN,
  originDatabaseUuid TEXT,
  originTrackId INTEGER,
  trackData BLOB,
  overviewWaveFormData BLOB,
  beatData BLOB,
  quickCues BLOB,
  loops BLOB,
  thirdPartySourceId INTEGER,
  streamingFlags INTEGER,
  explicitLyrics BOOLEAN,
  CONSTRAINT C_originDatabaseUuid_originTrackId UNIQUE (originDatabaseUuid, originTrackId),
  CONSTRAINT C_path UNIQUE (path),
  FOREIGN KEY (albumArtId) REFERENCES AlbumArt (id) ON DELETE RESTRICT
)
*/
type trackEntry struct {
	Id       int            `db:"id"`
	Title    sql.NullString `db:"title"`
	Album sql.NullString `db:"album"`
	Artist sql.NullString `db:"artist"`

	Length   sql.NullInt32  `db:"length"`
	BPM      sql.NullInt32  `db:"bpm"`
	Year     sql.NullInt32  `db:"year"`
	Path     sql.NullString `db:"path"`
	Filename sql.NullString `db:"filename"`
	Bitrate  sql.NullInt32  `db:"bitrate"`
	Size     sql.NullInt32  `db:"fileBytes"`

	Rating  sql.NullInt32 `db:"rating"`
	Created sql.NullTime `json:"dateCreated"`
	Added   sql.NullTime `json:"dateAdded"`

	OriginTrackId      sql.NullInt32  `db:"originTrackId"`
	OriginDatabaseUuid sql.NullString `db:"originDatabaseUuid"`
}

/*
CREATE TABLE Playlist (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  title TEXT,
  parentListId INTEGER,
  isPersisted BOOLEAN,
  nextListId INTEGER,
  lastEditTime DATETIME,
  isExplicitlyExported BOOLEAN,
  CONSTRAINT C_NAME_UNIQUE_FOR_PARENT UNIQUE (title, parentListId),
  CONSTRAINT C_NEXT_LIST_ID_UNIQUE_FOR_PARENT UNIQUE (parentListId, nextListId)
)
*/
type playlistEntry struct {
	Id               int            `db:"id"`
	Title            sql.NullString `db:"title"`
	ParentListId     sql.NullInt32  `db:"parentListId"`
	IsPersisted      bool           `db:"isPersisted"`
	NextListId       sql.NullInt32  `db:"nextListId"`
	LastEditTime     sql.NullTime   `db:"lastEditTime"`
	ExplicitExported bool           `db:"isExplicitlyExported"`
}

/*
CREATE TABLE PlaylistEntity (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  listId INTEGER,
  trackId INTEGER,
  databaseUuid TEXT,
  nextEntityId INTEGER,
  membershipReference INTEGER,
  CONSTRAINT C_NAME_UNIQUE_FOR_LIST UNIQUE (listId, databaseUuid, trackId),
  FOREIGN KEY (listId) REFERENCES Playlist (id) ON DELETE CASCADE
)
*/
type playListEntityEntry struct {
	Id                  int            `db:"id"`
	ListId              sql.NullInt32  `db:"listId"`
	TrackId             sql.NullInt32  `db:"trackId"`
	DatabaseUuid        sql.NullString `db:"databaseUuid"`
	NextEntityId        sql.NullInt32  `db:"nextEntityId"`
	MembershipReference sql.NullInt32  `db:"membershipReference"`
}

type infoEntry struct {
	Id                 int    `db:"id"`
	UUID               string `db:"uuid"`
	SchemaVersionMajor int    `db:"schemaVersionMajor"`
	SchemaVersionMinor int    `db:"schemaVersionMinor"`
	SchemaVersionPatch int    `db:"schemaVersionPatch"`
}
