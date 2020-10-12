/*
 *
 * options.go
 * pie
 *
 * Created by lintao on 2020/6/8 4:05 下午
 * Copyright © 2020-2020 LINTAO. All rights reserved.
 *
 */

package pie

import (
	"crypto/tls"
	"time"

	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type SessionOptions interface {
	// SetArrayFilters sets the value for the ArrayFilters field.
	SetArrayFilters(filters options.ArrayFilters) *Session

	// SetOrdered sets the value for the Ordered field.
	SetOrdered(ordered bool) *Session
	// SetBypassDocumentValidation sets the value for the BypassDocumentValidation field.
	SetBypassDocumentValidation(b bool) *Session

	// SetReturnDocument sets the value for the ReturnDocument field.
	SetReturnDocument(rd options.ReturnDocument) *Session

	// SetUpsert sets the value for the Upsert field.
	SetUpsert(b bool) *Session

	// SetCollation sets the value for the Collation field.
	SetCollation(collation *options.Collation) *Session

	// SetMaxTime sets the value for the MaxTime field.
	SetMaxTime(d time.Duration) *Session
	// SetProjection sets the value for the Projection field.
	SetProjection(projection interface{}) *Session

	// SetSort sets the value for the Sort field.
	SetSort(sort interface{}) *Session

	// SetHint sets the value for the Hint field.
	SetHint(hint interface{}) *Session
}

type Options struct {
	UpdateEmpty bool
}

func (o *Options) SetUpdateEmpty(e bool) {
	o.UpdateEmpty = e
}

func (d *Driver) SetOrdered(ordered bool) *Session {
	session := d.NewSession()
	return session.SetOrdered(ordered)
}

// SetArrayFilters sets the value for the ArrayFilters field.
func (d *Driver) SetArrayFilters(filters options.ArrayFilters) *Session {
	session := d.NewSession()
	return session.SetArrayFilters(filters)
}

// SetBypassDocumentValidation sets the value for the BypassDocumentValidation field.
func (d *Driver) SetBypassDocumentValidation(b bool) *Session {
	session := d.NewSession()
	return session.SetBypassDocumentValidation(b)
}

// SetReturnDocument sets the value for the ReturnDocument field.
func (d *Driver) SetReturnDocument(rd options.ReturnDocument) *Session {
	session := d.NewSession()
	return session.SetReturnDocument(rd)
}

// SetUpsert sets the value for the Upsert field.
func (d *Driver) SetUpsert(b bool) *Session {
	session := d.NewSession()
	return session.SetUpsert(b)
}

// SetCollation sets the value for the Collation field.
func (d *Driver) SetCollation(collation *options.Collation) *Session {
	session := d.NewSession()
	return session.SetCollation(collation)
}

// SetMaxTime sets the value for the MaxTime field.
func (d *Driver) SetMaxTime(t time.Duration) *Session {
	session := d.NewSession()
	return session.SetMaxTime(t)
}

// SetProjection sets the value for the Projection field.
func (d *Driver) SetProjection(projection interface{}) *Session {
	session := d.NewSession()
	return session.SetProjection(projection)
}

// SetSort sets the value for the Sort field.
func (d *Driver) SetSort(sort interface{}) *Session {
	session := d.NewSession()
	return session.SetSort(sort)
}

// SetHint sets the value for the Hint field.
func (d *Driver) SetHint(hint interface{}) *Session {
	session := d.NewSession()
	return session.SetHint(hint)
}

// SetURI parses the given URI and sets options accordingly. The URI can contain host names, IPv4/IPv6 literals, or
// an SRV record that will be resolved when the Client is created. When using an SRV record, TLS support is
// implictly enabled. Specify the "tls=false" URI option to override this.
//
// If the connection string contains any options that have previously been set, it will overwrite them. Options that
// correspond to multiple URI parameters, such as WriteConcern, will be completely overwritten if any of the query
// parameters are specified. If an option is set on ClientOptions after this method is called, that option will override
// any option applied via the connection string.
//
// If the URI format is incorrect or there are conflicing options specified in the URI an error will be recorded and
// can be retrieved by calling Validate.
//
// For more information about the URI format, see https://docs.mongodb.com/manual/reference/connection-string/. See
// mongo.Connect documentation for examples of using URIs for different Client configurations.
func (d *Driver) SetURI(uri string) {
	d.clientOpts = append(d.clientOpts, options.Client().ApplyURI(uri))
}

// SetAppName specifies an application name that is sent to the server when creating new connections. It is used by the
// server to log connection and profiling information (e.g. slow query logs). This can also be set through the "appName"
// URI option (e.g "appName=example_application"). The default is empty, meaning no app name will be sent.
func (d *Driver) SetAppName(s string) {
	d.clientOpts = append(d.clientOpts, options.Client().SetAppName(s))
}

// SetAuth specifies a Credential containing options for configuring authentication. See the options.Credential
// documentation for more information about Credential fields. The default is an empty Credential, meaning no
// authentication will be configured.
func (d *Driver) SetAuth(auth options.Credential) {
	d.clientOpts = append(d.clientOpts, options.Client().SetAuth(auth))
}

// SetCompressors sets the compressors that can be used when communicating with a server. Valid values are:
//
// 1. "snappy" - requires server version >= 3.4
//
// 2. "zlib" - requires server version >= 3.6
//
// 3. "zstd" - requires server version >= 4.2, and driver version >= 1.2.0 with cgo support enabled or driver version >= 1.3.0
//    without cgo
//
// To use compression, it must be enabled on the server as well. If this option is specified, the driver will perform a
// negotiation with the server to determine a common list of of compressors and will use the first one in that list when
// performing operations. See
// https://docs.mongodb.com/manual/reference/program/mongod/#cmdoption-mongod-networkmessagecompressors for more
// information about how to enable this feature on the server.
//
// This can also be set through the "compressors" URI option (e.g. "compressors=zstd,zlib,snappy"). The default is
// an empty slice, meaning no compression will be enabled.
func (d *Driver) SetCompressors(comps []string) {
	d.clientOpts = append(d.clientOpts, options.Client().SetCompressors(comps))
}

// SetConnectTimeout specifies a timeout that is used for creating connections to the server. If a custom Dialer is
// specified through SetDialer, this option must not be used. This can be set through SetURI with the
// "connectTimeoutMS" (e.g "connectTimeoutMS=30") option. If set to 0, no timeout will be used. The default is 30
// seconds.
func (d *Driver) SetConnectTimeout(t time.Duration) {
	d.clientOpts = append(d.clientOpts, options.Client().SetConnectTimeout(t))
}

// SetDialer specifies a custom ContextDialer to be used to create new connections to the server. The default is a
// net.Dialer instance with a 300 second keepalive time.
func (d *Driver) SetDialer(t options.ContextDialer) {
	d.clientOpts = append(d.clientOpts, options.Client().SetDialer(t))
}

// SetDirect specifies whether or not a direct connect should be made. To use this option, a URI with a single host must
// be specified through SetURI. If set to true, the driver will only connect to the host provided in the URI and will
// not discover other hosts in the cluster. This can also be set through the "connect" URI option with the following
// values:
//
// 1. "connect=direct" for direct connections
//
// 2. "connect=automatic" for automatic discovery.
//
// The default is false ("automatic" in the connection string).
func (d *Driver) SetDirect(b bool) {
	d.clientOpts = append(d.clientOpts, options.Client().SetDirect(b))
}

// SetHeartbeatInterval specifies the amount of time to wait between periodic background server checks. This can also be
// set through the "heartbeatIntervalMS" URI option (e.g. "heartbeatIntervalMS=10000"). The default is 10 seconds.
func (d *Driver) SetHeartbeatInterval(t time.Duration) {
	d.clientOpts = append(d.clientOpts, options.Client().SetHeartbeatInterval(t))
}

// SetHosts specifies a list of host names or IP addresses for servers in a cluster. Both IPv4 and IPv6 addresses are
// supported. IPv6 literals must be enclosed in '[]' following RFC-2732 syntax.
//
// Hosts can also be specified as a comma-separated list in a URI. For example, to include "localhost:27017" and
// "localhost:27018", a URI could be "mongodb://localhost:27017,localhost:27018". The default is ["localhost:27017"]
func (d *Driver) SetHosts(s []string) {
	d.clientOpts = append(d.clientOpts, options.Client().SetHosts(s))
}

// SetLocalThreshold specifies the width of the 'latency window': when choosing between multiple suitable servers for an
// operation, this is the acceptable non-negative delta between shortest and longest average round-trip times. A server
// within the latency window is selected randomly. This can also be set through the "localThresholdMS" URI option (e.g.
// "localThresholdMS=15000"). The default is 15 milliseconds.
func (d *Driver) SetLocalThreshold(t time.Duration) {
	d.clientOpts = append(d.clientOpts, options.Client().SetLocalThreshold(t))
}

// SetMaxConnIdleTime specifies the maximum amount of time that a connection will remain idle in a connection pool
// before it is removed from the pool and closed. This can also be set through the "maxIdleTimeMS" URI option (e.g.
// "maxIdleTimeMS=10000"). The default is 0, meaning a connection can remain unused indefinitely.
func (d *Driver) SetMaxConnIdleTime(t time.Duration) {
	d.clientOpts = append(d.clientOpts, options.Client().SetMaxConnIdleTime(t))
}

// SetMaxPoolSize specifies that maximum number of connections allowed in the driver's connection pool to each server.
// Requests to a server will block if this maximum is reached. This can also be set through the "maxPoolSize" URI option
// (e.g. "maxPoolSize=100"). The default is 100. If this is 0, it will be set to math.MaxInt64.
func (d *Driver) SetMaxPoolSize(u uint64) {
	d.clientOpts = append(d.clientOpts, options.Client().SetMaxPoolSize(u))
}

// SetMinPoolSize specifies the minimum number of connections allowed in the driver's connection pool to each server. If
// this is non-zero, each server's pool will be maintained in the background to ensure that the size does not fall below
// the minimum. This can also be set through the "minPoolSize" URI option (e.g. "minPoolSize=100"). The default is 0.
func (d *Driver) SetMinPoolSize(u uint64) {
	d.clientOpts = append(d.clientOpts, options.Client().SetMinPoolSize(u))
}

// SetPoolMonitor specifies a PoolMonitor to receive connection pool events. See the event.PoolMonitor documentation
// for more information about the structure of the monitor and events that can be received.
func (d *Driver) SetPoolMonitor(m *event.PoolMonitor) {
	d.clientOpts = append(d.clientOpts, options.Client().SetPoolMonitor(m))
}

// SetMonitor specifies a CommandMonitor to receive command events. See the event.CommandMonitor documentation for more
// information about the structure of the monitor and events that can be received.
func (d *Driver) SetMonitor(m *event.CommandMonitor) {
	d.clientOpts = append(d.clientOpts, options.Client().SetMonitor(m))
}

// SetReadConcern specifies the read concern to use for read operations. A read concern level can also be set through
// the "readConcernLevel" URI option (e.g. "readConcernLevel=majority"). The default is nil, meaning the server will use
// its configured default.
func (d *Driver) SetReadConcern(rc *readconcern.ReadConcern) {
	d.clientOpts = append(d.clientOpts, options.Client().SetReadConcern(rc))
}

// SetReadPreference specifies the read preference to use for read operations. This can also be set through the
// following URI options:
//
// 1. "readPreference" - Specifiy the read preference mode (e.g. "readPreference=primary").
//
// 2. "readPreferenceTags": Specify one or more read preference tags
// (e.g. "readPreferenceTags=region:south,datacenter:A").
//
// 3. "maxStalenessSeconds" (or "maxStaleness"): Specify a maximum replication lag for reads from secondaries in a
// replica set (e.g. "maxStalenessSeconds=10").
//
// The default is readpref.Primary(). See https://docs.mongodb.com/manual/core/read-preference/#read-preference for
// more information about read preferences.
func (d *Driver) SetReadPreference(rp *readpref.ReadPref) {
	d.clientOpts = append(d.clientOpts, options.Client().SetReadPreference(rp))
}

// SetRegistry specifies the BSON registry to use for BSON marshalling/unmarshalling operations. The default is
// bson.DefaultRegistry.
func (d *Driver) SetRegistry(registry *bsoncodec.Registry) {
	d.clientOpts = append(d.clientOpts, options.Client().SetRegistry(registry))
}

// SetReplicaSet specifies the replica set name for the cluster. If specified, the cluster will be treated as a replica
// set and the driver will automatically discover all servers in the set, starting with the nodes specified through
// SetURI or SetHosts. All nodes in the replica set must have the same replica set name, or they will not be
// considered as part of the set by the Client. This can also be set through the "replicaSet" URI option (e.g.
// "replicaSet=replset"). The default is empty.
func (d *Driver) SetReplicaSet(s string) {
	d.clientOpts = append(d.clientOpts, options.Client().SetReplicaSet(s))
}

// SetRetryWrites specifies whether supported write operations should be retried once on certain errors, such as network
// errors.
//
// Supported operations are InsertOne, UpdateOne, ReplaceOne, DeleteOne, FindOneAndDelete, FindOneAndReplace,
// FindOneAndDelete, InsertMany, and BulkWrite. Note that BulkWrite requests must not include UpdateManyModel or
// DeleteManyModel instances to be considered retryable. Unacknowledged writes will not be retried, even if this option
// is set to true.
//
// This option requires server version >= 3.6 and a replica set or sharded cluster and will be ignored for any other
// cluster type. This can also be set through the "retryWrites" URI option (e.g. "retryWrites=true"). The default is
// true.
func (d *Driver) SetRetryWrites(b bool) {
	d.clientOpts = append(d.clientOpts, options.Client().SetRetryWrites(b))
}

// SetRetryReads specifies whether supported read operations should be retried once on certain errors, such as network
// errors.
//
// Supported operations are Find, FindOne, Aggregate without a $out stage, Distinct, CountDocuments,
// EstimatedDocumentCount, Watch (for Client, Database, and Collection), ListCollections, and ListDatabases. Note that
// operations run through RunCommand are not retried.
//
// This option requires server version >= 3.6 and driver version >= 1.1.0. The default is true.
func (d *Driver) SetRetryReads(b bool) {
	d.clientOpts = append(d.clientOpts, options.Client().SetRetryReads(b))
}

// SetServerSelectionTimeout specifies how long the driver will wait to find an available, suitable server to execute an
// operation. This can also be set through the "serverSelectionTimeoutMS" URI option (e.g.
// "serverSelectionTimeoutMS=30000"). The default value is 30 seconds.
func (d *Driver) SetServerSelectionTimeout(t time.Duration) {
	d.clientOpts = append(d.clientOpts, options.Client().SetServerSelectionTimeout(t))
}

// SetSocketTimeout specifies how long the driver will wait for a socket read or write to return before returning a
// network error. This can also be set through the "socketTimeoutMS" URI option (e.g. "socketTimeoutMS=1000"). The
// default value is 0, meaning no timeout is used and socket operations can block indefinitely.
func (d *Driver) SetSocketTimeout(t time.Duration) {
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
func (d *Driver) SetTLSConfig(cfg *tls.Config) {
	d.clientOpts = append(d.clientOpts, options.Client().SetTLSConfig(cfg))
}

// SetWriteConcern specifies the write concern to use to for write operations. This can also be set through the following
// URI options:
//
// 1. "w": Specify the number of nodes in the cluster that must acknowledge write operations before the operation
// returns or "majority" to specify that a majority of the nodes must acknowledge writes. This can either be an integer
// (e.g. "w=10") or the string "majority" (e.g. "w=majority").
//
// 2. "wTimeoutMS": Specify how long write operations should wait for the correct number of nodes to acknowledge the
// operation (e.g. "wTimeoutMS=1000").
//
// 3. "journal": Specifies whether or not write operations should be written to an on-disk journal on the server before
// returning (e.g. "journal=true").
//
// The default is nil, meaning the server will use its configured default.
func (d *Driver) SetWriteConcern(wc *writeconcern.WriteConcern) {
	d.clientOpts = append(d.clientOpts, options.Client().SetWriteConcern(wc))
}

// SetZlibLevel specifies the level for the zlib compressor. This option is ignored if zlib is not specified as a
// compressor through SetURI or SetCompressors. Supported values are -1 through 9, inclusive. -1 tells the zlib
// library to use its default, 0 means no compression, 1 means best speed, and 9 means best compression.
// This can also be set through the "zlibCompressionLevel" URI option (e.g. "zlibCompressionLevel=-1"). Defaults to -1.
func (d *Driver) SetZlibLevel(level int) {
	d.clientOpts = append(d.clientOpts, options.Client().SetZlibLevel(level))
}

// SetZstdLevel sets the level for the zstd compressor. This option is ignored if zstd is not specified as a compressor
// through SetURI or SetCompressors. Supported values are 1 through 20, inclusive. 1 means best speed and 20 means
// best compression. This can also be set through the "zstdCompressionLevel" URI option. Defaults to 6.
func (d *Driver) SetZstdLevel(level int) {
	d.clientOpts = append(d.clientOpts, options.Client().SetZstdLevel(level))
}

// SetAutoEncryptionOptions specifies an AutoEncryptionOptions instance to automatically encrypt and decrypt commands
// and their results. See the options.AutoEncryptionOptions documentation for more information about the supported
// options.
func (d *Driver) SetAutoEncryptionOptions(opts *options.AutoEncryptionOptions) {
	d.clientOpts = append(d.clientOpts, options.Client().SetAutoEncryptionOptions(opts))
}
