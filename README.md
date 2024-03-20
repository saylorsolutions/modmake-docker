# Modmake Docker Plugin

[![Go Report Card](https://goreportcard.com/badge/github.com/saylorsolutions/modmake-docker)](https://goreportcard.com/report/github.com/saylorsolutions/modmake-docker)

This repo provides some base functionality for building, running, and interacting with Docker images and containers.
If you're not familiar with [Go](https://go.dev/), [Docker](https://www.docker.com/), or [Modmake](https://saylorsolutions.github.io/modmake/), then this probably isn't for you yet.
Check out those first and circle back here to find out how this can help you.

## API Docs

The API docs are generously hosted here:
https://pkg.go.dev/github.com/saylorsolutions/modmake-docker

## How to get it

Just like Modmake itself, this plugin can be added to your project as a dependency. 

> The Docker CLI is required to be available on the PATH to use this plugin.

```bash
go get github.com/saylorsolutions/modmake-docker@latest
```

Once that's done, this can be incorporated into any new or existing Modmake build.

## Reporting Bugs or Recommendations

If you find a bug in this project or have an idea for improvement, please [report an issue](https://github.com/saylorsolutions/modmake-docker/issues/new/choose) so I can get it fixed/added.

# Disclaimers and such

This project is ***NOT*** endorsed by or produced in conjunction with the owners, operators, project maintainers, directors, or other entities associated with Docker and related tooling or systems.

I do ***NOT*** own - in *any* capacity - the right to use or profit from Docker trademarks or branding.

This is ***NOT*** intended to be fully representative of the offering or capability of Docker and related tooling or systems in any way.

The Docker owners, operators, project maintainers, directors, or other entities associated with Docker and related tooling or systems are and continue to be in perpetuity, free to change the CLI interface on which this plugin depends in any way, with or without corresponding changes in this repository or documentation, and without notice.

This is *just* a Modmake-specific way to access the Docker CLI interface. The two projects are ***NOT*** related in any official way.

For official Docker tooling and information, see their website:
https://www.docker.com/

# Basic Usage

Building images and running containers are the most basic operations, and both are supported by this plugin.

Here's a step-by-step guide for basic usage.
There are many ways this can be customized to fit your needs.

### Build Example

Given a working tree like this.
```
.
├── cmd
│   └── my-app
│       └── main.go
├── go.mod
├── go.sum
└── modmake
    └── build.go
```

And an existing modmake build file in `./modmake/build.go`.

```go
package main

import (
	. "github.com/saylorsolutions/modmake"
)

func main() {
	b := NewBuild()

	b.Test().Does(Go().TestAll())
	b.Build().Does(Go().Build("./cmd/my-app"))

	b.Execute()
}
```

We can add a Dockerfile to the root of the project like this, allowing us to run our Go executable in a Docker container.

> This isn't an example of the *best* way to run Go applications in container, but it works for a simple example.

```dockerfile
FROM golang:1.22

WORKDIR /app
COPY . .
# Run the Modmake 'build' step
RUN go run ./modmake build

ENTRYPOINT ./my-app
```

Next, we need to update our Modmake build to include the step to build our image.

```diff
 
 import (
 	. "github.com/saylorsolutions/modmake"
+	. "github.com/saylorsolutions/modmake-docker"
 )
 
 func main() {
 	b := NewBuild()
 
+	// Make sure dependencies are present in container.
+	b.Generate().Does(Go().ModTidy())
 	b.Test().Does(Go().TestAll())
 	b.Build().Does(Go().Build("./cmd/my-app"))
 
+	pkgDocker := NewStep("package-docker", "Packages a docker image ot test-project").
+		Does(Docker().Build("test-project:latest", ""))
+	b.AddStep(pkgDocker)
+
 	b.Execute()
 }
```

A call to `go mod tidy` was added with the default `generate` step to make sure the required dependencies are available for `my-app` during image build.

Once this is done, we can run the new `package-docker` step to build our Docker image with this command.
```bash
go run ./modmake package-docker
```

### Run Example

Working off of the build example, we can add another step to run the built image.
A dependency is expressed on the `package-docker` step to make sure it's always using the latest version of the image.

```diff
 		Does(Docker().Build("test-project:latest", ""))
 	b.AddStep(pkgDocker)
 
+	runDocker := NewStep("run-docker", "Runs the Docker container built from package-docker").
+		Does(Docker().Run("test-project:latest").RemoveAfterExit()).
+		// Make this depend on 'run-docker'
+		DependsOn(pkgDocker)
+	b.AddStep(runDocker)
+
 	b.Execute()
 }
```

We can run the new step with this command.

```bash
go run ./modmake run-docker
```

We'll see the image build output too since there's a dependency on the `package-docker` step.

### Composable Sub-Build

There's a lot more that can be done with this foundation, including much more interesting configuration options for both `Docker().Build` and `Docker().Run`.

One big improvement over the previous examples would be to add Docker steps to their own sub-build.
This provides an easy demarcation point for what are Docker-specific steps, and what are the base steps of the build.
The Modmake build script below will create the sub-build and import it into the parent build under the `docker` heading.

```golang
package main

import (
 	. "github.com/saylorsolutions/modmake"
	. "github.com/saylorsolutions/modmake-docker"
)

func main() {
	// It's good practice to add variables that make the build easily configurable.
	var (
		// Uses Modmake's PathString for portable paths.
		// All Modmake PathStrings expect forward slashes regardless of OS.
		appPath = Path("./cmd/my-app")
		// Uses Modmake's F strings for environment variable interpolation.
		// This will default to 'test-project:local' if 'IMAGE_TAG' is empty/undefined.
		imageName = F("test-project:${IMAGE_TAG:local}")
	)

	b := NewBuild()

	b.Test().Does(Go().TestAll())
	// Build the application
	b.Build().Does(Go().Build(appPath))

	// Create a Docker sub-build.
	docker := NewBuild()
	// Do 'go mod tidy' to ensure dependencies are available during image build.
	docker.Generate().Does(Go().ModTidy())
	// Build the image using the default 'build' step.
	docker.Build().Does(Docker().Build(imageName, ""))
	// Run the image after building it with a custom step.
	docker.NewStep("run", "Runs my-app in container").
		Does(Docker().Run(imageName).RemoveAfterExit()).
		DependsOn(docker.Step("build"))
	// Import the sub-build into the base build with the prefix 'docker:'.
	b.Import("docker", docker)
	// Make sure that tests pass before building the image.
	// This also means that tests must pass before running 'docker:run'
	b.Step("docker:build").DependsOn("test")

	b.Execute()
}
```

This provides a neater reference pattern for Docker steps.
Here's a table mapping the previous step name to the new step using the `docker` sub-build.

| Previous Step | Becomes |
| :-- | :-- |
| `package-docker` | `docker:build` |
| `run-docker` | `docker:run` |

## Summary

There's a lot that Modmake and this plugin can do.
Consider using it in your own builds, and let me know how it works out. 😊

If you find a bug in modmake-docker, please [report an issue](https://github.com/saylorsolutions/modmake-docker/issues/new/choose) so I can get it fixed.
