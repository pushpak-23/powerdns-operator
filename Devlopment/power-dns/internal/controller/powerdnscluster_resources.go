package controller

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	platformv1alpha1 "github.com/example/powerdns-platform/api/v1alpha1"
)

const (
	defaultPowerDNSImage       = "ghcr.io/example/powerdns-authoritative:4.8"
	defaultPostgreSQLImage     = "public.ecr.aws/bitnami/postgresql:16"
	defaultPowerDNSAdminImage  = "ghcr.io/example/powerdns-admin:latest"
	defaultPowerDNSPort        = 53
	defaultPowerDNSAPIHealth   = 8081
	defaultPostgreSQLPort      = 5432
	defaultPowerDNSAdminPort   = 8080
	defaultPowerDNSDatabase    = "powerdns"
	defaultPowerDNSUser        = "powerdns"
	defaultPowerDNSAPITokenKey = "api-key"
	defaultDatabasePasswordKey = "postgres-password"
	defaultAdminSecretKeyName  = "secret-key"
)

func (r *PowerDNSClusterReconciler) reconcileWorkloads(ctx context.Context, instance *platformv1alpha1.PowerDNSCluster) error {
	resources := []client.Object{
		r.desiredPostgreSQLService(instance),
		r.desiredPostgreSQLStatefulSet(instance),
		r.desiredPowerDNSService(instance),
		r.desiredPowerDNSDeployment(instance),
	}
	if instance.Spec.Admin.Enabled {
		resources = append(resources, r.desiredPowerDNSAdminService(instance), r.desiredPowerDNSAdminDeployment(instance))
	}

	for _, obj := range resources {
		if err := controllerutil.SetControllerReference(instance, obj, r.Scheme); err != nil {
			return err
		}
		if err := r.applyObject(ctx, obj); err != nil {
			return err
		}
	}

	return nil
}

func (r *PowerDNSClusterReconciler) readyReplicas(ctx context.Context, instance *platformv1alpha1.PowerDNSCluster) int32 {
	ready := int32(0)

	postgres := &appsv1.StatefulSet{}
	if err := r.Get(ctx, types.NamespacedName{Name: postgresName(instance), Namespace: instance.Namespace}, postgres); err == nil {
		if postgres.Status.ReadyReplicas > 0 {
			ready++
		}
	}

	pdns := &appsv1.Deployment{}
	if err := r.Get(ctx, types.NamespacedName{Name: powerDNSName(instance), Namespace: instance.Namespace}, pdns); err == nil {
		ready += pdns.Status.ReadyReplicas
	}

	admin := &appsv1.Deployment{}
	if instance.Spec.Admin.Enabled {
		if err := r.Get(ctx, types.NamespacedName{Name: adminName(instance), Namespace: instance.Namespace}, admin); err == nil {
			ready += admin.Status.ReadyReplicas
		}
	}

	return ready
}

func (r *PowerDNSClusterReconciler) desiredReadyReplicas(instance *platformv1alpha1.PowerDNSCluster) int32 {
	desired := int32(1) + max32(instance.Spec.Replicas, 1)
	if instance.Spec.Admin.Enabled {
		desired += max32(instance.Spec.Admin.Replicas, 1)
	}
	return desired
}

func (r *PowerDNSClusterReconciler) applyObject(ctx context.Context, obj client.Object) error {
	desired := obj.DeepCopyObject().(client.Object)
	key := client.ObjectKeyFromObject(obj)
	current := desired.DeepCopyObject().(client.Object)
	if err := r.Get(ctx, key, current); err != nil {
		if !apierrors.IsNotFound(err) {
			return err
		}
		return r.Create(ctx, desired)
	}

	switch typedDesired := desired.(type) {
	case *corev1.Service:
		typedCurrent := current.(*corev1.Service)
		typedCurrent.Spec = typedDesired.Spec
		typedCurrent.Labels = typedDesired.Labels
		typedCurrent.Annotations = typedDesired.Annotations
		return r.Update(ctx, typedCurrent)
	case *appsv1.Deployment:
		typedCurrent := current.(*appsv1.Deployment)
		typedCurrent.Spec = typedDesired.Spec
		typedCurrent.Labels = typedDesired.Labels
		typedCurrent.Annotations = typedDesired.Annotations
		return r.Update(ctx, typedCurrent)
	case *appsv1.StatefulSet:
		typedCurrent := current.(*appsv1.StatefulSet)
		typedCurrent.Spec = typedDesired.Spec
		typedCurrent.Labels = typedDesired.Labels
		typedCurrent.Annotations = typedDesired.Annotations
		return r.Update(ctx, typedCurrent)
	default:
		return fmt.Errorf("unsupported object type %T", obj)
	}
}

func postgresName(instance *platformv1alpha1.PowerDNSCluster) string {
	return instance.Name + "-postgresql"
}

func powerDNSName(instance *platformv1alpha1.PowerDNSCluster) string {
	return instance.Name + "-powerdns"
}

func adminName(instance *platformv1alpha1.PowerDNSCluster) string {
	return instance.Name + "-powerdns-admin"
}

func baseLabels(instance *platformv1alpha1.PowerDNSCluster, component string) map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":       instance.Name,
		"app.kubernetes.io/instance":    instance.Name,
		"app.kubernetes.io/managed-by":  "powerdns-platform-operator",
		"app.kubernetes.io/component":   component,
		"app.kubernetes.io/part-of":     "powerdns-platform",
	}
}

func (r *PowerDNSClusterReconciler) desiredPostgreSQLService(instance *platformv1alpha1.PowerDNSCluster) *corev1.Service {
	labels := baseLabels(instance, "postgresql")
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      postgresName(instance),
			Namespace: instance.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: "None",
			Selector:  labels,
			Ports: []corev1.ServicePort{{
				Name:       "postgresql",
				Port:       defaultPostgreSQLPort,
				TargetPort: intstr.FromInt(defaultPostgreSQLPort),
				Protocol:   corev1.ProtocolTCP,
			}},
		},
	}
}

func (r *PowerDNSClusterReconciler) desiredPostgreSQLStatefulSet(instance *platformv1alpha1.PowerDNSCluster) *appsv1.StatefulSet {
	labels := baseLabels(instance, "postgresql")
	image := instance.Spec.PostgreSQL.Image
	if image == "" {
		image = defaultPostgreSQLImage
	}
	policy := corev1.PullIfNotPresent
	if instance.Spec.PostgreSQL.ImagePullPolicy != "" {
		policy = corev1.PullPolicy(instance.Spec.PostgreSQL.ImagePullPolicy)
	}
	storage := instance.Spec.PostgreSQL.StorageSize
	if storage == "" {
		storage = "20Gi"
	}
	passwordSecret := instance.Spec.PostgreSQL.CredentialsSecretRef
	if passwordSecret == "" {
		passwordSecret = instance.Name + "-postgres-credentials"
	}

	return &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      postgresName(instance),
			Namespace: instance.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: postgresName(instance),
			Replicas:    int32Ptr(1),
			Selector: &metav1.LabelSelector{MatchLabels: labels},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: labels},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:            "postgresql",
						Image:           image,
						ImagePullPolicy: policy,
						Ports: []corev1.ContainerPort{{
							Name:          "postgresql",
							ContainerPort:  defaultPostgreSQLPort,
							Protocol:       corev1.ProtocolTCP,
						}},
						Env: []corev1.EnvVar{{Name: "ALLOW_EMPTY_PASSWORD", Value: "no"}, {
							Name: "POSTGRESQL_USERNAME", Value: defaultPowerDNSUser,
						}, {
							Name: "POSTGRESQL_DATABASE", Value: defaultPowerDNSDatabase,
						}, {
							Name: "POSTGRESQL_PASSWORD", ValueFrom: secretKeySelector(passwordSecret, defaultDatabasePasswordKey),
						}},
						VolumeMounts: []corev1.VolumeMount{{Name: "data", MountPath: "/bitnami/postgresql"}},
					}},
				},
			},
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{{
				ObjectMeta: metav1.ObjectMeta{Name: "data"},
				Spec: corev1.PersistentVolumeClaimSpec{
					AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
					Resources: corev1.VolumeResourceRequirements{Requests: corev1.ResourceList{
						corev1.ResourceStorage: resource.MustParse(storage),
					}},
				},
			}},
		},
	}
}

func (r *PowerDNSClusterReconciler) desiredPowerDNSService(instance *platformv1alpha1.PowerDNSCluster) *corev1.Service {
	labels := baseLabels(instance, "powerdns")
	apiPort := powerDNSAPIPort(instance)
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      powerDNSName(instance),
			Namespace: instance.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Type:     serviceTypeOrDefault(instance.Spec.Service.Type),
					Ports: []corev1.ServicePort{{Name: "dns-tcp", Port: defaultPowerDNSPort, TargetPort: intstr.FromInt(defaultPowerDNSPort), Protocol: corev1.ProtocolTCP}, {
						Name: "dns-udp", Port: defaultPowerDNSPort, TargetPort: intstr.FromInt(defaultPowerDNSPort), Protocol: corev1.ProtocolUDP}, {
						Name: "api", Port: apiPort, TargetPort: intstr.FromInt(int(apiPort)), Protocol: corev1.ProtocolTCP}},
		},
	}
}

func (r *PowerDNSClusterReconciler) desiredPowerDNSDeployment(instance *platformv1alpha1.PowerDNSCluster) *appsv1.Deployment {
	labels := baseLabels(instance, "powerdns")
	apiPort := powerDNSAPIPort(instance)
	image := instance.Spec.Image
	if image == "" {
		image = defaultPowerDNSImage
	}
	policy := corev1.PullIfNotPresent
	if instance.Spec.ImagePullPolicy != "" {
		policy = corev1.PullPolicy(instance.Spec.ImagePullPolicy)
	}
	postgresSecret := instance.Spec.PostgreSQL.CredentialsSecretRef
	if postgresSecret == "" {
		postgresSecret = instance.Name + "-postgres-credentials"
	}
	apiSecret := instance.Spec.PowerDNS.APIKeySecretRef
	if apiSecret == "" {
		apiSecret = instance.Name + "-powerdns-api"
	}

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: powerDNSName(instance), Namespace: instance.Namespace, Labels: labels},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(max32(instance.Spec.Replicas, 1)),
			Selector: &metav1.LabelSelector{MatchLabels: labels},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: labels},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:            "powerdns",
						Image:           image,
						ImagePullPolicy: policy,
						Ports: []corev1.ContainerPort{{Name: "dns-tcp", ContainerPort: defaultPowerDNSPort, Protocol: corev1.ProtocolTCP}, {Name: "dns-udp", ContainerPort: defaultPowerDNSPort, Protocol: corev1.ProtocolUDP}, {Name: "api", ContainerPort: apiPort, Protocol: corev1.ProtocolTCP}},
						Env: []corev1.EnvVar{{Name: "PDNS_GPGSQL_HOST", Value: postgresName(instance)}, {Name: "PDNS_GPGSQL_DBNAME", Value: defaultPowerDNSDatabase}, {Name: "PDNS_GPGSQL_USER", Value: defaultPowerDNSUser}, {Name: "PDNS_GPGSQL_PASSWORD", ValueFrom: secretKeySelector(postgresSecret, defaultDatabasePasswordKey)}, {Name: "PDNS_API_KEY", ValueFrom: secretKeySelector(apiSecret, defaultPowerDNSAPITokenKey)}, {Name: "PDNS_WEBSERVER_PORT", Value: fmt.Sprintf("%d", apiPort)}},
					}},
				},
			},
		},
	}
}

func (r *PowerDNSClusterReconciler) desiredPowerDNSAdminService(instance *platformv1alpha1.PowerDNSCluster) *corev1.Service {
	labels := baseLabels(instance, "powerdns-admin")
	adminPort := powerDNSAdminPort(instance)
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{Name: adminName(instance), Namespace: instance.Namespace, Labels: labels},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Type:     serviceTypeOrDefault("ClusterIP"),
			Ports: []corev1.ServicePort{{Name: "http", Port: adminPort, TargetPort: intstr.FromInt(int(adminPort)), Protocol: corev1.ProtocolTCP}},
		},
	}
}

func (r *PowerDNSClusterReconciler) desiredPowerDNSAdminDeployment(instance *platformv1alpha1.PowerDNSCluster) *appsv1.Deployment {
	labels := baseLabels(instance, "powerdns-admin")
	adminPort := powerDNSAdminPort(instance)
	if !instance.Spec.Admin.Enabled {
		return &appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Name: adminName(instance), Namespace: instance.Namespace}}
	}
	image := instance.Spec.Admin.Image
	if image == "" {
		image = defaultPowerDNSAdminImage
	}
	policy := corev1.PullIfNotPresent
	if instance.Spec.Admin.ImagePullPolicy != "" {
		policy = corev1.PullPolicy(instance.Spec.Admin.ImagePullPolicy)
	}
	postgresSecret := instance.Spec.PostgreSQL.CredentialsSecretRef
	if postgresSecret == "" {
		postgresSecret = instance.Name + "-postgres-credentials"
	}
	apiSecret := instance.Spec.PowerDNS.APIKeySecretRef
	if apiSecret == "" {
		apiSecret = instance.Name + "-powerdns-api"
	}
	databaseURISecret := instance.Spec.Admin.DatabaseURISecretRef
	if databaseURISecret == "" {
		databaseURISecret = instance.Name + "-powerdns-admin-db"
	}
	adminSecret := instance.Spec.Admin.SecretKeySecretRef
	if adminSecret == "" {
		adminSecret = instance.Name + "-powerdns-admin"
	}
	serviceHost := instance.Spec.Admin.APIURL
	if serviceHost == "" {
		serviceHost = fmt.Sprintf("http://%s.%s.svc.cluster.local:%d", powerDNSName(instance), instance.Namespace, powerDNSAPIPort(instance))
	}

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: adminName(instance), Namespace: instance.Namespace, Labels: labels},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(max32(instance.Spec.Admin.Replicas, 1)),
			Selector: &metav1.LabelSelector{MatchLabels: labels},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: labels},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:            "powerdns-admin",
						Image:           image,
						ImagePullPolicy: policy,
						Ports: []corev1.ContainerPort{{Name: "http", ContainerPort: adminPort, Protocol: corev1.ProtocolTCP}},
						Env: []corev1.EnvVar{{Name: "SECRET_KEY", ValueFrom: secretKeySelector(adminSecret, defaultAdminSecretKeyName)}, {Name: "PDNS_API_URL", Value: serviceHost}, {Name: "PDNS_API_TOKEN", ValueFrom: secretKeySelector(apiSecret, defaultPowerDNSAPITokenKey)}, {Name: "SQLALCHEMY_DATABASE_URI", ValueFrom: secretKeySelector(databaseURISecret, "uri")}},
					}},
				},
			},
		},
	}
}

func int32Ptr(value int32) *int32 { return &value }

func max32(a, b int32) int32 {
	if a > b { return a }
	return b
}

func serviceTypeOrDefault(value string) corev1.ServiceType {
	if value == "" { return corev1.ServiceTypeClusterIP }
	return corev1.ServiceType(value)
}

func powerDNSAPIPort(instance *platformv1alpha1.PowerDNSCluster) int32 {
	if instance.Spec.PowerDNS.WebServerPort > 0 {
		return instance.Spec.PowerDNS.WebServerPort
	}
	return defaultPowerDNSAPIHealth
}

func powerDNSAdminPort(instance *platformv1alpha1.PowerDNSCluster) int32 {
	if instance.Spec.Admin.ServicePort > 0 {
		return instance.Spec.Admin.ServicePort
	}
	return 80
}

func secretKeySelector(secretName, key string) *corev1.EnvVarSource {
	return &corev1.EnvVarSource{SecretKeyRef: &corev1.SecretKeySelector{LocalObjectReference: corev1.LocalObjectReference{Name: secretName}, Key: key}}
}
