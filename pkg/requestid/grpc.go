package requestid

import "google.golang.org/grpc/metadata"

const MetaName = "x-request-id"

func WithGRPCMetadata(md metadata.MD, requestID string) metadata.MD {
	md.Set(MetaName, requestID)

	return md
}

func FromGRPCMetadata(md metadata.MD) string {
	requestIDs := md.Get(MetaName)

	if len(requestIDs) > 0 {
		return requestIDs[0]
	}

	return ""
}
