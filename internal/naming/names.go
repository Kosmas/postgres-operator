/*
 Copyright 2021 Crunchy Data Solutions, Inc.
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

 http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package naming

import (
	"fmt"
	"hash/fnv"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crunchydata/postgres-operator/pkg/apis/postgres-operator.crunchydata.com/v1beta1"
)

const (
	// ContainerDatabase is the name of the container running PostgreSQL and
	// supporting tools: Patroni, pgBackRest, etc.
	ContainerDatabase = "database"

	// ContainerPGBouncer is the name of a container running PgBouncer.
	ContainerPGBouncer = "pgbouncer"
	// ContainerPGBouncerConfig is the name of a container supporting PgBouncer.
	ContainerPGBouncerConfig = "pgbouncer-config"

	// ContainerPostgresStartup is the name of the initialization container
	// that prepares the filesystem for PostgreSQL.
	ContainerPostgresStartup = "postgres-startup"

	// ContainerClientCertInit is the name of the initialization container that is responsible
	// for copying and setting proper permissions on the client certificate and key
	ContainerClientCertInit = ContainerDatabase + "-client-cert-init"
	// ContainerClientCertCopy is the name of the container that is responsible for copying and
	// setting proper permissions on the client certificate and key after initialization whenever
	// there is a change in the certificates or key
	ContainerClientCertCopy = "replication-cert-copy"
	// ContainerNSSWrapperInit is the name of the init container utilized to configure support
	// for the nss_wrapper
	ContainerNSSWrapperInit = "nss-wrapper-init"

	// ContainerPGMonitorExporter is the name of a container running postgres_exporter
	ContainerPGMonitorExporter = "exporter"
)

const (
	PortExporter   = "exporter"
	PortPGBouncer  = "pgbouncer"
	PortPostgreSQL = "postgres"
)

const (
	// RootCertSecret is the default root certificate secret name
	RootCertSecret = "pgo-root-cacert" /* #nosec */
	// ClusterCertSecret is the default cluster leaf certificate secret name
	ClusterCertSecret = "%s-cluster-cert" /* #nosec */
)

const (
	// CertVolume is the name of the Certificate volume and volume mount in a
	// PostgreSQL instance Pod
	CertVolume = "cert-volume"

	// CertMountPath is the path for mounting the postgrescluster certificates
	// and key
	CertMountPath = "/pgconf/tls"

	// ReplicationDirectory is the directory at CertMountPath where the replication
	// certificates and key are mounted
	ReplicationDirectory = "/replication"

	// ReplicationTmp is the directory where the replication certificates and key can
	// have the proper permissions set due to:
	// https://github.com/kubernetes/kubernetes/issues/57923
	ReplicationTmp = "/tmp/replication"

	// ReplicationCert is the secret key to the postgrescluster's
	// replication/rewind user's client certificate
	ReplicationCert = "tls.crt"

	// ReplicationCertPath is the path to the postgrescluster's replication/rewind
	// user's client certificate
	ReplicationCertPath = "replication/tls.crt"

	// ReplicationPrivateKey is the secret key to the postgrescluster's
	// replication/rewind user's client private key
	ReplicationPrivateKey = "tls.key"

	// ReplicationPrivateKeyPath is the path to the postgrescluster's
	// replication/rewind user's client private key
	ReplicationPrivateKeyPath = "replication/tls.key"

	// ReplicationCACert is the key name of the postgrescluster's replication/rewind
	// user's client CA certificate
	// Note: when using auto-generated certificates, this will be identical to the
	// server CA cert
	ReplicationCACert = "ca.crt"

	// ReplicationCACertPath is the path to the postgrescluster's replication/rewind
	// user's client CA certificate
	ReplicationCACertPath = "replication/ca.crt"
)

const (
	// PGBackRestRepoContainerName is the name assigned to the container used to run pgBackRest and
	// SSH
	PGBackRestRepoContainerName = "pgbackrest"

	// PGBackRestRestoreContainerName is the name assigned to the container used to run pgBackRest
	// restores
	PGBackRestRestoreContainerName = "pgbackrest-restore"

	// PGBackRestRepoName is the name used for a pgbackrest repository
	PGBackRestRepoName = "%s-pgbackrest-repo-%s"

	// PGBackRestSSHVolume is the name the SSH volume used when configuring SSH in a pgBackRest Pod
	PGBackRestSSHVolume = "ssh"

	// suffix used with postgrescluster name for associated configmap.
	// for instance, if the cluster is named 'mycluster', the
	// configmap will be named 'mycluster-pgbackrest-config'
	cmNameSuffix = "%s-pgbackrest-config"

	// suffix used with postgrescluster name for associated configmap.
	// for instance, if the cluster is named 'mycluster', the
	// configmap will be named 'mycluster-ssh-config'
	sshCMNameSuffix = "%s-ssh-config"

	// suffix used with postgrescluster name for associated secret.
	// for instance, if the cluster is named 'mycluster', the
	// secret will be named 'mycluster-ssh'
	sshSecretNameSuffix = "%s-ssh"
)

// AsObjectKey converts the ObjectMeta API type to a client.ObjectKey.
// When you have a client.Object, use client.ObjectKeyFromObject() instead.
func AsObjectKey(m metav1.ObjectMeta) client.ObjectKey {
	return client.ObjectKey{Namespace: m.Namespace, Name: m.Name}
}

// ClusterConfigMap returns the ObjectMeta necessary to lookup
// cluster's shared ConfigMap.
func ClusterConfigMap(cluster *v1beta1.PostgresCluster) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Namespace: cluster.Namespace,
		Name:      cluster.Name + "-config",
	}
}

// ClusterInstanceRBAC returns the ObjectMeta necessary to lookup the
// ServiceAccount, Role, and RoleBinding for cluster's PostgreSQL instances.
func ClusterInstanceRBAC(cluster *v1beta1.PostgresCluster) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Namespace: cluster.Namespace,
		Name:      cluster.Name + "-instance",
	}
}

// ClusterPGBouncer returns the ObjectMeta necessary to lookup the ConfigMap,
// Deployment, Secret, or Service that is cluster's PgBouncer proxy.
func ClusterPGBouncer(cluster *v1beta1.PostgresCluster) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Namespace: cluster.Namespace,
		Name:      cluster.Name + "-pgbouncer",
	}
}

// ClusterPodService returns the ObjectMeta necessary to lookup the Service
// that is responsible for the network identity of Pods.
func ClusterPodService(cluster *v1beta1.PostgresCluster) metav1.ObjectMeta {
	// The hyphen below ensures that the DNS name will not be interpreted as a
	// top-level domain. Partially qualified requests for "{pod}.{cluster}-pods"
	// should not leave the Kubernetes cluster, and if they do they are less
	// likely to resolve.
	return metav1.ObjectMeta{
		Namespace: cluster.Namespace,
		Name:      cluster.Name + "-pods",
	}
}

// ClusterPrimaryService returns the ObjectMeta necessary to lookup the Service
// that exposes the PostgreSQL primary instance.
func ClusterPrimaryService(cluster *v1beta1.PostgresCluster) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Namespace: cluster.Namespace,
		Name:      cluster.Name + "-primary",
	}
}

// GenerateInstance returns a random name for a member of cluster and set.
func GenerateInstance(
	cluster *v1beta1.PostgresCluster, set *v1beta1.PostgresInstanceSetSpec,
) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Namespace: cluster.Namespace,
		Name:      cluster.Name + "-" + set.Name + "-" + rand.String(4),
	}
}

// GenerateStartupInstance returns a stable name that's shaped like
// GenerateInstance above. The stable name is based on a four character
// hash of the cluster name and instance set name
func GenerateStartupInstance(
	cluster *v1beta1.PostgresCluster, set *v1beta1.PostgresInstanceSetSpec,
) metav1.ObjectMeta {
	// Calculate a stable name that's shaped like GenerateInstance above.
	// hash.Hash.Write never returns an error: https://pkg.go.dev/hash#Hash.
	hash := fnv.New32()
	_, _ = hash.Write([]byte(cluster.Name + set.Name))
	suffix := rand.SafeEncodeString(fmt.Sprint(hash.Sum32()))[:4]

	return metav1.ObjectMeta{
		Namespace: cluster.Namespace,
		Name:      cluster.Name + "-" + set.Name + "-" + suffix,
	}
}

// InstanceConfigMap returns the ObjectMeta necessary to lookup
// instance's shared ConfigMap.
func InstanceConfigMap(instance metav1.Object) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Namespace: instance.GetNamespace(),
		Name:      instance.GetName() + "-config",
	}
}

// InstanceCertificates returns the ObjectMeta necessary to lookup the Secret
// containing instance's certificates.
func InstanceCertificates(instance metav1.Object) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Namespace: instance.GetNamespace(),
		Name:      instance.GetName() + "-certs",
	}
}

// InstancePostgresDataVolume returns the ObjectMeta for the PostgreSQL data
// volume for instance.
func InstancePostgresDataVolume(instance *appsv1.StatefulSet) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Namespace: instance.GetNamespace(),
		Name:      instance.GetName() + "-pgdata",
	}
}

// InstancePostgresWALVolume returns the ObjectMeta for the PostgreSQL WAL
// volume for instance.
func InstancePostgresWALVolume(instance *appsv1.StatefulSet) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Namespace: instance.GetNamespace(),
		Name:      instance.GetName() + "-pgwal",
	}
}

// MonitoringUserSecret returns ObjectMeta necessary to lookup the Secret
// containing authentication credentials for monitoring tools.
func MonitoringUserSecret(cluster *v1beta1.PostgresCluster) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Namespace: cluster.Namespace,
		Name:      cluster.Name + "-monitoring",
	}
}

// ReplicationClientCertSecret returns ObjectMeta necessary to lookup the Secret
// containing the Patroni client authentication certificate information.
func ReplicationClientCertSecret(cluster *v1beta1.PostgresCluster) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Namespace: cluster.Namespace,
		Name:      cluster.Name + "-replication-cert",
	}
}

// PatroniDistributedConfiguration returns the ObjectMeta necessary to lookup
// the DCS created by Patroni for cluster. This same name is used for both
// ConfigMap and Endpoints. See Patroni DCS "config_path".
func PatroniDistributedConfiguration(cluster *v1beta1.PostgresCluster) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Namespace: cluster.Namespace,
		Name:      PatroniScope(cluster) + "-config",
	}
}

// PatroniLeaderConfigMap returns the ObjectMeta necessary to lookup the
// ConfigMap created by Patroni for the leader election of cluster.
// See Patroni DCS "leader_path".
func PatroniLeaderConfigMap(cluster *v1beta1.PostgresCluster) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Namespace: cluster.Namespace,
		Name:      PatroniScope(cluster) + "-leader",
	}
}

// PatroniLeaderEndpoints returns the ObjectMeta necessary to lookup the
// Endpoints created by Patroni for the leader election of cluster.
// See Patroni DCS "leader_path".
func PatroniLeaderEndpoints(cluster *v1beta1.PostgresCluster) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Namespace: cluster.Namespace,
		Name:      PatroniScope(cluster),
	}
}

// PatroniScope returns the "scope" Patroni uses for cluster.
func PatroniScope(cluster *v1beta1.PostgresCluster) string {
	return cluster.Name + "-ha"
}

// PatroniTrigger returns the ObjectMeta necessary to lookup the ConfigMap or
// Endpoints Patroni creates for cluster to initiate a controlled change of the
// leader. See Patroni DCS "failover_path".
func PatroniTrigger(cluster *v1beta1.PostgresCluster) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Namespace: cluster.Namespace,
		Name:      PatroniScope(cluster) + "-failover",
	}
}

// PGBackRestConfig returns the ObjectMeta for a pgBackRest ConfigMap
func PGBackRestConfig(cluster *v1beta1.PostgresCluster) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Namespace: cluster.GetNamespace(),
		Name:      fmt.Sprintf(cmNameSuffix, cluster.GetName()),
	}
}

// PGBackRestBackupJob returns the ObjectMeta for the pgBackRest backup Job utilized
// to create replicas using pgBackRest
func PGBackRestBackupJob(cluster *v1beta1.PostgresCluster) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      cluster.GetName() + "-backup-" + rand.String(4),
		Namespace: cluster.GetNamespace(),
	}
}

// PGBackRestCronJob returns the ObjectMeta for a pgBackRest CronJob
func PGBackRestCronJob(cluster *v1beta1.PostgresCluster, backuptype, repoName string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Namespace: cluster.GetNamespace(),
		Name:      cluster.Name + "-pgbackrest-" + repoName + "-" + backuptype,
	}
}

// PGBackRestRestoreJob returns the ObjectMeta for a pgBackRest restore Job
func PGBackRestRestoreJob(cluster *v1beta1.PostgresCluster) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Namespace: cluster.GetNamespace(),
		Name:      cluster.Name + "-pgbackrest-restore",
	}
}

// PGBackRestRBAC returns the ObjectMeta necessary to lookup the ServiceAccount, Role, and
// RoleBinding for pgBackRest Jobs
func PGBackRestRBAC(cluster *v1beta1.PostgresCluster) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Namespace: cluster.Namespace,
		Name:      cluster.Name + "-pgbackrest",
	}
}

// PGBackRestRepoVolume returns the ObjectMeta for a pgBackRest repository volume
func PGBackRestRepoVolume(cluster *v1beta1.PostgresCluster,
	repoName string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      fmt.Sprintf("%s-%s", cluster.GetName(), repoName),
		Namespace: cluster.GetNamespace(),
	}
}

// PGBackRestSSHConfig returns the ObjectMeta for a pgBackRest SSHD ConfigMap
func PGBackRestSSHConfig(cluster *v1beta1.PostgresCluster) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      fmt.Sprintf(sshCMNameSuffix, cluster.GetName()),
		Namespace: cluster.GetNamespace(),
	}
}

// PGBackRestSSHSecret returns the ObjectMeta for a pgBackRest SSHD Secret
func PGBackRestSSHSecret(cluster *v1beta1.PostgresCluster) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      fmt.Sprintf(sshSecretNameSuffix, cluster.GetName()),
		Namespace: cluster.GetNamespace(),
	}
}

// DeprecatedPostgresUserSecret returns the ObjectMeta necessary to lookup the
// old Secret containing the default Postgres user and connection information.
// Use PostgresUserSecret instead.
func DeprecatedPostgresUserSecret(cluster *v1beta1.PostgresCluster) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Namespace: cluster.Namespace,
		Name:      cluster.Name + "-pguser",
	}
}

// PostgresUserSecret returns the ObjectMeta necessary to lookup a Secret
// containing a PostgreSQL user and its connection information.
func PostgresUserSecret(cluster *v1beta1.PostgresCluster, username string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Namespace: cluster.Namespace,
		Name:      cluster.Name + "-pguser-" + username,
	}
}

// PostgresTLSSecret returns the ObjectMeta necessary to lookup the Secret
// containing the default Postgres TLS certificates and key
func PostgresTLSSecret(cluster *v1beta1.PostgresCluster) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Namespace: cluster.Namespace,
		Name:      cluster.Name + "-cluster-cert",
	}
}
