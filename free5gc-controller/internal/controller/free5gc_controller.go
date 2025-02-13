/*
Copyright 2025.

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

package controller

import (
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	k8sstartkubernetescomv1 "k8s.startkubernetes.com/api/v1"
)

// Free5gcReconciler reconciles a Free5gc object
type Free5gcReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// Convertisseur pour pointer vers int32
func int32Ptr(i int) *int32 {
	val := int32(i)
	return &val
}

// +kubebuilder:rbac:groups=k8s.startkubernetes.com.my.domain,resources=free5gcs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=k8s.startkubernetes.com.my.domain,resources=free5gcs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=k8s.startkubernetes.com.my.domain,resources=free5gcs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Free5gc object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.20.0/pkg/reconcile
func (r *Free5gcReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	free5gc := &k8sstartkubernetescomv1.Free5gc{}
	err := r.Get(ctx, req.NamespacedName, free5gc)
	if err != nil {
		if errors.IsNotFound(err) {
			logger.Info("Free5GC resource not found. Ignoring since object might be deleted")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Failed to get Free5GC instance")
		return ctrl.Result{}, err
	}

	// Déploiement des services
	if err := r.deployAMF(ctx, free5gc); err != nil {
		logger.Error(err, "Failed to deploy AMF")
		return ctrl.Result{}, err
	}

	if err := r.deploySMF(ctx, free5gc); err != nil {
		logger.Error(err, "Failed to deploy SMF")
		return ctrl.Result{}, err
	}

	if err := r.deployAUSF(ctx, free5gc); err != nil {
		logger.Error(err, "Failed to deploy AUSF")
		return ctrl.Result{}, err
	}

	if err := r.deployPCF(ctx, free5gc); err != nil {
		logger.Error(err, "Failed to deploy PCF")
		return ctrl.Result{}, err
	}

	if err := r.deployNRF(ctx, free5gc); err != nil {
		logger.Error(err, "Failed to deploy NRF")
		return ctrl.Result{}, err
	}

	if err := r.deployUPF(ctx, free5gc); err != nil {
		logger.Error(err, "Failed to deploy UPF")
		return ctrl.Result{}, err
	}

	if err := r.deployUDM(ctx, free5gc); err != nil {
		logger.Error(err, "Failed to deploy UDM")
		return ctrl.Result{}, err
	}

	if err := r.deployUDR(ctx, free5gc); err != nil {
		logger.Error(err, "Failed to deploy UDR")
		return ctrl.Result{}, err
	}

	if err := r.deployNSSF(ctx, free5gc); err != nil {
		logger.Error(err, "Failed to deploy NSSF")
		return ctrl.Result{}, err
	}

	logger.Info("Reconciliation successful")

	return ctrl.Result{}, nil
}

// Déploiement du service AMF
func (r *Free5gcReconciler) deployAMF(ctx context.Context, free5gc *k8sstartkubernetescomv1.Free5gc) error {
	logger := log.FromContext(ctx)

	deploymentName := fmt.Sprintf("%s-amf", free5gc.Name)
	foundDeployment := &appsv1.Deployment{}
	err := r.Get(ctx, client.ObjectKey{Namespace: free5gc.Namespace, Name: deploymentName}, foundDeployment)
	if err != nil && errors.IsNotFound(err) {
		logger.Info("Creating new AMF deployment", "DeploymentName", deploymentName)

		AMFDeployment := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      deploymentName,
				Namespace: free5gc.Namespace,
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: int32Ptr(free5gc.Spec.AMF.Replicas),
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"app": "amf"},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{"app": "amf"},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "amf",
								Image: free5gc.Spec.AMF.Image,
								Ports: []corev1.ContainerPort{
									{
										ContainerPort: int32(free5gc.Spec.AMF.Network.Port),
									},
								},
							},
						},
					},
				},
			},
		}

		if err := controllerutil.SetControllerReference(free5gc, AMFDeployment, r.Scheme); err != nil {
			return err
		}
		return r.Create(ctx, AMFDeployment)
	} else if err != nil {
		return err
	}

	logger.Info("AMF Deployment already exists, skipping creation")
	return nil
}

// Déploiement du service SMF
func (r *Free5gcReconciler) deploySMF(ctx context.Context, free5gc *k8sstartkubernetescomv1.Free5gc) error {
	logger := log.FromContext(ctx)

	deploymentName := fmt.Sprintf("%s-smf", free5gc.Name)
	foundDeployment := &appsv1.Deployment{}
	err := r.Get(ctx, client.ObjectKey{Namespace: free5gc.Namespace, Name: deploymentName}, foundDeployment)
	if err != nil && errors.IsNotFound(err) {
		logger.Info("Creating new SMF deployment", "DeploymentName", deploymentName)

		SMFDeployment := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      deploymentName,
				Namespace: free5gc.Namespace,
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: int32Ptr(free5gc.Spec.SMF.Replicas),
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"app": "smf"},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{"app": "smf"},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "smf",
								Image: free5gc.Spec.SMF.Image,
								Ports: []corev1.ContainerPort{
									{
										ContainerPort: int32(free5gc.Spec.SMF.Network.Port),
									},
								},
							},
						},
					},
				},
			},
		}

		if err := controllerutil.SetControllerReference(free5gc, SMFDeployment, r.Scheme); err != nil {
			return err
		}
		return r.Create(ctx, SMFDeployment)
	} else if err != nil {
		return err
	}

	logger.Info("SMF Deployment already exists, skipping creation")
	return nil
}

// Déploiement du service AUSF
func (r *Free5gcReconciler) deployAUSF(ctx context.Context, free5gc *k8sstartkubernetescomv1.Free5gc) error {
	logger := log.FromContext(ctx)

	deploymentName := fmt.Sprintf("%s-ausf", free5gc.Name)
	foundDeployment := &appsv1.Deployment{}
	err := r.Get(ctx, client.ObjectKey{Namespace: free5gc.Namespace, Name: deploymentName}, foundDeployment)
	if err != nil && errors.IsNotFound(err) {
		logger.Info("Creating new AUSF deployment", "DeploymentName", deploymentName)

		AUSFDeployment := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      deploymentName,
				Namespace: free5gc.Namespace,
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: int32Ptr(free5gc.Spec.AUSF.Replicas),
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"app": "ausf"},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{"app": "ausf"},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "ausf",
								Image: free5gc.Spec.AUSF.Image,
								Ports: []corev1.ContainerPort{
									{
										ContainerPort: int32(free5gc.Spec.AUSF.Network.Port),
									},
								},
							},
						},
					},
				},
			},
		}

		if err := controllerutil.SetControllerReference(free5gc, AUSFDeployment, r.Scheme); err != nil {
			return err
		}
		return r.Create(ctx, AUSFDeployment)
	} else if err != nil {
		return err
	}

	logger.Info("AUSF Deployment already exists, skipping creation")
	return nil
}

// Déploiement du service PCF
func (r *Free5gcReconciler) deployPCF(ctx context.Context, free5gc *k8sstartkubernetescomv1.Free5gc) error {
	logger := log.FromContext(ctx)

	deploymentName := fmt.Sprintf("%s-pcf", free5gc.Name)
	foundDeployment := &appsv1.Deployment{}
	err := r.Get(ctx, client.ObjectKey{Namespace: free5gc.Namespace, Name: deploymentName}, foundDeployment)
	if err != nil && errors.IsNotFound(err) {
		logger.Info("Creating new PCF deployment", "DeploymentName", deploymentName)

		PCFDeployment := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      deploymentName,
				Namespace: free5gc.Namespace,
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: int32Ptr(free5gc.Spec.PCF.Replicas),
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"app": "pcf"},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{"app": "pcf"},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "pcf",
								Image: free5gc.Spec.PCF.Image,
								Ports: []corev1.ContainerPort{
									{
										ContainerPort: int32(free5gc.Spec.PCF.Network.Port),
									},
								},
							},
						},
					},
				},
			},
		}

		if err := controllerutil.SetControllerReference(free5gc, PCFDeployment, r.Scheme); err != nil {
			return err
		}
		return r.Create(ctx, PCFDeployment)
	} else if err != nil {
		return err
	}

	logger.Info("PCF Deployment already exists, skipping creation")
	return nil
}

// Déploiement du service NRF
func (r *Free5gcReconciler) deployNRF(ctx context.Context, free5gc *k8sstartkubernetescomv1.Free5gc) error {
	logger := log.FromContext(ctx)

	deploymentName := fmt.Sprintf("%s-nrf", free5gc.Name)
	foundDeployment := &appsv1.Deployment{}
	err := r.Get(ctx, client.ObjectKey{Namespace: free5gc.Namespace, Name: deploymentName}, foundDeployment)
	if err != nil && errors.IsNotFound(err) {
		logger.Info("Creating new NRF deployment", "DeploymentName", deploymentName)

		NRFDeployment := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      deploymentName,
				Namespace: free5gc.Namespace,
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: int32Ptr(free5gc.Spec.NRF.Replicas),
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"app": "nrf"},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{"app": "nrf"},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "nrf",
								Image: free5gc.Spec.NRF.Image,
								Ports: []corev1.ContainerPort{
									{
										ContainerPort: int32(free5gc.Spec.NRF.Network.Port),
									},
								},
							},
						},
					},
				},
			},
		}

		if err := controllerutil.SetControllerReference(free5gc, NRFDeployment, r.Scheme); err != nil {
			return err
		}
		return r.Create(ctx, NRFDeployment)
	} else if err != nil {
		return err
	}

	logger.Info("NRF Deployment already exists, skipping creation")
	return nil
}

// Déploiement du service UPF
func (r *Free5gcReconciler) deployUPF(ctx context.Context, free5gc *k8sstartkubernetescomv1.Free5gc) error {
	logger := log.FromContext(ctx)

	deploymentName := fmt.Sprintf("%s-upf", free5gc.Name)
	foundDeployment := &appsv1.Deployment{}
	err := r.Get(ctx, client.ObjectKey{Namespace: free5gc.Namespace, Name: deploymentName}, foundDeployment)
	if err != nil && errors.IsNotFound(err) {
		logger.Info("Creating new UPF deployment", "DeploymentName", deploymentName)

		UPFDeployment := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      deploymentName,
				Namespace: free5gc.Namespace,
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: int32Ptr(free5gc.Spec.UPF.Replicas),
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"app": "upf"},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{"app": "upf"},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "upf",
								Image: free5gc.Spec.UPF.Image,
								Ports: []corev1.ContainerPort{
									{
										ContainerPort: int32(free5gc.Spec.UPF.Network.Port),
									},
								},
							},
						},
					},
				},
			},
		}

		if err := controllerutil.SetControllerReference(free5gc, UPFDeployment, r.Scheme); err != nil {
			return err
		}
		return r.Create(ctx, UPFDeployment)
	} else if err != nil {
		return err
	}

	logger.Info("UPF Deployment already exists, skipping creation")
	return nil
}

// Déploiement du service UDM
func (r *Free5gcReconciler) deployUDM(ctx context.Context, free5gc *k8sstartkubernetescomv1.Free5gc) error {
	logger := log.FromContext(ctx)

	deploymentName := fmt.Sprintf("%s-udm", free5gc.Name)
	foundDeployment := &appsv1.Deployment{}
	err := r.Get(ctx, client.ObjectKey{Namespace: free5gc.Namespace, Name: deploymentName}, foundDeployment)
	if err != nil && errors.IsNotFound(err) {
		logger.Info("Creating new UDM deployment", "DeploymentName", deploymentName)

		UDMDeployment := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      deploymentName,
				Namespace: free5gc.Namespace,
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: int32Ptr(free5gc.Spec.UDM.Replicas),
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"app": "udm"},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{"app": "udm"},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "udm",
								Image: free5gc.Spec.UDM.Image,
								Ports: []corev1.ContainerPort{
									{
										ContainerPort: int32(free5gc.Spec.UDM.Network.Port),
									},
								},
							},
						},
					},
				},
			},
		}

		if err := controllerutil.SetControllerReference(free5gc, UDMDeployment, r.Scheme); err != nil {
			return err
		}
		return r.Create(ctx, UDMDeployment)
	} else if err != nil {
		return err
	}

	logger.Info("UDM Deployment already exists, skipping creation")
	return nil
}

// Déploiement du service UDR
func (r *Free5gcReconciler) deployUDR(ctx context.Context, free5gc *k8sstartkubernetescomv1.Free5gc) error {
	logger := log.FromContext(ctx)

	deploymentName := fmt.Sprintf("%s-udr", free5gc.Name)
	foundDeployment := &appsv1.Deployment{}
	err := r.Get(ctx, client.ObjectKey{Namespace: free5gc.Namespace, Name: deploymentName}, foundDeployment)
	if err != nil && errors.IsNotFound(err) {
		logger.Info("Creating new UDR deployment", "DeploymentName", deploymentName)

		UDRDeployment := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      deploymentName,
				Namespace: free5gc.Namespace,
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: int32Ptr(free5gc.Spec.UDR.Replicas),
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"app": "udr"},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{"app": "udr"},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "udr",
								Image: free5gc.Spec.UDR.Image,
								Ports: []corev1.ContainerPort{
									{
										ContainerPort: int32(free5gc.Spec.UDR.Network.Port),
									},
								},
							},
						},
					},
				},
			},
		}

		if err := controllerutil.SetControllerReference(free5gc, UDRDeployment, r.Scheme); err != nil {
			return err
		}
		return r.Create(ctx, UDRDeployment)
	} else if err != nil {
		return err
	}

	logger.Info("UDR Deployment already exists, skipping creation")
	return nil
}

// Déploiement du service NSSF
func (r *Free5gcReconciler) deployNSSF(ctx context.Context, free5gc *k8sstartkubernetescomv1.Free5gc) error {
	logger := log.FromContext(ctx)

	deploymentName := fmt.Sprintf("%s-nssf", free5gc.Name)
	foundDeployment := &appsv1.Deployment{}
	err := r.Get(ctx, client.ObjectKey{Namespace: free5gc.Namespace, Name: deploymentName}, foundDeployment)
	if err != nil && errors.IsNotFound(err) {
		logger.Info("Creating new NSSF deployment", "DeploymentName", deploymentName)

		NSSFDeployment := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      deploymentName,
				Namespace: free5gc.Namespace,
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: int32Ptr(free5gc.Spec.NSSF.Replicas),
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"app": "nssf"},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{"app": "nssf"},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "nssf",
								Image: free5gc.Spec.NSSF.Image,
								Ports: []corev1.ContainerPort{
									{
										ContainerPort: int32(free5gc.Spec.NSSF.Network.Port),
									},
								},
							},
						},
					},
				},
			},
		}

		if err := controllerutil.SetControllerReference(free5gc, NSSFDeployment, r.Scheme); err != nil {
			return err
		}
		return r.Create(ctx, NSSFDeployment)
	} else if err != nil {
		return err
	}

	logger.Info("NSSF Deployment already exists, skipping creation")
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *Free5gcReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&k8sstartkubernetescomv1.Free5gc{}).
		Owns(&appsv1.Deployment{}).
		Named("free5gc").
		Complete(r)
}
