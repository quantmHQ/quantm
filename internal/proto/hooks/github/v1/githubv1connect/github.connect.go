// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: hooks/github/v1/github.proto

package githubv1connect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	v1 "go.breu.io/quantm/internal/proto/hooks/github/v1"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect.IsAtLeastVersion1_13_0

const (
	// GithubServiceName is the fully-qualified name of the GithubService service.
	GithubServiceName = "hooks.github.v1.GithubService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// GithubServiceInstallProcedure is the fully-qualified name of the GithubService's Install RPC.
	GithubServiceInstallProcedure = "/hooks.github.v1.GithubService/Install"
)

// These variables are the protoreflect.Descriptor objects for the RPCs defined in this package.
var (
	githubServiceServiceDescriptor       = v1.File_hooks_github_v1_github_proto.Services().ByName("GithubService")
	githubServiceInstallMethodDescriptor = githubServiceServiceDescriptor.Methods().ByName("Install")
)

// GithubServiceClient is a client for the hooks.github.v1.GithubService service.
type GithubServiceClient interface {
	// complete installation github app hook.
	Install(context.Context, *connect.Request[v1.InstallRequest]) (*connect.Response[emptypb.Empty], error)
}

// NewGithubServiceClient constructs a client for the hooks.github.v1.GithubService service. By
// default, it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses,
// and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the
// connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewGithubServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) GithubServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &githubServiceClient{
		install: connect.NewClient[v1.InstallRequest, emptypb.Empty](
			httpClient,
			baseURL+GithubServiceInstallProcedure,
			connect.WithSchema(githubServiceInstallMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
	}
}

// githubServiceClient implements GithubServiceClient.
type githubServiceClient struct {
	install *connect.Client[v1.InstallRequest, emptypb.Empty]
}

// Install calls hooks.github.v1.GithubService.Install.
func (c *githubServiceClient) Install(ctx context.Context, req *connect.Request[v1.InstallRequest]) (*connect.Response[emptypb.Empty], error) {
	return c.install.CallUnary(ctx, req)
}

// GithubServiceHandler is an implementation of the hooks.github.v1.GithubService service.
type GithubServiceHandler interface {
	// complete installation github app hook.
	Install(context.Context, *connect.Request[v1.InstallRequest]) (*connect.Response[emptypb.Empty], error)
}

// NewGithubServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewGithubServiceHandler(svc GithubServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	githubServiceInstallHandler := connect.NewUnaryHandler(
		GithubServiceInstallProcedure,
		svc.Install,
		connect.WithSchema(githubServiceInstallMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	return "/hooks.github.v1.GithubService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case GithubServiceInstallProcedure:
			githubServiceInstallHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedGithubServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedGithubServiceHandler struct{}

func (UnimplementedGithubServiceHandler) Install(context.Context, *connect.Request[v1.InstallRequest]) (*connect.Response[emptypb.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("hooks.github.v1.GithubService.Install is not implemented"))
}
