package controller

import (
	"path"
	"regexp"

	csidriverv1alpha1 "github.com/openshift/csi-operator2/pkg/apis/csidriver/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// generateServiceAccount prepares a ServiceAccount that will be used by all pods (controller + daemon set) with
// CSI drivers and its sidecar containers.
func (h *Handler) generateServiceAccount(cr *csidriverv1alpha1.CSIDriverDeployment) *v1.ServiceAccount {
	scName := cr.Name

	sc := &v1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: cr.Namespace,
			Name:      scName,
		},
	}
	h.addOwnerLabels(&sc.ObjectMeta, cr)
	h.addOwner(&sc.ObjectMeta, cr)

	return sc
}

// generateClusterRoleBinding prepares a ClusterRoleBinding that gives a ServiceAccount privileges needed by
// sidecar containers.
func (h *Handler) generateClusterRoleBinding(cr *csidriverv1alpha1.CSIDriverDeployment, serviceAccount *v1.ServiceAccount) *rbacv1.ClusterRoleBinding {
	crbName := h.uniqueGlobalName(cr)
	crb := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: crbName,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      serviceAccount.Name,
				Namespace: serviceAccount.Namespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     h.config.ClusterRoleName,
		},
	}
	h.addOwnerLabels(&crb.ObjectMeta, cr)
	h.addOwner(&crb.ObjectMeta, cr)
	return crb
}

// generateLeaderElectionRoleBinding prepares a RoleBinding that gives a ServiceAccount privileges needed by
// attacher and provisioner leader election.
func (h *Handler) generateLeaderElectionRoleBinding(cr *csidriverv1alpha1.CSIDriverDeployment, serviceAccount *v1.ServiceAccount) *rbacv1.RoleBinding {
	rbName := "leader-election-" + cr.Name
	rb := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: cr.Namespace,
			Name:      rbName,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      serviceAccount.Name,
				Namespace: serviceAccount.Namespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     h.config.LeaderElectionClusterRoleName,
		},
	}
	h.addOwnerLabels(&rb.ObjectMeta, cr)
	h.addOwner(&rb.ObjectMeta, cr)
	return rb
}

// generateDaemonSet prepares a DaemonSet with CSI driver and driver registrar sidecar containers.
func (h *Handler) generateDaemonSet(cr *csidriverv1alpha1.CSIDriverDeployment, serviceAccount *v1.ServiceAccount) *appsv1.DaemonSet {
	dsName := cr.Name + "-node"

	labels := map[string]string{
		"csidriver.storage.okd.io/daemonset": dsName,
	}

	// Prepare DS.Spec.PodSpec
	podSpec := cr.Spec.DriverPerNodeTemplate.DeepCopy()
	if podSpec.Labels == nil {
		podSpec.Labels = labels
	} else {
		for k, v := range labels {
			podSpec.Labels[k] = v
		}
	}

	// Don't overwrite user's ServiceAccount
	if podSpec.Spec.ServiceAccountName == "" {
		podSpec.Spec.ServiceAccountName = serviceAccount.Name
	}

	// Path to the CSI driver socket in the driver container
	csiDriverSocketPath := cr.Spec.DriverSocket
	csiDriverSocketFileName := path.Base(csiDriverSocketPath)
	csiDriverSocketDirectory := path.Dir(csiDriverSocketPath)

	// Path to the CSI driver socket in the driver registrar container
	registrarSocketDirectory := "/csi"
	registrarSocketPath := path.Join(registrarSocketDirectory, csiDriverSocketFileName)

	// Path to the CSI driver socket from kubelet point of view
	kubeletSocketDirectory := path.Join(h.config.KubeletRootDir, "plugins", sanitizeDriverName(cr.Spec.DriverName))
	kubeletSocketPath := path.Join(registrarSocketDirectory, csiDriverSocketFileName)

	// Path to the kubelet dynamic registration directory
	kubeletRegistrationDirectory := path.Join(h.config.KubeletRootDir, "plugins")

	bTrue := true
	// Add CSI Registrar sidecar
	registrarImage := *h.config.DefaultImages.DriverRegistrarImage
	if cr.Spec.ContainerImages != nil && cr.Spec.ContainerImages.DriverRegistrarImage != nil {
		registrarImage = *cr.Spec.ContainerImages.DriverRegistrarImage
	}
	registrar := v1.Container{
		Name:  "csi-driver-registrar",
		Image: registrarImage,
		Args: []string{
			"--v=5",
			"--csi-address=$(ADDRESS)",
			// TODO: enable when 1.12 is rebased
			// "--kubelet-registration-path=$(DRIVER_REG_SOCK_PATH)",
		},
		SecurityContext: &v1.SecurityContext{
			Privileged: &bTrue,
		},
		Env: []v1.EnvVar{
			{
				Name:  "ADDRESS",
				Value: registrarSocketPath,
			},
			{
				Name:  "DRIVER_REG_SOCK_PATH",
				Value: kubeletSocketPath,
			},
			{
				Name: "KUBE_NODE_NAME",
				ValueFrom: &v1.EnvVarSource{
					FieldRef: &v1.ObjectFieldSelector{
						FieldPath: "spec.nodeName",
					},
				},
			},
		},
		VolumeMounts: []v1.VolumeMount{
			{
				Name:      "csi-driver",
				MountPath: registrarSocketDirectory,
			},
			{
				Name:      "registration-dir",
				MountPath: "/registration",
			},
		},
	}
	podSpec.Spec.Containers = append(podSpec.Spec.Containers, registrar)

	// Add volumes
	typeDir := v1.HostPathDirectory
	typeDirOrCreate := v1.HostPathDirectoryOrCreate
	volumes := []v1.Volume{
		{
			Name: "registration-dir",
			VolumeSource: v1.VolumeSource{
				HostPath: &v1.HostPathVolumeSource{
					Path: kubeletRegistrationDirectory,
					Type: &typeDir,
				},
			},
		},
		{
			Name: "csi-driver",
			VolumeSource: v1.VolumeSource{
				HostPath: &v1.HostPathVolumeSource{
					Path: kubeletSocketDirectory,
					Type: &typeDirOrCreate,
				},
			},
		},
	}
	podSpec.Spec.Volumes = append(podSpec.Spec.Volumes, volumes...)

	// Patch the driver container with the volume for CSI driver socket
	volumeMount := v1.VolumeMount{
		Name:      "csi-driver",
		MountPath: csiDriverSocketDirectory,
	}
	driverContainer := &podSpec.Spec.Containers[0]
	driverContainer.VolumeMounts = append(driverContainer.VolumeMounts, volumeMount)

	// Create the DaemonSet
	updateStrategy := appsv1.OnDeleteDaemonSetStrategyType
	if cr.Spec.NodeUpdateStrategy == csidriverv1alpha1.CSIDeploymentUpdateStrategyRolling {
		updateStrategy = appsv1.RollingUpdateDaemonSetStrategyType
	}
	ds := &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: cr.Namespace,
			Name:      dsName,
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: *podSpec,
			UpdateStrategy: appsv1.DaemonSetUpdateStrategy{
				Type: updateStrategy,
			},
		},
	}
	h.addOwnerLabels(&ds.ObjectMeta, cr)
	h.addOwner(&ds.ObjectMeta, cr)

	return ds
}

// generateDeployment prepares a Deployment with CSI driver and attacher and provisioner sidecar containers.
func (h *Handler) generateDeployment(cr *csidriverv1alpha1.CSIDriverDeployment, serviceAccount *v1.ServiceAccount) *appsv1.Deployment {
	dName := cr.Name + "-controller"

	labels := map[string]string{
		"csidriver.storage.okd.io/deployment": dName,
	}

	// Prepare the pod template
	podSpec := cr.Spec.DriverControllerTemplate.DeepCopy()
	if podSpec.Labels == nil {
		podSpec.Labels = labels
	} else {
		for k, v := range labels {
			podSpec.Labels[k] = v
		}
	}

	if podSpec.Spec.ServiceAccountName == "" {
		podSpec.Spec.ServiceAccountName = serviceAccount.Name
	}

	// Add sidecars

	// Path to the CSI driver socket in the driver container
	csiDriverSocketPath := cr.Spec.DriverSocket
	csiDriverSocketFileName := path.Base(csiDriverSocketPath)
	csiDriverSocketDirectory := path.Dir(csiDriverSocketPath)

	// Path to the CSI driver socket in the sidecar containers
	sidecarSocketDirectory := "/csi"
	sidecarSocketPath := path.Join(sidecarSocketDirectory, csiDriverSocketFileName)

	provisionerImage := *h.config.DefaultImages.ProvisionerImage
	if cr.Spec.ContainerImages != nil && cr.Spec.ContainerImages.ProvisionerImage != nil {
		provisionerImage = *cr.Spec.ContainerImages.ProvisionerImage
	}
	provisioner := v1.Container{
		Name:  "csi-provisioner",
		Image: provisionerImage,
		Args: []string{
			"--v=5",
			"--csi-address=$(ADDRESS)",
			"--provisioner=" + cr.Spec.DriverName,
			// TODO: add leader election parameters
		},
		Env: []v1.EnvVar{
			{
				Name:  "ADDRESS",
				Value: sidecarSocketPath,
			},
		},
		VolumeMounts: []v1.VolumeMount{
			{
				Name:      "csi-driver",
				MountPath: "/csi",
			},
		},
	}
	podSpec.Spec.Containers = append(podSpec.Spec.Containers, provisioner)

	attacherImage := *h.config.DefaultImages.AttacherImage
	if cr.Spec.ContainerImages != nil && cr.Spec.ContainerImages.AttacherImage != nil {
		attacherImage = *cr.Spec.ContainerImages.AttacherImage
	}
	attacher := v1.Container{
		Name:  "csi-attacher",
		Image: attacherImage,
		Args: []string{
			"--v=5",
			"--csi-address=$(ADDRESS)",
			// TODO: add leader election parameters
		},
		Env: []v1.EnvVar{
			{
				Name:  "ADDRESS",
				Value: sidecarSocketPath,
			},
		},
		VolumeMounts: []v1.VolumeMount{
			{
				Name:      "csi-driver",
				MountPath: "/csi",
			},
		},
	}
	podSpec.Spec.Containers = append(podSpec.Spec.Containers, attacher)

	// Add volumes
	volumes := []v1.Volume{
		{
			Name: "csi-driver",
			VolumeSource: v1.VolumeSource{
				EmptyDir: &v1.EmptyDirVolumeSource{},
			},
		},
	}
	podSpec.Spec.Volumes = append(podSpec.Spec.Volumes, volumes...)

	// Set selector to infra nodes only
	if podSpec.Spec.NodeSelector == nil {
		podSpec.Spec.NodeSelector = h.config.InfrastructureNodeSelector
	}

	// Patch the driver container with the volume for CSI driver socket
	volumeMount := v1.VolumeMount{
		Name:      "csi-driver",
		MountPath: csiDriverSocketDirectory,
	}
	driverContainer := &podSpec.Spec.Containers[0]
	driverContainer.VolumeMounts = append(driverContainer.VolumeMounts, volumeMount)

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: cr.Namespace,
			Name:      dName,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: *podSpec,
			Replicas: &h.config.DeploymentReplicas,
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RollingUpdateDeploymentStrategyType,
			},
		},
	}
	h.addOwnerLabels(&deployment.ObjectMeta, cr)
	h.addOwner(&deployment.ObjectMeta, cr)

	return deployment
}

// generateStorageClass prepares a StorageClass from given template
func (h *Handler) generateStorageClass(cr *csidriverv1alpha1.CSIDriverDeployment, template *csidriverv1alpha1.StorageClassTemplate) *storagev1.StorageClass {
	expectedSC := &storagev1.StorageClass{
		// ObjectMeta will be filled below
		Provisioner:          cr.Spec.DriverName,
		Parameters:           template.Parameters,
		ReclaimPolicy:        template.ReclaimPolicy,
		MountOptions:         template.MountOptions,
		AllowVolumeExpansion: template.AllowVolumeExpansion,
		VolumeBindingMode:    template.VolumeBindingMode,
		AllowedTopologies:    template.AllowedTopologies,
	}
	template.ObjectMeta.DeepCopyInto(&expectedSC.ObjectMeta)
	h.addOwnerLabels(&expectedSC.ObjectMeta, cr)
	h.addOwner(&expectedSC.ObjectMeta, cr)
	if template.Default != nil && *template.Default == true {
		expectedSC.Annotations = map[string]string{
			"storageclass.kubernetes.io/is-default-class": "true",
		}
	}
	return expectedSC
}

// sanitizeDriverName sanitizes CSI driver name to be usable as a directory name. All dangerous characters are replaced
// by '-'.
func sanitizeDriverName(driver string) string {
	re := regexp.MustCompile("[^a-zA-Z0-9-.]")
	name := re.ReplaceAllString(driver, "-")
	return name
}

// a CSIDriverDeployment (as OwnerReference does not work there) and may be used to limit Watch() in future.
func (h *Handler) addOwnerLabels(meta *metav1.ObjectMeta, cr *csidriverv1alpha1.CSIDriverDeployment) bool {
	changed := false
	if meta.Labels == nil {
		meta.Labels = map[string]string{}
		changed = true
	}
	if v, exists := meta.Labels["csidriver.storage.okd.io/owner-namespace"]; !exists || v != cr.Namespace {
		meta.Labels["csidriver.storage.okd.io/owner-namespace"] = cr.Namespace
		changed = true
	}
	if v, exists := meta.Labels["csidriver.storage.okd.io/owner-name"]; !exists || v != cr.Name {
		meta.Labels["csidriver.storage.okd.io/owner-name"] = cr.Name
		changed = true
	}

	return changed
}

func (h *Handler) addOwner(meta *metav1.ObjectMeta, cr *csidriverv1alpha1.CSIDriverDeployment) {
	bTrue := true
	meta.OwnerReferences = []metav1.OwnerReference{
		{
			APIVersion: "csidriver.storage.okd.io/v1alpha1",
			Kind:       "CSIDriverDeployment",
			Name:       cr.Name,
			UID:        cr.UID,
			Controller: &bTrue,
		},
	}
}

func (h *Handler) uniqueGlobalName(i *csidriverv1alpha1.CSIDriverDeployment) string {
	return "csidriverdeployment-" + string(i.UID)
}