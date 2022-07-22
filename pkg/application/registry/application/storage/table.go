/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package storage

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metav1beta1 "k8s.io/apimachinery/pkg/apis/meta/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"tkestack.io/tke/api/application"
	"tkestack.io/tke/pkg/util/printers"
)

// AddHandlers adds print handlers for default TKE types dealing with internal versions.
// Refer kubernetes/pkg/printers/internalversion/printers.go:78
func AddHandlers(h printers.PrintHandler) {
	appColumnDefinitions := []metav1beta1.TableColumnDefinition{
		{Name: "Name", Type: "string", Format: "name", Description: "Name must be unique within a namespace. Is required when creating resources, although some resources may allow a client to request the generation of an appropriate name automatically. Name is primarily intended for creation idempotence and configuration definition. Cannot be updated. More info: http://kubernetes.io/docs/user-guide/identifiers#names"},
		{Name: "ChartName", Type: "string", Description: "ChartName is the name of the chart."},
		{Name: "ChartVersion", Type: "string", Description: "ChartVersion is the version of the chart."},
		{Name: "Status", Type: "string", Description: "Phase the release is in, one of ('ChartFetched', 'ChartFetchFailed', 'Installing', 'Upgrading', 'Succeeded', 'RollingBack', 'RolledBack', 'RollbackFailed')"},
		{Name: "Age", Type: "string", Description: "CreationTimestamp is a timestamp representing the server time when this object was created. It is not guaranteed to be set in happens-before order across separate operations. Clients may not set this value. It is represented in RFC3339 form and is in UTC.\n\nPopulated by the system. Read-only. Null for lists. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata"},
	}
	h.TableHandler(appColumnDefinitions, printAppList)
	h.TableHandler(appColumnDefinitions, printApp)
}

func printAppList(appList *application.AppList, options printers.PrintOptions) ([]metav1.TableRow, error) {
	rows := make([]metav1.TableRow, 0, len(appList.Items))
	for i := range appList.Items {
		r, err := printApp(&appList.Items[i], options)
		if err != nil {
			return nil, err
		}
		rows = append(rows, r...)
	}
	return rows, nil
}

func printApp(app *application.App, options printers.PrintOptions) ([]metav1.TableRow, error) {
	row := metav1.TableRow{
		Object: runtime.RawExtension{Object: app},
	}
	row.Cells = append(row.Cells, app.Name, app.Spec.Chart.ChartName, app.Spec.Chart.ChartVersion, app.Status.Phase, printers.TranslateTimestampSince(app.CreationTimestamp))
	return []metav1beta1.TableRow{row}, nil
}
