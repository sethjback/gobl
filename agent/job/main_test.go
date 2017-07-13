package job

import (
	"io"
	"os"
	"testing"

	"google.golang.org/grpc"

	"github.com/sethjback/gobl/files"
	"github.com/sethjback/gobl/goblgrpc"
	context "golang.org/x/net/context"
)

// tgrpc implements the goblgprc CoordinatorClient interface for mocking
type tgrpc struct {
	fclient goblgrpc.Coordinator_FileClient
	rclient goblgrpc.Coordinator_RestoreClient
}

func (c *tgrpc) File(ctx context.Context, opts ...grpc.CallOption) (goblgrpc.Coordinator_FileClient, error) {
	if c.fclient == nil {
		c.fclient = &fclient{}
	}
	return c.fclient, nil
}

func (c *tgrpc) Restore(ctx context.Context, in *goblgrpc.RestoreRequest, opts ...grpc.CallOption) (goblgrpc.Coordinator_RestoreClient, error) {
	if c.rclient == nil {
		c.rclient = &rclient{count: -1, jobID: in.Id}
	}
	return c.rclient, nil
}

func (c *tgrpc) State(ctx context.Context, in *goblgrpc.JobState, opts ...grpc.CallOption) (*goblgrpc.ReturnMessage, error) {
	return nil, nil
}

// fclient implements Coordinator_FileClient for mocking
type fclient struct {
	grpc.ClientStream
}

func (f *fclient) Send(*goblgrpc.FileRequest) error {
	return nil
}

func (f *fclient) CloseAndRecv() (*goblgrpc.ReturnMessage, error) {
	return nil, nil
}

// rclient implements Coordinator_FileClient for mocking
type rclient struct {
	grpc.ClientStream
	files []*files.File
	count int
	jobID string
}

func (r *rclient) Recv() (*goblgrpc.FileRequest, error) {
	r.count++
	if r.count > len(r.files) {
		return nil, io.EOF
	}
	f := r.files[r.count-1]

	return &goblgrpc.FileRequest{
		JobId: r.jobID,
		File: &goblgrpc.File{
			Signature: &goblgrpc.Signature{Path: f.Path, Hash: f.Hash, Modifications: f.Modifications},
			Meta:      &goblgrpc.Meta{Uid: int32(f.UID), Gid: int32(f.GID), Mode: f.Mode},
		},
	}, nil
}

func TestMain(m *testing.M) {
	e := buildDirectoryTree()
	if e != nil {
		os.Exit(-1)
	}
	code := m.Run()
	cleanUpDirectoryTree()
	os.Exit(code)
}
