/*
Copyright 2022 dgsfor@gmail.com.

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

package controllers

import (
	"context"
	"encoding/json"
	"hwcloud-aom-operator/hwsinger"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	itmonkeyv1 "hwcloud-aom-operator/api/v1"
)

// AomGroupReconciler reconciles a AomGroup object
type AomGroupReconciler struct {
	client.Client
	Scheme  *runtime.Scheme
	Eventer record.EventRecorder
}

//+kubebuilder:rbac:groups=itmonkey.itmonkey.icu,resources=aomgroups,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=itmonkey.itmonkey.icu,resources=aomgroups/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=itmonkey.itmonkey.icu,resources=aomgroups/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the AomGroup object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.2/pkg/reconcile
func (r *AomGroupReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	// TODO(user): your logic here
	aomGroup := &itmonkeyv1.AomGroup{}
	// 1. 获取aomGroup
	if err := r.Get(ctx, req.NamespacedName, aomGroup); err != nil {
		r.Eventer.Eventf(aomGroup, corev1.EventTypeWarning, "Warning", "get aom group error")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// 2. 获取对应的deployment
	d := &v1.Deployment{}
	if err := r.Get(ctx, client.ObjectKey{Namespace: req.Namespace, Name: aomGroup.Spec.Deployment}, d); err != nil {
		if errors.IsNotFound(err) {
			r.Eventer.Eventf(aomGroup, corev1.EventTypeWarning, "Warning", "the deployment not exist")
			return ctrl.Result{}, err
		}
	}
	r.Eventer.Eventf(aomGroup, corev1.EventTypeNormal, "Normal", "binding deployment success")
	logger.Info("binding deployment success")

	// 3. set aom group reference
	if err := controllerutil.SetControllerReference(d, aomGroup, r.Scheme); err != nil {
		logger.Error(err, "set controller reference failed")
		r.Eventer.Eventf(aomGroup, corev1.EventTypeWarning, "Warning", "set controller reference failed")
		return ctrl.Result{}, err
	}
	r.Eventer.Eventf(aomGroup, corev1.EventTypeNormal, "Normal", "set controller reference success")
	logger.Info("set controller reference success")

	// 4. 获取该deployment对应的aom group配置
	aomGroupResult := hwsinger.GetAomGroupWithDeployment(req.Namespace, aomGroup.Spec.Deployment)
	result := new(hwsinger.AomGroupResult)
	// 如果格式化失败，首先打印结果信息，然后记录事件
	if err := json.Unmarshal(aomGroupResult.Data, &result); err != nil {
		logger.Error(err, string(aomGroupResult.Data))
		r.Eventer.Eventf(aomGroup, corev1.EventTypeWarning, "Warning", err.Error())
		return ctrl.Result{}, err
	}
	// 如果请求aom group数据错误，记录事件
	if aomGroupResult.Code != 200 {
		r.Eventer.Eventf(aomGroup, corev1.EventTypeWarning, "Warning", result.ErrorMessage)
		return ctrl.Result{}, nil
	}
	r.Eventer.Eventf(aomGroup, corev1.EventTypeNormal, "Normal", "get aom group data from hw success")
	logger.Info("get aom group data from hw success")

	// 5. 更新aom group cr，填充id、group_attribute等信息
	aomGroup.Spec.Id = result.Config.ID
	aomGroup.Spec.GroupAttribute.CooldownTime = result.Config.CooldownTime
	aomGroup.Spec.GroupAttribute.MaxInstances = result.Config.MaxInstances
	aomGroup.Spec.GroupAttribute.MinInstances = result.Config.MinInstances
	if err := r.Client.Update(ctx, aomGroup); err != nil {
		logger.Error(err, "update cr failed")
		r.Eventer.Eventf(aomGroup, corev1.EventTypeWarning, "Warning", "update cr failed")
		return ctrl.Result{}, err
	}
	r.Eventer.Eventf(aomGroup, corev1.EventTypeNormal, "Normal", "update cr success")
	logger.Info("update cr success")

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AomGroupReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&itmonkeyv1.AomGroup{}).
		Complete(r)
}
