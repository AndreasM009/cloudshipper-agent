package runner

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"k8s.io/apimachinery/pkg/labels"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
)

// K8sRunnerPod pod for runner
type K8sRunnerPod struct {
	errorChan chan error
}

// CreateAndWatchRunnerPod test
func CreateAndWatchRunnerPod(ctx context.Context, k8sclient *kubernetes.Clientset, natsChannelName string) (*K8sRunnerPod, error) {
	r := &K8sRunnerPod{
		errorChan: make(chan error, 1),
	}

	name := fmt.Sprintf("runner-%s", natsChannelName)

	pod := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Pod",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: map[string]string{"runner": natsChannelName},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				corev1.Container{
					Name:            "runnercntr",
					Image:           "m009/cs-agent-runner:latest",
					ImagePullPolicy: corev1.PullAlways,
					Env: []corev1.EnvVar{
						corev1.EnvVar{
							Name:  "ARTIFACTS_DIRECTORY",
							Value: "/artifacts",
						},
					},
					Command: []string{"./runner"},
					Args: []string{
						fmt.Sprintf("-s=%s", "nats://example-nats.default.svc.cluster.local:4222"),
						fmt.Sprintf("-c=%s", natsChannelName),
					},
				},
			},
			RestartPolicy: corev1.RestartPolicyNever,
		},
	}

	resp, err := k8sclient.CoreV1().Pods("cs-agent").Create(ctx, pod, metav1.CreateOptions{})
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("RunnerPod status:", resp.Status.Phase)

	w, err := k8sclient.CoreV1().Pods("cs-agent").Watch(ctx, metav1.ListOptions{
		Watch:           true,
		ResourceVersion: resp.ResourceVersion,
		FieldSelector:   fields.SelectorFromSet(fields.Set{"metadata.name": name}).String(),
		LabelSelector:   labels.SelectorFromSet(labels.Set{"runner": natsChannelName}).String(),
	})

	if err != nil {
		log.Panic(err)
	}

	deletePod := func() {
		k8sclient.CoreV1().Pods("cs-agent").Delete(context.Background(), name, metav1.DeleteOptions{})
	}

	loop := func() {
		for {
			select {
			case <-ctx.Done():
				w.Stop()
				deletePod()
				r.errorChan <- nil
				return
			case events, ok := <-w.ResultChan():
				if ok {
					if resp, ok := events.Object.(*corev1.Pod); ok {
						if resp.Status.Phase == corev1.PodFailed {
							w.Stop()
							deletePod()
							r.errorChan <- errors.New("runner pod failed")
							return
						}
					}
				}
			}
		}
	}

	// wait until pod is running
	for {
		select {
		case <-ctx.Done():
			deletePod()
			return nil, nil
		case events, ok := <-w.ResultChan():
			if !ok {
				deletePod()
				log.Panic()
			}

			resp := events.Object.(*corev1.Pod)
			fmt.Println("RunnerPod status:", resp.Status.Phase)

			if resp.Status.Phase == corev1.PodFailed {
				w.Stop()
				deletePod()
				return nil, errors.New("RunnerPod failed")
			} else if resp.Status.Phase == corev1.PodSucceeded {
				w.Stop()
				return r, nil
			} else if resp.Status.Phase == corev1.PodRunning {
				go loop()
				return r, nil
			}

		case <-time.After(20 * time.Second):
			fmt.Println("timeout to wait for runner pod active")
			w.Stop()
			deletePod()
			return nil, errors.New("timeout wait for runner pod active")
		}
	}
}

// Done channel
func (r *K8sRunnerPod) Error() <-chan error {
	return r.errorChan
}
