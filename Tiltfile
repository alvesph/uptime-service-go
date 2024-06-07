allow_k8s_contexts('k3s')

local_resource(
  'uptime-sercice',
  cmd='CGO_ENABLED=0 GOOS=linux go build -o ./main',
  deps=['main.go'],
  labels=['uptime-service-api']
)

docker_build(
  'recrutaz/uptime-service-api',
  './',
  dockerfile='./Dockerfile'
)

k8s_yaml('k8s/deployment.yaml')