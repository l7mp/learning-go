// Package ci runs the lab tests against a real Kubernetes cluster.
//
// The harness owns the whole run: it builds the students' applications into container
// images, brings up a k3d cluster, imports the images into it, runs every lab in turn, and
// tears the cluster down again. Nothing outside "go test" is needed, so a student can run
// the same checks the CI runs.
package ci

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/envfuncs"
	"sigs.k8s.io/e2e-framework/support"
	"sigs.k8s.io/e2e-framework/support/k3d"

	"labci/pkg/image"
	"labci/pkg/labenv"
)

// clusterOpts shape the cluster to look like the Minikube the labs are written for.
//
// Traefik has to go. k3d installs it by default, and its ServiceLB pod claims host port 80
// for the whole life of the cluster; the SplitDim Service the labs ask students to write
// asks for port 80 too, so its own ServiceLB pod could never be scheduled and the Service
// would wait for an external address that is never coming.
var clusterOpts = []support.ClusterOpts{
	k3d.WithArgs("--k3s-arg", "--disable=traefik@server:*"),
}

// clusterName is the k3d cluster the labs share. k3d reuses a cluster of this name if one
// is already running, so repeated local runs do not pay for a fresh cluster every time.
const clusterName = "labs"

// labApps are the applications the labs package into images. resilient and transactionlog
// are libraries these vendor in, never images of their own.
var labApps = []string{"helloworld", "splitdim", "kvstore"}

var testenv env.Environment

func TestMain(m *testing.M) {
	// With a cluster of your own, KUBECONFIG selects it and you have loaded the lab
	// images into it yourself, so the harness has nothing left to do.
	if os.Getenv("LABS_EXTERNAL_CLUSTER") != "" {
		fmt.Fprintf(os.Stderr, "LABS_EXTERNAL_CLUSTER is set: testing against the cluster "+
			"KUBECONFIG names, building no images and creating no cluster\n")
		os.Exit(m.Run())
	}

	root, err := labenv.FindRoot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	code := filepath.Join(root, "99-labs", "code")

	apps := buildable(code)
	if len(apps) == 0 {
		fmt.Fprintf(os.Stderr, "none of the lab applications can be built yet: none of "+
			"%v under %s has both a go.mod and a deploy/Dockerfile\n", labApps, code)
		os.Exit(1)
	}

	archives, err := os.MkdirTemp("", "labs-images-*")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	setup := []env.Func{
		buildImages(code, apps, archives),
		envfuncs.CreateClusterWithOpts(k3d.NewProvider(), clusterName, clusterOpts...),
		exportKubeconfig,
	}
	for _, app := range apps {
		setup = append(setup,
			envfuncs.LoadImageArchiveToCluster(clusterName, filepath.Join(archives, app+".tar")))
	}

	finish := []env.Func{
		func(ctx context.Context, _ *envconf.Config) (context.Context, error) {
			return ctx, os.RemoveAll(archives)
		},
	}

	// Keeping the cluster skips the single largest fixed cost of the next run, which is
	// worth having while working on one lab. CI wants the machine left clean, so tearing
	// down is the default.
	if os.Getenv("LABS_KEEP_CLUSTER") == "" {
		finish = append(finish, envfuncs.DestroyCluster(clusterName))
	} else {
		fmt.Fprintf(os.Stderr, "LABS_KEEP_CLUSTER is set: leaving the %q cluster running, "+
			"remove it with \"k3d cluster delete %s\"\n", clusterName, clusterName)
	}

	testenv = env.New()
	testenv.Setup(setup...)
	testenv.Finish(finish...)

	os.Exit(testenv.Run(m))
}

// buildable returns the lab applications that can be built into an image: the ones the
// student has already started, having both a module and a Dockerfile.
//
// Skipping what is not there yet is what lets somebody halfway through the course test the
// lab they are on. It does not need to know which lab that is: the code is cumulative, so
// whatever builds today also serves the earlier labs.
func buildable(code string) []string {
	apps := []string{}

	for _, app := range labApps {
		if _, err := os.Stat(filepath.Join(code, app, "go.mod")); err != nil {
			continue
		}
		if _, err := os.Stat(filepath.Join(code, app, "deploy", "Dockerfile")); err != nil {
			continue
		}

		apps = append(apps, app)
	}

	return apps
}

// buildImages builds each application from the student's own Dockerfile and saves it as an
// archive the cluster can import.
func buildImages(code string, apps []string, archives string) env.Func {
	return func(ctx context.Context, _ *envconf.Config) (context.Context, error) {
		b, err := image.NewBuilder()
		if err != nil {
			return ctx, err
		}
		defer b.Close()

		if err := b.Ping(ctx); err != nil {
			return ctx, err
		}

		for _, app := range apps {
			fmt.Fprintf(os.Stderr, "building localhost/%s:latest\n", app)

			dir, cleanup, err := image.Vendor(ctx, code, app)
			if err != nil {
				return ctx, err
			}

			err = b.Build(ctx, dir, "deploy/Dockerfile", "localhost/"+app+":latest",
				filepath.Join(archives, app+".tar"))
			cleanup()

			if err != nil {
				return ctx, err
			}
		}

		return ctx, nil
	}
}

// exportKubeconfig publishes the cluster's kubeconfig so that the lab tests, and the
// students' own test suites the labs run as subprocesses, all talk to it.
func exportKubeconfig(ctx context.Context, cfg *envconf.Config) (context.Context, error) {
	return ctx, os.Setenv("KUBECONFIG", cfg.KubeconfigFile())
}
