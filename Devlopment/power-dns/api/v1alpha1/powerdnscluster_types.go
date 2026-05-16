package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	PowerDNSClusterPhaseReady    = "Ready"
	PowerDNSClusterPhaseDegraded = "Degraded"
	PowerDNSClusterPhaseReconciling = "Reconciling"
)

// PowerDNSCluster defines an enterprise PowerDNS authoritative deployment plus PostgreSQL backend policy.
type PowerDNSCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PowerDNSClusterSpec   `json:"spec,omitempty"`
	Status PowerDNSClusterStatus `json:"status,omitempty"`
}

type PowerDNSClusterSpec struct {
	Replicas int32 `json:"replicas,omitempty"`
	Image string `json:"image,omitempty"`
	ImagePullPolicy string `json:"imagePullPolicy,omitempty"`
	Service PowerDNSServiceSpec `json:"service,omitempty"`
	PowerDNS PowerDNSSpec `json:"powerdns,omitempty"`
	PostgreSQL PostgreSQLSpec `json:"postgresql,omitempty"`
	Admin PowerDNSAdminSpec `json:"admin,omitempty"`
	DNSSEC DNSSECSpec `json:"dnssec,omitempty"`
	Backups BackupSpec `json:"backups,omitempty"`
	Monitoring MonitoringSpec `json:"monitoring,omitempty"`
	Security SecuritySpec `json:"security,omitempty"`
	Scheduling SchedulingSpec `json:"scheduling,omitempty"`
	Scaling ScalingSpec `json:"scaling,omitempty"`
	OpenStack OpenStackSpec `json:"openstack,omitempty"`
}

type PowerDNSServiceSpec struct {
	Type string `json:"type,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	Labels map[string]string `json:"labels,omitempty"`
}

type PowerDNSSpec struct {
	APIEnabled bool `json:"apiEnabled,omitempty"`
	APIKeySecretRef string `json:"apiKeySecretRef,omitempty"`
	WebServerAddress string `json:"webServerAddress,omitempty"`
	WebServerPort int32 `json:"webServerPort,omitempty"`
	AllowAXFR bool `json:"allowAxfr,omitempty"`
	AllowIXFR bool `json:"allowIxfr,omitempty"`
	EnableDNSSEC bool `json:"enableDnssec,omitempty"`
	ZoneTransferACL []string `json:"zoneTransferAcl,omitempty"`
	Recursor string `json:"recursor,omitempty"`
}

type PostgreSQLSpec struct {
	Mode string `json:"mode,omitempty"`
	Image string `json:"image,omitempty"`
	ImagePullPolicy string `json:"imagePullPolicy,omitempty"`
	ConnectionPooler bool `json:"connectionPooler,omitempty"`
	Provider string `json:"provider,omitempty"`
	StorageClass string `json:"storageClass,omitempty"`
	StorageSize string `json:"storageSize,omitempty"`
	CredentialsSecretRef string `json:"credentialsSecretRef,omitempty"`
	SyncReplication bool `json:"syncReplication,omitempty"`
	WALArchive bool `json:"walArchive,omitempty"`
	PITR bool `json:"pitr,omitempty"`
	BackupSchedule string `json:"backupSchedule,omitempty"`
	ReadReplicas int32 `json:"readReplicas,omitempty"`
	TLS bool `json:"tls,omitempty"`
}

type DNSSECSpec struct {
	Enabled bool `json:"enabled,omitempty"`
	KSKSecretRef string `json:"kskSecretRef,omitempty"`
	ZSKSecretRef string `json:"zskSecretRef,omitempty"`
	Algorithm string `json:"algorithm,omitempty"`
	KeyRollInterval string `json:"keyRollInterval,omitempty"`
}

type PowerDNSAdminSpec struct {
	Enabled bool `json:"enabled,omitempty"`
	Image string `json:"image,omitempty"`
	ImagePullPolicy string `json:"imagePullPolicy,omitempty"`
	Replicas int32 `json:"replicas,omitempty"`
	ServicePort int32 `json:"servicePort,omitempty"`
	IngressHost string `json:"ingressHost,omitempty"`
	SecretKeySecretRef string `json:"secretKeySecretRef,omitempty"`
	APIURL string `json:"apiUrl,omitempty"`
	APITokenSecretRef string `json:"apiTokenSecretRef,omitempty"`
}

type BackupSpec struct {
	Enabled bool `json:"enabled,omitempty"`
	Schedule string `json:"schedule,omitempty"`
	Retention string `json:"retention,omitempty"`
	Target string `json:"target,omitempty"`
	Verify bool `json:"verify,omitempty"`
}

type MonitoringSpec struct {
	ServiceMonitor bool `json:"serviceMonitor,omitempty"`
	PrometheusRules bool `json:"prometheusRules,omitempty"`
	GrafanaDashboards bool `json:"grafanaDashboards,omitempty"`
	OpenTelemetry bool `json:"openTelemetry,omitempty"`
}

type SecuritySpec struct {
	TLS bool `json:"tls,omitempty"`
	mTLS bool `json:"mtls,omitempty"`
	ReadOnlyRootFilesystem bool `json:"readOnlyRootFilesystem,omitempty"`
	RunAsNonRoot bool `json:"runAsNonRoot,omitempty"`
	PodSecurityStandard string `json:"podSecurityStandard,omitempty"`
	NetworkPolicies bool `json:"networkPolicies,omitempty"`
	SecretEncryption bool `json:"secretEncryption,omitempty"`
}

type SchedulingSpec struct {
	AntiAffinity bool `json:"antiAffinity,omitempty"`
	TopologySpreadConstraints bool `json:"topologySpreadConstraints,omitempty"`
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	Tolerations []corev1Toleration `json:"tolerations,omitempty"`
	Zones []string `json:"zones,omitempty"`
}

type corev1Toleration struct {
	Key string `json:"key,omitempty"`
	Operator string `json:"operator,omitempty"`
	Value string `json:"value,omitempty"`
	Effect string `json:"effect,omitempty"`
}

type ScalingSpec struct {
	HorizontalAutoscaler bool `json:"horizontalAutoscaler,omitempty"`
	VerticalAutoscaler bool `json:"verticalAutoscaler,omitempty"`
	MinReplicas int32 `json:"minReplicas,omitempty"`
	MaxReplicas int32 `json:"maxReplicas,omitempty"`
}

type OpenStackSpec struct {
	DesignateEnabled bool `json:"designateEnabled,omitempty"`
	KeystoneAuth bool `json:"keystoneAuth,omitempty"`
	TenantIsolation bool `json:"tenantIsolation,omitempty"`
	ReverseDNS bool `json:"reverseDns,omitempty"`
	MultiRegion bool `json:"multiRegion,omitempty"`
}

type PowerDNSClusterStatus struct {
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	Phase string `json:"phase,omitempty"`
	ReadyReplicas int32 `json:"readyReplicas,omitempty"`
	Conditions []metav1.Condition `json:"conditions,omitempty"`
	Endpoints PowerDNSEndpointsStatus `json:"endpoints,omitempty"`
	LastBackupTime *metav1.Time `json:"lastBackupTime,omitempty"`
}

type PowerDNSEndpointsStatus struct {
	API string `json:"api,omitempty"`
	DNS string `json:"dns,omitempty"`
	PostgreSQL string `json:"postgresql,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=powerdnsclusters,scope=Namespaced,shortName=pdns
type PowerDNSClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta  `json:"metadata,omitempty"`
	Items           []PowerDNSCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PowerDNSCluster{}, &PowerDNSClusterList{})
}

func (in *PowerDNSCluster) DeepCopyObject() runtime.Object {
	out := new(PowerDNSCluster)
	*out = *in
	out.ObjectMeta = *in.ObjectMeta.DeepCopy()
	return out
}

func (in *PowerDNSClusterList) DeepCopyObject() runtime.Object {
	out := new(PowerDNSClusterList)
	*out = *in
	return out
}

func (in *PowerDNSCluster) GetObjectKind() schema.ObjectKind {
	return &in.TypeMeta
}

func (in *PowerDNSClusterList) GetObjectKind() schema.ObjectKind {
	return &in.TypeMeta
}