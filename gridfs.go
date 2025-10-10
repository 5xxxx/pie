package pie

import (
	"context"
	"fmt"
	"io"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// GridFS GridFS file storage wrapper
type GridFS struct {
	engine *Engine
	bucket *mongo.GridFSBucket
}

// GridFSFile GridFS file information
type GridFSFile struct {
	ID          bson.ObjectID `bson:"_id"`
	Filename    string        `bson:"filename"`
	Length      int64         `bson:"length"`
	ChunkSize   int32         `bson:"chunkSize"`
	UploadDate  time.Time     `bson:"uploadDate"`
	MD5         string        `bson:"md5"`
	ContentType string        `bson:"contentType,omitempty"`
	Metadata    bson.M        `bson:"metadata,omitempty"`
}

// GridFSUploadOptions GridFS upload options
type GridFSUploadOptions struct {
	ChunkSizeBytes int32
	Metadata       bson.M
	ContentType    string
}

// GridFSDownloadOptions GridFS download options
type GridFSDownloadOptions struct {
	Revision int32
}

// GridFSQuery GridFS query options
type GridFSQuery struct {
	Filter bson.D
	Sort   bson.D
	Limit  int64
	Skip   int64
}

// NewGridFS create GridFS instance
func NewGridFS(engine *Engine) *GridFS {
	bucket := engine.database.GridFSBucket()
	return &GridFS{
		engine: engine,
		bucket: bucket,
	}
}

// Bucket get GridFS bucket
func (g *GridFS) Bucket() *mongo.GridFSBucket {
	return g.bucket
}

// Upload upload file
func (g *GridFS) Upload(ctx context.Context, filename string, source io.Reader, opts *GridFSUploadOptions) (bson.ObjectID, error) {
	start := time.Now()

	var uploadOpts []options.Lister[options.GridFSUploadOptions]
	if opts != nil {
		if opts.ChunkSizeBytes > 0 {
			uploadOpts = append(uploadOpts, options.GridFSUpload().SetChunkSizeBytes(opts.ChunkSizeBytes))
		}
		if opts.Metadata != nil {
			uploadOpts = append(uploadOpts, options.GridFSUpload().SetMetadata(opts.Metadata))
		}
	}

	fileID, err := g.bucket.UploadFromStream(ctx, filename, source, uploadOpts...)

	// Record query log
	if g.engine.queryLogger.IsEnabled() {
		g.engine.queryLogger.Log(&LogEntry{
			Timestamp:  start,
			Collection: "fs.files",
			Operation:  "upload",
			Filter:     bson.D{{Key: "filename", Value: filename}},
			Duration:   time.Since(start),
			Error:      err,
		})
	}

	return fileID, err
}

// UploadFromBytes upload file from byte array
func (g *GridFS) UploadFromBytes(ctx context.Context, filename string, data []byte, opts *GridFSUploadOptions) (bson.ObjectID, error) {
	return g.Upload(ctx, filename, &bytesReader{data: data}, opts)
}

// Download download file
func (g *GridFS) Download(ctx context.Context, fileID bson.ObjectID, dest io.Writer) (int64, error) {
	start := time.Now()

	bytesWritten, err := g.bucket.DownloadToStream(ctx, fileID, dest)

	// Record query log
	if g.engine.queryLogger.IsEnabled() {
		g.engine.queryLogger.Log(&LogEntry{
			Timestamp:  start,
			Collection: "fs.files",
			Operation:  "download",
			Filter:     bson.D{{Key: "_id", Value: fileID}},
			Duration:   time.Since(start),
			Error:      err,
		})
	}

	return bytesWritten, err
}

// DownloadToBytes download file to byte array
func (g *GridFS) DownloadToBytes(ctx context.Context, fileID bson.ObjectID) ([]byte, error) {
	var buf bytesBuffer
	_, err := g.Download(ctx, fileID, &buf)
	return buf.Bytes(), err
}

// DownloadByName download file by filename
func (g *GridFS) DownloadByName(ctx context.Context, filename string, dest io.Writer, opts *GridFSDownloadOptions) (int64, error) {
	start := time.Now()

	var downloadOpts []options.Lister[options.GridFSNameOptions]
	if opts != nil {
		downloadOpts = append(downloadOpts, options.GridFSName().SetRevision(opts.Revision))
	}

	bytesWritten, err := g.bucket.DownloadToStreamByName(ctx, filename, dest, downloadOpts...)

	// Record query log
	if g.engine.queryLogger.IsEnabled() {
		g.engine.queryLogger.Log(&LogEntry{
			Timestamp:  start,
			Collection: "fs.files",
			Operation:  "downloadByName",
			Filter:     bson.D{{Key: "filename", Value: filename}},
			Duration:   time.Since(start),
			Error:      err,
		})
	}

	return bytesWritten, err
}

// Find find files
func (g *GridFS) Find(ctx context.Context, query *GridFSQuery) (*Cursor[GridFSFile], error) {
	start := time.Now()

	var findOpts []options.Lister[options.GridFSFindOptions]
	if query != nil {
		if query.Sort != nil {
			findOpts = append(findOpts, options.GridFSFind().SetSort(query.Sort))
		}
		if query.Limit > 0 {
			findOpts = append(findOpts, options.GridFSFind().SetLimit(int32(query.Limit)))
		}
		if query.Skip > 0 {
			findOpts = append(findOpts, options.GridFSFind().SetSkip(int32(query.Skip)))
		}
	}

	filter := bson.D{}
	if query != nil && query.Filter != nil {
		filter = query.Filter
	}

	cursor, err := g.bucket.Find(ctx, filter, findOpts...)
	if err != nil {
		return nil, err
	}

	// Record query log
	if g.engine.queryLogger.IsEnabled() {
		g.engine.queryLogger.Log(&LogEntry{
			Timestamp:  start,
			Collection: "fs.files",
			Operation:  "find",
			Filter:     filter,
			Duration:   time.Since(start),
			Error:      err,
		})
	}

	return &Cursor[GridFSFile]{cursor: cursor}, nil
}

// FindOne find single file
func (g *GridFS) FindOne(ctx context.Context, query *GridFSQuery) (*GridFSFile, error) {
	cursor, err := g.Find(ctx, query)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		var file GridFSFile
		if err := cursor.Decode(&file); err != nil {
			return nil, err
		}
		return &file, nil
	}

	return nil, fmt.Errorf("file not found")
}

// Delete delete file
func (g *GridFS) Delete(ctx context.Context, fileID bson.ObjectID) error {
	start := time.Now()

	err := g.bucket.Delete(ctx, fileID)

	// Record query log
	if g.engine.queryLogger.IsEnabled() {
		g.engine.queryLogger.Log(&LogEntry{
			Timestamp:  start,
			Collection: "fs.files",
			Operation:  "delete",
			Filter:     bson.D{{Key: "_id", Value: fileID}},
			Duration:   time.Since(start),
			Error:      err,
		})
	}

	return err
}

// DeleteByName delete file by filename
func (g *GridFS) DeleteByName(ctx context.Context, filename string) error {
	// First find the file
	file, err := g.FindOne(ctx, &GridFSQuery{
		Filter: bson.D{{Key: "filename", Value: filename}},
	})
	if err != nil {
		return err
	}

	return g.Delete(ctx, file.ID)
}

// Rename rename file
func (g *GridFS) Rename(ctx context.Context, fileID bson.ObjectID, newFilename string) error {
	start := time.Now()

	err := g.bucket.Rename(ctx, fileID, newFilename)

	// Record query log
	if g.engine.queryLogger.IsEnabled() {
		g.engine.queryLogger.Log(&LogEntry{
			Timestamp:  start,
			Collection: "fs.files",
			Operation:  "rename",
			Filter:     bson.D{{Key: "_id", Value: fileID}},
			Duration:   time.Since(start),
			Error:      err,
		})
	}

	return err
}

// Drop delete GridFS bucket
func (g *GridFS) Drop(ctx context.Context) error {
	start := time.Now()

	err := g.bucket.Drop(ctx)

	// Record query log
	if g.engine.queryLogger.IsEnabled() {
		g.engine.queryLogger.Log(&LogEntry{
			Timestamp:  start,
			Collection: "fs.files",
			Operation:  "drop",
			Filter:     bson.D{},
			Duration:   time.Since(start),
			Error:      err,
		})
	}

	return err
}

// Helper types for GridFS operations

// bytesReader byte array reader
type bytesReader struct {
	data []byte
	pos  int
}

func (br *bytesReader) Read(p []byte) (n int, err error) {
	if br.pos >= len(br.data) {
		return 0, io.EOF
	}
	n = copy(p, br.data[br.pos:])
	br.pos += n
	return n, nil
}

// bytesBuffer byte array buffer
type bytesBuffer struct {
	data []byte
}

func (bb *bytesBuffer) Write(p []byte) (n int, err error) {
	bb.data = append(bb.data, p...)
	return len(p), nil
}

func (bb *bytesBuffer) Bytes() []byte {
	return bb.data
}
