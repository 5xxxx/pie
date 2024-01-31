/*
 *
 * options.go
 * driver
 *
 * Created by lintao on 2024/1/25 15:20
 * Copyright © 2020-2024 LINTAO. All rights reserved.
 *
 */

package pie

import (
	"crypto/tls"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"time"
)

// ClientOptions represents the options for configuring a client session.
type ClientOptions interface {
	// SetArrayFilters sets the value for the ArrayFilters field.
	SetArrayFilters(filters options.ArrayFilters) Session
	// SetOrdered sets the value for the Ordered field.
	SetOrdered(ordered bool) Session
	// SetBypassDocumentValidation sets the value for the BypassDocumentValidation field.
	SetBypassDocumentValidation(b bool) Session
	// SetReturnDocument sets the value for the ReturnDocument field.
	SetReturnDocument(rd options.ReturnDocument) Session

	// SetUpsert sets the value for the Upsert field.
	SetUpsert(b bool) Session

	// SetCollation sets the value for the Collation field.
	SetCollation(collation *options.Collation) Session

	// SetMaxTime sets the value for the MaxTime field.
	SetMaxTime(d time.Duration) Session
	// SetProjection sets the value for the Projection field.
	SetProjection(projection any) Session

	// SetSort sets the value for the Sort field.
	SetSort(sort any) Session

	// SetHint sets the value for the Hint field.
	SetHint(hint any) Session
}

type Options struct {
	UpdateEmpty bool
}

// SetUpdateEmpty sets the value for the UpdateEmpty field.
func (o *Options) SetUpdateEmpty(e bool) {
	o.UpdateEmpty = e
}

// SetOrdered sets the value for the Ordered field.
func (d *defaultClient) SetOrdered(ordered bool) Session {
	return d.NewSession().SetOrdered(ordered)
}

// SetArrayFilters sets the value for the ArrayFilters field.
func (d *defaultClient) SetArrayFilters(filters options.ArrayFilters) Session {
	return d.NewSession().SetArrayFilters(filters)
}

// SetBypassDocumentValidation sets the value for the BypassDocumentValidation field.
func (d *defaultClient) SetBypassDocumentValidation(b bool) Session {
	return d.NewSession().SetBypassDocumentValidation(b)
}

// SetReturnDocument sets the value for the ReturnDocument field.
func (d *defaultClient) SetReturnDocument(rd options.ReturnDocument) Session {
	return d.NewSession().SetReturnDocument(rd)
}

// SetUpsert sets the value for the Upsert field.
func (d *defaultClient) SetUpsert(b bool) Session {
	return d.NewSession().SetUpsert(b)
}

// SetCollation sets the value for the Collation field.
func (d *defaultClient) SetCollation(collation *options.Collation) Session {
	return d.NewSession().SetCollation(collation)
}

// SetMaxTime sets the maximum amount of time for a command to run before an error is returned.
// The time.Duration parameter represents the maximum time duration in which a command can run
// without returning an error.
func (d *defaultClient) SetMaxTime(t time.Duration) Session {
	return d.NewSession().SetMaxTime(t)
}

// SetProjection sets the value for the Projection field.
func (d *defaultClient) SetProjection(projection any) Session {
	return d.NewSession().SetProjection(projection)
}

// SetSort sets the value for the sort field of the session created by defaultClient.
// It takes an argument "sort" of type any which represents the sorting mechanism to be applied.
// The method returns a Session object which is created using the NewSession method.
// The SetSort method is responsible for setting the sort value of the session and returning the modified session.
// Example usage:
//
//	client := &defaultClient{}
//	sortValue := "ascending"
//	session := client.SetSort(sortValue)
//	session.DoSomething()
//	...
//	session.DoSomethingElse()
func (d *defaultClient) SetSort(sort any) Session {
	return d.NewSession().SetSort(sort)
}

// SetHint sets the hint for the session.
// It takes an any argument representing the hint.
// The method creates a new session using the defaultClient's NewSession method.
// It then sets the hint for the session using the SetHint method on the session.
// Finally, it returns the session with the new hint set.
func (d *defaultClient) SetHint(hint any) Session {
	return d.NewSession().SetHint(hint)
}

// SetURI sets the URI for the defaultClient instance.
// The URI is applied to the clientOpts field using the options.Client().ApplyURI() method.
func (d *defaultClient) SetURI(uri string) {
	d.clientOpts = append(d.clientOpts, options.Client().ApplyURI(uri))
}

// SetAppName sets the value for the AppName field.
func (d *defaultClient) SetAppName(s string) {
	d.clientOpts = append(d.clientOpts, options.Client().SetAppName(s))
}

// SetAuth sets the value for the Auth field.
func (d *defaultClient) SetAuth(auth options.Credential) {
	d.clientOpts = append(d.clientOpts, options.Client().SetAuth(auth))
}

// SetCompressors sets the compressors that can be used when communicating with a server. Valid values are:
//
// 1. "snappy" - requires server version >= 3.4
//
// 2. "zlib" - requires server version >= 3.6
//
//  3. "zstd" - requires server version >= 4.2, and internal version >= 1.2.0 with cgo support enabled or internal version >= 1.3.0
//     without cgo
//
// To use compression, it must be enabled on the server as well. If this option is specified, the internal will perform a
// negotiation with the server to determine a common list of of compressors and will use the first one in that list when
// performing operations. See
// https://docs.mongodb.com/manual/reference/program/mongod/#cmdoption-mongod-networkmessagecompressors for more
// information about how to enable this feature on the server.
//
// This can also be set through the "compressors" URI option (e.g. "compressors=zstd,zlib,snappy"). The default is
// an empty slice, meaning no compression will be enabled.
func (d *defaultClient) SetCompressors(comps []string) {
	d.clientOpts = append(d.clientOpts, options.Client().SetCompressors(comps))
}

// SetConnectTimeout sets the value for the ConnectTimeout field in the client options.
func (d *defaultClient) SetConnectTimeout(t time.Duration) {
	d.clientOpts = append(d.clientOpts, options.Client().SetConnectTimeout(t))
}

// SetDialer sets the value for the ContextDialer field.
func (d *defaultClient) SetDialer(t options.ContextDialer) {
	d.clientOpts = append(d.clientOpts, options.Client().SetDialer(t))
}

// SetDirect sets the value for the Direct field.
func (d *defaultClient) SetDirect(b bool) {
	d.clientOpts = append(d.clientOpts, options.Client().SetDirect(b))
}

// SetHeartbeatInterval sets the value for the HeartbeatInterval field.
func (d *defaultClient) SetHeartbeatInterval(t time.Duration) {
	d.clientOpts = append(d.clientOpts, options.Client().SetHeartbeatInterval(t))
}

// SetHosts sets the value for the Hosts field in the options.Client() object.
// It appends the given slice of strings (s) to the clientOpts slice using the options.Client().SetHosts() function.
func (d *defaultClient) SetHosts(s []string) {
	d.clientOpts = append(d.clientOpts, options.Client().SetHosts(s))
}

// SetLocalThreshold sets the value for the LocalThreshold field.
func (d *defaultClient) SetLocalThreshold(t time.Duration) {
	d.clientOpts = append(d.clientOpts, options.Client().SetLocalThreshold(t))
}

// SetMaxConnIdleTime sets the value for the MaxConnIdleTime field.
func (d *defaultClient) SetMaxConnIdleTime(t time.Duration) {
	d.clientOpts = append(d.clientOpts, options.Client().SetMaxConnIdleTime(t))
}

// SetMaxPoolSize sets the value for the MaxPoolSize field.
func (d *defaultClient) SetMaxPoolSize(u uint64) {
	d.clientOpts = append(d.clientOpts, options.Client().SetMaxPoolSize(u))
}

// SetMinPoolSize sets the value for the MinPoolSize field.
func (d *defaultClient) SetMinPoolSize(u uint64) {
	d.clientOpts = append(d.clientOpts, options.Client().SetMinPoolSize(u))
}

// SetPoolMonitor sets the value for the PoolMonitor field.
func (d *defaultClient) SetPoolMonitor(m *event.PoolMonitor) {
	d.clientOpts = append(d.clientOpts, options.Client().SetPoolMonitor(m))
}

// SetMonitor sets the value for the Monitor field.
func (d *defaultClient) SetMonitor(m *event.CommandMonitor) {
	d.clientOpts = append(d.clientOpts, options.Client().SetMonitor(m))
}

// SetReadConcern appends the given read concern to the client options.
func (d *defaultClient) SetReadConcern(rc *readconcern.ReadConcern) {
	d.clientOpts = append(d.clientOpts, options.Client().SetReadConcern(rc))
}

// SetReadPreference sets the value for the ReadPreference field.
func (d *defaultClient) SetReadPreference(rp *readpref.ReadPref) {
	d.clientOpts = append(d.clientOpts, options.Client().SetReadPreference(rp))
}

// SetRegistry sets the value for the Registry field of the defaultClient struct.
func (d *defaultClient) SetRegistry(registry *bsoncodec.Registry) {
	d.clientOpts = append(d.clientOpts, options.Client().SetRegistry(registry))
}

// SetReplicaSet sets the value for the ReplicaSet field.
func (d *defaultClient) SetReplicaSet(s string) {
	d.clientOpts = append(d.clientOpts, options.Client().SetReplicaSet(s))
}

// SetRetryWrites sets the value for the RetryWrites field.
func (d *defaultClient) SetRetryWrites(b bool) {
	d.clientOpts = append(d.clientOpts, options.Client().SetRetryWrites(b))
}

// SetRetryReads sets the value for the RetryReads field.
func (d *defaultClient) SetRetryReads(b bool) {
	d.clientOpts = append(d.clientOpts, options.Client().SetRetryReads(b))
}

// SetServerSelectionTimeout sets the value for the ServerSelectionTimeout field.
func (d *defaultClient) SetServerSelectionTimeout(t time.Duration) {
	d.clientOpts = append(d.clientOpts, options.Client().SetServerSelectionTimeout(t))
}

// SetSocketTimeout sets the value for the SocketTimeout field.
func (d *defaultClient) SetSocketTimeout(t time.Duration) {
	d.clientOpts = append(d.clientOpts, options.Client().SetSocketTimeout(t))
}

// SetTLSConfig specifies a tls.Config instance to use use to configure TLS on all connections created to the cluster.
// This can also be set through the following URI options:
//
// 1. "tls" (or "ssl"): Specify if TLS should be used (e.g. "tls=true").
//
// 2. Either "tlsCertificateKeyFile" (or "sslClientCertificateKeyFile") or a combination of "tlsCertificateFile" and
// "tlsPrivateKeyFile". The "tlsCertificateKeyFile" option specifies a path to the client certificate and private key,
// which must be concatenated into one file. The "tlsCertificateFile" and "tlsPrivateKey" combination specifies separate
// paths to the client certificate and private key, respectively. Note that if "tlsCertificateKeyFile" is used, the
// other two options must not be specified.
//
// 3. "tlsCertificateKeyFilePassword" (or "sslClientCertificateKeyPassword"): Specify the password to decrypt the client
// private key file (e.g. "tlsCertificateKeyFilePassword=password").
//
// 4. "tlsCaFile" (or "sslCertificateAuthorityFile"): Specify the path to a single or bundle of certificate authorities
// to be considered trusted when making a TLS connection (e.g. "tlsCaFile=/path/to/caFile").
//
// 5. "tlsInsecure" (or "sslInsecure"): Specifies whether or not certificates and hostnames received from the server
// should be validated. If true (e.g. "tlsInsecure=true"), the TLS library will accept any certificate presented by the
// server and any host name in that certificate. Note that setting this to true makes TLS susceptible to
// man-in-the-middle attacks and should only be done for testing.
//
// The default is nil, meaning no TLS will be enabled.
func (d *defaultClient) SetTLSConfig(cfg *tls.Config) {
	d.clientOpts = append(d.clientOpts, options.Client().SetTLSConfig(cfg))
}

// SetWriteConcern sets the value for the WriteConcern field.
func (d *defaultClient) SetWriteConcern(wc *writeconcern.WriteConcern) {
	d.clientOpts = append(d.clientOpts, options.Client().SetWriteConcern(wc))
}

// SetZlibLevel sets the value for the ZlibLevel field of the options.Client.
func (d *defaultClient) SetZlibLevel(level int) {
	d.clientOpts = append(d.clientOpts, options.Client().SetZlibLevel(level))
}

// SetZstdLevel sets the value for the ZstdLevel field.
func (d *defaultClient) SetZstdLevel(level int) {
	d.clientOpts = append(d.clientOpts, options.Client().SetZstdLevel(level))
}

// SetAutoEncryptionOptions sets the value for the AutoEncryptionOptions field.
func (d *defaultClient) SetAutoEncryptionOptions(opts *options.AutoEncryptionOptions) {
	d.clientOpts = append(d.clientOpts, options.Client().SetAutoEncryptionOptions(opts))
}
