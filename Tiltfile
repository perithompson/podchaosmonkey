# -*- mode: Python -*-

# For more on Extension
load('ext://restart_process', 'docker_build_with_restart')

IMG="localhost:50949/podchaosmonkey:latest"

def yaml():
  return local('cd config/manager; kustomize edit set image controller=' + IMG + '; cd ../..' + '; kustomize build config/default')

def set_image():
  return local('cd config/manager && kustomize edit set image controller={}'.format(IMG))

docker_build(
  ref='localhost:50949/podchaosmonkey:latest',
  context='.',
  dockerfile='./Dockerfile',
  live_update=[sync('./bin/manager', '/')],
  entrypoint=['/manager'],
)

compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o ./bin/manager main.go'

local_resource(
  'monkey-go-compile',
  compile_cmd,
  deps=['./main.go', './api/v1alpha1/monkey_types.go', './controllers', './internal', './pkg'],
)

set_image()
k8s_yaml(kustomize('config/default'))